package main

import (
	"context"
	"encoding/json"
	clientdelivery "eshop/internal/clients/delivery"
	clientpayment "eshop/internal/clients/payment"
	clientstore "eshop/internal/clients/store"
	dbClient "eshop/internal/db"
	orderDB "eshop/internal/db/order"
	response "eshop/internal/request"
	requestorder "eshop/internal/request/order"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"
)

const (
	dbName  = "Order"
	colName = "Order"
)

var db *orderDB.DB

var clientStore *clientstore.Client
var clientDelivery *clientdelivery.Client
var clientPayment *clientpayment.Client

func main() {
	uri, ok := os.LookupEnv("DB")
	if !ok {
		slog.Error("no db uri")
		os.Exit(1)
	}
	slog.Info("DB", "uri", uri)
	dbClient := dbClient.New(uri, dbName, colName)
	db = orderDB.New(dbClient)

	uriStore, ok := os.LookupEnv("STORE")
	if !ok {
		slog.Error("no STORE uri")
		os.Exit(1)
	}
	clientStore = clientstore.NewClient(uriStore)

	uriDelivery, ok := os.LookupEnv("DELIVERY")
	if !ok {
		slog.Error("no DELIVERY uri")
		os.Exit(1)
	}
	clientDelivery = clientdelivery.NewClient(uriDelivery)

	uriPay, ok := os.LookupEnv("PAY")
	if !ok {
		slog.Error("no PAY uri")
		os.Exit(1)
	}
	clientPayment = clientpayment.NewClient(uriPay)

	muxRouter := mux.NewRouter()
	muxRouter.HandleFunc("/", index).Methods(http.MethodGet)
	muxRouter.HandleFunc("/create", create).Methods(http.MethodPost)

	srv := &http.Server{
		Handler:      muxRouter,
		Addr:         "0.0.0.0:8001",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}

func index(w http.ResponseWriter, req *http.Request) {
	_, err := fmt.Fprint(w, "I'm service order")
	if err != nil {
		slog.Error("index write")
	}
}

func create(w http.ResponseWriter, r *http.Request) {
	slog.Info("call create order")
	data, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var order requestorder.RequestOrder
	err = json.Unmarshal(data, &order)
	if err != nil {
		slog.Error("unmarshal create order", "err", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// логика
	orderId, orderDBData, err := db.AddNewOrder(r.Context(), order.Items)
	if err != nil {
		slog.Error("addNewOrder db", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := db.AddLog(r.Context(), orderId, fmt.Sprintf("order created")); err != nil {
		slog.Error("Save log 'created'", "err", err)
	}

	ok, msg, err := clientStore.Reserve(orderDBData.Items)
	if err != nil {
		slog.Error("Save log reserve", "err", err)
		// не получилось зарезервировать, ошибка сети или 500
		if err := db.AddLog(r.Context(), orderId, fmt.Sprintf("reserving fail: %s", err.Error())); err != nil {
			slog.Error("Save log reserve", "err", err)
		}
		if err := db.ChangeStatus(r.Context(), orderId, orderDB.Status_Cancel); err != nil {
			slog.Error("Change status to cancel fail", "err", err)
		}
		response.Common{
			Status:  false,
			Message: fmt.Sprintf("Not reserve, order cancel"),
		}.WriteResponse(w)
		return
	}

	if !ok {
		// не получилось зарезервировать, логическая ошибка
		if err := db.AddLog(r.Context(), orderId, fmt.Sprintf("reserving fail: %s", msg)); err != nil {
			slog.Error("Save log 'not reserve'", "err", err)
		}
		if err := db.ChangeStatus(r.Context(), orderId, orderDB.Status_Cancel); err != nil {
			slog.Error("Change status to cancel", "err", err)
		}
		response.Common{
			Status:  false,
			Message: msg,
		}.WriteResponse(w)
		return
	}

	// зарезервировали товар

	if err := db.ChangeStatus(r.Context(), orderId, orderDB.Status_ReservedItems); err != nil {
		slog.Error("Change status to reserved", "err", err)
	}

	// бронируем слот для доставки
	if err := db.AddLog(r.Context(), orderId, fmt.Sprintf("reserving delivery slot")); err != nil {
		slog.Error("Save log 'reserving delivery slot'", "err", err)
	}
	reserve, err := clientDelivery.Reserve()
	if err != nil {
		if err := db.AddLog(r.Context(), orderId, fmt.Sprintf("reserving delivery slot faild")); err != nil {
			slog.Error("Save log 'reserving delivery slot faild'", "err", err)
		}
		Rollback(orderId, w)
		return
	}

	if reserve.Status == false {
		if err := db.AddLog(r.Context(), orderId, fmt.Sprintf("reserving delivery slot faild: %s", reserve.Message)); err != nil {
			slog.Error("Save log 'reserving delivery slot faild'", "err", err)
		}
		Rollback(orderId, w)
		return
	}

	// успешно зарезервировали слот доставки
	if err := db.ChangeStatus(r.Context(), orderId, orderDB.Status_ReserveDeliverSlot); err != nil {
		slog.Error("Change status to reserveDelivery", "err", err)
	}
	if err = db.SetSlot(r.Context(), orderId, reserve.Slot); err != nil {
		slog.Error("saving slot to order", "err", err)
	}

	// оплачиваем
	_, err = clientPayment.Pay()
	if err != nil {
		if err := db.AddLog(r.Context(), orderId, fmt.Sprintf("pay fail")); err != nil {
			slog.Error("Save log 'pay fail'", "err", err)
		}
		Rollback(orderId, w)
		return
	}

	// успешно оплатили
	if err := db.AddLog(r.Context(), orderId, fmt.Sprintf("pay success")); err != nil {
		slog.Error("Save log 'pay success'", "err", err)
	}
	if err := db.ChangeStatus(r.Context(), orderId, orderDB.Status_Paid); err != nil {
		slog.Error("Change status to paid", "err", err)
	}

	var logOrder []string
	getOrder, _ := db.GetOrder(r.Context(), orderId)
	for _, s := range getOrder.Log {
		logOrder = append(logOrder, s)
	}
	response.Common{
		Status:  true,
		Message: fmt.Sprintf("Order created :%s", orderId),
		Log:     logOrder,
	}.WriteResponse(w)

}

// Rollback получает заказ с текущим статусом и откатывает его в зад
func Rollback(orderId string, w http.ResponseWriter) {
	ctx := context.TODO()
	order, err := db.GetOrder(ctx, orderId)
	if err != nil {
		slog.Error("Rollback get order", "err", err)
		return
	}

	if order.Status == orderDB.Status_ReserveDeliverSlot {
		err := db.AddLog(ctx, orderId, "try rollback delivery slot")
		if err != nil {
			slog.Error("save logOrder rollback delivery slot", "err", err)
		}

		err = clientDelivery.Rollback(order.Slot)
		if err != nil {
			db.AddLog(ctx, orderId, fmt.Sprintf("try rollback delivery slot fail: %s", err.Error()))
			return
		}
		if err = db.ChangeStatus(ctx, orderId, orderDB.Status_ReservedItems); err != nil {
			slog.Error("rollback change slot", "err", err)
			return
		}
		order.Status = orderDB.Status_ReservedItems
	}

	if order.Status == orderDB.Status_ReservedItems {
		err := db.AddLog(ctx, orderId, "try rollback reserve items")
		if err != nil {
			slog.Error("save logOrder rollback reserve items", "err", err)
		}
		err = clientStore.Rollback(order.Items)
		if err != nil {
			db.AddLog(ctx, orderId, fmt.Sprintf("try rollback reserve items fail: %s", err.Error()))
		}
		db.ChangeStatus(ctx, orderId, orderDB.Status_Cancel)
		order.Status = orderDB.Status_Cancel
	}

	var logOrder []string
	getOrder, _ := db.GetOrder(ctx, orderId)
	for _, s := range getOrder.Log {
		logOrder = append(logOrder, s)
	}
	response.Common{
		Status:  false,
		Message: fmt.Sprintf("order canceled"),
		Log:     logOrder,
	}.WriteResponse(w)
}

package main

import (
	"encoding/json"
	"errors"
	dbClient "eshop/internal/db"
	deliveryDB "eshop/internal/db/delivery"
	response "eshop/internal/request"
	deliveryRequest "eshop/internal/request/delivery"
	"fmt"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"
)

const (
	dbName  = "Delivery"
	colName = "Delivery"
)

var db *deliveryDB.DB

func main() {
	uri, ok := os.LookupEnv("DB")
	if !ok {
		slog.Error("no db uri")
		os.Exit(1)
	}
	slog.Info("DB", "uri", uri)

	dbClient := dbClient.New(uri, dbName, colName)

	db = deliveryDB.New(dbClient)

	muxRouter := mux.NewRouter()
	muxRouter.HandleFunc("/", index).Methods(http.MethodGet)
	muxRouter.HandleFunc("/reserve_slot", reserveSlot).Methods(http.MethodPost)
	muxRouter.HandleFunc("/rollback", rollback).Methods(http.MethodPost)

	srv := &http.Server{
		Handler:      muxRouter,
		Addr:         "0.0.0.0:8004",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}

func index(w http.ResponseWriter, req *http.Request) {
	_, err := fmt.Fprint(w, "i'm payment service")
	if err != nil {
		slog.Error("index write")
	}
}

// ReserveSlot бронируем слот
func reserveSlot(w http.ResponseWriter, r *http.Request) {
	slog.Info("call reserve delivery slot")

	// в тупую перебираем слоты, если находим то они заняты
	for i := 10; i < 20; i++ {
		one := db.Collection.FindOne(r.Context(), bson.D{{"slot", i}})
		slog.Info("find free slot", "slot", i)
		if one.Err() != nil {
			// если документ не найдет, то слот свободен
			if errors.As(one.Err(), &mongo.ErrNoDocuments) {
				resp := deliveryRequest.Response{Slot: int32(i), Status: true, Message: "ok"}

				// занимаем слот(пофиг на ошибку уже)
				db.Collection.InsertOne(r.Context(), deliveryDB.SlotDB{
					Slot: int32(i),
				})

				resp.WriteResponse(w)
				return
			}

			slog.Error("Find slot", "err", one.Err())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	// свободных слотов нет, выходим
	resp := deliveryRequest.Response{Slot: int32(0), Status: false, Message: "There are no free slots"}
	resp.WriteResponse(w)
}

// Rollback удаляем из базы занятый слот
func rollback(w http.ResponseWriter, r *http.Request) {
	slog.Info("call rollback delivery")

	data, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var request deliveryRequest.Rollback
	err = json.Unmarshal(data, &request)
	if err != nil {
		slog.Error("unmarshal rollback delivery", "err", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	slog.Info("rollback..", "slot", request.Slot)

	result := db.Collection.FindOneAndDelete(r.Context(), bson.D{{"slot", request.Slot}})
	if result.Err() != nil {
		if !errors.As(result.Err(), &mongo.ErrNoDocuments) {
			slog.Error("rollback FindOneAndDelete", "err", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	indent, err := json.MarshalIndent(response.Common{
		Status:  true,
		Message: "Ok",
	}, "	", "")
	if err != nil {
		slog.Error("MarshalIndent", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Write(indent)

}

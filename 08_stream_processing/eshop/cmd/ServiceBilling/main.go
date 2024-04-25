package main

import (
	"context"
	dbClient "eshop/internal/db"
	billingDB "eshop/internal/db/billing"
	"eshop/internal/natsclient"
	request "eshop/internal/request"
	"eshop/internal/schemagen/BillingPB"
	"eshop/internal/schemagen/NatsPB"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/mux"
	"github.com/nats-io/nats.go/jetstream"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"
)

const (
	dbName  = "Billing"
	colName = "Billing"
)

var db *billingDB.DB
var nats *natsclient.Client

func main() {
	uriNats, ok := os.LookupEnv("NATS")
	if !ok {
		slog.Error("no NATS uri")
		os.Exit(1)
	}
	slog.Info("NATS", "uri", uriNats)
	nats = natsclient.NewClient(uriNats)

	uri, ok := os.LookupEnv("DB")
	if !ok {
		slog.Error("no db uri")
		os.Exit(1)
	}
	slog.Info("DB", "uri", uri)

	dbClient := dbClient.New(uri, dbName, colName)
	db = billingDB.New(dbClient, dbName, colName)

	muxRouter := mux.NewRouter()
	muxRouter.HandleFunc("/", index).Methods(http.MethodGet)
	muxRouter.HandleFunc("/CreateUser", CreateUser).Methods(http.MethodPost)
	muxRouter.HandleFunc("/DepositCash", DepositCash).Methods(http.MethodPost)

	muxRouter.Use(func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			slog.Info("Call", "uri", r.RequestURI)
			handler.ServeHTTP(w, r)
		})
	})

	srv := &http.Server{
		Handler:      muxRouter,
		Addr:         "0.0.0.0:8002",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFunc()

	consumer, err := nats.Js.Consumer(ctx, natsclient.Pays, natsclient.Pays)
	if err != nil {
		slog.Error("nats consumer pays", "err", err)
		os.Exit(1)
	}
	_, err = consumer.Consume(handlerPays)
	if err != nil {
		slog.Error("nats consume pays", "err", err)
		os.Exit(1)
	}

	log.Fatal(srv.ListenAndServe())
}

func index(w http.ResponseWriter, req *http.Request) {
	_, err := fmt.Fprint(w, "i'm billing service")
	if err != nil {
		slog.Error("index write")
	}
}

// DepositCash пополняем счет
func DepositCash(w http.ResponseWriter, r *http.Request) {
	var depositCashData BillingPB.DepositCashRequest
	if ok := request.HttpReqToObj(w, r, &depositCashData); ok != true {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if ok := db.DepositCash(r.Context(), depositCashData.GetUser(), depositCashData.GetSum()); !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	slog.Info("Deposit balance success", "user", depositCashData.GetUser())
	request.Common{Status: true}.WriteResponse(w)
}

// CreateUser создаем пользователя
func CreateUser(w http.ResponseWriter, r *http.Request) {
	var createUserData BillingPB.CreateUserRequest
	if ok := request.HttpReqToObj(w, r, &createUserData); ok != true {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if ok := db.CreateUser(r.Context(), createUserData.GetUser()); !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	slog.Info("user created", "user", createUserData.GetUser())
	request.Common{Status: true}.WriteResponse(w)
}

func handlerPays(msg jetstream.Msg) {
	order := NatsPB.PayOrder{}
	if err := proto.Unmarshal(msg.Data(), &order); err != nil {
		slog.Error("handlerPays unmarshal", "err", err)
		return
	}
	defer msg.Ack()

	slog.Info("new pay", "msg", &order)
	ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFunc()
	pay := db.Pay(ctx, order.GetUser(), order.GetSum())

	marshal, err := proto.Marshal(&NatsPB.Notify{
		User:   order.GetUser(),
		Result: pay,
	})
	if err != nil {
		slog.Error("marshal notify", "err", err)
		return
	}

	if _, err := nats.Js.Publish(ctx, natsclient.Notify, marshal); err != nil {
		slog.Error("public notify", "err", err)
		return
	}

}

package main

import (
	"eshop/internal/natsclient"
	"eshop/internal/request"
	"eshop/internal/schemagen/OrderPB"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/mux"
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

// var db *orderDB.DB
var nats *natsclient.Client

func main() {
	/*	uri, ok := os.LookupEnv("DB")
		if !ok {
			slog.Error("no db uri")
			os.Exit(1)
		}
		slog.Info("DB", "uri", uri)
		dbClient := dbClient.New(uri, dbName, colName)
		db = orderDB.New(dbClient, dbName, colName)*/

	uriNats, ok := os.LookupEnv("NATS")
	if !ok {
		slog.Error("no NATS uri")
		os.Exit(1)
	}
	slog.Info("NATS", "uri", uriNats)
	nats = natsclient.NewClient(uriNats)

	muxRouter := mux.NewRouter()
	muxRouter.HandleFunc("/", index).Methods(http.MethodGet)
	muxRouter.HandleFunc("/Order", order).Methods(http.MethodPost)
	muxRouter.Use(func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			slog.Info("Call", "uri", r.RequestURI)
			handler.ServeHTTP(w, r)
		})
	})

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

func order(w http.ResponseWriter, r *http.Request) {
	var orderData OrderPB.CreateRequest
	if ok := request.HttpReqToObj(w, r, &orderData); ok != true {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	bOrderData, err := proto.Marshal(&orderData)
	if err != nil {
		slog.Error("order marshal", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if _, err := nats.Js.Publish(r.Context(), natsclient.Pays, bOrderData); err != nil {
		slog.Error("order: publish pay", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	slog.Info("Order created", "user", orderData.GetUser())
	request.Common{Status: true}.WriteResponse(w)
}

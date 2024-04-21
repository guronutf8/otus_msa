package main

import (
	"encoding/json"
	dbClient "eshop/internal/db"
	paymentDB "eshop/internal/db/payment"
	response "eshop/internal/request"
	payRequest "eshop/internal/request/pay"
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
	dbName  = "Payment"
	colName = "Payment"
)

var db *paymentDB.DB

func main() {
	uri, ok := os.LookupEnv("DB")
	if !ok {
		slog.Error("no db uri")
		os.Exit(1)
	}
	slog.Info("DB", "uri", uri)

	dbClient := dbClient.New(uri, dbName, colName)

	db = paymentDB.New(dbClient)

	muxRouter := mux.NewRouter()
	muxRouter.HandleFunc("/", index).Methods(http.MethodGet)
	muxRouter.HandleFunc("/pay", pay).Methods(http.MethodPost)
	//muxRouter.HandleFunc("/rollback", rollback).Methods(http.MethodPost)

	srv := &http.Server{
		Handler:      muxRouter,
		Addr:         "0.0.0.0:8003",
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

// Pay оплачиваем заказ
func pay(w http.ResponseWriter, r *http.Request) {
	slog.Info("call pay")
	data, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Error("read body request", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var pay payRequest.Pay
	err = json.Unmarshal(data, &pay)
	if err != nil {
		slog.Error("unmarshal pay request", "err", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	indent, err := json.MarshalIndent(response.Common{
		Status:  true,
		Message: "Ok",
	}, "	", "")
	if err != nil {
		slog.Error("MarshalIndent", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	if _, err := w.Write(indent); err != nil {
		slog.Error("write response", "err", err)
	}

}

/*func rollback(w http.ResponseWriter, r *http.Request) {
	slog.Info("call rollback pay")
	indent, err := json.MarshalIndent(response.Common{
		Status:  true,
		Message: "Ok",
	}, "	", "")
	if err != nil {
		slog.Error("MarshalIndent", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Write(indent)

}*/

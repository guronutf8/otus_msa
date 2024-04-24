package main

import (
	"context"
	dbClient "eshop/internal/db"
	notifyDB "eshop/internal/db/notify"
	"eshop/internal/natsclient"
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
	dbName  = "Notify"
	colName = "Notify"
)

var db *notifyDB.DB
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
	db = notifyDB.New(dbClient, dbName, colName)

	muxRouter := mux.NewRouter()
	muxRouter.HandleFunc("/List", index).Methods(http.MethodGet)

	muxRouter.Use(func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			slog.Info("Call", "uri", r.RequestURI)
			handler.ServeHTTP(w, r)
		})
	})

	srv := &http.Server{
		Handler:      muxRouter,
		Addr:         "0.0.0.0:8003",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFunc()

	consumer, err := nats.Js.Consumer(ctx, natsclient.Notify, natsclient.Notify)
	if err != nil {
		slog.Error("nats consumer panotifyys", "err", err)
		os.Exit(1)
	}
	_, err = consumer.Consume(handlerNotify)
	if err != nil {
		slog.Error("nats consume notify", "err", err)
		os.Exit(1)
	}

	log.Fatal(srv.ListenAndServe())
}

func index(w http.ResponseWriter, req *http.Request) {
	_, err := fmt.Fprint(w, "i'm notify service\n\n")
	if err != nil {
		slog.Error("index write")
		return
	}

	rows := db.GetLog(req.Context())
	for _, row := range rows {
		fmt.Fprint(w, fmt.Sprintf("user %s result %t \r\n", row.User, row.Result))
	}
}

// handlerNotify
func handlerNotify(msg jetstream.Msg) {
	notify := NatsPB.Notify{}
	if err := proto.Unmarshal(msg.Data(), &notify); err != nil {
		slog.Error("handler notify unmarshal", "err", err)
		return
	}
	defer msg.Ack()

	slog.Info("new notify", "msg", &notify)
	ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFunc()
	db.SaveNotify(ctx, notify.GetUser(), notify.GetResult())
}

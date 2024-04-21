package main

import (
	"encoding/json"
	"errors"
	dbClient "eshop/internal/db"
	storeDB "eshop/internal/db/store"
	response "eshop/internal/request"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	//requestorder "eshop/internal/fields/order"
	storeRequest "eshop/internal/request/store"
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
	dbName  = "Store"
	colName = "Store"
)

var db *storeDB.DB

func main() {
	// DB=mongodb://root:root@localhost:27017/?authMechanism=SCRAM-SHA-1
	uri, ok := os.LookupEnv("DB")
	if !ok {
		slog.Error("no db uri")
		os.Exit(1)
	}
	slog.Info("DB", "uri", uri)

	dbClient := dbClient.New(uri, dbName, colName)
	check, err := dbClient.Check()
	if err != nil {
		slog.Error("CheckDB", "err", err)
		os.Exit(1)
	}
	if !check {
		itemsInit := []interface{}{
			storeDB.ItemDB{
				Title: "mouse",
				Count: 30,
			},
			storeDB.ItemDB{
				Title: "keyboard",
				Count: 30,
			},
			storeDB.ItemDB{
				Title: "monitor",
				Count: 30,
			},
		}
		dbClient.Init(itemsInit)
	}

	db = storeDB.New(dbClient)

	muxRouter := mux.NewRouter()
	muxRouter.HandleFunc("/", index).Methods(http.MethodGet)
	muxRouter.HandleFunc("/reserve", reserve).Methods(http.MethodPost)
	muxRouter.HandleFunc("/rollback", rollback).Methods(http.MethodPost)

	srv := &http.Server{
		Handler:      muxRouter,
		Addr:         "0.0.0.0:8002",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}

func index(w http.ResponseWriter, req *http.Request) {
	_, err := fmt.Fprint(w, "i'm store service")
	if err != nil {
		slog.Error("index write")
	}
}

// Reserve резервирует товар, просто отнимая от остатка
func reserve(w http.ResponseWriter, r *http.Request) {
	slog.Info("call reserve")
	data, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Error("read body request", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var itemsRequest storeRequest.Reserve
	err = json.Unmarshal(data, &itemsRequest)
	if err != nil {
		slog.Error("unmarshal create order", "err", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// старый остаток, перед изменением
	oldValues := map[storeRequest.Item]storeDB.ItemDB{}
	for _, itemRequest := range itemsRequest.Items {
		slog.Info("check itemRequest count", "itemRequest", itemRequest.Title, "count", itemRequest.Count)
		one := db.Collection.FindOne(r.Context(), bson.D{{"title", itemRequest.Title}, {"count", bson.D{{"$gte", itemRequest.Count}}}})

		if one.Err() != nil {
			if errors.As(one.Err(), &mongo.ErrNoDocuments) {
				responseJson, err := json.MarshalIndent(response.Common{
					Status:  false,
					Message: fmt.Sprintf("Item '%s' not enough in stock, need %d", itemRequest.Title, itemRequest.Count),
				}, "", "")
				if err != nil {
					slog.Error("marshal", "err", err)
					w.WriteHeader(http.StatusOK)
					return
				}
				if _, err = w.Write(responseJson); err != nil {
					slog.Error("write responseJson", "err", err)
				}
				return
			}
			slog.Error("FindOne", "err", one.Err())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		itemDB := storeDB.ItemDB{}
		if err := one.Decode(&itemDB); err != nil {
			slog.Error("decode find one", "err", one.Err())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		oldValues[itemRequest] = itemDB
	}

	for itemRequest, itemDB := range oldValues {

		objectID, err := primitive.ObjectIDFromHex(itemDB.Id)
		if err != nil {
			slog.Error("ObjectIDFromHex", "err", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		filter := bson.D{{"_id", objectID}, {"count", itemDB.Count}}

		res := db.Collection.FindOneAndUpdate(r.Context(), filter, bson.D{{"$inc", bson.D{{"count", -itemRequest.Count}}}})
		if res.Err() != nil {
			slog.Error("FindOneAndUpdate", "err", res.Err())
			w.WriteHeader(http.StatusInternalServerError)
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
		w.Write(indent)

	}

}

func rollback(w http.ResponseWriter, r *http.Request) {
	slog.Info("call rollback reserve")
	data, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Error("read body request", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var itemsRequest storeRequest.Reserve
	err = json.Unmarshal(data, &itemsRequest)
	if err != nil {
		slog.Error("unmarshal create order", "err", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	for _, itemRequest := range itemsRequest.Items {
		filter := bson.D{{"title", itemRequest.Title}}

		res := db.Collection.FindOneAndUpdate(r.Context(), filter, bson.D{{"$inc", bson.D{{"count", itemRequest.Count}}}})
		if res.Err() != nil {
			slog.Error("FindOneAndUpdate", "err", res.Err())
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
		return
	}

	w.Write(indent)

}

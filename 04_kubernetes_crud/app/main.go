package main

import (
	"context"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
	"usercrud/internal/api"
)

func main() {

	ctx := context.Background()
	muxRouter := mux.NewRouter()
	server := api.NewServer()
	server.Init(ctx, muxRouter)

	srv := &http.Server{
		Handler: muxRouter,
		Addr:    "0.0.0.0:8000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}

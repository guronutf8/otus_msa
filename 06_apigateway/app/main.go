package main

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	muxprom "gitlab.com/msvechla/mux-prometheus/pkg/middleware"
	"log"
	"net/http"
	"time"
	"usercrud/internal/api"
)

func main() {
	recordMetrics()
	ctx := context.Background()
	muxRouter := mux.NewRouter()
	instrumentation := muxprom.NewCustomInstrumentation(true, "mux", "router", []float64{.005, .01, .025, .05, .1, .25, .5, 1, 1.5, 2.5, 3, 3.4, 5, 10}, nil, prometheus.DefaultRegisterer)
	//instrumentation := muxprom.NewCustomInstrumentation(true, "mux", "router",prometheus.DefBuckets, nil, prometheus.DefaultRegisterer)
	/*
		UseRouteTemplate:   true,
				Namespace:          "mux",
				Subsystem:          "router",
				ReqDurationBuckets: prometheus.DefBuckets,
				Registerer:         prometheus.DefaultRegisterer,*/
	muxRouter.Use(instrumentation.Middleware)

	muxRouter.Handle("/metrics", promhttp.Handler())
	fmt.Println("http://localhost/metrics")

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

func recordMetrics() {
	go func() {
		for {
			opsProcessed.Inc()
			time.Sleep(2 * time.Second)
		}
	}()
}

var (
	opsProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "myapp_processed_ops_total",
		Help: "The total number of processed events",
	})
)

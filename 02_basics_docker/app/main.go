package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
)

type HealthCheck struct {
	Status string `json:"status"`
}

var count = 0

const version = "v7"

func main() {
	fmt.Println("hello im " + version)
	http.HandleFunc("/health", health)
	http.HandleFunc("/health/", health)
	http.HandleFunc("/", index)
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		slog.Error("start http serve ", "err", err)
		os.Exit(1)
	}
}

func health(w http.ResponseWriter, r *http.Request) {
	healthCheck := HealthCheck{Status: "OK"}
	bytes, err := json.Marshal(healthCheck)
	if err != nil {
		if _, err := fmt.Fprintf(w, "{\"status\": \"FAIL\"}"); err != nil {
			slog.Error("health print status fail", "err", err)
		}
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(200)
	_, err = w.Write(bytes)
	if err != nil {
		slog.Error("health write", "err", err)
	}
}
func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<html>")
	fmt.Fprintf(w, "<head><title>otus dz 2</title></head>")
	fmt.Fprintf(w, "<body>")
	fmt.Fprintf(w, "Hello")
	fmt.Fprintf(w, "<br/>service verison %s", version)
	fmt.Fprintf(w, "<br/>host %s", r.Host)
	fmt.Fprintf(w, "<br/>RequestURI %s", r.RequestURI)
	fmt.Fprintf(w, "<br/>counter %d", count)
	fmt.Fprintf(w, "<br><a href=\"/health/\" >health</a>")
	fmt.Fprintf(w, "</body>")
	fmt.Fprintf(w, "<html>")
	fmt.Println("hello im %s, count %d", version, count)
	count++
}

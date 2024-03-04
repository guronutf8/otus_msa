package api

import (
	"fmt"
	"net/http"
	"os"
)

func (s *Server) IndexHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)

	fmt.Fprintf(w, "<html>")
	fmt.Fprintf(w, "<head><title>Домашка 6</title></head><body>\n")
	fmt.Fprintf(w, "<h1>Hello %s %s </h1>\n", r.Header.Get("X-First-Name"), r.Header.Get("X-Last-Name"))
	fmt.Fprintf(w, "<h2>You login %s</h2>\n", r.Header.Get("X-User"))
	fmt.Fprintf(w, "<h2>You email %s</h2>\n", r.Header.Get("X-Email"))
	fmt.Fprintf(w, "<hr>\n")

	hostname, err := os.Hostname()
	if err != nil {
		fmt.Fprintf(w, "<h1>hostname error: %s</h1>\n", err)
	} else {
		fmt.Fprintf(w, "<h1>host %s</h1>\n", hostname)
	}

	fmt.Fprintf(w, "<body></html>")
}

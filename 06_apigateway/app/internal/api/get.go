package api

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func (s *Server) userGetHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId, ok := vars["userId"]
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	fmt.Println(userId)
	w.WriteHeader(http.StatusOK)
}

func (s *Server) signinHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("v %d, call auth; method %s %s \n", v, r.Method, r.RequestURI)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "signin please app")
}

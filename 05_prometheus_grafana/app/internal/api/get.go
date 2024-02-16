package api

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
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

	get, err := s.db.Get(r.Context(), userId)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if get == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)

	intB, _ := json.Marshal(get)
	fmt.Println(string(intB))

	fmt.Fprintf(w, "%s", intB)

}

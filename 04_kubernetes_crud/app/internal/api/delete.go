package api

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func (s *Server) userDeleteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId, ok := vars["userId"]
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)

	fmt.Fprintf(w, "Delete user id: %v\n", userId)
}

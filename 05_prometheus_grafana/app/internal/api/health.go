package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type health struct {
	Status string
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	status := health{Status: "OK"}
	w.WriteHeader(http.StatusOK)

	jsonData, _ := json.Marshal(status)
	fmt.Println(string(jsonData))

	fmt.Fprintf(w, "%s", jsonData)

}

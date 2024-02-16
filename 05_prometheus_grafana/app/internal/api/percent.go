package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"usercrud/internal/entity"
)

func (s *Server) PercentPostHandler(w http.ResponseWriter, r *http.Request) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var per entity.Percent200
	err = json.Unmarshal(data, &per)
	if err != nil {
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	percent200 = per.Percent

	w.WriteHeader(http.StatusOK)
	fmt.Printf("Set percent: %s\n", percent200)
	fmt.Fprintf(w, "Set percent: %s\n", percent200)
}

func (s *Server) PercentGetHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	intB, _ := json.Marshal(entity.Percent200{percent200})
	fmt.Println(string(intB))

	fmt.Fprintf(w, "%s", intB)
}

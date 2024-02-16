package api

import (
	"fmt"
	"net/http"
)

var ApiEnabled = true

func (s *Server) SwitchApiHandler(w http.ResponseWriter, r *http.Request) {
	ApiEnabled = !ApiEnabled
	fmt.Printf("Api enable %t", ApiEnabled)

	fmt.Fprintf(w, "Api enable %t", ApiEnabled)

}

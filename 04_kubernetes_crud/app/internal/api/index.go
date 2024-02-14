package api

import (
	"fmt"
	"net/http"
)

func (s *Server) IndexHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)

	fmt.Fprintf(w, "<html>")
	fmt.Fprintf(w, "<head><title>Домашка 4</title></head><body>\n")
	fmt.Fprintf(w, "<h1>Hello!!!</h1>\n")
	fmt.Fprintf(w, "<a href=\"/documentation/yaml\">Swagger schema yaml</a><br>\n")
	fmt.Fprintf(w, "<a href=\"/documentation/json\">Swagger schema json</a><br>\n")
	fmt.Fprintf(w, "<a href=\"/health\">Health</a><br>\n")

	err := s.db.Client.Ping(r.Context(), nil)
	if err != nil {
		fmt.Fprintf(w, "<h3>DB connect fail: %s</h3>\n", err.Error())
	} else {
		fmt.Fprintf(w, "<h3>DB connect ok</h3>\n")
	}

	list, err := s.db.List(r.Context())
	if err != nil {
		fmt.Fprintf(w, "<h3>DB connect fail: %s</h3>\n", err.Error())
	}

	fmt.Fprintf(w, "<h3>Users</h3>\n")
	for _, user := range list {
		fmt.Fprintf(w, "<div><a href=\"/user/%s\">%s</a> %s %s %s %s %s</div>\n", user.Id, user.Id, user.Username, user.FirstName, user.LastName, user.Email, user.Phone)
	}

	fmt.Println(list)

	fmt.Fprintf(w, "<body></html>")
}

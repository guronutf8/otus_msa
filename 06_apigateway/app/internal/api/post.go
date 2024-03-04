package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"usercrud/internal/entity"
)

func (s *Server) userPostHandler(w http.ResponseWriter, r *http.Request) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var user entity.User
	err = json.Unmarshal(data, &user)
	if err != nil {
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	//user.Id = nil
	//id, err := s.db.Post(r.Context(), user)
	//if err != nil {
	//	fmt.Println(err.Error())
	//	w.WriteHeader(http.StatusInternalServerError)
	//	return
	//}

	w.WriteHeader(http.StatusOK)
	//fmt.Printf("Post user id: %s\n", id)
	//fmt.Fprintf(w, "Post user id: %s \n", id)
}

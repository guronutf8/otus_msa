package api

import (
	"math/rand"
	"net/http"
)

var percent200 = 90

func getRandCode() int {
	if rand.Intn(100) > percent200 {
		return http.StatusInternalServerError
	}

	return http.StatusOK
}

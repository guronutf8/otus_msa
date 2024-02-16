package api

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"sync"
	"time"
)

var memLoad = make(map[string]string)
var memLoadMB [][]byte

var mutex = sync.Mutex{}

func (s *Server) IndexHandlerMM(w http.ResponseWriter, r *http.Request) {
	st := make([]byte, 1014)
	mutex.Lock()
	memLoadMB = append(memLoadMB, st)
	mutex.Unlock()
	fmt.Printf("mm len %d \n", len(memLoadMB))
	w.WriteHeader(http.StatusOK)
}

func (s *Server) IndexHandlerS1(w http.ResponseWriter, r *http.Request) {
	if !ApiEnabled {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	time.Sleep(1 * time.Second)
	w.WriteHeader(getRandCode())
}
func (s *Server) IndexHandlerS2(w http.ResponseWriter, r *http.Request) {
	if !ApiEnabled {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	time.Sleep(2 * time.Second)
	w.WriteHeader(getRandCode())
}
func (s *Server) IndexHandlerS3(w http.ResponseWriter, r *http.Request) {
	if !ApiEnabled {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	time.Sleep(3 * time.Second)
	w.WriteHeader(getRandCode())
}
func (s *Server) IndexHandlerS25(w http.ResponseWriter, r *http.Request) {
	if !ApiEnabled {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	min := 0
	max := 6
	rnd := rand.Intn(max-min) + min

	time.Sleep(time.Duration(rnd) * time.Second)
	w.WriteHeader(getRandCode())
}

var counter = 0

func GetMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

func (s *Server) IndexHandler(w http.ResponseWriter, r *http.Request) {
	hash := GetMD5Hash(fmt.Sprintf("%f", rand.Float64()))
	for i := 0; i < 10000000; i++ {
		hash = GetMD5Hash(hash)
	}
	memLoad[hash] = GetMD5Hash(hash)
	counter++
	if !ApiEnabled {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(getRandCode())

	fmt.Fprintf(w, "<html>")
	fmt.Fprintf(w, "<head><title>Домашка 5</title></head><body>\n")
	fmt.Fprintf(w, "<h1>Hello!!!</h1>\n")
	fmt.Fprintf(w, "<h1>rand hash %s</h1>\n", hash)
	hostname, err := os.Hostname()
	if err != nil {
		fmt.Fprintf(w, "<h1>hostname error: %s</h1>\n", err)
	} else {
		fmt.Fprintf(w, "<h1>host %s</h1>\n", hostname)
	}
	fmt.Fprintf(w, "<h1>counter %d</h1>\n", counter)
	fmt.Fprintf(w, "<a href=\"/documentation/yaml\">Swagger schema yaml</a><br>\n")
	fmt.Fprintf(w, "<a href=\"/documentation/json\">Swagger schema json</a><br>\n")
	fmt.Fprintf(w, "<a href=\"/health\">Health</a><br>\n")

	err = s.db.Client.Ping(r.Context(), nil)
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

	//fmt.Println(list)

	fmt.Fprintf(w, "<body></html>")
}

package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"time"
)

const v = 3

var Users = make(map[string]struct{})
var Sessions = make(map[string]*Registration)

func main() {
	muxRouter := mux.NewRouter()
	muxRouter.HandleFunc("/auth", AuthHandler)
	//muxRouter.HandleFunc("/authF", AuthFHandler)
	//muxRouter.HandleFunc("/authS", AuthSHandler)
	muxRouter.HandleFunc("/signin", signinHandler)
	muxRouter.HandleFunc("/registration", RegistrationHandler)
	muxRouter.HandleFunc("/login", LoginHandler)
	muxRouter.HandleFunc("/logout", LogoutHandler)
	muxRouter.HandleFunc("/health", HealthHandler)
	muxRouter.HandleFunc("/changeEmail", ChangeEmailHandler)

	srv := &http.Server{
		Handler:      muxRouter,
		Addr:         "0.0.0.0:8001",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}

func AuthHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("v %d, call auth; method %s %s \n", v, r.Method, r.RequestURI)
	log.Printf("%s \n", v, r.Header)

	cookie, err := r.Cookie("session_id")
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	var sessionId string
	if cookie.Value != "" {
		sessionId = cookie.Value
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	registration, ok := Sessions[sessionId]
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	w.Header().Set("X-User", registration.Login)
	w.Header().Set("X-Email", registration.Email)
	w.Header().Set("X-First-Name", registration.FirstName)
	w.Header().Set("X-Last-Name", registration.LastName)

	http.SetCookie(w, &http.Cookie{Name: "session_id", Value: sessionId})
	w.WriteHeader(http.StatusOK)
}

type Registration struct {
	Login     string
	Pass      string
	Email     string
	FirstName string
	LastName  string
}

func GetMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

func RegistrationHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("call registration")

	var reg Registration
	err := json.NewDecoder(r.Body).Decode(&reg)
	if err != nil {
		log.Println("json decode", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if reg.Login == "" || reg.Pass == "" || reg.Email == "" || reg.FirstName == "" || reg.LastName == "" {
		log.Println("Field is nil")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if _, ok := Users[reg.Login]; ok == true {
		log.Println("user has been")

		http.Error(w, "user has been", http.StatusConflict)
		return
	}

	Users[reg.Login] = struct{}{}
	sessionId := GetMD5Hash(reg.Login + reg.Pass)
	Sessions[sessionId] = &reg
	cookie := http.Cookie{Name: "session_id", Value: sessionId}
	http.SetCookie(w, &cookie)
	w.WriteHeader(http.StatusOK)
	log.Printf("reg succ %s", reg.Login)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("v %d, call auth; method %s %s \n", v, r.Method, r.RequestURI)
	log.Printf("%s \n", v, r.Header)

	username, password, ok := r.BasicAuth()
	if !ok {
		http.Error(w, "hz", http.StatusUnauthorized)
	}

	for sessionId, registration := range Sessions {
		if registration.Login == username && registration.Pass == password {
			cookie := http.Cookie{Name: "session_id", Value: sessionId}
			http.SetCookie(w, &cookie)
			w.WriteHeader(http.StatusOK)
			return
		}
	}

	http.Error(w, "ЛОгин не правильный", http.StatusForbidden)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("v %d, call Logout; method %s %s \n", v, r.Method, r.RequestURI)
	log.Printf("%s \n", v, r.Header)

	cookie := http.Cookie{Name: "session_id", Value: ""}
	http.SetCookie(w, &cookie)
	w.WriteHeader(http.StatusOK)
}

func ChangeEmailHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("v %d, call ChangeEmail; method %s %s \n", v, r.Method, r.RequestURI)

	cookie, err := r.Cookie("session_id")
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	var sessionId string
	if cookie.Value != "" {
		sessionId = cookie.Value
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	registration, ok := Sessions[sessionId]
	if !ok {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	all, err := io.ReadAll(r.Body)
	registration.Email = string(all)

	w.WriteHeader(http.StatusOK)
}

func signinHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("v %d, call auth; method %s %s \n", v, r.Method, r.RequestURI)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "signin please auth")
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "{\"status\":\"OK\"}")
}

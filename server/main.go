package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/pborman/uuid"
	"net/http"
	"strconv"
	"time"
)

const cardsOnPage = 3
var mapUser map[int]UserProfile
var mapSession map[string]string
var mapLoginPwd map[string]LogPwdId
var UserCards []userCard
var MeetingCards []meetCard


func getUser(w http.ResponseWriter, r *http.Request) {
	idNeed, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	resUser, findUser := mapUser[idNeed]
	if !findUser {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resUser)
	w.WriteHeader(http.StatusOK)
}

func editUser(w http.ResponseWriter, r *http.Request) {
	session, er := r.Cookie("authToken")
	if er != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	login, ok := mapSession[session.Value]
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var change userChanges
	err := json.NewDecoder(r.Body).Decode(&change)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	fmt.Println(login, change.Field, change.Text)
	//ДОДЕЛАТЬ
}

func getPeoples(w http.ResponseWriter, r *http.Request)  {
	pageNeed, err := strconv.Atoi(mux.Vars(r)["pageNum"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if len(UserCards) < pageNeed * cardsOnPage {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	border := (pageNeed + 1) * cardsOnPage
	if len(UserCards) < border {
		border = len(UserCards)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(UserCards[pageNeed * cardsOnPage : border])
	w.WriteHeader(http.StatusOK)
}

func getMeetings(w http.ResponseWriter, r *http.Request) {
	pageNeed, err := strconv.Atoi(mux.Vars(r)["pageNum"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if len(MeetingCards) < pageNeed * cardsOnPage {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	border := (pageNeed + 1) * cardsOnPage
	if len(MeetingCards) < border {
		border = len(UserCards)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(MeetingCards[pageNeed * cardsOnPage : border])
	w.WriteHeader(http.StatusOK)
}

func signIn(w http.ResponseWriter, r *http.Request)  {
	var userData LogPwd
	err := json.NewDecoder(r.Body).Decode(&userData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if mapLoginPwd[userData.Login].Pwd != userData.Pwd {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	user := mapLoginPwd[userData.Login]

	token := uuid.New()
	expire := time.Now().Add(30 * 24 * time.Hour)
	cookie := http.Cookie{
		Name:       "authToken",
		Value:      token,
		Expires:    expire,
	}
	http.SetCookie(w, &cookie)
	mapSession[token] = user.Login
	w.WriteHeader(http.StatusOK)
}

func signOut(w http.ResponseWriter, r *http.Request)  {
	session, er := r.Cookie("authToken")
	if er != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	delete(mapSession, session.Value)
	expire := time.Now().Add(-24 * time.Hour)
	cookie := http.Cookie{
		Name:       "authToken",
		Value:      session.Value,
		Expires:    expire,
	}
	http.SetCookie(w, &cookie)
	w.WriteHeader(http.StatusOK)
}

func signUp(w http.ResponseWriter, r *http.Request)  {
	w.WriteHeader(http.StatusOK)
}

func runServer(address string) {
	router := mux.NewRouter()
	router.HandleFunc("/user/{id}", getUser).Methods(http.MethodGet) // Ready
	router.HandleFunc("/user", editUser).Methods(http.MethodPost) // Need To set Value
	router.HandleFunc("/people/{pageNum}", getPeoples).Methods(http.MethodGet) // Ready
	router.HandleFunc("/meetings/{pageNum}", getMeetings).Methods(http.MethodGet) // Ready

	router.HandleFunc("/signUp", signUp).Methods(http.MethodPost)
	router.HandleFunc("/signIn", signIn).Methods(http.MethodPost)
	router.HandleFunc("/signOut", signOut).Methods(http.MethodPost)

	server := &http.Server{
		Handler: router,
		Addr:    "127.0.0.1" + address,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	err := server.ListenAndServe()
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	FillTestData()
	runServer(":5000")
}

package main

import (
	"github.com/gorilla/mux"
	"net/http"
)

type UserProfile struct {
	ImgSrc    string 	`json:"img_src"`
	Name      string 	`json:"name"`
	City      string 	`json:"city"`
	Telegram  string 	`json:"telegram"`
	Vk        string 	`json:"vk"`
	Meetings  []meeting `json:"meetings"`
	Interest  string	`json:"interest"`
	Skills    string	`json:"skills"`
	Education string	`json:"education"`
	Job       string	`json:"job"`
	Aims      string	`json:"aims"`
}

type meeting struct {
	ImgSrc string	`json:"img_src"`
	Text   string 	`json:"text"`
}

type userCard struct {
	CardId    int      `json:"card_id"`
	ImgSrc    string   `json:"img_src"`
	Name      string   `json:"name"`
	Job       string   `json:"job"`
	Interests []string `json:"interests"`
	Skills    []string `json:"skills"`
}

type meetCard struct {
	ImgSrc		string		`json:"img_src"`
	CardId		int			`json:"card_id"`
	Text		string		`json:"text"`
	Labels		[]string	`json:"labels"`
	Title 		string		`json:"title"`
	Place 		string		`json:"place"`
	Date 		string		`json:"date"`
}

type userChanges struct {
	Field 		string	`json:"field"`
	Text 		string	`json:"text"`
}








func register(w http.ResponseWriter, r *http.Request)  {
	userLogin, err := mux.Vars(r)["login"]
	if !err {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userPassword, err := mux.Vars(r)["password"]
	if !err {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	mapLoginPwd[userLogin] = userPassword
	w.WriteHeader(http.StatusOK)
}

func login(w http.ResponseWriter, r *http.Request)  {
	userLogin, err := mux.Vars(r)["login"]
	if !err {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userPassword, err := mux.Vars(r)["password"]
	if !err {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if userPassword == mapLoginPwd[userLogin] {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusUnauthorized)
	}
}


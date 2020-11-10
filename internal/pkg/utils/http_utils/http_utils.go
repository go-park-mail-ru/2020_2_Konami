package http_utils

import (
	"encoding/json"
	"net/http"
	"time"
)

type ErrResponse struct {
	RespCode int    `json:"-"`
	ErrMsg   string `json:"error"`
}

func WriteJson(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func WriteError(w http.ResponseWriter, resp *ErrResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.RespCode)
	json.NewEncoder(w).Encode(resp)
}

func SetMonthCookie(w http.ResponseWriter, name, value string) {
	cookie := http.Cookie{
		Name:    name,
		Value:   value,
		Expires: time.Now().Add(30 * 24 * time.Hour),
	}
	http.SetCookie(w, &cookie)
}

func RemoveCookie(w http.ResponseWriter, name, value string) {
	expire := time.Now().AddDate(0, 0, -1)
	cookie := http.Cookie{
		Name:    name,
		Value:   value,
		Expires: expire,
	}
	http.SetCookie(w, &cookie)
}

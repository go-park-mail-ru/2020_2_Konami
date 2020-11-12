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
	_ = json.NewEncoder(w).Encode(resp)
}

func SetAuthCookie(w http.ResponseWriter, value string) {
	cookie := http.Cookie{
		Name:     "authToken",
		Value:    value,
		HttpOnly: true,
		Expires:  time.Now().Add(30 * 24 * time.Hour),
	}
	http.SetCookie(w, &cookie)
}

func RemoveAuthCookie(w http.ResponseWriter, value string) {
	expire := time.Now().AddDate(0, 0, -1)
	cookie := http.Cookie{
		Name:    "authToken",
		Value:   value,
		Expires: expire,
	}
	http.SetCookie(w, &cookie)
}

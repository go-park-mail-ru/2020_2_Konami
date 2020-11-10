package models

type Session struct {
	UserId int `json:"userId"`
	token  string
}

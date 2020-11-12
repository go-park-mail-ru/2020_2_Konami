package models

type Session struct {
	UserId int    `json:"userId"`
	Token  string `json:"-"`
}

package models

type Session struct {
	UserId int64  `json:"userId"`
	Token  string `json:"-"`
}

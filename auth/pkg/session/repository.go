package session

import "errors"

var ErrInvalidToken = errors.New("invalid token")
var ErrSessionNotFound = errors.New("session not found")

type Repository interface {
	GetUserId(token string) (userId int64, err error)
	CreateSession(userId int64) (token string, err error)
	RemoveSession(token string) error
}

//go:generate mockgen -source=repository.go -destination=./repositoty_mock.go -package=session
package session

import "errors"

var ErrInvalidToken = errors.New("invalid token")

type Repository interface {
	GetUserId(token string) (userId int, err error)
	CreateSession(userId int) (token string, err error)
	RemoveSession(token string) error
}

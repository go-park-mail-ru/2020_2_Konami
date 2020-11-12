//go:generate mockgen -source=usecase.go -destination=./usecase_mock.go -package=csrf
package csrf

import "errors"

var ErrExpiredToken = errors.New("token expired")

type UseCase interface {
	Create(sid string, timeStamp int64) (string, error)
	Check(sid string, inputToken string) (bool, error)
}

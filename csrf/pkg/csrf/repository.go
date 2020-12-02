//go:generate mockgen -source=repository.go -destination=./repositoty_mock.go -package=csrf
package csrf

type Repository interface {
	Add(token string, expire int64) error
	Validate(token string) error
}

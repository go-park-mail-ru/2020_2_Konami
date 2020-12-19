//go:generate mockgen -source=usecase.go -destination=./usecase_mock.go -package=profile
package profile

import (
	"errors"
	"io"
	"konami_backend/internal/pkg/models"
)

var ErrInvalidCredentials = errors.New("invalid credentials")

type UseCase interface {
	GetAll(params FilterParams) ([]models.ProfileCard, error)
	GetUserSubscriptions(params FilterParams) ([]models.ProfileCard, error)
	CreateSubscription(authorId int, targetId int) (int, error)
	RemoveSubscription(authorId int, targetId int) error
	GetProfile(reqAuthorId, userId int) (models.Profile, error)
	EditProfile(userId int, update models.ProfileUpdate) error
	UploadProfilePic(userId int, filename string, img io.Reader) error
	SignUp(cred models.Credentials) (userId int, err error)
	Validate(cred models.Credentials) (userId int, err error)
}

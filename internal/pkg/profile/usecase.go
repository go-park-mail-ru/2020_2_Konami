package profile

import (
	"errors"
	"io"
	"konami_backend/internal/pkg/models"
)

var ErrInvalidCredentials = errors.New("invalid credentials")

type UseCase interface {
	GetAll() ([]models.ProfileCard, error)
	GetProfile(userId int) (models.Profile, error)
	EditProfile(userId int, update models.ProfileUpdate) error
	UploadProfilePic(userId int, filename string, img io.Reader) error
	SignUp(cred models.Credentials) (userId int, err error)
	Validate(cred models.Credentials) (userId int, err error)
}

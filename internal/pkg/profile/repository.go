package profile

import (
	"errors"
	"konami_backend/internal/pkg/models"
)

var ErrUserNonExistent = errors.New("user non existent")

type Repository interface {
	GetAll() ([]models.ProfileCard, error)
	GetProfile(userId int) (models.Profile, error)
	EditProfile(userId int, update models.ProfileUpdate) error
	EditProfilePic(userId int, imgSrc string) error
	SignUp(cred models.Credentials) (userId int, err error)
	Validate(cred models.Credentials) (userId int, err error)
}

//go:generate mockgen -source=repository.go -destination=./repositoty_mock.go -package=profile
package profile

import (
	"errors"
	"konami_backend/internal/pkg/models"
)

var ErrUserNonExistent = errors.New("user non existent")

type Repository interface {
	GetAll() ([]models.ProfileCard, error)
	GetUserSubscriptions(userId int) ([]models.ProfileCard, error)
	CreateSubscription(authorId int, targetId int) (int, error)
	RemoveSubscription(authorId int, targetId int) error
	GetProfile(userId int) (models.Profile, error)
	EditProfile(update models.Profile) error
	EditProfilePic(userId int, imgSrc string) error
	Create(p models.Profile) (userId int, err error)
	GetCredentials(login string) (userId int, pwdHash string, err error)
	GetLabel(userId int) (models.ProfileLabel, error)
	GetTagSubscriptions(userId int) (tagIds []int, err error)
}

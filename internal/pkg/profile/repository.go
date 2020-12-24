//go:generate mockgen -source=repository.go -destination=./repositoty_mock.go -package=profile
package profile

import (
	"errors"
	"konami_backend/internal/pkg/models"
)

var ErrUserNonExistent = errors.New("user non existent")

type FilterParams struct {
	PrevId      int
	CountLimit  int
	ReqAuthorId int
}

type Repository interface {
	GetAll(params FilterParams) ([]models.ProfileCard, error)
	GetUserSubscriptionIds(params FilterParams) ([]int, error)
	GetUserSubscriptions(params FilterParams) ([]models.ProfileCard, error)
	CheckUserSubscription(authorId, targetId int) (bool, error)
	CreateSubscription(authorId int, targetId int) (int, error)
	RemoveSubscription(authorId int, targetId int) error
	GetProfile(reqAuthorId, userId int) (models.Profile, error)
	EditProfile(update models.Profile) error
	EditProfilePic(userId int, imgSrc string) error
	Create(p models.Profile) (userId int, err error)
	GetCredentials(login string) (userId int, pwdHash string, err error)
	GetLabel(userId int) (models.ProfileLabel, error)
	GetTagSubscriptions(userId int) (tagIds []int, err error)
}

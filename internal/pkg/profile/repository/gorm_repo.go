package repository

import (
	"github.com/jinzhu/gorm"
	"konami_backend/internal/pkg/models"
	"konami_backend/internal/pkg/profile"
)

type ProfileGormRepo struct {
	db *gorm.DB
}

func NewProfileGormRepo(db *gorm.DB) profile.Repository {
	return &ProfileGormRepo{db: db}
}

func (p ProfileGormRepo) GetAll() ([]models.ProfileCard, error) {
	panic("implement me")
}

func (p ProfileGormRepo) GetProfile(userId int) (models.Profile, error) {
	panic("implement me")
}

func (p ProfileGormRepo) EditProfile(userId int, update models.ProfileUpdate) error {
	panic("implement me")
}

func (p ProfileGormRepo) EditProfilePic(userId int, imgSrc string) error {
	panic("implement me")
}

func (p ProfileGormRepo) SignUp(cred models.Credentials) (userId int, err error) {
	panic("implement me")
}

func (p ProfileGormRepo) Validate(cred models.Credentials) (userId int, err error) {
	panic("implement me")
}

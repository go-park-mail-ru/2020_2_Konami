package repository

import (
	"github.com/jinzhu/gorm"
	"konami_backend/internal/pkg/session"
)

type SessionGormRepo struct {
	db *gorm.DB
}

func NewSessionGormRepo(db *gorm.DB) session.Repository {
	return &SessionGormRepo{db: db}
}

func (s SessionGormRepo) GetUserId(token string) (userId int, err error) {
	panic("implement me")
}

func (s SessionGormRepo) CreateSession(userId int) (token string, err error) {
	panic("implement me")
}

func (s SessionGormRepo) RemoveSession(token string) error {
	panic("implement me")
}

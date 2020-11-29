package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"konami_backend/internal/pkg/models"
	"konami_backend/internal/pkg/session"
)

type SessionGormRepo struct {
	db *gorm.DB
}

func NewSessionGormRepo(db *gorm.DB) session.Repository {
	return &SessionGormRepo{db: db}
}

type Session struct {
	Id     int `gorm:"primaryKey;autoIncrement;"`
	UserId int
	Token  string
}

func (s *Session) TableName() string {
	return "sessions"
}

func ToDbObject(data models.Session) Session {
	return Session{
		UserId: data.UserId,
		Token:  data.Token,
	}
}

func ToModel(obj Session) models.Session {
	return models.Session{
		UserId: obj.UserId,
		Token:  obj.Token,
	}

}

func (h SessionGormRepo) GetUserId(token string) (int, error) {
	var s Session
	db := h.db.
		Table("sessions").
		Where("Token = ?", token).
		First(&s)
	err := db.Error
	if err != nil {
		return 0, err
	}
	return s.UserId, nil
}

func (h SessionGormRepo) CreateSession(userId int) (string, error) {
	s := Session{
		UserId: userId,
		Token:  uuid.New().String(),
	}
	db := h.db.Create(&s)
	err := db.Error
	if err != nil {
		return "", err
	}
	return s.Token, nil
}

func (h SessionGormRepo) RemoveSession(token string) error {
	db := h.db.Delete(Session{}, "token = ?", token)
	return db.Error
}

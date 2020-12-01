package repository

import (
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	ModelsPkg "konami_backend/auth/pkg/models"
	"konami_backend/auth/pkg/session"
)

type SessionGormRepo struct {
	db *gorm.DB
}

func NewSessionGormRepo(db *gorm.DB) session.Repository {
	return &SessionGormRepo{db: db}
}

type Session struct {
	Id     int `gorm:"primaryKey;autoIncrement;"`
	UserId int64
	Token  string
}

func (s *Session) TableName() string {
	return "sessions"
}

func ToDbObject(data ModelsPkg.Session) Session {
	return Session{
		UserId: data.UserId,
		Token:  data.Token,
	}
}

func ToModel(obj Session) ModelsPkg.Session {
	return ModelsPkg.Session{
		UserId: obj.UserId,
		Token:  obj.Token,
	}

}

func (h SessionGormRepo) GetUserId(token string) (int64, error) {
	var s Session
	db := h.db.
		Table("sessions").
		Where("Token = ?", token).
		First(&s)
	err := db.Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, session.ErrSessionNotFound
	}
	if err != nil {
		return 0, err
	}
	return s.UserId, nil
}

func (h SessionGormRepo) CreateSession(userId int64) (string, error) {
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

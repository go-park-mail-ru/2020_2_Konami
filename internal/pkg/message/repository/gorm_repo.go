package repository

import (
	"gorm.io/gorm"
	"konami_backend/internal/pkg/message"
	"konami_backend/internal/pkg/models"
	"time"
)

type MessageGormRepo struct {
	db *gorm.DB
}

func NewMeetingGormRepo(db *gorm.DB) message.Repository {
	return &MessageGormRepo{db: db}
}

type Message struct {
	Id        int `gorm:"primaryKey;autoIncrement;"`
	AuthorId  int
	MeetingId string
	Text      string
	Timestamp time.Time
}

func (m *Message) TableName() string {
	return "messages"
}

func ToModel(obj Message) models.Message {
	return models.Message{
		Id:        obj.Id,
		AuthorId:  obj.AuthorId,
		MeetingId: obj.MeetingId,
		Text:      obj.Text,
		Timestamp: obj.Timestamp.Format("2006-01-02T15:04:05.000Z0700"),
	}
}

func ToDbObject(m models.Message) (Message, error) {
	res := Message{
		AuthorId:  m.AuthorId,
		MeetingId: m.MeetingId,
		Text:      m.Text,
		Timestamp: time.Time{},
	}
	layout := "2006-01-02T15:04:05.000Z0700"
	var err error
	res.Timestamp, err = time.Parse(layout, m.Timestamp)
	return res, err
}

func (h *MessageGormRepo) SaveMessage(message models.Message) (int, error) {
	m, err := ToDbObject(message)
	if err != nil {
		return 0, err
	}
	db := h.db.Create(&m)
	err = db.Error
	if err != nil {
		return 0, err
	}
	return m.Id, nil
}

func (h *MessageGormRepo) GetMessages(meetingId int) ([]models.Message, error) {
	var messages []Message
	bd := h.db.
		Where("MeetingId = ?", meetingId).
		Find(&messages)
	err := bd.Error
	if err != nil {
		return nil, err
	}
	res := make([]models.Message, len(messages))
	for i, msg := range messages {
		res[i] = ToModel(msg)
	}
	return res, nil
}

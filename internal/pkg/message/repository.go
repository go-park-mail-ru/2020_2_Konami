package message

import "konami_backend/internal/pkg/models"

type Repository interface {
	SaveMessage(message models.Message) (int, error)
	GetMessages(meetingId int) ([]models.Message, error)
}

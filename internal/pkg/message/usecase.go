package message

import "konami_backend/internal/pkg/models"

type UseCase interface {
	CreateMessage(message models.Message) (int, error)
	GetMessages(meetingId int) ([]models.Message, error)
}

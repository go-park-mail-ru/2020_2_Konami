//go:generate mockgen -source=repository.go -destination=./repositoty_mock.go -package=message
package message

import "konami_backend/internal/pkg/models"

type Repository interface {
	SaveMessage(message models.Message) (int, error)
	GetMessages(meetingId int) ([]models.Message, error)
}

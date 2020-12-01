//go:generate mockgen -source=usecase.go -destination=./usecase_mock.go -package=message
package message

import "konami_backend/internal/pkg/models"

type UseCase interface {
	CreateMessage(message models.Message) (int, error)
	GetMessages(meetingId int) ([]models.Message, error)
}

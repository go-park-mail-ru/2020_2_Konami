package usecase

import (
	"konami_backend/internal/pkg/message"
	"konami_backend/internal/pkg/models"
)

type MessageUseCase struct {
	repo message.Repository
}

func NewMessageUseCase(mRepo message.Repository) message.UseCase {
	return MessageUseCase{repo: mRepo}
}

func (u MessageUseCase) CreateMessage(message models.Message) (int, error) {
	return u.repo.SaveMessage(message)
}

func (u MessageUseCase) GetMessages(meetingId int) ([]models.Message, error) {
	return u.repo.GetMessages(meetingId)
}

package usecase

import (
	"konami_backend/internal/pkg/session"
)

type SessionUseCase struct {
	SessionRepo session.Repository
}

func NewSessionUseCase(SessionRepo session.Repository) session.UseCase {
	return &SessionUseCase{SessionRepo: SessionRepo}
}

func (s SessionUseCase) GetUserId(token string) (userId int, err error) {
	panic("implement me")
}

func (s SessionUseCase) CreateSession(userId int) (token string, err error) {
	panic("implement me")
}

func (s SessionUseCase) RemoveSession(token string) error {
	panic("implement me")
}

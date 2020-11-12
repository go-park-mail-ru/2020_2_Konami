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

func (uc SessionUseCase) GetUserId(token string) (userId int, err error) {
	return uc.SessionRepo.GetUserId(token)
}

func (uc SessionUseCase) CreateSession(userId int) (token string, err error) {
	return uc.SessionRepo.CreateSession(userId)
}

func (uc SessionUseCase) RemoveSession(token string) error {
	return uc.SessionRepo.RemoveSession(token)
}

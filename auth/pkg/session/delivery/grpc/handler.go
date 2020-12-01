package grpc

import (
	"context"
	"errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"konami_backend/auth/pkg/session"
	"konami_backend/proto/auth"
)

type SessionHandler struct {
	SessionUC session.UseCase
}

func NewSessionHandler(sessionUC session.UseCase) SessionHandler {
	return SessionHandler{SessionUC: sessionUC}
}

func (uc *SessionHandler) Create(_ context.Context, in *auth.Session) (*auth.SessionToken, error) {
	sid, err := uc.SessionUC.CreateSession(in.UserId)
	if err != nil {
		return nil, status.Error(codes.Aborted, err.Error())
	}
	return &auth.SessionToken{Token: sid}, nil
}

func (uc *SessionHandler) Check(_ context.Context, in *auth.SessionToken) (*auth.Session, error) {
	userId, err := uc.SessionUC.GetUserId(in.Token)
	if errors.Is(err, session.ErrSessionNotFound) {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	if err != nil {
		return nil, status.Error(codes.Aborted, err.Error())
	}
	return &auth.Session{
		UserId: userId,
	}, nil
}

func (uc *SessionHandler) Delete(_ context.Context, in *auth.SessionToken) (*auth.Nothing, error) {
	err := uc.SessionUC.RemoveSession(in.Token)
	if err != nil {
		return nil, status.Error(codes.Aborted, err.Error())
	}
	return &auth.Nothing{Dummy: true}, nil
}

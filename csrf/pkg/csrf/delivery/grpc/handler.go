package grpc

import (
	"context"
	"errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"konami_backend/csrf/pkg/csrf"
	loggerPkg "konami_backend/logger"
	csrfProto "konami_backend/proto/csrf"
)

type CsrfHandler struct {
	CsrfUC csrf.UseCase
	Logger *loggerPkg.Logger
}

func NewCsrfHandler(csrfUC csrf.UseCase, logger *loggerPkg.Logger) CsrfHandler {
	return CsrfHandler{CsrfUC: csrfUC, Logger: logger}
}

func (uc *CsrfHandler) Create(_ context.Context, in *csrfProto.CsrfData) (*csrfProto.CsrfToken, error) {
	sid := in.Sid
	tok, err := uc.CsrfUC.Create(sid, in.TimeStamp)
	if err != nil {
		uc.Logger.LogError("csrf/delivery/grpc", "Create", err)
		return nil, status.Error(codes.Aborted, err.Error())
	}
	return &csrfProto.CsrfToken{Sid: sid, Token: tok}, nil
}

func (uc *CsrfHandler) Check(_ context.Context, in *csrfProto.CsrfToken) (*csrfProto.IsValid, error) {
	isValid, err := uc.CsrfUC.Check(in.Sid, in.Token)
	if errors.Is(err, csrf.ErrExpiredToken) || (err == nil && !isValid) {
		return &csrfProto.IsValid{Value: false}, nil
	}
	if err != nil {
		uc.Logger.LogError("csrf/delivery/grpc", "Check", err)
		return nil, status.Error(codes.Aborted, err.Error())
	}
	return &csrfProto.IsValid{Value: isValid}, nil
}

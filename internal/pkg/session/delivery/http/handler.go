package http

import (
	"bytes"
	"errors"
	"konami_backend/internal/pkg/middleware"
	"konami_backend/internal/pkg/models"
	"konami_backend/internal/pkg/profile"
	"konami_backend/internal/pkg/session"
	hu "konami_backend/internal/pkg/utils/http_utils"
	"net/http"
)

type SessionHandler struct {
	SessionUC session.UseCase
	ProfileUC profile.UseCase
}

func (h *SessionHandler) GetUserId(w http.ResponseWriter, r *http.Request) {
	uId, ok := r.Context().Value(middleware.UserID).(int)
	if !ok {
		hu.WriteError(w, &hu.ErrResponse{RespCode: http.StatusUnauthorized})
		return
	}
	hu.WriteJson(w, struct {
		UserId int `json:"userId"`
	}{uId})
}

func (h *SessionHandler) LogIn(w http.ResponseWriter, r *http.Request) {
	var cred models.Credentials
	buf := new(bytes.Buffer)
	_, _ = buf.ReadFrom(r.Body)
	err := cred.UnmarshalJSON(buf.Bytes())
	if err != nil {
		hu.WriteError(w, &hu.ErrResponse{RespCode: http.StatusBadRequest})
		return
	}
	userId, err := h.ProfileUC.Validate(cred)
	if errors.Is(err, profile.ErrInvalidCredentials) || errors.Is(err, profile.ErrUserNonExistent) {
		hu.WriteError(w, &hu.ErrResponse{RespCode: http.StatusBadRequest, ErrMsg: "invalid credentials"})
		return
	}
	if err != nil {
		hu.WriteError(w, &hu.ErrResponse{RespCode: http.StatusInternalServerError})
		return
	}
	token, err := h.SessionUC.CreateSession(userId)
	hu.SetAuthCookie(w, token)
	w.WriteHeader(http.StatusCreated)
}

func (h *SessionHandler) LogOut(w http.ResponseWriter, r *http.Request) {
	token, ok := r.Context().Value(middleware.AuthToken).(string)
	if !ok {
		hu.WriteError(w, &hu.ErrResponse{RespCode: http.StatusUnauthorized})
		return
	}
	err := h.SessionUC.RemoveSession(token)
	if err != nil && err != session.ErrInvalidToken {
		hu.WriteError(w, &hu.ErrResponse{RespCode: http.StatusInternalServerError})
		return
	}
	hu.RemoveAuthCookie(w, token)
	w.WriteHeader(http.StatusOK)
}

package http

import (
	"encoding/json"
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
	token, err := r.Cookie("authToken")
	if err != nil {
		hu.WriteError(w, &hu.ErrResponse{RespCode: http.StatusUnauthorized})
		return
	}
	uId, err := h.SessionUC.GetUserId(token.Value)
	if err != nil {
		hu.WriteError(w, &hu.ErrResponse{RespCode: http.StatusUnauthorized})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	hu.WriteJson(w, struct{ userId int }{uId})
}

func (h *SessionHandler) LogIn(w http.ResponseWriter, r *http.Request) {
	var cred models.Credentials
	err := json.NewDecoder(r.Body).Decode(&cred)
	if err != nil {
		hu.WriteError(w, &hu.ErrResponse{RespCode: http.StatusBadRequest})
		return
	}
	userId, err := h.ProfileUC.Validate(cred)
	if err == profile.ErrUserNonExistent {
		hu.WriteError(w, &hu.ErrResponse{RespCode: http.StatusUnauthorized, ErrMsg: "invalid credentials"})
		return
	}
	if err != nil {
		hu.WriteError(w, &hu.ErrResponse{RespCode: http.StatusInternalServerError})
		return
	}
	token, err := h.SessionUC.CreateSession(userId)
	hu.SetMonthCookie(w, "authToken", token)
	w.WriteHeader(http.StatusOK)
}

func (h *SessionHandler) LogOut(w http.ResponseWriter, r *http.Request) {
	token, err := r.Cookie("authToken")
	if err != nil {
		w.WriteHeader(http.StatusOK)
		return
	}
	err = h.SessionUC.RemoveSession(token.Value)
	if err != nil && err != session.ErrInvalidToken {
		hu.WriteError(w, &hu.ErrResponse{RespCode: http.StatusInternalServerError})
		return
	}
	hu.RemoveCookie(w, "authToken", token.Value)
	w.WriteHeader(http.StatusOK)
}

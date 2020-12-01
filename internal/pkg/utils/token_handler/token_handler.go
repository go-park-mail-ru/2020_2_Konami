package token_handler

import (
	"context"
	"errors"
	"konami_backend/internal/pkg/middleware"
	"konami_backend/logger"
	"konami_backend/proto/csrf"
	"net/http"
	"time"
)

type TokenHandler struct {
	CsrfClient csrf.CsrfDispatcherClient
	Log        *logger.Logger
}

func (h *TokenHandler) GetCSRF(w http.ResponseWriter, r *http.Request) {
	authTok, ok := r.Context().Value(middleware.AuthToken).(string)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	csrfTok, err := h.CsrfClient.Create(context.Background(), &csrf.CsrfData{Sid: authTok, TimeStamp: time.Now().Unix()})

	if err != nil {
		h.Log.LogError("csrfDelivery", "GetCSRF", errors.New("CSRF creation failed"))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Access-Control-Expose-Headers", "Csrf-Token")
	w.Header().Set("Csrf-Token", csrfTok.Token)
	w.WriteHeader(http.StatusOK)
}

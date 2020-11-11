package http

import (
	"errors"
	"konami_backend/internal/pkg/csrf"
	"konami_backend/internal/pkg/middleware"
	"konami_backend/logger"
	"net/http"
	"time"
)

type CSRFHandler struct {
	CsrfUC csrf.UseCase
	Log    *logger.Logger
}

func (h *CSRFHandler) GetCSRF(w http.ResponseWriter, r *http.Request) {
	authTok, ok := r.Context().Value(middleware.AuthToken).(string)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	csrfTok, err := h.CsrfUC.Create(authTok, time.Now().Unix())
	if err != nil {
		h.Log.LogError("csrfDelivery", "GetCSRF", errors.New("CSRF creation failed"))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Access-Control-Expose-Headers", "Csrf-Token")
	w.Header().Set("Csrf-Token", csrfTok)
	w.WriteHeader(http.StatusOK)
}

package http

import (
	"github.com/golang/mock/gomock"
	"github.com/steinfletcher/apitest"
	"konami_backend/internal/pkg/middleware"
	"konami_backend/internal/pkg/session"
	"konami_backend/logger"
	"net/http"
	"os"
	"testing"
)

var testHandler SessionHandler

func init() {
	testHandler.Log = logger.NewLogger(os.Stdout)
}

func TestSessions(t *testing.T) {
	t.Run("Get-OK", func(t *testing.T) {
		key := "4"
		handler := middleware.SetMuxVars(testHandler.GetUserId, "userId", key)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := session.NewMockUseCase(ctrl)
		testHandler.SessionUC = m

		apitest.New("Get-OK").
			Handler(handler).
			Method("Get").
			URL("/api/me/").
			Expect(t).
			Status(http.StatusOK).
			End()
	})

	t.Run("Get-OK", func(t *testing.T) {
		key := "4"
		handler := middleware.SetMuxVars(testHandler.GetUserId, "userId", key)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := session.NewMockUseCase(ctrl)
		testHandler.SessionUC = m

		apitest.New("Get-OK").
			Handler(handler).
			Method("Get").
			URL("/api/me/").
			Expect(t).
			Status(http.StatusOK).
			End()
	})
}
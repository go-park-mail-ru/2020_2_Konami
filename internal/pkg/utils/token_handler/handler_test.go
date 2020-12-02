package token_handler

import (
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/steinfletcher/apitest"
	"konami_backend/csrf/pkg/csrf"
	"konami_backend/internal/pkg/middleware"
	"konami_backend/logger"
	"net/http"
	"os"
	"testing"
	"time"
)

var testHandler TokenHandler

func init() {
	testHandler.Log = logger.NewLogger(os.Stdout)
}

func TestSessions(t *testing.T) {
	t.Run("Get-OK", func(t *testing.T) {
		var args []middleware.RouteArgs
		args = append(args, middleware.RouteArgs{Key: "authToken", Value: "4234124"})
		handler := middleware.SetMuxVars(testHandler.GetCSRF, args)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := csrf.NewMockCsrfDispatcherClient(ctrl)
		testHandler.CsrfClient = m

		m.EXPECT().
			Create("4234124", time.Now().Unix()).
			Return("LOL", nil)

		apitest.New("Get-OK").
			Handler(handler).
			Method("Get").
			URL("/api/csrf/").
			Expect(t).
			Status(http.StatusOK).
			End()
	})

	t.Run("Get-BAD", func(t *testing.T) {
		var args []middleware.RouteArgs
		args = append(args, middleware.RouteArgs{Key: "authTken", Value: "4234124"})
		handler := middleware.SetMuxVars(testHandler.GetCSRF, args)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := csrf.NewMockCsrfDispatcherClient(ctrl)
		testHandler.CsrfClient = m

		apitest.New("Get-BAD").
			Handler(handler).
			Method("Get").
			URL("/api/csrf/").
			Expect(t).
			Status(http.StatusUnauthorized).
			End()
	})

	t.Run("Create-Bad", func(t *testing.T) {
		var args []middleware.RouteArgs
		args = append(args, middleware.RouteArgs{Key: "authToken", Value: "4234124"})
		handler := middleware.SetMuxVars(testHandler.GetCSRF, args)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := csrf.NewMockCsrfDispatcherClient(ctrl)
		testHandler.CsrfClient = m

		m.EXPECT().
			Create("4234124", time.Now().Unix()).
			Return("LOL", errors.New("ERROR"))

		apitest.New("Create-Bad").
			Handler(handler).
			Method("Get").
			URL("/api/csrf/").
			Expect(t).
			Status(http.StatusInternalServerError).
			End()
	})
}

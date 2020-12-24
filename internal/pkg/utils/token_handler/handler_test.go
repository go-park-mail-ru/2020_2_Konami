package token_handler

import (
	"github.com/golang/mock/gomock"
	"github.com/steinfletcher/apitest"
	"konami_backend/csrf/pkg/csrf"
	"konami_backend/internal/pkg/middleware"
	"konami_backend/logger"
	"net/http"
	"os"
	"testing"
)

var testHandler TokenHandler

func init() {
	testHandler.Log = logger.NewLogger(os.Stdout)
}

func TestSessions(t *testing.T) {
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
}

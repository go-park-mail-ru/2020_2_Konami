package http

import (
	"encoding/json"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/steinfletcher/apitest"
	"github.com/stretchr/testify/assert"
	"konami_backend/auth/pkg/session"
	"konami_backend/internal/pkg/middleware"
	"konami_backend/internal/pkg/models"
	"konami_backend/internal/pkg/profile"
	"net/http"
	"testing"
)

var testHandler SessionHandler

func TestSessions(t *testing.T) {
	t.Run("Get-OK", func(t *testing.T) {
		var args []middleware.RouteArgs
		args = append(args, middleware.RouteArgs{Key: "userId", Value: 4})

		handler := middleware.SetMuxVars(testHandler.GetUserId, args)

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

	t.Run("Get-BAD", func(t *testing.T) {
		var args []middleware.RouteArgs
		args = append(args, middleware.RouteArgs{Key: "useId", Value: 4})

		handler := middleware.SetMuxVars(testHandler.GetUserId, args)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := session.NewMockUseCase(ctrl)
		testHandler.SessionUC = m

		apitest.New("Get-OK").
			Handler(handler).
			Method("Get").
			URL("/api/me/").
			Expect(t).
			Status(http.StatusUnauthorized).
			End()
	})

	t.Run("LogIN", func(t *testing.T) {
		var args []middleware.RouteArgs
		handler := middleware.SetMuxVars(testHandler.LogIn, args)

		testCredit := models.Credentials{
			Login:    "qwerty",
			Password: "qwerty",
		}
		testCreditJSON, err := json.Marshal(testCredit)
		assert.NoError(t, err)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := session.NewMockUseCase(ctrl)
		testHandler.SessionUC = m

		n := profile.NewMockUseCase(ctrl)
		testHandler.ProfileUC = n

		n.EXPECT().
			Validate(testCredit).
			Return(1, nil)

		m.EXPECT().
			CreateSession(1).
			Return("lol", nil)

		apitest.New("LogIN").
			Handler(handler).
			Method("POST").
			URL("/login").
			Body(string(testCreditJSON)).
			Expect(t).
			Status(http.StatusCreated).
			End()
	})

	t.Run("LogIN-BadReq", func(t *testing.T) {
		var args []middleware.RouteArgs
		handler := middleware.SetMuxVars(testHandler.LogIn, args)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := session.NewMockUseCase(ctrl)
		testHandler.SessionUC = m

		apitest.New("LogIN").
			Handler(handler).
			Method("POST").
			URL("/login").
			Expect(t).
			Status(http.StatusBadRequest).
			End()
	})

	t.Run("LogIN-BadValidate", func(t *testing.T) {
		var args []middleware.RouteArgs
		handler := middleware.SetMuxVars(testHandler.LogIn, args)

		testCredit := models.Credentials{
			Login:    "qwerty",
			Password: "qwerty",
		}
		testCreditJSON, err := json.Marshal(testCredit)
		assert.NoError(t, err)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := session.NewMockUseCase(ctrl)
		testHandler.SessionUC = m

		n := profile.NewMockUseCase(ctrl)
		testHandler.ProfileUC = n

		n.EXPECT().
			Validate(testCredit).
			Return(0, profile.ErrInvalidCredentials)

		apitest.New("LogIN").
			Handler(handler).
			Method("POST").
			URL("/login").
			Body(string(testCreditJSON)).
			Expect(t).
			Status(http.StatusBadRequest).
			End()
	})

	t.Run("LogIN-InternalErr", func(t *testing.T) {
		var args []middleware.RouteArgs
		handler := middleware.SetMuxVars(testHandler.LogIn, args)

		testCredit := models.Credentials{
			Login:    "qwerty",
			Password: "qwerty",
		}
		testCreditJSON, err := json.Marshal(testCredit)
		assert.NoError(t, err)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := session.NewMockUseCase(ctrl)
		testHandler.SessionUC = m

		n := profile.NewMockUseCase(ctrl)
		testHandler.ProfileUC = n

		n.EXPECT().
			Validate(testCredit).
			Return(0, errors.New("ERROR"))

		apitest.New("LogIN").
			Handler(handler).
			Method("POST").
			URL("/login").
			Body(string(testCreditJSON)).
			Expect(t).
			Status(http.StatusInternalServerError).
			End()
	})

	t.Run("LogOut-Bad", func(t *testing.T) {
		var args []middleware.RouteArgs
		handler := middleware.SetMuxVars(testHandler.LogOut, args)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := session.NewMockUseCase(ctrl)
		testHandler.SessionUC = m

		apitest.New("LogOUT").
			Handler(handler).
			Method("POST").
			URL("/logout").
			Expect(t).
			Status(http.StatusUnauthorized).
			End()
	})

	t.Run("LogOut", func(t *testing.T) {
		var args []middleware.RouteArgs
		args = append(args, middleware.RouteArgs{
			Key:   middleware.AuthToken,
			Value: "Some_tok",
		})
		handler := middleware.SetMuxVars(testHandler.LogOut, args)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := session.NewMockUseCase(ctrl)
		testHandler.SessionUC = m

		m.EXPECT().
			RemoveSession("Some_tok").
			Return(nil)

		apitest.New("LogOUT").
			Handler(handler).
			Method("POST").
			URL("/logout").
			Expect(t).
			Status(http.StatusOK).
			End()
	})

	t.Run("LogOut", func(t *testing.T) {
		var args []middleware.RouteArgs
		args = append(args, middleware.RouteArgs{
			Key:   middleware.AuthToken,
			Value: "Some_tok",
		})
		handler := middleware.SetMuxVars(testHandler.LogOut, args)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := session.NewMockUseCase(ctrl)
		testHandler.SessionUC = m

		m.EXPECT().
			RemoveSession("Some_tok").
			Return(errors.New("ERROR"))

		apitest.New("LogOUT").
			Handler(handler).
			Method("POST").
			URL("/logout").
			Expect(t).
			Status(http.StatusInternalServerError).
			End()
	})
}

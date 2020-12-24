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

var testHandler ProfileHandler

func TestSessions(t *testing.T) {
	t.Run("GetUserSubscriptions", func(t *testing.T) {
		var args []middleware.RouteArgs

		handler := middleware.SetMuxVars(testHandler.GetUserSubscriptions, args)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		p := profile.NewMockUseCase(ctrl)
		testHandler.ProfileUC = p

		p.EXPECT().GetUserSubscriptions(profile.FilterParams{ReqAuthorId: -1}).Return([]models.ProfileCard{}, nil)

		apitest.New("Get-All-Ok").
			Handler(handler).
			Method("Get").
			URL("/people").
			Expect(t).
			Status(http.StatusOK).
			End()
	})

	t.Run("GetUserSubscriptionsErr1", func(t *testing.T) {
		var args []middleware.RouteArgs

		handler := middleware.SetMuxVars(testHandler.GetUserSubscriptions, args)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		p := profile.NewMockUseCase(ctrl)
		testHandler.ProfileUC = p

		p.EXPECT().GetUserSubscriptions(profile.FilterParams{ReqAuthorId: -1}).Return([]models.ProfileCard{}, errors.New("err"))

		apitest.New("Get-All-Ok").
			Handler(handler).
			Method("Get").
			URL("/people").
			Expect(t).
			Status(http.StatusNotFound).
			End()
	})

	t.Run("Get-All-Ok", func(t *testing.T) {
		var args []middleware.RouteArgs
		//args = append(args, middleware.RouteArgs{Key: "userId", Value: 4})

		handler := middleware.SetMuxVars(testHandler.GetPeople, args)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		p := profile.NewMockUseCase(ctrl)
		testHandler.ProfileUC = p

		p.EXPECT().GetAll(profile.FilterParams{ReqAuthorId: -1}).Return([]models.ProfileCard{}, nil)

		apitest.New("Get-All-Ok").
			Handler(handler).
			Method("Get").
			URL("/people").
			Expect(t).
			Status(http.StatusOK).
			End()
	})

	t.Run("Get-All-Bad", func(t *testing.T) {
		var args []middleware.RouteArgs
		//args = append(args, middleware.RouteArgs{Key: "userId", Value: 4})

		handler := middleware.SetMuxVars(testHandler.GetPeople, args)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		p := profile.NewMockUseCase(ctrl)
		testHandler.ProfileUC = p

		p.EXPECT().GetAll(profile.FilterParams{ReqAuthorId: -1}).Return([]models.ProfileCard{}, errors.New("Error"))

		apitest.New("Get-All-Bad").
			Handler(handler).
			Method("Get").
			URL("/people").
			Expect(t).
			Status(http.StatusInternalServerError).
			End()
	})

	t.Run("Edit", func(t *testing.T) {
		var args []middleware.RouteArgs
		args = append(args, middleware.RouteArgs{Key: middleware.UserID, Value: 4})
		args = append(args, middleware.RouteArgs{Key: middleware.CSRFValid, Value: true})

		handler := middleware.SetMuxVars(testHandler.EditUser, args)

		testStr := "test"
		testUpd := models.ProfileUpdate{
			Name:        nil,
			Gender:      nil,
			City:        &testStr,
			Birthday:    nil,
			Telegram:    nil,
			Vk:          nil,
			MeetingTags: nil,
			Education:   nil,
			Job:         nil,
			Aims:        nil,
			Interests:   nil,
			Skills:      nil,
		}
		testUpdJSON, _ := json.Marshal(testUpd)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		p := profile.NewMockUseCase(ctrl)
		testHandler.ProfileUC = p

		p.EXPECT().EditProfile(4, testUpd)

		apitest.New("Edit").
			Handler(handler).
			Method("PATCH").
			URL("/user").
			Body(string(testUpdJSON)).
			Expect(t).
			Status(http.StatusOK).
			End()
	})

	t.Run("EditBad1", func(t *testing.T) {
		var args []middleware.RouteArgs

		handler := middleware.SetMuxVars(testHandler.EditUser, args)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		p := profile.NewMockUseCase(ctrl)
		testHandler.ProfileUC = p

		apitest.New("Edit").
			Handler(handler).
			Method("PATCH").
			URL("/user").
			Expect(t).
			Status(http.StatusUnauthorized).
			End()
	})

	t.Run("EditBad2", func(t *testing.T) {
		var args []middleware.RouteArgs
		args = append(args, middleware.RouteArgs{Key: middleware.UserID, Value: 4})

		handler := middleware.SetMuxVars(testHandler.EditUser, args)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		p := profile.NewMockUseCase(ctrl)
		testHandler.ProfileUC = p

		apitest.New("Edit").
			Handler(handler).
			Method("PATCH").
			URL("/user").
			Expect(t).
			Status(http.StatusUnauthorized).
			End()
	})

	t.Run("EditBad3", func(t *testing.T) {
		var args []middleware.RouteArgs
		args = append(args, middleware.RouteArgs{Key: middleware.UserID, Value: 4})
		args = append(args, middleware.RouteArgs{Key: middleware.CSRFValid, Value: true})

		handler := middleware.SetMuxVars(testHandler.EditUser, args)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		p := profile.NewMockUseCase(ctrl)
		testHandler.ProfileUC = p

		apitest.New("Edit").
			Handler(handler).
			Method("PATCH").
			URL("/user").
			Expect(t).
			Status(http.StatusBadRequest).
			End()
	})

	t.Run("Edit", func(t *testing.T) {
		var args []middleware.RouteArgs
		args = append(args, middleware.RouteArgs{Key: middleware.UserID, Value: 4})
		args = append(args, middleware.RouteArgs{Key: middleware.CSRFValid, Value: true})

		handler := middleware.SetMuxVars(testHandler.EditUser, args)

		testStr := "test"
		testUpd := models.ProfileUpdate{
			Name:        nil,
			Gender:      nil,
			City:        &testStr,
			Birthday:    nil,
			Telegram:    nil,
			Vk:          nil,
			MeetingTags: nil,
			Education:   nil,
			Job:         nil,
			Aims:        nil,
			Interests:   nil,
			Skills:      nil,
		}
		testUpdJSON, _ := json.Marshal(testUpd)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		p := profile.NewMockUseCase(ctrl)
		testHandler.ProfileUC = p

		p.EXPECT().EditProfile(4, testUpd).Return(errors.New("Err"))

		apitest.New("Edit").
			Handler(handler).
			Method("PATCH").
			URL("/user").
			Body(string(testUpdJSON)).
			Expect(t).
			Status(http.StatusInternalServerError).
			End()
	})

	t.Run("Get-User-Ok", func(t *testing.T) {
		var args []middleware.QueryArgs
		args = append(args, middleware.QueryArgs{Key: "userId", Value: "4"})

		handler := middleware.SetVars(testHandler.GetUser, args)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		p := profile.NewMockUseCase(ctrl)
		testHandler.ProfileUC = p

		p.EXPECT().GetProfile(-1, 4).Return(models.Profile{}, nil)

		apitest.New("Get-All-Ok").
			Handler(handler).
			Method("Get").
			URL("/user").
			Expect(t).
			Status(http.StatusOK).
			End()
	})

	t.Run("Get-User-Bad1", func(t *testing.T) {
		var args []middleware.QueryArgs

		handler := middleware.SetVars(testHandler.GetUser, args)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		p := profile.NewMockUseCase(ctrl)
		testHandler.ProfileUC = p

		apitest.New("Get-All-Ok").
			Handler(handler).
			Method("Get").
			URL("/user").
			Expect(t).
			Status(http.StatusNotFound).
			End()
	})

	t.Run("Get-User-Bad2", func(t *testing.T) {
		var args []middleware.QueryArgs
		args = append(args, middleware.QueryArgs{Key: "userId", Value: "4"})

		handler := middleware.SetVars(testHandler.GetUser, args)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		p := profile.NewMockUseCase(ctrl)
		testHandler.ProfileUC = p

		p.EXPECT().GetProfile(-1, 4).Return(models.Profile{}, errors.New("Err"))

		apitest.New("Get-All-Ok").
			Handler(handler).
			Method("Get").
			URL("/user").
			Expect(t).
			Status(http.StatusNotFound).
			End()
	})

	t.Run("UploadBad1", func(t *testing.T) {
		var args []middleware.RouteArgs
		args = append(args, middleware.RouteArgs{Key: middleware.UserID, Value: 4})
		args = append(args, middleware.RouteArgs{Key: middleware.CSRFValid, Value: true})

		handler := middleware.SetMuxVars(testHandler.UploadUserPic, args)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		p := profile.NewMockUseCase(ctrl)
		testHandler.ProfileUC = p

		apitest.New("Get-All-Ok").
			Handler(handler).
			Method("Get").
			URL("/people").
			Expect(t).
			Status(http.StatusBadRequest).
			End()
	})

	t.Run("UploadBad2", func(t *testing.T) {
		var args []middleware.RouteArgs
		args = append(args, middleware.RouteArgs{Key: middleware.UserID, Value: 4})

		handler := middleware.SetMuxVars(testHandler.UploadUserPic, args)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		p := profile.NewMockUseCase(ctrl)
		testHandler.ProfileUC = p

		apitest.New("Get-All-Ok").
			Handler(handler).
			Method("Get").
			URL("/people").
			Expect(t).
			Status(http.StatusUnauthorized).
			End()
	})

	t.Run("UploadBad3", func(t *testing.T) {
		var args []middleware.RouteArgs

		args = append(args, middleware.RouteArgs{Key: middleware.CSRFValid, Value: true})
		handler := middleware.SetMuxVars(testHandler.UploadUserPic, args)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		p := profile.NewMockUseCase(ctrl)
		testHandler.ProfileUC = p

		apitest.New("Get-All-Ok").
			Handler(handler).
			Method("Get").
			URL("/people").
			Expect(t).
			Status(http.StatusUnauthorized).
			End()
	})

	t.Run("SignUpBad2", func(t *testing.T) {
		var args []middleware.RouteArgs

		handler := middleware.SetMuxVars(testHandler.SignUp, args)

		testCred := models.Credentials{
			Login:    "qwerty",
			Password: "qwerty",
		}
		testUpdJSON, _ := json.Marshal(testCred)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		p := profile.NewMockUseCase(ctrl)
		testHandler.ProfileUC = p

		m := session.NewMockAuthCheckerClient(ctrl)
		testHandler.AuthClient = m

		p.EXPECT().Validate(testCred).Return(0, profile.ErrUserNonExistent)
		p.EXPECT().SignUp(testCred).Return(0, errors.New("Err"))

		apitest.New("Edit").
			Handler(handler).
			Method("PATCH").
			URL("/user").
			Body(string(testUpdJSON)).
			Expect(t).
			Status(http.StatusBadRequest).
			End()
	})

	t.Run("SignUpBad3", func(t *testing.T) {
		var args []middleware.RouteArgs

		handler := middleware.SetMuxVars(testHandler.SignUp, args)

		testCred := models.Credentials{
			Login:    "qwerty",
			Password: "qwerty",
		}
		testUpdJSON, _ := json.Marshal(testCred)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		p := profile.NewMockUseCase(ctrl)
		testHandler.ProfileUC = p

		m := session.NewMockAuthCheckerClient(ctrl)
		testHandler.AuthClient = m

		p.EXPECT().Validate(testCred).Return(0, nil)

		apitest.New("Edit").
			Handler(handler).
			Method("PATCH").
			URL("/user").
			Body(string(testUpdJSON)).
			Expect(t).
			Status(http.StatusConflict).
			End()
	})

	t.Run("SignUpBad4", func(t *testing.T) {
		var args []middleware.RouteArgs

		handler := middleware.SetMuxVars(testHandler.SignUp, args)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		p := profile.NewMockUseCase(ctrl)
		testHandler.ProfileUC = p

		m := session.NewMockAuthCheckerClient(ctrl)
		testHandler.AuthClient = m

		apitest.New("Edit").
			Handler(handler).
			Method("PATCH").
			URL("/user").
			Expect(t).
			Status(http.StatusBadRequest).
			End()
	})

	t.Run("Get-OK", func(t *testing.T) {
		var args []middleware.RouteArgs
		args = append(args, middleware.RouteArgs{Key: "userId", Value: 4})

		handler := middleware.SetMuxVars(testHandler.GetUserId, args)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := session.NewMockAuthCheckerClient(ctrl)
		testHandler.AuthClient = m

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

		m := session.NewMockAuthCheckerClient(ctrl)
		testHandler.AuthClient = m

		apitest.New("Get-OK").
			Handler(handler).
			Method("Get").
			URL("/api/me/").
			Expect(t).
			Status(http.StatusUnauthorized).
			End()
	})

	t.Run("LogIN-BadReq", func(t *testing.T) {
		var args []middleware.RouteArgs
		handler := middleware.SetMuxVars(testHandler.LogIn, args)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := session.NewMockAuthCheckerClient(ctrl)
		testHandler.AuthClient = m

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

		m := session.NewMockAuthCheckerClient(ctrl)
		testHandler.AuthClient = m

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

		m := session.NewMockAuthCheckerClient(ctrl)
		testHandler.AuthClient = m

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

		m := session.NewMockAuthCheckerClient(ctrl)
		testHandler.AuthClient = m

		apitest.New("LogOUT").
			Handler(handler).
			Method("POST").
			URL("/logout").
			Expect(t).
			Status(http.StatusUnauthorized).
			End()
	})

	t.Run("Get-OK", func(t *testing.T) {
		var args []middleware.RouteArgs
		args = append(args, middleware.RouteArgs{Key: "userId", Value: 4})

		handler := middleware.SetMuxVars(testHandler.GetUserId, args)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := session.NewMockAuthCheckerClient(ctrl)
		testHandler.AuthClient = m

		apitest.New("Get-OK").
			Handler(handler).
			Method("Get").
			URL("/api/me/").
			Expect(t).
			Status(http.StatusOK).
			End()
	})

	t.Run("CreateUserSubscription", func(t *testing.T) {
		var args []middleware.RouteArgs
		args = append(args, middleware.RouteArgs{Key: "userId", Value: 4})
		args = append(args, middleware.RouteArgs{Key: middleware.CSRFValid, Value: true})

		handler := middleware.SetMuxVars(testHandler.CreateUserSubscription, args)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := session.NewMockAuthCheckerClient(ctrl)
		testHandler.AuthClient = m

		p := profile.NewMockUseCase(ctrl)
		testHandler.ProfileUC = p

		p.EXPECT().CreateSubscription(4, 10).Return(1, nil)
		testCredit := models.UserSubscription{TargetId: 10}

		testCreditJSON, err := json.Marshal(testCredit)
		assert.NoError(t, err)

		apitest.New("Get-OK").
			Handler(handler).
			Method("Get").
			URL("/api/me/").
			Body(string(testCreditJSON)).
			Expect(t).
			Status(http.StatusOK).
			End()
	})

	t.Run("CreateUserSubscriptionErr1", func(t *testing.T) {
		var args []middleware.RouteArgs
		args = append(args, middleware.RouteArgs{Key: "userId", Value: 4})
		args = append(args, middleware.RouteArgs{Key: middleware.CSRFValid, Value: true})

		handler := middleware.SetMuxVars(testHandler.CreateUserSubscription, args)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := session.NewMockAuthCheckerClient(ctrl)
		testHandler.AuthClient = m

		p := profile.NewMockUseCase(ctrl)
		testHandler.ProfileUC = p

		p.EXPECT().CreateSubscription(4, 10).Return(1, errors.New("err"))
		testCredit := models.UserSubscription{TargetId: 10}

		testCreditJSON, err := json.Marshal(testCredit)
		assert.NoError(t, err)

		apitest.New("Get-OK").
			Handler(handler).
			Method("Get").
			URL("/api/me/").
			Body(string(testCreditJSON)).
			Expect(t).
			Status(http.StatusInternalServerError).
			End()
	})

	t.Run("CreateUserSubscriptionErr2", func(t *testing.T) {
		var args []middleware.RouteArgs
		args = append(args, middleware.RouteArgs{Key: "userId", Value: 4})
		args = append(args, middleware.RouteArgs{Key: middleware.CSRFValid, Value: true})

		handler := middleware.SetMuxVars(testHandler.CreateUserSubscription, args)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := session.NewMockAuthCheckerClient(ctrl)
		testHandler.AuthClient = m

		p := profile.NewMockUseCase(ctrl)
		testHandler.ProfileUC = p

		p.EXPECT().CreateSubscription(4, 10).Return(1, profile.ErrUserNonExistent)
		testCredit := models.UserSubscription{TargetId: 10}

		testCreditJSON, err := json.Marshal(testCredit)
		assert.NoError(t, err)

		apitest.New("Get-OK").
			Handler(handler).
			Method("Get").
			URL("/api/me/").
			Body(string(testCreditJSON)).
			Expect(t).
			Status(http.StatusBadRequest).
			End()
	})

	t.Run("CreateUserSubscriptionErr3", func(t *testing.T) {
		var args []middleware.RouteArgs
		args = append(args, middleware.RouteArgs{Key: "userId", Value: 4})
		args = append(args, middleware.RouteArgs{Key: middleware.CSRFValid, Value: true})

		handler := middleware.SetMuxVars(testHandler.CreateUserSubscription, args)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := session.NewMockAuthCheckerClient(ctrl)
		testHandler.AuthClient = m

		p := profile.NewMockUseCase(ctrl)
		testHandler.ProfileUC = p

		apitest.New("Get-OK").
			Handler(handler).
			Method("Get").
			URL("/api/me/").
			Expect(t).
			Status(http.StatusBadRequest).
			End()
	})

	t.Run("CreateUserSubscriptionErr4", func(t *testing.T) {
		var args []middleware.RouteArgs
		args = append(args, middleware.RouteArgs{Key: "userId", Value: 4})

		handler := middleware.SetMuxVars(testHandler.CreateUserSubscription, args)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		apitest.New("Get-OK").
			Handler(handler).
			Method("Get").
			URL("/api/me/").
			Expect(t).
			Status(http.StatusUnauthorized).
			End()
	})

	t.Run("CreateUserSubscriptionErr5", func(t *testing.T) {
		var args []middleware.RouteArgs

		handler := middleware.SetMuxVars(testHandler.CreateUserSubscription, args)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		apitest.New("Get-OK").
			Handler(handler).
			Method("Get").
			URL("/api/me/").
			Expect(t).
			Status(http.StatusUnauthorized).
			End()
	})

	t.Run("RemoveUserSubscription", func(t *testing.T) {
		var args []middleware.RouteArgs
		args = append(args, middleware.RouteArgs{Key: "userId", Value: 4})
		args = append(args, middleware.RouteArgs{Key: middleware.CSRFValid, Value: true})

		handler := middleware.SetMuxVars(testHandler.RemoveUserSubscription, args)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := session.NewMockAuthCheckerClient(ctrl)
		testHandler.AuthClient = m

		p := profile.NewMockUseCase(ctrl)
		testHandler.ProfileUC = p

		p.EXPECT().RemoveSubscription(4, 10).Return(nil)

		testCredit := models.UserSubscription{TargetId: 10}

		testCreditJSON, err := json.Marshal(testCredit)
		assert.NoError(t, err)

		apitest.New("Get-OK").
			Handler(handler).
			Method("Get").
			URL("/api/me/").
			Body(string(testCreditJSON)).
			Expect(t).
			Status(http.StatusOK).
			End()
	})

	t.Run("RemoveUserSubscriptionErr1", func(t *testing.T) {
		var args []middleware.RouteArgs
		args = append(args, middleware.RouteArgs{Key: "userId", Value: 4})
		args = append(args, middleware.RouteArgs{Key: middleware.CSRFValid, Value: true})

		handler := middleware.SetMuxVars(testHandler.RemoveUserSubscription, args)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := session.NewMockAuthCheckerClient(ctrl)
		testHandler.AuthClient = m

		p := profile.NewMockUseCase(ctrl)
		testHandler.ProfileUC = p

		p.EXPECT().RemoveSubscription(4, 10).Return(errors.New("Err"))

		testCredit := models.UserSubscription{TargetId: 10}

		testCreditJSON, err := json.Marshal(testCredit)
		assert.NoError(t, err)

		apitest.New("Get-OK").
			Handler(handler).
			Method("Get").
			URL("/api/me/").
			Body(string(testCreditJSON)).
			Expect(t).
			Status(http.StatusInternalServerError).
			End()
	})

	t.Run("RemoveUserSubscriptionErr1", func(t *testing.T) {
		var args []middleware.RouteArgs
		args = append(args, middleware.RouteArgs{Key: "userId", Value: 4})
		args = append(args, middleware.RouteArgs{Key: middleware.CSRFValid, Value: true})

		handler := middleware.SetMuxVars(testHandler.RemoveUserSubscription, args)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := session.NewMockAuthCheckerClient(ctrl)
		testHandler.AuthClient = m

		p := profile.NewMockUseCase(ctrl)
		testHandler.ProfileUC = p

		p.EXPECT().RemoveSubscription(4, 10).Return(errors.New("Err"))

		testCredit := models.UserSubscription{TargetId: 10}

		testCreditJSON, err := json.Marshal(testCredit)
		assert.NoError(t, err)

		apitest.New("Get-OK").
			Handler(handler).
			Method("Get").
			URL("/api/me/").
			Body(string(testCreditJSON)).
			Expect(t).
			Status(http.StatusInternalServerError).
			End()
	})

	t.Run("RemoveUserSubscriptionErr2", func(t *testing.T) {
		var args []middleware.RouteArgs
		args = append(args, middleware.RouteArgs{Key: "userId", Value: 4})
		args = append(args, middleware.RouteArgs{Key: middleware.CSRFValid, Value: true})

		handler := middleware.SetMuxVars(testHandler.RemoveUserSubscription, args)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := session.NewMockAuthCheckerClient(ctrl)
		testHandler.AuthClient = m

		apitest.New("Get-OK").
			Handler(handler).
			Method("Get").
			URL("/api/me/").
			Expect(t).
			Status(http.StatusBadRequest).
			End()
	})

	t.Run("RemoveUserSubscriptionErr3", func(t *testing.T) {
		var args []middleware.RouteArgs
		args = append(args, middleware.RouteArgs{Key: "userId", Value: 4})

		handler := middleware.SetMuxVars(testHandler.RemoveUserSubscription, args)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := session.NewMockAuthCheckerClient(ctrl)
		testHandler.AuthClient = m

		apitest.New("Get-OK").
			Handler(handler).
			Method("Get").
			URL("/api/me/").
			Expect(t).
			Status(http.StatusUnauthorized).
			End()
	})

	t.Run("RemoveUserSubscriptionErr3", func(t *testing.T) {
		var args []middleware.RouteArgs

		handler := middleware.SetMuxVars(testHandler.RemoveUserSubscription, args)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := session.NewMockAuthCheckerClient(ctrl)
		testHandler.AuthClient = m

		apitest.New("Get-OK").
			Handler(handler).
			Method("Get").
			URL("/api/me/").
			Expect(t).
			Status(http.StatusUnauthorized).
			End()
	})
}

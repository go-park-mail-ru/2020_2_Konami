package http

import (
	"encoding/json"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/steinfletcher/apitest"
	"gorm.io/gorm"
	"konami_backend/internal/pkg/message"
	messageRepoPkg "konami_backend/internal/pkg/message/repository"
	messageUseCasePkg "konami_backend/internal/pkg/message/usecase"
	"konami_backend/internal/pkg/middleware"
	"konami_backend/internal/pkg/models"
	"net/http"
	"testing"
)

var testHandler MessageHandler

func TestSessions(t *testing.T) {
	t.Run("SendMes", func(t *testing.T) {
		var args []middleware.RouteArgs
		args = append(args, middleware.RouteArgs{Key: middleware.UserID, Value: 4})
		args = append(args, middleware.RouteArgs{Key: middleware.CSRFValid, Value: true})

		handler := middleware.SetMuxVars(testHandler.SendMessage, args)

		msg := &models.Message{
			Id:        2,
			AuthorId:  4,
			MeetingId: 4,
			Text:      "fdsfsfs",
			Timestamp: "fslnfslkfs",
		}
		testUpdJSON, _ := json.Marshal(msg)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := message.NewMockUseCase(ctrl)
		testHandler.MessageUC = m

		testHandler.MaxReqSize = 10000

		m.EXPECT().CreateMessage(*msg).Return(0, nil)

		apitest.New("Get-All-Ok").
			Handler(handler).
			Method("Get").
			URL("/people").
			Body(string(testUpdJSON)).
			Expect(t).
			Status(http.StatusCreated).
			End()
	})

	t.Run("SendMesBad1", func(t *testing.T) {
		var args []middleware.RouteArgs
		args = append(args, middleware.RouteArgs{Key: middleware.UserID, Value: 4})
		args = append(args, middleware.RouteArgs{Key: middleware.CSRFValid, Value: true})

		handler := middleware.SetMuxVars(testHandler.SendMessage, args)

		msg := &models.Message{
			Id:        2,
			AuthorId:  4,
			MeetingId: 4,
			Text:      "fdsfsfs",
			Timestamp: "fslnfslkfs",
		}
		testUpdJSON, _ := json.Marshal(msg)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := message.NewMockUseCase(ctrl)
		testHandler.MessageUC = m

		testHandler.MaxReqSize = 10000

		m.EXPECT().CreateMessage(*msg).Return(0, errors.New("err"))

		apitest.New("Get-All-Ok").
			Handler(handler).
			Method("Get").
			URL("/people").
			Body(string(testUpdJSON)).
			Expect(t).
			Status(http.StatusInternalServerError).
			End()
	})

	t.Run("SendMesBad2", func(t *testing.T) {
		var args []middleware.RouteArgs
		args = append(args, middleware.RouteArgs{Key: middleware.UserID, Value: 4})
		args = append(args, middleware.RouteArgs{Key: middleware.CSRFValid, Value: true})

		handler := middleware.SetMuxVars(testHandler.SendMessage, args)

		msg := &models.Message{
			Id:        2,
			AuthorId:  4,
			MeetingId: 4,
			Text:      "fdsfsfs",
			Timestamp: "fslnfslkfs",
		}
		testUpdJSON, _ := json.Marshal(msg)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := message.NewMockUseCase(ctrl)
		testHandler.MessageUC = m

		testHandler.MaxReqSize = 0

		apitest.New("Get-All-Ok").
			Handler(handler).
			Method("Get").
			URL("/people").
			Body(string(testUpdJSON)).
			Expect(t).
			Status(http.StatusBadRequest).
			End()
	})

	t.Run("SendMesBad3", func(t *testing.T) {
		var args []middleware.RouteArgs
		args = append(args, middleware.RouteArgs{Key: middleware.UserID, Value: 4})

		handler := middleware.SetMuxVars(testHandler.SendMessage, args)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := message.NewMockUseCase(ctrl)
		testHandler.MessageUC = m

		apitest.New("Get-All-Ok").
			Handler(handler).
			Method("Get").
			URL("/people").
			Expect(t).
			Status(http.StatusUnauthorized).
			End()
	})


	t.Run("SendMesBad4", func(t *testing.T) {
		var args []middleware.RouteArgs
		handler := middleware.SetMuxVars(testHandler.SendMessage, args)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := message.NewMockUseCase(ctrl)
		testHandler.MessageUC = m

		apitest.New("Get-All-Ok").
			Handler(handler).
			Method("Get").
			URL("/people").
			Expect(t).
			Status(http.StatusUnauthorized).
			End()
	})

	t.Run("SendMesBad3", func(t *testing.T) {
		db := &gorm.DB{}
		msgRepo := messageRepoPkg.NewMeetingGormRepo(db)
		msgUC := messageUseCasePkg.NewMessageUseCase(msgRepo)
		_ = NewMessageHandler(msgUC, nil, 0)
	})

	t.Run("GetMessage", func(t *testing.T) {
		var args []middleware.QueryArgs
		args = append(args, middleware.QueryArgs{Key: "meetId", Value: "4"})
		handler := middleware.SetVars(testHandler.GetMessages, args)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := message.NewMockUseCase(ctrl)
		testHandler.MessageUC = m

		m.EXPECT().GetMessages(4).Return([]models.Message{}, nil)
		apitest.New("Get-All-Ok").
			Handler(handler).
			Method("Get").
			URL("/people").
			Expect(t).
			Status(http.StatusOK).
			End()
	})

	t.Run("GetMessageBad1", func(t *testing.T) {
		var args []middleware.QueryArgs
		args = append(args, middleware.QueryArgs{Key: "meetId", Value: "4"})
		handler := middleware.SetVars(testHandler.GetMessages, args)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := message.NewMockUseCase(ctrl)
		testHandler.MessageUC = m

		m.EXPECT().GetMessages(4).Return([]models.Message{}, errors.New("Err"))
		apitest.New("Get-All-Ok").
			Handler(handler).
			Method("Get").
			URL("/people").
			Expect(t).
			Status(http.StatusInternalServerError).
			End()
	})

	t.Run("GetMessageBad1", func(t *testing.T) {
		var args []middleware.QueryArgs
		handler := middleware.SetVars(testHandler.GetMessages, args)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := message.NewMockUseCase(ctrl)
		testHandler.MessageUC = m

		apitest.New("Get-All-Ok").
			Handler(handler).
			Method("Get").
			URL("/people").
			Expect(t).
			Status(http.StatusBadRequest).
			End()
	})
}
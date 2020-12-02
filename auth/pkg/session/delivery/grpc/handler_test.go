package grpc

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"konami_backend/auth/pkg/session"
	"konami_backend/proto/auth"
	"testing"
)

var testHandler SessionHandler

func TestGRPC(t *testing.T) {
	t.Run("GRPCCreate", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := session.NewMockUseCase(ctrl)
		testHandler = NewSessionHandler(m)

		var id int64
		id = 123

		testStr := "TOK"
		m.EXPECT().CreateSession(id).Return(testStr, nil)
		_, _ = testHandler.Create(context.Background(), &auth.Session{UserId: id})
	})

	t.Run("GRPCCreateBad", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := session.NewMockUseCase(ctrl)
		testHandler = NewSessionHandler(m)

		var id int64
		id = 123

		testStr := "TOK"
		m.EXPECT().CreateSession(id).Return(testStr, errors.New("err"))
		_, _ = testHandler.Create(context.Background(), &auth.Session{UserId: id})
	})

	t.Run("GRPCCheak", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := session.NewMockUseCase(ctrl)
		testHandler = NewSessionHandler(m)

		var id int64
		id = 123

		testStr := "TOK"
		m.EXPECT().GetUserId(testStr).Return(id, nil)
		_, _ = testHandler. Check(context.Background(), &auth.SessionToken{Token: "TOK"})
	})

	t.Run("GRPCCheakBad1", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := session.NewMockUseCase(ctrl)
		testHandler = NewSessionHandler(m)

		var id int64
		id = 123

		testStr := "TOK"
		m.EXPECT().GetUserId(testStr).Return(id, session.ErrSessionNotFound)
		_, _ = testHandler. Check(context.Background(), &auth.SessionToken{Token: "TOK"})
	})

	t.Run("GRPCCheakBad1", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := session.NewMockUseCase(ctrl)
		testHandler = NewSessionHandler(m)

		var id int64
		id = 123

		testStr := "TOK"
		m.EXPECT().GetUserId(testStr).Return(id, errors.New("Err"))
		_, _ = testHandler. Check(context.Background(), &auth.SessionToken{Token: "TOK"})
	})

	t.Run("GRPCDelete", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := session.NewMockUseCase(ctrl)
		testHandler = NewSessionHandler(m)

		testStr := "TOK"
		m.EXPECT().RemoveSession(testStr).Return(nil)
		_, _ = testHandler.Delete(context.Background(), &auth.SessionToken{Token: "TOK"})
	})

	t.Run("GRPCDelete", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := session.NewMockUseCase(ctrl)
		testHandler = NewSessionHandler(m)

		testStr := "TOK"
		m.EXPECT().RemoveSession(testStr).Return(errors.New("Err"))
		_, _ = testHandler.Delete(context.Background(), &auth.SessionToken{Token: "TOK"})
	})

	t.Run("GRPCDelete", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := session.NewMockUseCase(ctrl)
		testHandler = NewSessionHandler(m)

		testStr := "TOK"
		m.EXPECT().RemoveSession(testStr).Return(nil)
		_, _ = testHandler.Delete(context.Background(), &auth.SessionToken{Token: "TOK"})
	})
}
package grpc

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/sirupsen/logrus"
	"konami_backend/csrf/pkg/csrf"
	loggerPkg "konami_backend/logger"
	csrfProto "konami_backend/proto/csrf"
	"os"
	"testing"
)

var testHandler CsrfHandler

func TestCSRF(t *testing.T) {
	t.Run("GRPCCreate", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		logger := loggerPkg.NewLogger(os.Stdout)
		logger.SetLevel(logrus.TraceLevel)

		m := csrf.NewMockUseCase(ctrl)
		testHandler = NewCsrfHandler(m, logger)

		data := csrfProto.CsrfData{
			Sid:       "040",
			TimeStamp: 404,
		}

		testStr := "TOK"
		m.EXPECT().Create(data.Sid, data.TimeStamp).Return(testStr, nil)
		_, _ = testHandler.Create(context.Background(), &data)
	})

	t.Run("GRPCCreateErr", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		logger := loggerPkg.NewLogger(os.Stdout)
		logger.SetLevel(logrus.TraceLevel)

		m := csrf.NewMockUseCase(ctrl)
		testHandler = NewCsrfHandler(m, logger)

		data := csrfProto.CsrfData{
			Sid:       "040",
			TimeStamp: 404,
		}

		testStr := "TOK"
		m.EXPECT().Create(data.Sid, data.TimeStamp).Return(testStr, csrf.ErrExpiredToken)
		_, _ = testHandler.Create(context.Background(), &data)
	})

	t.Run("GRPCCreateErr2", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		logger := loggerPkg.NewLogger(os.Stdout)
		logger.SetLevel(logrus.TraceLevel)

		m := csrf.NewMockUseCase(ctrl)
		testHandler = NewCsrfHandler(m, logger)

		data := csrfProto.CsrfData{
			Sid:       "040",
			TimeStamp: 404,
		}

		testStr := "TOK"
		m.EXPECT().Create(data.Sid, data.TimeStamp).Return(testStr, errors.New("Err"))
		_, _ = testHandler.Create(context.Background(), &data)
	})

	t.Run("GRPCCheck", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		logger := loggerPkg.NewLogger(os.Stdout)
		logger.SetLevel(logrus.TraceLevel)

		m := csrf.NewMockUseCase(ctrl)
		testHandler = NewCsrfHandler(m, logger)

		data := csrfProto.CsrfToken{
			Sid:   "123",
			Token: "321",
		}

		m.EXPECT().Check(data.Sid, data.Token).Return(true, nil)
		_, _ = testHandler.Check(context.Background(), &data)
	})

	t.Run("GRPCCheckBad1", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		logger := loggerPkg.NewLogger(os.Stdout)
		logger.SetLevel(logrus.TraceLevel)

		m := csrf.NewMockUseCase(ctrl)
		testHandler = NewCsrfHandler(m, logger)

		data := csrfProto.CsrfToken{
			Sid:   "123",
			Token: "321",
		}

		m.EXPECT().Check(data.Sid, data.Token).Return(false, nil)
		_, _ = testHandler.Check(context.Background(), &data)
	})

	t.Run("GRPCCheckBad2", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		logger := loggerPkg.NewLogger(os.Stdout)
		logger.SetLevel(logrus.TraceLevel)

		m := csrf.NewMockUseCase(ctrl)
		testHandler = NewCsrfHandler(m, logger)

		data := csrfProto.CsrfToken{
			Sid:   "123",
			Token: "321",
		}

		m.EXPECT().Check(data.Sid, data.Token).Return(false, errors.New("Err"))
		_, _ = testHandler.Check(context.Background(), &data)
	})
}

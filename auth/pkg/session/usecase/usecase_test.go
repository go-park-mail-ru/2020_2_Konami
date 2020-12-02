package usecase

import (
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"konami_backend/auth/pkg/session"

	"testing"
)

func TestTag(t *testing.T) {
	t.Run("TestOnUsedToken", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		tagRepo := session.NewMockRepository(ctrl)
		ta := NewSessionUseCase(tagRepo)

		tagRepo.EXPECT().GetUserId("gg")
		_, err := ta.GetUserId("gg")
		assert.NoError(t, err)

		tagRepo.EXPECT().CreateSession(134)
		_, err = ta.CreateSession(134)
		assert.NoError(t, err)

		tagRepo.EXPECT().RemoveSession("ggor")
		err = ta.RemoveSession("ggor")
		assert.NoError(t, err)
	})
}

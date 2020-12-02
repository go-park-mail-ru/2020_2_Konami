package usecase

import (
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"konami_backend/internal/pkg/message"
	"konami_backend/internal/pkg/models"
	"testing"
)

func TestTag(t *testing.T) {
	t.Run("TestOnUsedToken", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		tagRepo := message.NewMockRepository(ctrl)

		ta := NewMessageUseCase(tagRepo)

		tagRepo.EXPECT().GetMessages(0)
		_, err := ta.GetMessages(0)
		assert.NoError(t, err)

		gg := models.Message{
			Id:        0,
			AuthorId:  0,
			MeetingId: 0,
			Text:      "",
			Timestamp: "",
		}

		tagRepo.EXPECT().SaveMessage(gg)
		_, err = ta.CreateMessage(gg)
		assert.NoError(t, err)
	})
}

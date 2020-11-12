package usecase

import (
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"konami_backend/internal/pkg/tag"
	"testing"
)

func TestTag(t *testing.T) {
	t.Run("TestOnUsedToken", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		tagRepo := tag.NewMockRepository(ctrl)

		ta := NewTagUseCase(tagRepo)

		tagRepo.EXPECT().CreateTag("gg")
		_, err := ta.CreateTag("gg")
		assert.NoError(t, err)

		tagRepo.EXPECT().FilterTags("ggo")
		_, err = ta.FilterTags("ggo")
		assert.NoError(t, err)

		tagRepo.EXPECT().GetOrCreateTag("ggor")
		_, err = ta.GetOrCreateTag("ggor")
		assert.NoError(t, err)

		tagRepo.EXPECT().GetOrCreateTag("ww")
		_, err = ta.GetOrCreateTag("ww")
		assert.NoError(t, err)

		tagRepo.EXPECT().GetTagById(1)
		_, err = ta.GetTagById(1)
		assert.NoError(t, err)

		tagRepo.EXPECT().GetOrCreateTag("ww")
		_, err = ta.GetOrCreateTag("ww")
		assert.NoError(t, err)

		tagRepo.EXPECT().GetTagByName("wrrw")
		_, err = ta.GetTagByName("wrrw")
		assert.NoError(t, err)
	})
}

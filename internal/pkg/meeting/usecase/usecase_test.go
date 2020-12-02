package usecase

import (
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"konami_backend/internal/pkg/meeting"
	"konami_backend/internal/pkg/models"
	"konami_backend/internal/pkg/tag"
	uploadsHandlerPkg "konami_backend/internal/pkg/utils/uploads_handler"
	"testing"
)

func TestTag(t *testing.T) {
	t.Run("TestOnUsedToken", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mRep := meeting.NewMockRepository(ctrl)
		tagRep := tag.NewMockRepository(ctrl)

		uploadsHandler := uploadsHandlerPkg.NewUploadsHandler("uploadsDir")

		uc := NewMeetingUseCase(mRep, uploadsHandler, tagRep, "test", "test")

		mRep.EXPECT().GetMeeting(1, 1, true).
			Return(models.MeetingDetails{}, nil)
		_, err := uc.GetMeeting(1, 1, true)
		assert.NoError(t, err)

		mRep.EXPECT().UpdateMeeting(1, models.MeetingUpdate{}).
			Return(nil)
		err = uc.UpdateMeeting(1, models.MeetingUpdate{})
		assert.NoError(t, err)

		mRep.EXPECT().GetNextMeetings(meeting.FilterParams{}).
			Return([]models.Meeting{}, nil)
		_, err = uc.GetNextMeetings(meeting.FilterParams{})
		assert.NoError(t, err)

		mRep.EXPECT().GetTopMeetings(meeting.FilterParams{}).
			Return([]models.Meeting{}, nil)
		_, err = uc.GetTopMeetings(meeting.FilterParams{})
		assert.NoError(t, err)

		mRep.EXPECT().FilterLiked(meeting.FilterParams{}).
			Return([]models.Meeting{}, nil)
		_, err = uc.FilterLiked(meeting.FilterParams{})
		assert.NoError(t, err)

		mRep.EXPECT().FilterRegistered(meeting.FilterParams{}).
			Return([]models.Meeting{}, nil)
		_, err = uc.FilterRegistered(meeting.FilterParams{})
		assert.NoError(t, err)

		mRep.EXPECT().FilterRecommended(meeting.FilterParams{}).
			Return([]models.Meeting{}, nil)
		_, err = uc.FilterRecommended(meeting.FilterParams{})
		assert.NoError(t, err)

		mRep.EXPECT().FilterTagged(meeting.FilterParams{}, 1).
			Return([]models.Meeting{}, nil)
		_, err = uc.FilterTagged(meeting.FilterParams{}, 1)
		assert.NoError(t, err)

		mRep.EXPECT().FilterSimilar(meeting.FilterParams{}, 1).
			Return([]models.Meeting{}, nil)
		_, err = uc.FilterSimilar(meeting.FilterParams{}, 1)
		assert.NoError(t, err)

		mRep.EXPECT().SearchMeetings(meeting.FilterParams{}, "LOL",1).
			Return([]models.Meeting{}, nil)
		_, err = uc.SearchMeetings(meeting.FilterParams{}, "LOL", 1)
		assert.NoError(t, err)

		testStr := "Some data"
		testNumber := 6

		testModel := models.MeetingData{
			Address:   &testStr,
			City:      &testStr,
			Start:     &testStr,
			End:       &testStr,
			Text:      &testStr,
			Tags:      []string{"tag", "error"},
			Title:     &testStr,
			Photo:     nil,
			Seats:     &testNumber,
			SeatsLeft: &testNumber,
		}

		tagRep.EXPECT().GetOrCreateTag("tag").
			Return(models.Tag{}, nil)

		tagRep.EXPECT().GetOrCreateTag("error").
			Return(models.Tag{}, errors.New("err"))

		_, err = uc.CreateMeeting(3, testModel)
		assert.Error(t, err)

		testModel.Photo = &testStr
		_, err = uc.CreateMeeting(3, testModel)
		assert.Error(t, err)

		testModel.Title = nil
		_, err = uc.CreateMeeting(3, testModel)
		assert.Error(t, err)
	})
}

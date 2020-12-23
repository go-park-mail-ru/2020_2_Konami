package http

import (
	"encoding/json"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/steinfletcher/apitest"
	"konami_backend/internal/pkg/meeting"
	"konami_backend/internal/pkg/middleware"
	"konami_backend/internal/pkg/models"
	"net/http"
	"testing"
	"time"
)

var testHandler MeetingHandler

func TestSessions(t *testing.T) {
	t.Run("GetMeetingsList", func(t *testing.T) {
		var args []middleware.QueryArgs
		args = append(args, middleware.QueryArgs{Key: "start", Value: "2006-01-02"})
		args = append(args, middleware.QueryArgs{Key: "end", Value: "2007-01-02"})
		args = append(args, middleware.QueryArgs{Key: "prevId", Value: "3"})
		args = append(args, middleware.QueryArgs{Key: "limit", Value: "10"})

		var args2 []middleware.RouteArgs
		args2 = append(args2, middleware.RouteArgs{Key: middleware.UserID, Value: 4})
		handler := middleware.SetVarsAndMux(testHandler.GetMeetingsList, args, args2)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := meeting.NewMockUseCase(ctrl)
		testHandler.MeetingUC = m

		layout := "2006-01-02"
		time1, _ := time.Parse(layout, "2006-01-02")
		time2, _ := time.Parse(layout, "2007-01-02")

		m.EXPECT().GetNextMeetings(meeting.FilterParams{
			StartDate:  time1,
			EndDate:    time2,
			PrevId:     3,
			CountLimit: 10,
			UserId:     4,
			PrevLikes:  MaxLikes,
			PrevStart:  time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
		}).Return([]models.Meeting{}, nil)

		apitest.New("GetMeetingsList").
			Handler(handler).
			Method("Get").
			URL("/user").
			Expect(t).
			Status(http.StatusOK).
			End()
	})

	t.Run("GetMeetingsListError", func(t *testing.T) {
		var args []middleware.QueryArgs
		args = append(args, middleware.QueryArgs{Key: "start", Value: "2006-01-02"})
		args = append(args, middleware.QueryArgs{Key: "end", Value: "2007-01-02"})
		args = append(args, middleware.QueryArgs{Key: "prevId", Value: "3"})
		args = append(args, middleware.QueryArgs{Key: "limit", Value: "10"})

		var args2 []middleware.RouteArgs
		args2 = append(args2, middleware.RouteArgs{Key: middleware.UserID, Value: 4})
		handler := middleware.SetVarsAndMux(testHandler.GetMeetingsList, args, args2)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := meeting.NewMockUseCase(ctrl)
		testHandler.MeetingUC = m

		layout := "2006-01-02"
		time1, _ := time.Parse(layout, "2006-01-02")
		time2, _ := time.Parse(layout, "2007-01-02")

		m.EXPECT().GetNextMeetings(meeting.FilterParams{
			StartDate:  time1,
			EndDate:    time2,
			PrevId:     3,
			CountLimit: 10,
			UserId:     4,
			PrevLikes:  MaxLikes,
			PrevStart:  time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
		}).Return([]models.Meeting{}, errors.New("New err"))

		apitest.New("GetMeetingsList").
			Handler(handler).
			Method("Get").
			URL("/user").
			Expect(t).
			Status(http.StatusInternalServerError).
			End()
	})

	t.Run("GetUserMeet", func(t *testing.T) {
		var args []middleware.QueryArgs
		args = append(args, middleware.QueryArgs{Key: "start", Value: "2006-01-02"})
		args = append(args, middleware.QueryArgs{Key: "end", Value: "2007-01-02"})
		args = append(args, middleware.QueryArgs{Key: "prevId", Value: "3"})
		args = append(args, middleware.QueryArgs{Key: "limit", Value: "10"})

		var args2 []middleware.RouteArgs
		args2 = append(args2, middleware.RouteArgs{Key: middleware.UserID, Value: 4})
		handler := middleware.SetVarsAndMux(testHandler.GetUserMeetingsList, args, args2)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := meeting.NewMockUseCase(ctrl)
		testHandler.MeetingUC = m

		layout := "2006-01-02"
		time1, _ := time.Parse(layout, "2006-01-02")
		time2, _ := time.Parse(layout, "2007-01-02")

		m.EXPECT().FilterRegistered(meeting.FilterParams{
			StartDate:  time1,
			EndDate:    time2,
			PrevId:     3,
			CountLimit: 10,
			UserId:     4,
			PrevLikes:  MaxLikes,
			PrevStart:  time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
		}).Return([]models.Meeting{}, nil)

		apitest.New("GetMeetingsList").
			Handler(handler).
			Method("Get").
			URL("/user").
			Expect(t).
			Status(http.StatusOK).
			End()
	})

	t.Run("GetUserMeetErr1", func(t *testing.T) {
		var args []middleware.QueryArgs
		args = append(args, middleware.QueryArgs{Key: "start", Value: "2006-01-02"})
		args = append(args, middleware.QueryArgs{Key: "end", Value: "2007-01-02"})
		args = append(args, middleware.QueryArgs{Key: "prevId", Value: "3"})
		args = append(args, middleware.QueryArgs{Key: "limit", Value: "10"})

		var args2 []middleware.RouteArgs
		args2 = append(args2, middleware.RouteArgs{Key: middleware.UserID, Value: 4})
		handler := middleware.SetVarsAndMux(testHandler.GetUserMeetingsList, args, args2)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := meeting.NewMockUseCase(ctrl)
		testHandler.MeetingUC = m

		layout := "2006-01-02"
		time1, _ := time.Parse(layout, "2006-01-02")
		time2, _ := time.Parse(layout, "2007-01-02")

		m.EXPECT().FilterRegistered(meeting.FilterParams{
			StartDate:  time1,
			EndDate:    time2,
			PrevId:     3,
			CountLimit: 10,
			UserId:     4,
			PrevLikes:  MaxLikes,
			PrevStart:  time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
		}).Return([]models.Meeting{}, errors.New("Err"))

		apitest.New("GetMeetingsList").
			Handler(handler).
			Method("Get").
			URL("/user").
			Expect(t).
			Status(http.StatusInternalServerError).
			End()
	})

	t.Run("GetUserMeetErr2", func(t *testing.T) {
		var args []middleware.QueryArgs
		args = append(args, middleware.QueryArgs{Key: "start", Value: "2006-01-02"})
		args = append(args, middleware.QueryArgs{Key: "end", Value: "2007-01-02"})
		args = append(args, middleware.QueryArgs{Key: "prevId", Value: "3"})
		args = append(args, middleware.QueryArgs{Key: "limit", Value: "10"})

		var args2 []middleware.RouteArgs
		args2 = append(args2, middleware.RouteArgs{Key: middleware.UserID, Value: -1})
		handler := middleware.SetVarsAndMux(testHandler.GetUserMeetingsList, args, args2)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := meeting.NewMockUseCase(ctrl)
		testHandler.MeetingUC = m

		apitest.New("GetMeetingsList").
			Handler(handler).
			Method("Get").
			URL("/user").
			Expect(t).
			Status(http.StatusUnauthorized).
			End()
	})
	////
	t.Run("GetUserMeet", func(t *testing.T) {
		var args []middleware.QueryArgs
		args = append(args, middleware.QueryArgs{Key: "start", Value: "2006-01-02"})
		args = append(args, middleware.QueryArgs{Key: "end", Value: "2007-01-02"})
		args = append(args, middleware.QueryArgs{Key: "prevId", Value: "3"})
		args = append(args, middleware.QueryArgs{Key: "limit", Value: "10"})

		var args2 []middleware.RouteArgs
		args2 = append(args2, middleware.RouteArgs{Key: middleware.UserID, Value: 4})
		handler := middleware.SetVarsAndMux(testHandler.GetFavMeetingsList, args, args2)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := meeting.NewMockUseCase(ctrl)
		testHandler.MeetingUC = m

		layout := "2006-01-02"
		time1, _ := time.Parse(layout, "2006-01-02")
		time2, _ := time.Parse(layout, "2007-01-02")

		m.EXPECT().FilterLiked(meeting.FilterParams{
			StartDate:  time1,
			EndDate:    time2,
			PrevId:     3,
			CountLimit: 10,
			UserId:     4,
			PrevLikes:  MaxLikes,
			PrevStart:  time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
		}).Return([]models.Meeting{}, nil)

		apitest.New("GetMeetingsList").
			Handler(handler).
			Method("Get").
			URL("/user").
			Expect(t).
			Status(http.StatusOK).
			End()
	})

	t.Run("GetUserMeetErr1", func(t *testing.T) {
		var args []middleware.QueryArgs
		args = append(args, middleware.QueryArgs{Key: "start", Value: "2006-01-02"})
		args = append(args, middleware.QueryArgs{Key: "end", Value: "2007-01-02"})
		args = append(args, middleware.QueryArgs{Key: "prevId", Value: "3"})
		args = append(args, middleware.QueryArgs{Key: "limit", Value: "10"})

		var args2 []middleware.RouteArgs
		args2 = append(args2, middleware.RouteArgs{Key: middleware.UserID, Value: 4})
		handler := middleware.SetVarsAndMux(testHandler.GetFavMeetingsList, args, args2)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := meeting.NewMockUseCase(ctrl)
		testHandler.MeetingUC = m

		layout := "2006-01-02"
		time1, _ := time.Parse(layout, "2006-01-02")
		time2, _ := time.Parse(layout, "2007-01-02")

		m.EXPECT().FilterLiked(meeting.FilterParams{
			StartDate:  time1,
			EndDate:    time2,
			PrevId:     3,
			CountLimit: 10,
			UserId:     4,
			PrevLikes:  MaxLikes,
			PrevStart:  time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
		}).Return([]models.Meeting{}, errors.New("Err"))

		apitest.New("GetMeetingsList").
			Handler(handler).
			Method("Get").
			URL("/user").
			Expect(t).
			Status(http.StatusInternalServerError).
			End()
	})

	t.Run("GetUserMeetErr2", func(t *testing.T) {
		var args []middleware.QueryArgs
		args = append(args, middleware.QueryArgs{Key: "start", Value: "2006-01-02"})
		args = append(args, middleware.QueryArgs{Key: "end", Value: "2007-01-02"})
		args = append(args, middleware.QueryArgs{Key: "prevId", Value: "3"})
		args = append(args, middleware.QueryArgs{Key: "limit", Value: "10"})

		var args2 []middleware.RouteArgs
		args2 = append(args2, middleware.RouteArgs{Key: middleware.UserID, Value: -1})
		handler := middleware.SetVarsAndMux(testHandler.GetFavMeetingsList, args, args2)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := meeting.NewMockUseCase(ctrl)
		testHandler.MeetingUC = m

		apitest.New("GetMeetingsList").
			Handler(handler).
			Method("Get").
			URL("/user").
			Expect(t).
			Status(http.StatusUnauthorized).
			End()
	})

	t.Run("GetUserMeet", func(t *testing.T) {
		var args []middleware.QueryArgs
		args = append(args, middleware.QueryArgs{Key: "start", Value: "2006-01-02"})
		args = append(args, middleware.QueryArgs{Key: "end", Value: "2007-01-02"})
		args = append(args, middleware.QueryArgs{Key: "prevId", Value: "3"})
		args = append(args, middleware.QueryArgs{Key: "limit", Value: "10"})

		var args2 []middleware.RouteArgs
		args2 = append(args2, middleware.RouteArgs{Key: middleware.UserID, Value: 4})
		handler := middleware.SetVarsAndMux(testHandler.GetTopMeetingsList, args, args2)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := meeting.NewMockUseCase(ctrl)
		testHandler.MeetingUC = m

		layout := "2006-01-02"
		time1, _ := time.Parse(layout, "2006-01-02")
		time2, _ := time.Parse(layout, "2007-01-02")

		m.EXPECT().GetTopMeetings(meeting.FilterParams{
			StartDate:  time1,
			EndDate:    time2,
			PrevId:     3,
			CountLimit: 10,
			UserId:     4,
			PrevLikes:  MaxLikes,
			PrevStart:  time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
		}).Return([]models.Meeting{}, nil)

		apitest.New("GetMeetingsList").
			Handler(handler).
			Method("Get").
			URL("/user").
			Expect(t).
			Status(http.StatusOK).
			End()
	})

	t.Run("GetUserMeetErr1", func(t *testing.T) {
		var args []middleware.QueryArgs
		args = append(args, middleware.QueryArgs{Key: "start", Value: "2006-01-02"})
		args = append(args, middleware.QueryArgs{Key: "end", Value: "2007-01-02"})
		args = append(args, middleware.QueryArgs{Key: "prevId", Value: "3"})
		args = append(args, middleware.QueryArgs{Key: "limit", Value: "10"})

		var args2 []middleware.RouteArgs
		args2 = append(args2, middleware.RouteArgs{Key: middleware.UserID, Value: 4})
		handler := middleware.SetVarsAndMux(testHandler.GetTopMeetingsList, args, args2)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := meeting.NewMockUseCase(ctrl)
		testHandler.MeetingUC = m

		layout := "2006-01-02"
		time1, _ := time.Parse(layout, "2006-01-02")
		time2, _ := time.Parse(layout, "2007-01-02")

		m.EXPECT().GetTopMeetings(meeting.FilterParams{
			StartDate:  time1,
			EndDate:    time2,
			PrevId:     3,
			CountLimit: 10,
			UserId:     4,
			PrevLikes:  MaxLikes,
			PrevStart:  time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
		}).Return([]models.Meeting{}, errors.New("Err"))

		apitest.New("GetMeetingsList").
			Handler(handler).
			Method("Get").
			URL("/user").
			Expect(t).
			Status(http.StatusInternalServerError).
			End()
	})
	////
	t.Run("GetUserMeet", func(t *testing.T) {
		var args []middleware.QueryArgs
		args = append(args, middleware.QueryArgs{Key: "start", Value: "2006-01-02"})
		args = append(args, middleware.QueryArgs{Key: "end", Value: "2007-01-02"})
		args = append(args, middleware.QueryArgs{Key: "prevId", Value: "3"})
		args = append(args, middleware.QueryArgs{Key: "limit", Value: "10"})

		var args2 []middleware.RouteArgs
		args2 = append(args2, middleware.RouteArgs{Key: middleware.UserID, Value: 4})
		handler := middleware.SetVarsAndMux(testHandler.GetRecommendedList, args, args2)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := meeting.NewMockUseCase(ctrl)
		testHandler.MeetingUC = m

		layout := "2006-01-02"
		time1, _ := time.Parse(layout, "2006-01-02")
		time2, _ := time.Parse(layout, "2007-01-02")

		m.EXPECT().FilterRecommended(meeting.FilterParams{
			StartDate:  time1,
			EndDate:    time2,
			PrevId:     3,
			CountLimit: 10,
			UserId:     4,
			PrevLikes:  MaxLikes,
			PrevStart:  time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
		}).Return([]models.Meeting{}, nil)

		apitest.New("GetMeetingsList").
			Handler(handler).
			Method("Get").
			URL("/user").
			Expect(t).
			Status(http.StatusOK).
			End()
	})

	t.Run("GetUserMeetErr1", func(t *testing.T) {
		var args []middleware.QueryArgs
		args = append(args, middleware.QueryArgs{Key: "start", Value: "2006-01-02"})
		args = append(args, middleware.QueryArgs{Key: "end", Value: "2007-01-02"})
		args = append(args, middleware.QueryArgs{Key: "prevId", Value: "3"})
		args = append(args, middleware.QueryArgs{Key: "limit", Value: "10"})

		var args2 []middleware.RouteArgs
		args2 = append(args2, middleware.RouteArgs{Key: middleware.UserID, Value: 4})
		handler := middleware.SetVarsAndMux(testHandler.GetRecommendedList, args, args2)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := meeting.NewMockUseCase(ctrl)
		testHandler.MeetingUC = m

		layout := "2006-01-02"
		time1, _ := time.Parse(layout, "2006-01-02")
		time2, _ := time.Parse(layout, "2007-01-02")

		m.EXPECT().FilterRecommended(meeting.FilterParams{
			StartDate:  time1,
			EndDate:    time2,
			PrevId:     3,
			CountLimit: 10,
			UserId:     4,
			PrevLikes:  MaxLikes,
			PrevStart:  time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
		}).Return([]models.Meeting{}, errors.New("Err"))

		apitest.New("GetMeetingsList").
			Handler(handler).
			Method("Get").
			URL("/user").
			Expect(t).
			Status(http.StatusInternalServerError).
			End()
	})

	t.Run("GetUserMeetErr2", func(t *testing.T) {
		var args []middleware.QueryArgs
		args = append(args, middleware.QueryArgs{Key: "start", Value: "2006-01-02"})
		args = append(args, middleware.QueryArgs{Key: "end", Value: "2007-01-02"})
		args = append(args, middleware.QueryArgs{Key: "prevId", Value: "3"})
		args = append(args, middleware.QueryArgs{Key: "limit", Value: "10"})

		var args2 []middleware.RouteArgs
		args2 = append(args2, middleware.RouteArgs{Key: middleware.UserID, Value: -1})
		handler := middleware.SetVarsAndMux(testHandler.GetRecommendedList, args, args2)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := meeting.NewMockUseCase(ctrl)
		testHandler.MeetingUC = m

		apitest.New("GetMeetingsList").
			Handler(handler).
			Method("Get").
			URL("/user").
			Expect(t).
			Status(http.StatusUnauthorized).
			End()
	})
	////
	t.Run("GetUserMeet", func(t *testing.T) {
		var args []middleware.QueryArgs
		args = append(args, middleware.QueryArgs{Key: "start", Value: "2006-01-02"})
		args = append(args, middleware.QueryArgs{Key: "end", Value: "2007-01-02"})
		args = append(args, middleware.QueryArgs{Key: "prevId", Value: "3"})
		args = append(args, middleware.QueryArgs{Key: "limit", Value: "10"})
		args = append(args, middleware.QueryArgs{Key: "tag", Value: "banana"})

		var args2 []middleware.RouteArgs
		args2 = append(args2, middleware.RouteArgs{Key: middleware.UserID, Value: 4})
		handler := middleware.SetVarsAndMux(testHandler.GetTaggedMeetings, args, args2)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := meeting.NewMockUseCase(ctrl)
		testHandler.MeetingUC = m

		layout := "2006-01-02"
		time1, _ := time.Parse(layout, "2006-01-02")
		time2, _ := time.Parse(layout, "2007-01-02")

		m.EXPECT().FilterTagged(meeting.FilterParams{
			StartDate:  time1,
			EndDate:    time2,
			PrevId:     3,
			CountLimit: 10,
			UserId:     4,
			PrevLikes:  MaxLikes,
			PrevStart:  time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
		}, []string{"banana"}).Return([]models.Meeting{}, nil)

		apitest.New("GetMeetingsList").
			Handler(handler).
			Method("Get").
			URL("/user").
			Expect(t).
			Status(http.StatusOK).
			End()
	})

	t.Run("GetUserMeetErr1", func(t *testing.T) {
		var args []middleware.QueryArgs
		args = append(args, middleware.QueryArgs{Key: "start", Value: "2006-01-02"})
		args = append(args, middleware.QueryArgs{Key: "end", Value: "2007-01-02"})
		args = append(args, middleware.QueryArgs{Key: "prevId", Value: "3"})
		args = append(args, middleware.QueryArgs{Key: "limit", Value: "10"})
		args = append(args, middleware.QueryArgs{Key: "tag", Value: "banana"})

		var args2 []middleware.RouteArgs
		args2 = append(args2, middleware.RouteArgs{Key: middleware.UserID, Value: 4})
		handler := middleware.SetVarsAndMux(testHandler.GetTaggedMeetings, args, args2)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := meeting.NewMockUseCase(ctrl)
		testHandler.MeetingUC = m

		layout := "2006-01-02"
		time1, _ := time.Parse(layout, "2006-01-02")
		time2, _ := time.Parse(layout, "2007-01-02")

		m.EXPECT().FilterTagged(meeting.FilterParams{
			StartDate:  time1,
			EndDate:    time2,
			PrevId:     3,
			CountLimit: 10,
			UserId:     4,
			PrevLikes:  MaxLikes,
			PrevStart:  time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
		}, []string{"banana"}).Return([]models.Meeting{}, errors.New("Err"))

		apitest.New("GetMeetingsList").
			Handler(handler).
			Method("Get").
			URL("/user").
			Expect(t).
			Status(http.StatusInternalServerError).
			End()
	})

	t.Run("GetUserMeetErr2", func(t *testing.T) {
		var args []middleware.QueryArgs
		args = append(args, middleware.QueryArgs{Key: "start", Value: "2006-01-02"})
		args = append(args, middleware.QueryArgs{Key: "end", Value: "2007-01-02"})
		args = append(args, middleware.QueryArgs{Key: "prevId", Value: "3"})
		args = append(args, middleware.QueryArgs{Key: "limit", Value: "10"})

		var args2 []middleware.RouteArgs
		args2 = append(args2, middleware.RouteArgs{Key: middleware.UserID, Value: -1})
		handler := middleware.SetVarsAndMux(testHandler.GetTaggedMeetings, args, args2)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := meeting.NewMockUseCase(ctrl)
		testHandler.MeetingUC = m

		apitest.New("GetMeetingsList").
			Handler(handler).
			Method("Get").
			URL("/user").
			Expect(t).
			Status(http.StatusBadRequest).
			End()
	})

	t.Run("GetUserMeet", func(t *testing.T) {
		var args []middleware.QueryArgs
		args = append(args, middleware.QueryArgs{Key: "start", Value: "2006-01-02"})
		args = append(args, middleware.QueryArgs{Key: "end", Value: "2007-01-02"})
		args = append(args, middleware.QueryArgs{Key: "prevId", Value: "3"})
		args = append(args, middleware.QueryArgs{Key: "limit", Value: "10"})
		args = append(args, middleware.QueryArgs{Key: "meetId", Value: "15"})

		var args2 []middleware.RouteArgs
		args2 = append(args2, middleware.RouteArgs{Key: middleware.UserID, Value: 4})
		handler := middleware.SetVarsAndMux(testHandler.GetAkinMeetings, args, args2)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := meeting.NewMockUseCase(ctrl)
		testHandler.MeetingUC = m

		layout := "2006-01-02"
		time1, _ := time.Parse(layout, "2006-01-02")
		time2, _ := time.Parse(layout, "2007-01-02")

		m.EXPECT().FilterSimilar(meeting.FilterParams{
			StartDate:  time1,
			EndDate:    time2,
			PrevId:     3,
			CountLimit: 10,
			UserId:     4,
			PrevLikes:  MaxLikes,
			PrevStart:  time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
		}, 15).Return([]models.Meeting{}, nil)

		apitest.New("GetMeetingsList").
			Handler(handler).
			Method("Get").
			URL("/user").
			Expect(t).
			Status(http.StatusOK).
			End()
	})

	t.Run("GetUserMeetErr1", func(t *testing.T) {
		var args []middleware.QueryArgs
		args = append(args, middleware.QueryArgs{Key: "start", Value: "2006-01-02"})
		args = append(args, middleware.QueryArgs{Key: "end", Value: "2007-01-02"})
		args = append(args, middleware.QueryArgs{Key: "prevId", Value: "3"})
		args = append(args, middleware.QueryArgs{Key: "limit", Value: "10"})
		args = append(args, middleware.QueryArgs{Key: "meetId", Value: "15"})

		var args2 []middleware.RouteArgs
		args2 = append(args2, middleware.RouteArgs{Key: middleware.UserID, Value: 4})
		handler := middleware.SetVarsAndMux(testHandler.GetAkinMeetings, args, args2)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := meeting.NewMockUseCase(ctrl)
		testHandler.MeetingUC = m

		layout := "2006-01-02"
		time1, _ := time.Parse(layout, "2006-01-02")
		time2, _ := time.Parse(layout, "2007-01-02")

		m.EXPECT().FilterSimilar(meeting.FilterParams{
			StartDate:  time1,
			EndDate:    time2,
			PrevId:     3,
			CountLimit: 10,
			UserId:     4,
			PrevLikes:  MaxLikes,
			PrevStart:  time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
		}, 15).Return([]models.Meeting{}, errors.New("Err"))

		apitest.New("GetMeetingsList").
			Handler(handler).
			Method("Get").
			URL("/user").
			Expect(t).
			Status(http.StatusInternalServerError).
			End()
	})

	t.Run("GetUserMeetErr2", func(t *testing.T) {
		var args []middleware.QueryArgs
		args = append(args, middleware.QueryArgs{Key: "start", Value: "2006-01-02"})
		args = append(args, middleware.QueryArgs{Key: "end", Value: "2007-01-02"})
		args = append(args, middleware.QueryArgs{Key: "prevId", Value: "3"})
		args = append(args, middleware.QueryArgs{Key: "limit", Value: "10"})

		var args2 []middleware.RouteArgs
		args2 = append(args2, middleware.RouteArgs{Key: middleware.UserID, Value: -1})
		handler := middleware.SetVarsAndMux(testHandler.GetAkinMeetings, args, args2)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := meeting.NewMockUseCase(ctrl)
		testHandler.MeetingUC = m

		apitest.New("GetMeetingsList").
			Handler(handler).
			Method("Get").
			URL("/user").
			Expect(t).
			Status(http.StatusBadRequest).
			End()
	})

	t.Run("GetUserMeet", func(t *testing.T) {
		var args []middleware.QueryArgs

		var args2 []middleware.RouteArgs
		args2 = append(args2, middleware.RouteArgs{Key: middleware.UserID, Value: 4})
		args2 = append(args2, middleware.RouteArgs{Key: middleware.CSRFValid, Value: true})
		handler := middleware.SetVarsAndMux(testHandler.CreateMeeting, args, args2)

		strTest := "lulz i wanna sleeeeeep"
		testHandler.MaxReqSize = 10000
		testUpd := &models.MeetingData{
			Address:   &strTest,
			City:      nil,
			Start:     nil,
			End:       nil,
			Text:      nil,
			Tags:      nil,
			Title:     nil,
			Photo:     nil,
			Seats:     nil,
			SeatsLeft: nil,
		}
		testUpdJSON, _ := json.Marshal(testUpd)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := meeting.NewMockUseCase(ctrl)
		testHandler.MeetingUC = m

		m.EXPECT().CreateMeeting(4, *testUpd).
			Return(4, nil)

		apitest.New("GetMeetingsList").
			Handler(handler).
			Method("Get").
			URL("/user").
			Body(string(testUpdJSON)).
			Expect(t).
			Status(http.StatusCreated).
			End()
	})

	t.Run("GetUserMeetErr1", func(t *testing.T) {
		var args []middleware.QueryArgs

		var args2 []middleware.RouteArgs
		args2 = append(args2, middleware.RouteArgs{Key: middleware.UserID, Value: 4})
		args2 = append(args2, middleware.RouteArgs{Key: middleware.CSRFValid, Value: true})
		handler := middleware.SetVarsAndMux(testHandler.CreateMeeting, args, args2)

		strTest := "lulz i wanna sleeeeeep"
		testHandler.MaxReqSize = 10000
		testUpd := &models.MeetingData{
			Address:   &strTest,
			City:      nil,
			Start:     nil,
			End:       nil,
			Text:      nil,
			Tags:      nil,
			Title:     nil,
			Photo:     nil,
			Seats:     nil,
			SeatsLeft: nil,
		}
		testUpdJSON, _ := json.Marshal(testUpd)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := meeting.NewMockUseCase(ctrl)
		testHandler.MeetingUC = m

		m.EXPECT().CreateMeeting(4, *testUpd).
			Return(4, errors.New("err"))

		apitest.New("GetMeetingsList").
			Handler(handler).
			Method("Get").
			URL("/user").
			Body(string(testUpdJSON)).
			Expect(t).
			Status(http.StatusInternalServerError).
			End()
	})

	t.Run("GetUserMeetErr2", func(t *testing.T) {
		var args []middleware.QueryArgs

		var args2 []middleware.RouteArgs
		args2 = append(args2, middleware.RouteArgs{Key: middleware.UserID, Value: 4})
		args2 = append(args2, middleware.RouteArgs{Key: middleware.CSRFValid, Value: true})
		handler := middleware.SetVarsAndMux(testHandler.CreateMeeting, args, args2)

		strTest := "lulz i wanna sleeeeeep"
		testHandler.MaxReqSize = 0
		testUpd := &models.MeetingData{
			Address:   &strTest,
			City:      nil,
			Start:     nil,
			End:       nil,
			Text:      nil,
			Tags:      nil,
			Title:     nil,
			Photo:     nil,
			Seats:     nil,
			SeatsLeft: nil,
		}
		testUpdJSON, _ := json.Marshal(testUpd)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := meeting.NewMockUseCase(ctrl)
		testHandler.MeetingUC = m

		apitest.New("GetMeetingsList").
			Handler(handler).
			Method("Get").
			URL("/user").
			Body(string(testUpdJSON)).
			Expect(t).
			Status(http.StatusBadRequest).
			End()
	})

	t.Run("GetUserMeetErr2", func(t *testing.T) {
		var args []middleware.QueryArgs

		var args2 []middleware.RouteArgs
		args2 = append(args2, middleware.RouteArgs{Key: middleware.UserID, Value: 4})
		handler := middleware.SetVarsAndMux(testHandler.CreateMeeting, args, args2)

		strTest := "lulz i wanna sleeeeeep"
		testHandler.MaxReqSize = 0
		testUpd := &models.MeetingData{
			Address:   &strTest,
			City:      nil,
			Start:     nil,
			End:       nil,
			Text:      nil,
			Tags:      nil,
			Title:     nil,
			Photo:     nil,
			Seats:     nil,
			SeatsLeft: nil,
		}
		testUpdJSON, _ := json.Marshal(testUpd)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := meeting.NewMockUseCase(ctrl)
		testHandler.MeetingUC = m

		apitest.New("GetMeetingsList").
			Handler(handler).
			Method("Get").
			URL("/user").
			Body(string(testUpdJSON)).
			Expect(t).
			Status(http.StatusUnauthorized).
			End()
	})

	t.Run("GetUserMeetErr3", func(t *testing.T) {
		var args []middleware.QueryArgs

		var args2 []middleware.RouteArgs
		handler := middleware.SetVarsAndMux(testHandler.CreateMeeting, args, args2)

		strTest := "lulz i wanna sleeeeeep"
		testHandler.MaxReqSize = 0
		testUpd := &models.MeetingData{
			Address:   &strTest,
			City:      nil,
			Start:     nil,
			End:       nil,
			Text:      nil,
			Tags:      nil,
			Title:     nil,
			Photo:     nil,
			Seats:     nil,
			SeatsLeft: nil,
		}
		testUpdJSON, _ := json.Marshal(testUpd)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := meeting.NewMockUseCase(ctrl)
		testHandler.MeetingUC = m

		apitest.New("GetMeetingsList").
			Handler(handler).
			Method("Get").
			URL("/user").
			Body(string(testUpdJSON)).
			Expect(t).
			Status(http.StatusUnauthorized).
			End()
	})

	t.Run("GetUserMeet", func(t *testing.T) {
		var args []middleware.QueryArgs
		args = append(args, middleware.QueryArgs{
			Key:   "meetId",
			Value: "1",
		})

		var args2 []middleware.RouteArgs
		args2 = append(args2, middleware.RouteArgs{Key: middleware.UserID, Value: 4})
		handler := middleware.SetVarsAndMux(testHandler.GetMeeting, args, args2)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := meeting.NewMockUseCase(ctrl)
		testHandler.MeetingUC = m

		m.EXPECT().GetMeeting(1, 4, true).
			Return(models.MeetingDetails{}, nil)

		apitest.New("GetMeetingsList").
			Handler(handler).
			Method("Get").
			URL("/user").
			Expect(t).
			Status(http.StatusOK).
			End()
	})

	t.Run("GetUserMeet", func(t *testing.T) {
		var args []middleware.QueryArgs
		args = append(args, middleware.QueryArgs{
			Key:   "meetId",
			Value: "1",
		})

		var args2 []middleware.RouteArgs
		handler := middleware.SetVarsAndMux(testHandler.GetMeeting, args, args2)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := meeting.NewMockUseCase(ctrl)
		testHandler.MeetingUC = m

		m.EXPECT().GetMeeting(1, -1, false).
			Return(models.MeetingDetails{}, errors.New("Err"))

		apitest.New("GetMeetingsList").
			Handler(handler).
			Method("Get").
			URL("/user").
			Expect(t).
			Status(http.StatusBadRequest).
			End()
	})

	t.Run("GetUserMeet", func(t *testing.T) {
		var args []middleware.QueryArgs

		var args2 []middleware.RouteArgs
		handler := middleware.SetVarsAndMux(testHandler.GetMeeting, args, args2)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := meeting.NewMockUseCase(ctrl)
		testHandler.MeetingUC = m

		apitest.New("GetMeetingsList").
			Handler(handler).
			Method("Get").
			URL("/user").
			Expect(t).
			Status(http.StatusNotFound).
			End()
	})

	t.Run("GetUserMeet", func(t *testing.T) {
		var args []middleware.QueryArgs

		var args2 []middleware.RouteArgs
		args2 = append(args2, middleware.RouteArgs{Key: middleware.UserID, Value: 4})
		args2 = append(args2, middleware.RouteArgs{Key: middleware.CSRFValid, Value: true})
		handler := middleware.SetVarsAndMux(testHandler.UpdateMeeting, args, args2)

		testUpd := &models.MeetingUpdate{
			MeetId: 1,
			Fields: nil,
		}

		testUpdJSON, _ := json.Marshal(testUpd)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := meeting.NewMockUseCase(ctrl)
		testHandler.MeetingUC = m

		m.EXPECT().UpdateMeeting(4, *testUpd).Return(nil)

		apitest.New("GetMeetingsList").
			Handler(handler).
			Method("Get").
			URL("/user").
			Body(string(testUpdJSON)).
			Expect(t).
			Status(http.StatusOK).
			End()
	})

	t.Run("GetUserMeet", func(t *testing.T) {
		var args []middleware.QueryArgs

		var args2 []middleware.RouteArgs
		args2 = append(args2, middleware.RouteArgs{Key: middleware.UserID, Value: 4})
		args2 = append(args2, middleware.RouteArgs{Key: middleware.CSRFValid, Value: true})
		handler := middleware.SetVarsAndMux(testHandler.UpdateMeeting, args, args2)

		testUpd := &models.MeetingUpdate{
			MeetId: 1,
			Fields: nil,
		}

		testUpdJSON, _ := json.Marshal(testUpd)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := meeting.NewMockUseCase(ctrl)
		testHandler.MeetingUC = m

		m.EXPECT().UpdateMeeting(4, *testUpd).Return(errors.New("Err"))

		apitest.New("GetMeetingsList").
			Handler(handler).
			Method("Get").
			URL("/user").
			Body(string(testUpdJSON)).
			Expect(t).
			Status(http.StatusBadRequest).
			End()
	})

	t.Run("GetUserMeet", func(t *testing.T) {
		var args []middleware.QueryArgs

		var args2 []middleware.RouteArgs
		args2 = append(args2, middleware.RouteArgs{Key: middleware.UserID, Value: 4})
		args2 = append(args2, middleware.RouteArgs{Key: middleware.CSRFValid, Value: true})
		handler := middleware.SetVarsAndMux(testHandler.UpdateMeeting, args, args2)

		testUpd := &models.MeetingUpdate{
			MeetId: 1,
			Fields: nil,
		}

		testUpdJSON, _ := json.Marshal(testUpd)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := meeting.NewMockUseCase(ctrl)
		testHandler.MeetingUC = m

		m.EXPECT().UpdateMeeting(4, *testUpd).Return(meeting.ErrMeetingNotFound)

		apitest.New("GetMeetingsList").
			Handler(handler).
			Method("Get").
			URL("/user").
			Body(string(testUpdJSON)).
			Expect(t).
			Status(http.StatusNotFound).
			End()
	})

	t.Run("GetUserMeet", func(t *testing.T) {
		var args []middleware.QueryArgs

		var args2 []middleware.RouteArgs
		args2 = append(args2, middleware.RouteArgs{Key: middleware.UserID, Value: 4})
		args2 = append(args2, middleware.RouteArgs{Key: middleware.CSRFValid, Value: true})
		handler := middleware.SetVarsAndMux(testHandler.UpdateMeeting, args, args2)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := meeting.NewMockUseCase(ctrl)
		testHandler.MeetingUC = m

		apitest.New("GetMeetingsList").
			Handler(handler).
			Method("Get").
			URL("/user").
			Expect(t).
			Status(http.StatusBadRequest).
			End()
	})

	t.Run("GetUserMeet", func(t *testing.T) {
		var args []middleware.QueryArgs

		var args2 []middleware.RouteArgs
		args2 = append(args2, middleware.RouteArgs{Key: middleware.UserID, Value: 4})
		handler := middleware.SetVarsAndMux(testHandler.UpdateMeeting, args, args2)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := meeting.NewMockUseCase(ctrl)
		testHandler.MeetingUC = m

		apitest.New("GetMeetingsList").
			Handler(handler).
			Method("Get").
			URL("/user").
			Expect(t).
			Status(http.StatusUnauthorized).
			End()
	})

	t.Run("GetUserMeet", func(t *testing.T) {
		var args []middleware.QueryArgs

		var args2 []middleware.RouteArgs
		handler := middleware.SetVarsAndMux(testHandler.UpdateMeeting, args, args2)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := meeting.NewMockUseCase(ctrl)
		testHandler.MeetingUC = m

		apitest.New("GetMeetingsList").
			Handler(handler).
			Method("Get").
			URL("/user").
			Expect(t).
			Status(http.StatusUnauthorized).
			End()
	})

	t.Run("GetMeetingsList", func(t *testing.T) {
		var args []middleware.QueryArgs
		args = append(args, middleware.QueryArgs{Key: "start", Value: "2006-01-02"})
		args = append(args, middleware.QueryArgs{Key: "end", Value: "2007-01-02"})
		args = append(args, middleware.QueryArgs{Key: "prevId", Value: "3"})

		args = append(args, middleware.QueryArgs{Key: "query", Value: "test"})
		args = append(args, middleware.QueryArgs{Key: "limit", Value: "-5"})

		var args2 []middleware.RouteArgs
		handler := middleware.SetVarsAndMux(testHandler.SearchMeetings, args, args2)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := meeting.NewMockUseCase(ctrl)
		testHandler.MeetingUC = m

		layout := "2006-01-02"
		time1, _ := time.Parse(layout, "2006-01-02")
		time2, _ := time.Parse(layout, "2007-01-02")

		m.EXPECT().SearchMeetings(meeting.FilterParams{
			StartDate:  time1,
			EndDate:    time2,
			PrevId:     3,
			CountLimit: 10,
			UserId:     -1,
			PrevLikes:  MaxLikes,
			PrevStart:  time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
		}, "test", -1).Return([]models.Meeting{}, nil)

		apitest.New("GetMeetingsList").
			Handler(handler).
			Method("Get").
			URL("/user").
			Expect(t).
			Status(http.StatusOK).
			End()
	})

	t.Run("GetMeetingsList", func(t *testing.T) {
		var args []middleware.QueryArgs
		args = append(args, middleware.QueryArgs{Key: "start", Value: "2006-01-02"})
		args = append(args, middleware.QueryArgs{Key: "end", Value: "2007-01-02"})
		args = append(args, middleware.QueryArgs{Key: "prevId", Value: "3"})

		args = append(args, middleware.QueryArgs{Key: "query", Value: "test"})
		args = append(args, middleware.QueryArgs{Key: "limit", Value: "-5"})

		var args2 []middleware.RouteArgs
		handler := middleware.SetVarsAndMux(testHandler.SearchMeetings, args, args2)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := meeting.NewMockUseCase(ctrl)
		testHandler.MeetingUC = m

		layout := "2006-01-02"
		time1, _ := time.Parse(layout, "2006-01-02")
		time2, _ := time.Parse(layout, "2007-01-02")

		m.EXPECT().SearchMeetings(meeting.FilterParams{
			StartDate:  time1,
			EndDate:    time2,
			PrevId:     3,
			CountLimit: 10,
			UserId:     -1,
			PrevLikes:  MaxLikes,
			PrevStart:  time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
		}, "test", -1).Return([]models.Meeting{}, errors.New("Err"))

		apitest.New("GetMeetingsList").
			Handler(handler).
			Method("Get").
			URL("/user").
			Expect(t).
			Status(http.StatusInternalServerError).
			End()
	})

	t.Run("GetMeetingsList", func(t *testing.T) {
		var args []middleware.QueryArgs
		args = append(args, middleware.QueryArgs{Key: "start", Value: "2006-01-02"})
		args = append(args, middleware.QueryArgs{Key: "end", Value: "2007-01-02"})
		args = append(args, middleware.QueryArgs{Key: "prevId", Value: "3"})

		args = append(args, middleware.QueryArgs{Key: "limit", Value: "-5"})

		var args2 []middleware.RouteArgs
		handler := middleware.SetVarsAndMux(testHandler.SearchMeetings, args, args2)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := meeting.NewMockUseCase(ctrl)
		testHandler.MeetingUC = m

		apitest.New("GetMeetingsList").
			Handler(handler).
			Method("Get").
			URL("/user").
			Expect(t).
			Status(http.StatusBadRequest).
			End()
	})

	t.Run("GetMeetingsList", func(t *testing.T) {
		var args []middleware.QueryArgs
		args = append(args, middleware.QueryArgs{Key: "limit", Value: "-5"})

		var args2 []middleware.RouteArgs
		handler := middleware.SetVarsAndMux(testHandler.SearchMeetings, args, args2)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := meeting.NewMockUseCase(ctrl)
		testHandler.MeetingUC = m

		apitest.New("GetMeetingsList").
			Handler(handler).
			Method("Get").
			URL("/user").
			Expect(t).
			Status(http.StatusBadRequest).
			End()
	})

	t.Run("GetSubsMeetingsList", func(t *testing.T) {
		var args []middleware.QueryArgs
		args = append(args, middleware.QueryArgs{Key: "start", Value: "2006-01-02"})
		args = append(args, middleware.QueryArgs{Key: "end", Value: "2007-01-02"})
		args = append(args, middleware.QueryArgs{Key: "prevId", Value: "3"})
		args = append(args, middleware.QueryArgs{Key: "limit", Value: "10"})

		var args2 []middleware.RouteArgs
		args2 = append(args2, middleware.RouteArgs{Key: middleware.UserID, Value: 4})
		handler := middleware.SetVarsAndMux(testHandler.GetSubsMeetingsList, args, args2)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := meeting.NewMockUseCase(ctrl)
		testHandler.MeetingUC = m

		layout := "2006-01-02"
		time1, _ := time.Parse(layout, "2006-01-02")
		time2, _ := time.Parse(layout, "2007-01-02")

		m.EXPECT().FilterSubsRegistered(meeting.FilterParams{
			StartDate:  time1,
			EndDate:    time2,
			PrevId:     3,
			CountLimit: 10,
			UserId:     4,
			PrevLikes:  MaxLikes,
			PrevStart:  time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
		}).Return([]models.Meeting{}, nil)

		apitest.New("GetMeetingsList").
			Handler(handler).
			Method("Get").
			URL("/user").
			Expect(t).
			Status(http.StatusOK).
			End()
	})

	t.Run("GetSubsMeetingsListErr1", func(t *testing.T) {
		var args []middleware.QueryArgs
		args = append(args, middleware.QueryArgs{Key: "start", Value: "2006-01-02"})
		args = append(args, middleware.QueryArgs{Key: "end", Value: "2007-01-02"})
		args = append(args, middleware.QueryArgs{Key: "prevId", Value: "3"})
		args = append(args, middleware.QueryArgs{Key: "limit", Value: "10"})

		var args2 []middleware.RouteArgs
		args2 = append(args2, middleware.RouteArgs{Key: middleware.UserID, Value: 4})
		handler := middleware.SetVarsAndMux(testHandler.GetSubsMeetingsList, args, args2)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := meeting.NewMockUseCase(ctrl)
		testHandler.MeetingUC = m

		layout := "2006-01-02"
		time1, _ := time.Parse(layout, "2006-01-02")
		time2, _ := time.Parse(layout, "2007-01-02")

		m.EXPECT().FilterSubsRegistered(meeting.FilterParams{
			StartDate:  time1,
			EndDate:    time2,
			PrevId:     3,
			CountLimit: 10,
			UserId:     4,
			PrevLikes:  MaxLikes,
			PrevStart:  time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
		}).Return([]models.Meeting{}, errors.New("Err"))

		apitest.New("GetMeetingsList").
			Handler(handler).
			Method("Get").
			URL("/user").
			Expect(t).
			Status(http.StatusInternalServerError).
			End()
	})

	t.Run("GetSubsMeetingsListErr1", func(t *testing.T) {
		var args []middleware.QueryArgs
		var args2 []middleware.RouteArgs
		args2 = append(args2, middleware.RouteArgs{Key: middleware.UserID, Value: -1})
		handler := middleware.SetVarsAndMux(testHandler.GetSubsMeetingsList, args, args2)

		apitest.New("GetMeetingsList").
			Handler(handler).
			Method("Get").
			URL("/user").
			Expect(t).
			Status(http.StatusUnauthorized).
			End()
	})
	//
	t.Run("GetSubsMeetingsList", func(t *testing.T) {
		var args []middleware.QueryArgs
		args = append(args, middleware.QueryArgs{Key: "start", Value: "2006-01-02"})
		args = append(args, middleware.QueryArgs{Key: "end", Value: "2007-01-02"})
		args = append(args, middleware.QueryArgs{Key: "prevId", Value: "3"})
		args = append(args, middleware.QueryArgs{Key: "limit", Value: "10"})

		var args2 []middleware.RouteArgs
		args2 = append(args2, middleware.RouteArgs{Key: middleware.UserID, Value: 4})
		handler := middleware.SetVarsAndMux(testHandler.GetSubsFavMeetingsList, args, args2)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := meeting.NewMockUseCase(ctrl)
		testHandler.MeetingUC = m

		layout := "2006-01-02"
		time1, _ := time.Parse(layout, "2006-01-02")
		time2, _ := time.Parse(layout, "2007-01-02")

		m.EXPECT().FilterSubsLiked(meeting.FilterParams{
			StartDate:  time1,
			EndDate:    time2,
			PrevId:     3,
			CountLimit: 10,
			UserId:     4,
			PrevLikes:  MaxLikes,
			PrevStart:  time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
		}).Return([]models.Meeting{}, nil)

		apitest.New("GetMeetingsList").
			Handler(handler).
			Method("Get").
			URL("/user").
			Expect(t).
			Status(http.StatusOK).
			End()
	})

	t.Run("GetSubsMeetingsListErr1", func(t *testing.T) {
		var args []middleware.QueryArgs
		args = append(args, middleware.QueryArgs{Key: "start", Value: "2006-01-02"})
		args = append(args, middleware.QueryArgs{Key: "end", Value: "2007-01-02"})
		args = append(args, middleware.QueryArgs{Key: "prevId", Value: "3"})
		args = append(args, middleware.QueryArgs{Key: "limit", Value: "10"})

		var args2 []middleware.RouteArgs
		args2 = append(args2, middleware.RouteArgs{Key: middleware.UserID, Value: 4})
		handler := middleware.SetVarsAndMux(testHandler.GetSubsFavMeetingsList, args, args2)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := meeting.NewMockUseCase(ctrl)
		testHandler.MeetingUC = m

		layout := "2006-01-02"
		time1, _ := time.Parse(layout, "2006-01-02")
		time2, _ := time.Parse(layout, "2007-01-02")

		m.EXPECT().FilterSubsLiked(meeting.FilterParams{
			StartDate:  time1,
			EndDate:    time2,
			PrevId:     3,
			CountLimit: 10,
			UserId:     4,
			PrevLikes:  MaxLikes,
			PrevStart:  time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
		}).Return([]models.Meeting{}, errors.New("Err"))

		apitest.New("GetMeetingsList").
			Handler(handler).
			Method("Get").
			URL("/user").
			Expect(t).
			Status(http.StatusInternalServerError).
			End()
	})

	t.Run("GetSubsMeetingsListErr1", func(t *testing.T) {
		var args []middleware.QueryArgs
		var args2 []middleware.RouteArgs
		args2 = append(args2, middleware.RouteArgs{Key: middleware.UserID, Value: -1})
		handler := middleware.SetVarsAndMux(testHandler.GetSubsFavMeetingsList, args, args2)

		apitest.New("GetMeetingsList").
			Handler(handler).
			Method("Get").
			URL("/user").
			Expect(t).
			Status(http.StatusUnauthorized).
			End()
	})

}

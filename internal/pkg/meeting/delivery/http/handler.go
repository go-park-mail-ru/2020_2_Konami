package http

import (
	"bytes"
	"errors"
	"konami_backend/internal/pkg/meeting"
	"konami_backend/internal/pkg/middleware"
	"konami_backend/internal/pkg/models"
	"konami_backend/internal/pkg/session"
	hu "konami_backend/internal/pkg/utils/http_utils"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type MeetingHandler struct {
	MeetingUC  meeting.UseCase
	SessionUC  session.UseCase
	MaxReqSize int64
}

const DefCountLimit = 10

func GetQueryParams(r *http.Request) meeting.FilterParams {
	var res meeting.FilterParams
	var err error
	layout := "2006-01-02"
	res.StartDate, err = time.Parse(layout, r.URL.Query().Get("start"))
	if err != nil {
		res.StartDate = time.Now()
	}
	res.EndDate, err = time.Parse(layout, r.URL.Query().Get("end"))
	if err != nil {
		res.EndDate = res.StartDate.AddDate(100, 0, 0)
	}
	res.PrevId, err = strconv.Atoi(r.URL.Query().Get("prevId"))
	if err != nil {
		res.PrevId = 0
	}
	res.CountLimit, err = strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil || res.CountLimit <= 0 {
		res.CountLimit = DefCountLimit
	}
	var ok bool
	res.UserId, ok = r.Context().Value(middleware.UserID).(int)
	if !ok {
		res.UserId = -1
	}
	return res
}

func (h *MeetingHandler) GetMeetingsList(w http.ResponseWriter, r *http.Request) {
	params := GetQueryParams(r)
	var meets []models.Meeting
	var err error
	meets, err = h.MeetingUC.GetNextMeetings(params)
	if err != nil {
		hu.WriteError(w, &hu.ErrResponse{RespCode: http.StatusInternalServerError})
		return
	}
	hu.WriteJson(w, meets)
}

func (h *MeetingHandler) GetUserMeetingsList(w http.ResponseWriter, r *http.Request) {
	params := GetQueryParams(r)
	if params.UserId == -1 {
		hu.WriteError(w, &hu.ErrResponse{RespCode: http.StatusUnauthorized})
		return
	}
	var meets []models.Meeting
	var err error
	meets, err = h.MeetingUC.FilterRegistered(params)
	if err != nil {
		hu.WriteError(w, &hu.ErrResponse{RespCode: http.StatusInternalServerError})
		return
	}
	hu.WriteJson(w, meets)
}

func (h *MeetingHandler) GetFavMeetingsList(w http.ResponseWriter, r *http.Request) {
	params := GetQueryParams(r)
	if params.UserId == -1 {
		hu.WriteError(w, &hu.ErrResponse{RespCode: http.StatusUnauthorized})
		return
	}
	var meets []models.Meeting
	var err error
	meets, err = h.MeetingUC.FilterLiked(params)
	if err != nil {
		hu.WriteError(w, &hu.ErrResponse{RespCode: http.StatusInternalServerError})
		return
	}
	hu.WriteJson(w, meets)
}

func (h *MeetingHandler) GetTopMeetingsList(w http.ResponseWriter, r *http.Request) {
	params := GetQueryParams(r)
	var meets []models.Meeting
	var err error
	meets, err = h.MeetingUC.GetTopMeetings(params)
	if err != nil {
		hu.WriteError(w, &hu.ErrResponse{RespCode: http.StatusInternalServerError})
		return
	}
	hu.WriteJson(w, meets)
}

func (h *MeetingHandler) GetRecommendedList(w http.ResponseWriter, r *http.Request) {
	params := GetQueryParams(r)
	if params.UserId == -1 {
		hu.WriteError(w, &hu.ErrResponse{RespCode: http.StatusUnauthorized})
		return
	}
	var meets []models.Meeting
	var err error
	meets, err = h.MeetingUC.FilterRecommended(params)
	if err != nil {
		hu.WriteError(w, &hu.ErrResponse{RespCode: http.StatusInternalServerError})
		return
	}
	hu.WriteJson(w, meets)
}

func (h *MeetingHandler) GetTaggedMeetings(w http.ResponseWriter, r *http.Request) {
	params := GetQueryParams(r)
	tagId, err := strconv.Atoi(r.URL.Query().Get("tagId"))
	if err != nil || tagId < 0 {
		hu.WriteError(w, &hu.ErrResponse{RespCode: http.StatusBadRequest})
		return
	}
	var meets []models.Meeting
	meets, err = h.MeetingUC.FilterTagged(params, tagId)
	if err != nil {
		hu.WriteError(w, &hu.ErrResponse{RespCode: http.StatusInternalServerError})
		return
	}
	hu.WriteJson(w, meets)
}

func (h *MeetingHandler) GetAkinMeetings(w http.ResponseWriter, r *http.Request) {
	params := GetQueryParams(r)
	meetId, err := strconv.Atoi(r.URL.Query().Get("meetId"))
	if err != nil || meetId < 0 {
		hu.WriteError(w, &hu.ErrResponse{RespCode: http.StatusBadRequest})
		return
	}
	var meets []models.Meeting
	meets, err = h.MeetingUC.FilterSimilar(params, meetId)
	if err != nil {
		hu.WriteError(w, &hu.ErrResponse{RespCode: http.StatusInternalServerError})
		return
	}
	hu.WriteJson(w, meets)
}

func (h *MeetingHandler) CreateMeeting(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value(middleware.UserID).(int)
	if !ok {
		hu.WriteError(w, &hu.ErrResponse{RespCode: http.StatusUnauthorized})
		return
	}
	tokenValid, ok := r.Context().Value(middleware.CSRFValid).(bool)
	if !ok || !tokenValid {
		hu.WriteError(w, &hu.ErrResponse{RespCode: http.StatusUnauthorized, ErrMsg: "Invalid CSRF token"})
		return
	}
	mData := &models.MeetingData{}
	buf := new(bytes.Buffer)
	_, _ = buf.ReadFrom(http.MaxBytesReader(w, r.Body, h.MaxReqSize))
	err := mData.UnmarshalJSON(buf.Bytes())
	if err != nil {
		hu.WriteError(w, &hu.ErrResponse{RespCode: http.StatusBadRequest})
		return
	}
	_, err = h.MeetingUC.CreateMeeting(userId, *mData)
	if err != nil {
		hu.WriteError(w, &hu.ErrResponse{RespCode: http.StatusInternalServerError})
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *MeetingHandler) GetMeeting(w http.ResponseWriter, r *http.Request) {
	meetId, err := strconv.Atoi(r.URL.Query().Get("meetId"))
	if err != nil {
		hu.WriteError(w, &hu.ErrResponse{RespCode: http.StatusNotFound})
		return
	}
	userId, ok := r.Context().Value(middleware.UserID).(int)
	var meet models.MeetingDetails
	if !ok {
		meet, err = h.MeetingUC.GetMeeting(meetId, -1, false)
	} else {
		meet, err = h.MeetingUC.GetMeeting(meetId, userId, true)
	}
	if err != nil {
		hu.WriteError(w, &hu.ErrResponse{RespCode: http.StatusBadRequest})
		return
	}
	hu.WriteJson(w, meet)
}

func (h *MeetingHandler) UpdateMeeting(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value(middleware.UserID).(int)
	if !ok {
		hu.WriteError(w, &hu.ErrResponse{RespCode: http.StatusUnauthorized})
		return
	}
	tokenValid, ok := r.Context().Value(middleware.CSRFValid).(bool)
	if !ok || !tokenValid {
		hu.WriteError(w, &hu.ErrResponse{RespCode: http.StatusUnauthorized, ErrMsg: "Invalid CSRF token"})
		return
	}
	update := &models.MeetingUpdate{}
	buf := new(bytes.Buffer)
	_, _ = buf.ReadFrom(r.Body)
	err := update.UnmarshalJSON(buf.Bytes())

	if err != nil {
		hu.WriteError(w, &hu.ErrResponse{RespCode: http.StatusBadRequest})
		return
	}
	err = h.MeetingUC.UpdateMeeting(userId, *update)
	if errors.Is(err, meeting.ErrMeetingNotFound) {
		hu.WriteError(w, &hu.ErrResponse{RespCode: http.StatusNotFound})
	}
	if err != nil {
		hu.WriteError(w, &hu.ErrResponse{RespCode: http.StatusBadRequest})
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *MeetingHandler) SearchMeetings(w http.ResponseWriter, r *http.Request) {
	params := GetQueryParams(r)
	searchQuery := strings.TrimSpace(r.URL.Query().Get("query"))
	if searchQuery == "" {
		hu.WriteError(w, &hu.ErrResponse{RespCode: http.StatusBadRequest})
		return
	}
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil || limit <= 0 {
		limit = -1
	}
	var meets []models.Meeting
	meets, err = h.MeetingUC.SearchMeetings(params, searchQuery, limit)
	if err != nil {
		hu.WriteError(w, &hu.ErrResponse{RespCode: http.StatusInternalServerError})
		return
	}
	hu.WriteJson(w, meets)
}

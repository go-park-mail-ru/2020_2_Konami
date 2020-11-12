package http

import (
	"encoding/json"
	"errors"
	"konami_backend/internal/pkg/meeting"
	"konami_backend/internal/pkg/middleware"
	"konami_backend/internal/pkg/models"
	"konami_backend/internal/pkg/session"
	hu "konami_backend/internal/pkg/utils/http_utils"
	"net/http"
	"strconv"
)

type MeetingHandler struct {
	MeetingUC  meeting.UseCase
	SessionUC  session.UseCase
	MaxReqSize int64
}

func (h *MeetingHandler) GetMeetingsList(w http.ResponseWriter, r *http.Request) {
	todayOnly := r.URL.Query().Get("today") == "true"
	tomorrowOnly := r.URL.Query().Get("tomorrow") == "true"
	myOnly := r.URL.Query().Get("mymeetings") == "true"
	favOnly := r.URL.Query().Get("favorites") == "true"

	var meets []models.Meeting
	var ok bool
	var err error
	userId, ok := r.Context().Value(middleware.UserID).(int)
	if !ok {
		if myOnly || favOnly {
			hu.WriteError(w, &hu.ErrResponse{RespCode: http.StatusUnauthorized})
			return
		}
		userId = -1
	}
	if todayOnly {
		meets, err = h.MeetingUC.FilterToday(userId)
	} else if tomorrowOnly {
		meets, err = h.MeetingUC.FilterTomorrow(userId)
	} else if myOnly {
		meets, err = h.MeetingUC.FilterRegistered(userId)
	} else if favOnly {
		meets, err = h.MeetingUC.FilterLiked(userId)
	} else {
		meets, err = h.MeetingUC.GetAll(userId)
	}

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
	err := json.NewDecoder(http.MaxBytesReader(w, r.Body, h.MaxReqSize)).Decode(&mData)
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
	var meet models.Meeting
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
	err := json.NewDecoder(r.Body).Decode(&update)
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

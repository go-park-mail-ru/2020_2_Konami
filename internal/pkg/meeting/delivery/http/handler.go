package http

import (
	"encoding/json"
	"konami_backend/internal/pkg/meeting"
	"konami_backend/internal/pkg/models"
	"konami_backend/internal/pkg/session"
	hu "konami_backend/internal/pkg/utils/http_utils"
	"net/http"
	"strconv"
)

type MeetingHandler struct {
	MeetingUC meeting.UseCase
	SessionUC session.UseCase
}

func (h *MeetingHandler) GetMeetingsList(w http.ResponseWriter, r *http.Request) {
	todayOnly := r.URL.Query().Get("today") == "true"
	tomorrowOnly := r.URL.Query().Get("tomorrow") == "true"
	myOnly := r.URL.Query().Get("mymeetings") == "true"
	favOnly := r.URL.Query().Get("favorites") == "true"

	var meets []models.MeetingCard
	var err error
	var userId int
	if myOnly || favOnly {
		token, err := r.Cookie("authToken")
		if err != nil {
			hu.WriteError(w, &hu.ErrResponse{RespCode: http.StatusUnauthorized})
			return
		}
		userId, err = h.SessionUC.GetUserId(token.Value)
		if err != nil {
			hu.WriteError(w, &hu.ErrResponse{RespCode: http.StatusUnauthorized})
			return
		}
	}

	if todayOnly {
		meets, err = h.MeetingUC.FilterToday()
	} else if tomorrowOnly {
		meets, err = h.MeetingUC.FilterTomorrow()
	} else if myOnly {
		meets, err = h.MeetingUC.FilterRegistered(userId)
	} else if favOnly {
		meets, err = h.MeetingUC.FilterLiked(userId)
	} else {
		meets, err = h.MeetingUC.GetAll()
	}

	if err != nil {
		hu.WriteError(w, &hu.ErrResponse{RespCode: http.StatusInternalServerError})
		return
	}
	hu.WriteJson(w, meets)
}

func (h *MeetingHandler) CreateMeeting(w http.ResponseWriter, r *http.Request) {
	token, err := r.Cookie("authToken")
	if err != nil {
		hu.WriteError(w, &hu.ErrResponse{RespCode: http.StatusUnauthorized})
		return
	}
	userId, err := h.SessionUC.GetUserId(token.Value)
	if err != nil {
		hu.WriteError(w, &hu.ErrResponse{RespCode: http.StatusUnauthorized})
		return
	}
	mData := &models.MeetingData{}
	err = json.NewDecoder(http.MaxBytesReader(w, r.Body, 10*1024*1024)).Decode(&mData)
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
	token, err := r.Cookie("authToken")
	if err != nil {
		hu.WriteError(w, &hu.ErrResponse{RespCode: http.StatusUnauthorized})
		return
	}
	userId, err := h.SessionUC.GetUserId(token.Value)
	if err != nil {
		hu.WriteError(w, &hu.ErrResponse{RespCode: http.StatusUnauthorized})
		return
	}
	meet, err := h.MeetingUC.GetMeeting(meetId, userId)
	if err != nil {
		hu.WriteError(w, &hu.ErrResponse{RespCode: http.StatusInternalServerError})
		return
	}
	hu.WriteJson(w, meet)
}

func (h *MeetingHandler) UpdateMeeting(w http.ResponseWriter, r *http.Request) {
	panic("implement me")
}

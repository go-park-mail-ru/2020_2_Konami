package http

import (
	"bytes"
	"context"
	"errors"
	sessionPkg "konami_backend/auth/pkg/session"
	"konami_backend/internal/pkg/middleware"
	"konami_backend/internal/pkg/models"
	"konami_backend/internal/pkg/profile"
	hu "konami_backend/internal/pkg/utils/http_utils"
	"konami_backend/proto/auth"
	"net/http"
	"strconv"
)

type ProfileHandler struct {
	ProfileUC  profile.UseCase
	AuthClient auth.AuthCheckerClient
	MaxReqSize int64
}

func GetQueryParams(r *http.Request) profile.FilterParams {
	var res profile.FilterParams
	var err error
	res.PrevId, err = strconv.Atoi(r.URL.Query().Get("prevId"))
	if err != nil {
		res.PrevId = 0
	}
	res.CountLimit, err = strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil || res.CountLimit <= 0 {
		res.CountLimit = 0
	}
	var ok bool
	res.ReqAuthorId, ok = r.Context().Value(middleware.UserID).(int)
	if !ok {
		res.ReqAuthorId = -1
	}
	return res
}

func (h *ProfileHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	var creds models.Credentials
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(r.Body)
	if err == nil {
		err = creds.UnmarshalJSON(buf.Bytes())
	}
	if err != nil {
		hu.WriteError(w, &hu.ErrResponse{RespCode: http.StatusBadRequest})
		return
	}
	_, err = h.ProfileUC.Validate(creds)
	if err != profile.ErrUserNonExistent {
		hu.WriteError(w, &hu.ErrResponse{RespCode: http.StatusConflict, ErrMsg: "login has already been taken"})
		return
	}
	userId, err := h.ProfileUC.SignUp(creds)
	if err != nil {
		hu.WriteError(w, &hu.ErrResponse{RespCode: http.StatusBadRequest})
		return
	}
	session, err := h.AuthClient.Create(context.Background(), &auth.Session{UserId: int64(userId)})
	if err != nil {
		hu.WriteError(w, &hu.ErrResponse{RespCode: http.StatusInternalServerError})
		return
	}
	hu.SetAuthCookie(w, session.Token)
	w.WriteHeader(http.StatusCreated)
}

func (h *ProfileHandler) GetUserId(w http.ResponseWriter, r *http.Request) {
	uId, ok := r.Context().Value(middleware.UserID).(int)
	if !ok {
		hu.WriteError(w, &hu.ErrResponse{RespCode: http.StatusUnauthorized})
		return
	}
	hu.WriteJson(w, struct {
		UserId int `json:"userId"`
	}{uId})
}

func (h *ProfileHandler) LogIn(w http.ResponseWriter, r *http.Request) {
	var creds models.Credentials
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(r.Body)
	if err == nil {
		err = creds.UnmarshalJSON(buf.Bytes())
	}
	if err != nil {
		hu.WriteError(w, &hu.ErrResponse{RespCode: http.StatusBadRequest})
		return
	}
	userId, err := h.ProfileUC.Validate(creds)
	if errors.Is(err, profile.ErrInvalidCredentials) || errors.Is(err, profile.ErrUserNonExistent) {
		hu.WriteError(w, &hu.ErrResponse{RespCode: http.StatusBadRequest, ErrMsg: "invalid credentials"})
		return
	}
	if err != nil {
		hu.WriteError(w, &hu.ErrResponse{RespCode: http.StatusInternalServerError})
		return
	}
	session, err := h.AuthClient.Create(context.Background(), &auth.Session{UserId: int64(userId)})
	if err != nil {
		hu.WriteError(w, &hu.ErrResponse{RespCode: http.StatusInternalServerError})
		return
	}
	hu.SetAuthCookie(w, session.Token)
	w.WriteHeader(http.StatusCreated)
}

func (h *ProfileHandler) LogOut(w http.ResponseWriter, r *http.Request) {
	token, ok := r.Context().Value(middleware.AuthToken).(string)
	if !ok {
		hu.WriteError(w, &hu.ErrResponse{RespCode: http.StatusUnauthorized})
		return
	}
	_, err := h.AuthClient.Delete(context.Background(), &auth.SessionToken{Token: token})
	if err != nil && err != sessionPkg.ErrInvalidToken {
		hu.WriteError(w, &hu.ErrResponse{RespCode: http.StatusInternalServerError})
		return
	}
	hu.RemoveAuthCookie(w, token)
	w.WriteHeader(http.StatusOK)
}

func (h *ProfileHandler) UploadUserPic(w http.ResponseWriter, r *http.Request) {
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
	err := r.ParseMultipartForm(h.MaxReqSize)
	if err != nil {
		hu.WriteError(w, &hu.ErrResponse{RespCode: http.StatusBadRequest, ErrMsg: "invalid multipart form"})
		return
	}
	file, header, err := r.FormFile("fileToUpload")
	if err != nil {
		hu.WriteError(w, &hu.ErrResponse{RespCode: http.StatusBadRequest, ErrMsg: "invalid form file"})
		return
	}
	defer file.Close()
	err = h.ProfileUC.UploadProfilePic(userId, header.Filename, file)
	if err != nil {
		hu.WriteError(w, &hu.ErrResponse{RespCode: http.StatusBadRequest, ErrMsg: "invalid image file"})
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *ProfileHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	targetId, err := strconv.Atoi(r.URL.Query().Get("userId"))
	if err != nil {
		hu.WriteError(w, &hu.ErrResponse{RespCode: http.StatusNotFound})
		return
	}
	userId, ok := r.Context().Value(middleware.UserID).(int)
	if !ok {
		userId = -1
	}
	prof, err := h.ProfileUC.GetProfile(userId, targetId)
	if err != nil {
		hu.WriteError(w, &hu.ErrResponse{RespCode: http.StatusNotFound})
		return
	}
	hu.WriteJson(w, prof)
}

func (h *ProfileHandler) GetPeople(w http.ResponseWriter, r *http.Request) {
	params := GetQueryParams(r)
	users, err := h.ProfileUC.GetAll(params)
	if err != nil {
		hu.WriteError(w, &hu.ErrResponse{RespCode: http.StatusInternalServerError})
		return
	}
	hu.WriteJson(w, users)
}

func (h *ProfileHandler) GetUserSubscriptions(w http.ResponseWriter, r *http.Request) {
	params := GetQueryParams(r)
	users, err := h.ProfileUC.GetUserSubscriptions(params)
	if err != nil {
		hu.WriteError(w, &hu.ErrResponse{RespCode: http.StatusNotFound})
		return
	}
	hu.WriteJson(w, users)
}

func (h *ProfileHandler) CreateUserSubscription(w http.ResponseWriter, r *http.Request) {
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
	sub := &models.UserSubscription{}
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(r.Body)
	if err == nil {
		err = sub.UnmarshalJSON(buf.Bytes())
	}
	if err != nil {
		hu.WriteError(w, &hu.ErrResponse{RespCode: http.StatusBadRequest})
		return
	}
	_, err = h.ProfileUC.CreateSubscription(userId, sub.TargetId)
	if err == profile.ErrUserNonExistent {
		hu.WriteError(w, &hu.ErrResponse{RespCode: http.StatusBadRequest})
		return
	}
	if err != nil {
		hu.WriteError(w, &hu.ErrResponse{RespCode: http.StatusInternalServerError})
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *ProfileHandler) RemoveUserSubscription(w http.ResponseWriter, r *http.Request) {
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
	sub := &models.UserSubscription{}
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(r.Body)
	if err == nil {
		err = sub.UnmarshalJSON(buf.Bytes())
	}
	if err != nil {
		hu.WriteError(w, &hu.ErrResponse{RespCode: http.StatusBadRequest})
		return
	}
	err = h.ProfileUC.RemoveSubscription(userId, sub.TargetId)
	if err != nil {
		hu.WriteError(w, &hu.ErrResponse{RespCode: http.StatusInternalServerError})
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *ProfileHandler) EditUser(w http.ResponseWriter, r *http.Request) {
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
	update := &models.ProfileUpdate{}
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(r.Body)
	if err == nil {
		err = update.UnmarshalJSON(buf.Bytes())
	}
	if err != nil {
		hu.WriteError(w, &hu.ErrResponse{RespCode: http.StatusBadRequest})
		return
	}
	err = h.ProfileUC.EditProfile(userId, *update)
	if err != nil {
		hu.WriteError(w, &hu.ErrResponse{RespCode: http.StatusInternalServerError})
		return
	}
	w.WriteHeader(http.StatusOK)
}

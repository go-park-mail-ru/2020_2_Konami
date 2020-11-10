package http

import (
	"encoding/json"
	"konami_backend/internal/pkg/models"
	"konami_backend/internal/pkg/profile"
	"konami_backend/internal/pkg/session"
	hu "konami_backend/internal/pkg/utils/http_utils"
	"net/http"
	"strconv"
)

type ProfileHandler struct {
	ProfileUC profile.UseCase
	SessionUC session.UseCase
}

func (h *ProfileHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	var creds models.Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
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
	token, err := h.SessionUC.CreateSession(userId)
	if err != nil {
		hu.WriteError(w, &hu.ErrResponse{RespCode: http.StatusInternalServerError})
		return
	}
	hu.SetMonthCookie(w, "authToken", token)
	w.WriteHeader(http.StatusOK)
}

func (h *ProfileHandler) UploadUserPic(w http.ResponseWriter, r *http.Request) {
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
	err = r.ParseMultipartForm(10 * 1024 * 1024)
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
	userId, err := strconv.Atoi(r.URL.Query().Get("userId"))
	if err != nil {
		hu.WriteError(w, &hu.ErrResponse{RespCode: http.StatusNotFound})
		return
	}
	prof, err := h.ProfileUC.GetProfile(userId)
	if err != nil {
		hu.WriteError(w, &hu.ErrResponse{RespCode: http.StatusNotFound})
		return
	}
	hu.WriteJson(w, prof)
}

func (h *ProfileHandler) GetPeople(w http.ResponseWriter, r *http.Request) {
	users, err := h.ProfileUC.GetAll()
	if err != nil {
		hu.WriteError(w, &hu.ErrResponse{RespCode: http.StatusInternalServerError})
		return
	}
	hu.WriteJson(w, users)
}

func (h *ProfileHandler) EditUser(w http.ResponseWriter, r *http.Request) {
	panic("implement me")
}

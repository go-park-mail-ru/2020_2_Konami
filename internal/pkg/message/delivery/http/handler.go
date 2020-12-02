package http

import (
	"bytes"
	"github.com/gorilla/websocket"
	"konami_backend/internal/pkg/message"
	"konami_backend/internal/pkg/middleware"
	"konami_backend/internal/pkg/models"
	hu "konami_backend/internal/pkg/utils/http_utils"
	"konami_backend/logger"
	"net/http"
	"strconv"
)

type MessageHandler struct {
	MessageUC  message.UseCase
	Log        *logger.Logger
	MaxReqSize int64
	clients    map[*websocket.Conn]bool
	msgChan    chan *models.Message
	upgrader   websocket.Upgrader
}

func NewMessageHandler(messageUC message.UseCase, log *logger.Logger, maxReqSize int64) MessageHandler {
	return MessageHandler{
		MessageUC:  messageUC,
		Log:        log,
		MaxReqSize: maxReqSize,
		clients:    make(map[*websocket.Conn]bool),
		msgChan:    make(chan *models.Message),
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
}

func (h *MessageHandler) SendMessage(w http.ResponseWriter, r *http.Request) {
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
	msg := &models.Message{}
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(http.MaxBytesReader(w, r.Body, h.MaxReqSize))
	if err == nil {
		err = msg.UnmarshalJSON(buf.Bytes())
	}
	if err != nil {
		hu.WriteError(w, &hu.ErrResponse{RespCode: http.StatusBadRequest})
		return
	}
	msg.AuthorId = userId
	_, err = h.MessageUC.CreateMessage(*msg)
	if err != nil {
		hu.WriteError(w, &hu.ErrResponse{RespCode: http.StatusInternalServerError})
		return
	}
	w.WriteHeader(http.StatusCreated)
	go h.PublishMsg(msg)
}

func (h *MessageHandler) GetMessages(w http.ResponseWriter, r *http.Request) {
	mId, err := strconv.Atoi(r.URL.Query().Get("meetId"))
	if err != nil {
		hu.WriteError(w, &hu.ErrResponse{RespCode: http.StatusBadRequest})
		return
	}
	messages, err := h.MessageUC.GetMessages(mId)
	if err != nil {
		hu.WriteError(w, &hu.ErrResponse{RespCode: http.StatusInternalServerError})
		return
	}
	hu.WriteJson(w, messages)
}

func (h *MessageHandler) Upgrade(w http.ResponseWriter, r *http.Request) {
	ws, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		h.Log.LogError("message/delivery/http", "Upgrade", err)
	} else {
		h.clients[ws] = true
	}
}

func (h *MessageHandler) PublishMsg(msg *models.Message) {
	h.msgChan <- msg
}

func (h *MessageHandler) ServeWS() {
	for {
		msg := <-h.msgChan
		resp := struct {
			Payload *models.Message `json:"payload"`
			MsgType string          `json:"type"`
		}{Payload: msg, MsgType: "chatMessage"}
		for client := range h.clients {
			err := client.WriteJSON(resp)
			if err != nil {
				err = client.Close()
				delete(h.clients, client)
				if err != nil {
					h.Log.LogError("message/delivery/http", "ServeWS", err)
				}
			}
		}
	}
}

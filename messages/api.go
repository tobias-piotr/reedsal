package messages

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"reedsal/api"
	"reedsal/users"

	"github.com/go-chi/chi/v5"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/jmoiron/sqlx"
)

var upgrader = websocket.Upgrader{}

type MessageHandler struct {
	Srv   *MessageService
	Redis *redis.Client
}

func NewMessageHandler(db *sqlx.DB, redisClient *redis.Client) *MessageHandler {
	return &MessageHandler{NewMessageService(NewMessageRepository(db), users.NewUserRepository(db)), redisClient}
}

func (h MessageHandler) Router() chi.Router {
	r := chi.NewRouter()

	r.Use(api.JWTMiddleware)
	r.Get("/{id}", h.HandleConversation)

	return r
}

func (h MessageHandler) HandleConversation(w http.ResponseWriter, r *http.Request) {
	// Check recipient
	recipientID := chi.URLParam(r, "id")
	if _, err := uuid.Parse(recipientID); err != nil {
		api.RespondWithError(w, api.NewAPIError(http.StatusBadRequest, "Invalid user id", nil))
		return
	}
	exists, err := h.Srv.GetRecipientExistence(recipientID)
	if err != nil {
		api.RespondWithError(w, err)
		return
	}
	if !exists {
		api.RespondWithError(w, api.NewAPIError(http.StatusNotFound, "User does not exist", nil))
		return
	}

	// Get messages between user and recipient
	senderID := r.Context().Value(api.SubKey).(string)
	msgs, err := h.Srv.GetConversation(senderID, recipientID)
	if err != nil {
		api.RespondWithError(w, err)
		return
	}

	// Upgrade to WS
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.Error("Upgrade failed", "err", err)
		return
	}

	// Make sure the connection and Redis sub are both closed
	closeWS := make(chan bool, 1)
	defer func() {
		closeWS <- true
		conn.Close()
	}()

	// Send conversation
	api.SendDataMessage(conn, msgs)

	ctx := context.Background()

	go func() {
		sub := h.Redis.Subscribe(ctx, fmt.Sprintf("%v-%v", recipientID, senderID))
		// Wait for the signal and close the sub (channel loop will be stopped)
		go func() {
			<-closeWS
			sub.Close()
		}()
		for msg := range sub.Channel() {
			dm := api.DataMessage{}
			json.Unmarshal([]byte(msg.Payload), &dm.Data)
			conn.WriteJSON(dm)
		}
	}()

	for {
		// Read message
		_, message, err := conn.ReadMessage()
		if err != nil {
			slog.Error("Reading message failed", "err", err)
			return
		}

		// Validate data
		data := MessageCreatePayload{senderID, recipientID, string(message)}
		verr := data.Validate()
		if verr != nil {
			api.SendErrorMessage(conn, verr)
			continue
		}

		// Create message
		msg, err := h.Srv.CreateMessage(&data)
		if err != nil {
			api.SendErrorMessage(conn, err)
			continue
		}

		// Publish message via Redis pubsub
		go func() {
			msgBytes, _ := json.Marshal(msg)
			h.Redis.Publish(ctx, fmt.Sprintf("%v-%v", senderID, recipientID), msgBytes)
		}()

		// Respond
		if err != nil {
			api.SendErrorMessage(conn, err)
			continue
		}
		api.SendDataMessage(conn, msg)
	}
}

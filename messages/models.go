package messages

import (
	"reedsal/api"
	"time"

	"github.com/google/uuid"
)

type MessageCreatePayload struct {
	Sender    string `json:"sender"`
	Recipient string `json:"recipient"`
	Content   string `json:"content"`
}

func (p MessageCreatePayload) Validate() *api.ValidationError {
	if len(p.Content) < 1 {
		return &api.ValidationError{Details: api.ValidationDetails{"content": "Message content cannot be empty"}}
	}
	return nil
}

type Message struct {
	ID        uuid.UUID `json:"id"`
	Sender    uuid.UUID `json:"sender"`
	Recipient uuid.UUID `json:"recipient"`
	Content   string    `json:"content"`
	Datetime  time.Time `json:"datetime"`
}

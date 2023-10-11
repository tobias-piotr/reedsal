package messages

import "github.com/jmoiron/sqlx"

type MessageRepository struct {
	DB *sqlx.DB
}

func NewMessageRepository(db *sqlx.DB) *MessageRepository {
	return &MessageRepository{db}
}

func (r MessageRepository) CreateMessage(data *MessageCreatePayload) (*Message, error) {
	query := `
		INSERT INTO messages (sender, recipient, content)
		VALUES (:sender, :recipient, :content)
		RETURNING id, sender, recipient, content, datetime
	`
	query, args, err := r.DB.BindNamed(query, data)
	if err != nil {
		return nil, err
	}

	var msg Message
	err = r.DB.QueryRowx(query, args...).StructScan(&msg)
	return &msg, err
}

func (r MessageRepository) GetConversation(senderID string, recipientID string) (*[]Message, error) {
	query := `
		SELECT id, sender, recipient, content, datetime
		FROM messages
		WHERE (sender=$1 AND recipient=$2) OR (sender=$2 AND recipient=$1)
		ORDER BY datetime DESC;
	`
	msgs := []Message{}
	err := r.DB.Select(&msgs, query, senderID, recipientID)
	return &msgs, err
}

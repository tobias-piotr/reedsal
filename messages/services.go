package messages

import (
	"reedsal/users"
)

type MessageService struct {
	Repo     *MessageRepository
	UserRepo *users.UserRepository
}

func NewMessageService(repo *MessageRepository, userRepo *users.UserRepository) *MessageService {
	return &MessageService{repo, userRepo}
}

func (s MessageService) CreateMessage(data *MessageCreatePayload) (*Message, error) {
	return s.Repo.CreateMessage(data)
}

func (s MessageService) GetRecipientExistence(id string) (bool, error) {
	return s.UserRepo.GetUserExistence("id", id)
}

func (s MessageService) GetConversation(senderID string, recipientID string) (*[]Message, error) {
	return s.Repo.GetConversation(senderID, recipientID)
}

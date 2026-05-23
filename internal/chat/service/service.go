package chat_service

import (
	"context"

	"github.com/egorkto/Chat-go/internal/domain"
)

type ChatService struct {
	usersStorage UsersStorage
	msgStorage   MessagesStorage
}

type UsersStorage interface {
	GetUserByID(
		ctx context.Context,
		id int,
	) (domain.User, error)
}

type MessagesStorage interface {
	CreateMessage(
		ctx context.Context,
		msg domain.Message,
	) (domain.Message, error)
	GetMessages(
		ctx context.Context,
		limit *int,
		offset *int,
	) ([]domain.Message, error)
}

func New(usersStorage UsersStorage, msgStorage MessagesStorage) ChatService {
	return ChatService{
		usersStorage: usersStorage,
		msgStorage:   msgStorage,
	}
}

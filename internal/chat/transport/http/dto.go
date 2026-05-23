package chat_transport_http

import (
	"time"

	"github.com/egorkto/Chat-go/internal/domain"
)

type MessageDTO struct {
	ID      int       `json:"id"`
	Version int       `json:"vesion"`
	Sender  UserDTO   `json:"sender"`
	Text    string    `json:"text"`
	SendAt  time.Time `json:"send_at"`
}

type UserDTO struct {
	ID       int    `json:"id"`
	Version  int    `json:"version"`
	Login    string `json:"login"`
	FullName string `json:"full_name"`
}

func dtoFromDomains(msgs []domain.Message) []MessageDTO {
	dto := make([]MessageDTO, len(msgs))

	for i, msg := range msgs {
		dto[i] = dtoFromDomain(msg)
	}

	return dto
}

func dtoFromDomain(msg domain.Message) MessageDTO {
	return MessageDTO{
		ID:      msg.ID,
		Version: msg.Version,
		Sender: UserDTO{
			ID:       msg.Sender.ID(),
			Version:  msg.Sender.Version(),
			Login:    msg.Sender.Login(),
			FullName: msg.Sender.FullName(),
		},
		SendAt: msg.SendTime,
	}
}

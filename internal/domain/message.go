package domain

import "time"

type Message struct {
	ID       int
	Version  int
	Sender   User
	Text     string
	SendTime time.Time
}

func NewMessage(
	ID int,
	version int,
	sender User,
	text string,
	sendTime time.Time,
) Message {
	return Message{
		ID:       ID,
		Version:  version,
		Sender:   sender,
		Text:     text,
		SendTime: sendTime,
	}
}

func NewUninitializedMessage(
	sender User,
	text string,
	sendTime time.Time,
) Message {
	return NewMessage(UninitializedID, UninitializedVersion, sender, text, sendTime)
}

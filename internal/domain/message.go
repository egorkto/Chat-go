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

func (m Message) Validate() error {
	valErrors := make(map[string]string)

	if err := m.Sender.Validate(); err != nil {
		valErrors["sender"] = "invalid sender: " + err.Error()
	}

	textLen := len([]rune(m.Text))
	if textLen < 1 || textLen > 2048 {
		valErrors["text"] = "text must be between 1 and 2048 symbols"
	}

	if len(valErrors) > 0 {
		return ValidationError{Errors: valErrors}
	}

	return nil
}

package entity

import (
	"bytes"
	"time"
)

type message struct {
	SenderID UserID    `json:"sender_id"`
	Text     string    `json:"text"`
	SentAt   time.Time `json:"sent_at"`
}

func NewMesssage(senderID UserID, content []byte) *message {
	return &message{
		SenderID: senderID,
		Text:     bytes.NewBuffer(content).String(),
		SentAt:   time.Now(),
	}
}

package entity

import (
	"bytes"
	"time"
)

type Message struct {
	SenderID   UserID    `json:"sender_id"`
	SenderName string    `json:"sender_name"`
	Text       string    `json:"text"`
	SentAt     time.Time `json:"sent_at"`
}

func NewMesssage(senderID UserID, senderName string, content []byte) *Message {
	return &Message{
		SenderID:   senderID,
		SenderName: senderName,
		Text:       bytes.NewBuffer(content).String(),
		SentAt:     time.Now(),
	}
}

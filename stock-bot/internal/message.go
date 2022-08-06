package internal

import (
	"time"
)

type Message struct {
	// SenderID UserID    `json:"sender_id"`
	Text   string    `json:"text"`
	SentAt time.Time `json:"sent_at"`
}

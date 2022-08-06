package entity

import "github.com/google/uuid"

type UserID = uuid.UUID

type user struct {
	ID       UserID
	Nick     string
	Password string
}

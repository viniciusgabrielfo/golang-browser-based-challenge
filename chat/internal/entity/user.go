package entity

import "github.com/google/uuid"

type UserID = uuid.UUID
type User struct {
	ID       UserID
	Nick     string
	Password string
}

func NewUser(nick, password string) *User {
	return &User{
		ID:       uuid.New(),
		Nick:     nick,
		Password: password,
	}
}

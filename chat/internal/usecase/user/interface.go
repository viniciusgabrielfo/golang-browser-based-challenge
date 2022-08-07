package user

import "github.com/viniciusgabrielfo/golang-browser-based-challenge/chat/internal/entity"

type Repository interface {
	Create(user *entity.User) error
	Get(id entity.UserID) (*entity.User, error)
	GetByNick(nick string) (*entity.User, error)
}

type Service interface {
	CreateUser(nick, password string) error
	GetUser(id entity.UserID) (*entity.User, error)
	Auth(nick, password string) (entity.UserID, error)
}

package repository

import (
	"github.com/viniciusgabrielfo/golang-browser-based-challenge/chat/internal/entity"
)

type userInMemory struct {
	m map[entity.UserID]*entity.User
}

func NewUserInMemory() *userInMemory {
	var m = map[entity.UserID]*entity.User{}
	return &userInMemory{
		m: m,
	}
}

func (r *userInMemory) Create(e *entity.User) error {
	r.m[e.ID] = e
	return nil
}

func (r *userInMemory) Get(id entity.UserID) (*entity.User, error) {
	if r.m[id] == nil {
		return nil, entity.ErrNotFoundEntity
	}
	return r.m[id], nil
}

func (r *userInMemory) GetByNick(nick string) (*entity.User, error) {

	var user *entity.User
	for _, j := range r.m {
		if j.Nick == nick {
			user = j
		}
	}

	if user == nil {
		return nil, entity.ErrNotFoundEntity
	}

	return user, nil
}

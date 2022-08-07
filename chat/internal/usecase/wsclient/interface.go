package wsclient

import "github.com/viniciusgabrielfo/golang-browser-based-challenge/chat/internal/entity"

type Service interface {
	GetClient() *entity.Client
	Read()
	Write()
}

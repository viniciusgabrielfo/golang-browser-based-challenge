package chatroom

import (
	"github.com/viniciusgabrielfo/golang-browser-based-challenge/chat/internal/entity"
	"go.uber.org/zap"
)

type service struct {
	chatroom *entity.Chatroom
	logger   *zap.SugaredLogger
}

func NewService(chatroom *entity.Chatroom, logger *zap.SugaredLogger) *service {
	return &service{
		chatroom: chatroom,
		logger:   logger,
	}
}

func (s *service) Start() {
	for {
		select {
		case Client := <-s.chatroom.GetRegisterChan():
			s.chatroom.Clients[Client] = true
			s.logger.Infow("new client registered", "active_clients", len(s.chatroom.Clients))
		case Client := <-s.chatroom.GetUnregisterChan():
			if _, ok := s.chatroom.Clients[Client]; ok {
				delete(s.chatroom.Clients, Client)
				s.logger.Infow("a client was unregistered", "active_clients", len(s.chatroom.Clients))
				close(Client.GetOutboundChan())
			}
		case message := <-s.chatroom.GetBroadcastChan():
			for Client := range s.chatroom.Clients {
				Client.GetOutboundChan() <- message
			}
		}
	}
}

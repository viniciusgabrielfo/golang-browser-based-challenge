package wsclient

import (
	"github.com/gorilla/websocket"
	"github.com/viniciusgabrielfo/golang-browser-based-challenge/chat/internal/entity"
	"go.uber.org/zap"
)

type service struct {
	client *entity.Client
	logger *zap.SugaredLogger
}

func NewService(client *entity.Client, logger *zap.SugaredLogger) *service {
	return &service{
		client: client,
		logger: logger,
	}
}

func (s *service) GetClient() *entity.Client {
	return s.client
}

func (s *service) Read() {
	defer s.client.Websocket.Close()

	for {
		_, bMsg, err := s.client.Websocket.ReadMessage()
		if err != nil {
			if !websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
				s.logger.Error(err)
			}
			break
		}

		s.client.Chatroom.GetBroadcastChan() <- entity.NewMesssage(s.client.User.ID, s.client.User.Nick, bMsg)
	}
}

func (s *service) Write() {
	defer s.client.Websocket.Close()

	for msg := range s.client.GetOutboundChan() {

		if err := s.client.Websocket.WriteJSON(msg); err != nil {
			s.logger.Error(err)
			break
		}
	}
}

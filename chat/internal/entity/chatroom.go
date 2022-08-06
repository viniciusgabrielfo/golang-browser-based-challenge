package entity

import (
	"net/http"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

type chatroom struct {
	clients    map[*client]bool
	register   chan *client
	unregister chan *client
	broadcast  chan *message

	logger *zap.SugaredLogger
}

func NewChatroom(logger *zap.SugaredLogger) *chatroom {
	return &chatroom{
		clients:    make(map[*client]bool),
		register:   make(chan *client),
		unregister: make(chan *client),
		broadcast:  make(chan *message),
		logger:     logger,
	}
}

func (c *chatroom) Start() {
	for {
		select {
		case client := <-c.register:
			c.clients[client] = true
			c.logger.Infow("new client registered", "active_clients", len(c.clients))
		case client := <-c.unregister:
			if _, ok := c.clients[client]; ok {
				delete(c.clients, client)
				c.logger.Infow("a client was unregistered", "active_clients", len(c.clients))
				close(client.outbound)
			}
		case message := <-c.broadcast:
			for client := range c.clients {
				client.outbound <- message
			}
		}
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (c *chatroom) HandleWebSocketConn(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		c.logger.Error(err)
		return
	}

	client := NewClient(conn, nil, c)

	c.register <- client
	defer func() {
		c.unregister <- client
	}()

	go client.Write()
	client.Read()
}

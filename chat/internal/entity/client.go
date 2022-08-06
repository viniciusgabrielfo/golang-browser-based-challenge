package entity

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type client struct {
	websocket *websocket.Conn
	user      *user
	chatRoom  *chatroom
	outbound  chan *message
}

func NewClient(wsConn *websocket.Conn, user *user, chatroom *chatroom) *client {
	return &client{
		websocket: wsConn,
		user:      user,
		chatRoom:  chatroom,
		outbound:  make(chan *message),
	}
}

func (c *client) Read() {
	defer c.websocket.Close()

	for {
		_, bMsg, err := c.websocket.ReadMessage()
		fmt.Println(string(bMsg))
		if err != nil {
			if !websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
				c.chatRoom.logger.Error(err)
			}
			break
		}

		c.chatRoom.broadcast <- NewMesssage(uuid.UUID{}, bMsg)
	}
}

func (c *client) Write() {
	defer c.websocket.Close()

	for msg := range c.outbound {
		if err := c.websocket.WriteJSON(msg); err != nil {
			c.chatRoom.logger.Error(err)
			break
		}
	}
}

package entity

import (
	"github.com/gorilla/websocket"
)

type Client struct {
	Websocket *websocket.Conn
	User      *User
	Chatroom  *Chatroom
	outbound  chan *Message
}

func NewClient(wsConn *websocket.Conn, user *User, chatroom *Chatroom) *Client {
	return &Client{
		Websocket: wsConn,
		User:      user,
		Chatroom:  chatroom,
		outbound:  make(chan *Message),
	}
}

func (c *Client) GetOutboundChan() chan *Message {
	return c.outbound
}

// func (c *Client) Read() {
// 	defer c.websocket.Close()

// 	for {
// 		_, bMsg, err := c.websocket.ReadMessage()
// 		if err != nil {
// 			if !websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
// 				// TODO: handle this error correctly
// 				log.Println(err)
// 			}
// 			break
// 		}

// 		c.chatRoom.broadcast <- NewMesssage(uuid.UUID{}, bMsg)
// 	}
// }

// func (c *Client) Write() {
// 	defer c.websocket.Close()

// 	for msg := range c.outbound {
// 		if err := c.websocket.WriteJSON(msg); err != nil {
// 			// TODO: handle this error correctly
// 			log.Println(err)
// 			break
// 		}
// 	}
// }

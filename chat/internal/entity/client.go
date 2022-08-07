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

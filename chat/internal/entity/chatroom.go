package entity

type Chatroom struct {
	Clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan *Message
}

func NewChatroom() *Chatroom {
	return &Chatroom{
		Clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan *Message),
	}
}

func (c *Chatroom) GetBroadcastChan() chan *Message {
	return c.broadcast
}

func (c *Chatroom) GetRegisterChan() chan *Client {
	return c.register
}

func (c *Chatroom) GetUnregisterChan() chan *Client {
	return c.unregister
}

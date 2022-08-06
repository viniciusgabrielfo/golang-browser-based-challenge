package internal

import (
	"log"
	"net/url"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

type webSocketClient struct {
	url        url.URL
	conn       *websocket.Conn
	dispatcher chan []byte

	logger *zap.SugaredLogger
}

func NewWebSocketClient(scheme, host, path string, dispatcher chan []byte, logger *zap.SugaredLogger) (*webSocketClient, error) {
	client := &webSocketClient{
		url:        url.URL{Scheme: scheme, Host: host, Path: path},
		dispatcher: dispatcher,
		logger:     logger,
	}

	if err := client.connect(); err != nil {
		return nil, err
	}

	return client, nil
}

func (c *webSocketClient) connect() error {
	if c.conn != nil {
		return nil
	}

	conn, _, err := websocket.DefaultDialer.Dial(c.url.String(), nil)
	if err != nil {
		return err
	}

	c.logger.Infof("successfull connection with websocket on %s", c.url.String())
	c.conn = conn

	go c.listen()

	return nil
}

func (c *webSocketClient) listen() {
	c.logger.Info("starting listen websocket messages")
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			c.logger.Errorf("error when read message from websocket: %w", err)
			return
		}
		log.Printf("recv: %s", message)

		// c.dispatcher <- message
	}
}

func (c *webSocketClient) Close() error {
	return c.conn.Close()
}

func (c *webSocketClient) Write(msg string) {
	if err := c.conn.WriteMessage(websocket.TextMessage, []byte(msg)); err != nil {
		c.logger.Error("error when try to send message to websocket: %w", err)
	}
}

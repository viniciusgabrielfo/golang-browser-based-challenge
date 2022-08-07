package internal

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

type webSocketClient struct {
	url        url.URL
	conn       *websocket.Conn
	dispatcher chan *Message

	// nick and password used to enable stock-bot to enter in a chatroom
	nick     string
	password string

	log *zap.SugaredLogger
}

func NewWebSocketClient(scheme, host, path, botNick, botPass string, dispatcher chan *Message, log *zap.SugaredLogger) (*webSocketClient, error) {
	client := &webSocketClient{
		url:        url.URL{Scheme: scheme, Host: host, Path: path},
		dispatcher: dispatcher,
		nick:       botNick,
		password:   botPass,
		log:        log,
	}

	if err := client.connect(); err != nil {
		return nil, err
	}

	log.Infof("successfull connection with websocket on %s", client.url.String())

	return client, nil
}

func (c *webSocketClient) connect() error {
	if c.conn != nil {
		return nil
	}

	jwtToken, err := c.auth()
	if err != nil {
		return err
	}

	header := http.Header{}
	header.Add("Authorization", "bearer "+jwtToken)

	conn, _, err := websocket.DefaultDialer.Dial(c.url.String(), header)
	if err != nil {
		return err
	}

	c.conn = conn

	go c.listen()

	return nil
}

func (c *webSocketClient) auth() (string, error) {
	authData := struct {
		Nick     string `json:"nick"`
		Password string `json:"password"`
	}{
		Nick:     c.nick,
		Password: c.password,
	}

	authJSON, _ := json.Marshal(authData)

	req, err := http.NewRequest("POST", "http://"+c.url.Host+"/auth", bytes.NewBuffer(authJSON))
	if err != nil {
		return "", nil
	}

	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", nil
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", nil
	}

	return string(body), nil
}

func (c *webSocketClient) listen() {
	c.log.Info("starting listen websocket messages")
	for {
		var msg *Message
		if err := c.conn.ReadJSON(&msg); err != nil {
			c.log.Errorw("error when try to read message from websocket", "error", err)
			continue
		}

		c.dispatcher <- msg
	}
}

func (c *webSocketClient) Close() error {
	if c.conn == nil {
		return errors.New("no connection open")
	}

	c.conn.WriteMessage(websocket.CloseMessage, []byte(""))
	if err := c.conn.Close(); err != nil {
		return nil
	}
	c.conn = nil

	return nil
}

func (c *webSocketClient) Write(msg string) error {
	if err := c.conn.WriteMessage(websocket.TextMessage, []byte(msg)); err != nil {
		c.log.Errorw("error when try to send message to websocket", "error", err)
		return err
	}

	return nil
}

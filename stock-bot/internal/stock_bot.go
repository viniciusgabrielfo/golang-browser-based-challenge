package internal

import (
	"fmt"
	"strings"
)

type stockBot struct {
	// writer is a function used to stock-bot send your messages trigger by commands
	writer func(msg string) error
	// listener is a channel used to stock-bot receive messages to proccess
	listener chan *Message
	// validCommands stores bot valid commands (key) and if accepts parameters (val)
	validCommands map[string]bool
}

func NewStockBot(writer func(msg string) error, listener chan *Message) *stockBot {
	return &stockBot{
		writer:        writer,
		listener:      listener,
		validCommands: map[string]bool{"stock": true},
	}
}

func (c *stockBot) Start() {
	for {
		if client := <-c.listener; c.isCommand(client.Text) {
			c.proccessCommand(client.Text)
		}
	}
}

func (c *stockBot) isCommand(msg string) bool {
	return msg[0:1] == "/"
}

func (c *stockBot) proccessCommand(msg string) {
	var parameter string
	command := msg[1:]

	if strings.Contains(command, "=") {
		commandSplitted := strings.Split(msg[1:], "=")

		command = commandSplitted[0]
		parameter = commandSplitted[1]
	}

	if acceptParemeter, ok := c.validCommands[command]; ok {
		if parameter != "" && !acceptParemeter {
			c.writer(fmt.Sprintf("the command %s does not accept parameters", command))
		}

		c.writer("one day I will work")

		return
	}

	c.writer(fmt.Sprintf("command %s not found", command))
}

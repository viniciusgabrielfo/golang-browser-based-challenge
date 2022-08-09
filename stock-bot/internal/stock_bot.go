package internal

import (
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/viniciusgabrielfo/golang-browser-based-challenge/stock-bot/pkg/rabbitmq"
)

const stooqGetStockURL = "https://stooq.com/q/l/?s=%s&f=sd2t2ohlcv&h&e=csv"

const (
	MessageCommandNeedParameterError     = "Command '/stock' need a parameter to run. For example: '/stock=parameter'. Please try again."
	MessageCommandParameterNotFoundError = "Command '/stock' need a no empety parameter. For example: '/stock=parameter'. Please try again."
)

type stockBot struct {
	// nick registered for stock-bot
	nick string
	// listener is a channel used to stock-bot receive messages to proccess
	listener chan *ChatMessage
	// validCommands stores bot valid commands (key) and if accepts parameters (val)
	validCommands map[string]func(string, rabbitmq.Producer)
	// used to send stock-bot messages to chatroom
	producer rabbitmq.Producer

	ctx       context.Context
	ctxCancel context.CancelFunc
}

func NewStockBot(nick string, producer rabbitmq.Producer) *stockBot {
	stockBot := &stockBot{
		nick:     nick,
		listener: make(chan *ChatMessage),
		producer: producer,
	}

	stockBot.validCommands = map[string]func(string, rabbitmq.Producer){"stock": stockBot.proccessStockCommand}

	stockBot.ctx, stockBot.ctxCancel = context.WithCancel(context.Background())

	return stockBot
}

func (c *stockBot) GetListener() chan *ChatMessage {
	return c.listener
}

func (c *stockBot) Start() {
	for {
		select {
		case client := <-c.listener:
			if c.isCommand(client.Text) {
				c.proccessCommand(client.Text)
			}
		case <-c.ctx.Done():
			return
		}
	}
}

func (c *stockBot) Stop() {
	c.ctxCancel()
	close(c.listener)
}

func (c *stockBot) isCommand(msg string) bool {
	return msg[0:1] == "/"
}

func (c *stockBot) proccessCommand(msg string) {
	command := msg[1:]

	if strings.Contains(command, "=") {
		commandSplitted := strings.Split(msg[1:], "=")
		command = commandSplitted[0]
	}

	if proccessFunc, ok := c.validCommands[command]; ok {
		proccessFunc(msg[1:], c.producer)
		return
	}

	c.producer.Send(c.nick, fmt.Sprintf("command '%s' not found", command))
}

func (c *stockBot) proccessStockCommand(command string, producer rabbitmq.Producer) {
	commandSplitted := strings.Split(command, "=")

	if len(commandSplitted) < 2 {
		producer.Send(c.nick, MessageCommandNeedParameterError)
		return
	}

	stockName := commandSplitted[1]
	if stockName == "" {
		producer.Send(c.nick, MessageCommandParameterNotFoundError)
		return
	}

	req, err := http.Get(fmt.Sprintf(stooqGetStockURL, stockName))
	if err != nil {
		log.Fatal(err)
	}

	defer req.Body.Close()

	csvReader := csv.NewReader(req.Body)
	data, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	var stockQuote StockQuote
	for i, line := range data {
		if i == 0 {
			// skip header
			continue
		}

		for column, value := range line {
			if column == 0 {
				stockQuote.Symbol = value
			} else if column == 6 {
				parsedQuote, _ := strconv.ParseFloat(value, 64)
				stockQuote.Quote = parsedQuote
			}
		}
	}

	if err := producer.Send(c.nick, stockQuote.String()); err != nil {
		log.Fatal(err)
	}
}

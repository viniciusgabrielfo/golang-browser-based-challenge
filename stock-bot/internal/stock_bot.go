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

type stockBot struct {
	// listener is a channel used to stock-bot receive messages to proccess
	listener chan *Message
	// validCommands stores bot valid commands (key) and if accepts parameters (val)
	validCommands map[string]func(string, rabbitmq.Producer)
	// used to send stock-bot messages to chatroom
	producer rabbitmq.Producer

	ctx       context.Context
	ctxCancel context.CancelFunc
}

func NewStockBot(producer rabbitmq.Producer) *stockBot {
	stockBot := &stockBot{
		listener:      make(chan *Message),
		validCommands: map[string]func(string, rabbitmq.Producer){"stock": proccessStockCommand},
		producer:      producer,
	}

	stockBot.ctx, stockBot.ctxCancel = context.WithCancel(context.Background())

	return stockBot
}

func (c *stockBot) GetListener() chan *Message {
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

	fmt.Println(command)
	if proccessFunc, ok := c.validCommands[command]; ok {
		proccessFunc(msg[1:], c.producer)
		return
	}

	c.producer.Send(fmt.Sprintf("command '%s' not found", command))
}

const stooqGetStockURL = "https://stooq.com/q/l/?s=%s&f=sd2t2ohlcv&h&e=csv"

const (
	MessageCommandNeedParameterError     = "Command '/stock' need a parameter to run. For example: '/stock=parameter'. Please try again."
	MessageCommandParameterNotFoundError = "Command '/stock' need a no empety parameter. For example: '/stock=parameter'. Please try again."
)

func proccessStockCommand(command string, producer rabbitmq.Producer) {
	commandSplitted := strings.Split(command, "=")

	if len(commandSplitted) < 2 {
		producer.Send(MessageCommandNeedParameterError)
		return
	}

	stockName := commandSplitted[1]
	if stockName == "" {
		producer.Send(MessageCommandParameterNotFoundError)
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

	if err := producer.Send(stockQuote.String()); err != nil {
		log.Fatal(err)
	}
}

type StockQuote struct {
	Symbol string
	Quote  float64
}

func (s *StockQuote) String() string {
	return fmt.Sprintf("%s quote is $%.2f per share", s.Symbol, s.Quote)
}

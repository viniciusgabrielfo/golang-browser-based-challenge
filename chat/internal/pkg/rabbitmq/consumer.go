package rabbitmq

import (
	"fmt"
	"log"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/viniciusgabrielfo/golang-browser-based-challenge/chat/internal/entity"
)

type Consumer interface {
	Start(dispatcher chan *entity.Message) error
	Close() error
}

type consumer struct {
	channel *amqp.Channel
	queue   string
	tag     string
	done    chan error
}

func newConsumer(channel *amqp.Channel, consumerName, queue string) *consumer {
	return &consumer{
		channel: channel,
		queue:   queue,
		tag:     consumerName,
		done:    make(chan error),
	}
}

func (c *consumer) Start(dispatcher chan *entity.Message) error {
	deliveries, err := c.channel.Consume(
		c.queue,
		c.tag,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		defer func() {
			c.done <- nil
		}()

		for d := range deliveries {
			dispatcher <- entity.NewMesssage(uuid.UUID{}, d.Headers["nick"].(string), d.Body) // TODO: try to get correct UUID
		}
	}()

	return nil
}

func (c *consumer) Close() error {
	if err := c.channel.Cancel(c.tag, true); err != nil {
		return fmt.Errorf("consumer cancel failed: %s", err)
	}

	return <-c.done
}

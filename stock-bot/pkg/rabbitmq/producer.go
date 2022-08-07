package rabbitmq

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Producer interface {
	Send(msg string) error
}

type producer struct {
	channel  *amqp.Channel
	exchange string
	routing  string
}

func newProducer(channel *amqp.Channel, exchange string, routing string) *producer {
	return &producer{
		channel:  channel,
		exchange: exchange,
		routing:  routing,
	}
}

func (p *producer) Send(msg string) error {
	if err := p.channel.Publish(
		p.exchange,
		p.routing,
		false,
		false,
		amqp.Publishing{
			Headers:         amqp.Table{},
			ContentType:     "text/plain",
			ContentEncoding: "",
			Body:            []byte(msg),
			DeliveryMode:    amqp.Transient,
			Priority:        0,
		},
	); err != nil {
		return fmt.Errorf("exchange Publish: %w", err)
	}

	return nil
}

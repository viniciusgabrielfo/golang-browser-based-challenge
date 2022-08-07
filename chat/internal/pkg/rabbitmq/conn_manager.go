package rabbitmq

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type AmqpConfig struct {
	User     string
	Password string
	Host     string
}

func (s *AmqpConfig) String() string {
	return fmt.Sprintf("amqp://%s:....@%s", s.User, s.Host)
}

type AmqpConnManager struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	config  AmqpConfig
}

func NewAmqpConnManager(config AmqpConfig) (*AmqpConnManager, error) {
	connManager := &AmqpConnManager{
		config: config,
	}

	conn, err := connManager.dial()
	if err != nil {
		return nil, err
	}

	connManager.conn = conn

	channel, err := connManager.conn.Channel()
	if err != nil {
		return nil, err
	}

	connManager.channel = channel

	return connManager, nil
}

func (m *AmqpConnManager) Close() error {
	if err := m.channel.Close(); err != nil {
		return err
	}

	return m.conn.Close()
}

func (m *AmqpConnManager) dial() (*amqp.Connection, error) {
	url := fmt.Sprintf("amqp://%s:%s@%s", m.config.User, m.config.Password, m.config.Host)
	return amqp.Dial(url)
}

func (m *AmqpConnManager) ExchangeDeclare(exchangeName, kind string) error {
	return m.channel.ExchangeDeclare(exchangeName,
		kind,
		true,
		false,
		false,
		false,
		nil,
	)
}

func (m *AmqpConnManager) QueueDeclare(queueName string) error {
	_, err := m.channel.QueueDeclare(queueName,
		true,
		false,
		false,
		false,
		nil,
	)

	return err
}

func (m *AmqpConnManager) QeueBind(queueName, key, exchange string) error {
	return m.channel.QueueBind(
		queueName,
		key,
		exchange,
		false,
		nil,
	)
}

func (m *AmqpConnManager) CreateConsumer(name, queueName string) Consumer {
	return newConsumer(m.channel, name, queueName)
}

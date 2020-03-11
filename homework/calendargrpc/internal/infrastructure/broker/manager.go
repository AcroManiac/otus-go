package broker

import (
	"fmt"
	"io"

	"github.com/pkg/errors"
	"github.com/streadway/amqp"
)

const (
	exchangeName = "otus.homework"
	queueName    = "senderQueue"
	routingKey   = "senderQueue.events"
)

type Manager struct {
	Protocol string
	User     string
	Password string
	Host     string
	Port     int
	Conn     *amqp.Connection
	Ch       *amqp.Channel
	Que      amqp.Queue
	wr       io.Writer
}

func NewManager(protocol, user, password, host string, port int) *Manager {
	return &Manager{
		Protocol: protocol,
		User:     user,
		Password: password,
		Host:     host,
		Port:     port,
	}
}

func (m *Manager) Open() error {
	var err error
	connUrl := fmt.Sprintf("%s://%s:%s@%s:%d/", m.Protocol, m.User, m.Password, m.Host, m.Port)

	// Open connection to broker
	m.Conn, err = amqp.Dial(connUrl)
	if err != nil {
		return errors.Wrap(err, "failed to connect to RabbitMQ")
	}

	// Open channel
	m.Ch, err = m.Conn.Channel()
	if err != nil {
		return errors.Wrap(err, "failed to open a channel")
	}

	// Open exchange
	err = m.Ch.ExchangeDeclare(
		exchangeName, // name
		"topic",      // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		return errors.Wrap(err, "failed to declare an exchange")
	}

	// Create queue
	m.Que, err = m.Ch.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when unused
		true,      // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return errors.Wrap(err, "failed to declare a queue")
	}

	// Binding queue to exchange
	err = m.Ch.QueueBind(
		m.Que.Name,   // queue name
		routingKey,   // routing key
		exchangeName, // exchange
		false,
		nil)
	if err != nil {
		return errors.Wrap(err, "failed to bind a queue")
	}

	return nil
}

func (m *Manager) GetWriter() io.Writer {
	if m.wr == nil {
		// Create broker writer
		m.wr = NewAmqpWriter(m.Ch)
	}
	return m.wr
}

func (m *Manager) Close() error {
	// Close queue
	if _, err := m.Ch.QueueDelete(m.Que.Name, false, false, true); err != nil {
		return errors.Wrap(err, "error deleting queue")
	}

	// Close channel
	if m.Ch != nil {
		if err := m.Ch.Close(); err != nil {
			return errors.Wrap(err, "failed to close channel")
		}
	}
	// Close connection
	if m.Conn != nil {
		if err := m.Conn.Close(); err != nil {
			return errors.Wrap(err, "failed closing connection")
		}
	}
	return nil
}

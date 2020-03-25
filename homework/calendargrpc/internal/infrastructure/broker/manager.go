package broker

import (
	"context"
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
	wr       io.WriteCloser
	rd       io.ReadCloser
}

func NewManager(protocol, user, password, host string, port int) *Manager {
	connUrl := fmt.Sprintf("%s://%s:%s@%s:%d/", protocol, user, password, host, port)

	// Open connection to broker
	conn, err := amqp.Dial(connUrl)
	if err != nil {
		return nil
	}

	return &Manager{
		Protocol: protocol,
		User:     user,
		Password: password,
		Host:     host,
		Port:     port,
		Conn:     conn,
	}
}

func (m *Manager) GetWriter() io.Writer {
	if m.wr == nil {
		// Create broker writer
		m.wr = NewAmqpWriter(m.Conn)
	}
	return m.wr
}

func (m *Manager) GetReader(ctx context.Context) io.Reader {
	if m.rd == nil {
		// Create broker reader
		m.rd = NewAmqpReader(ctx, m.Conn)
	}
	return m.rd
}

func (m *Manager) Close() error {
	// Close i/o channels
	if m.rd != nil {
		if err := m.rd.Close(); err != nil {
			return errors.Wrap(err, "failed closing reader")
		}
	}
	if m.wr != nil {
		if err := m.wr.Close(); err != nil {
			return errors.Wrap(err, "failed closing writer")
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

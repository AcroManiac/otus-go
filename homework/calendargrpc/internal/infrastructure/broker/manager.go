package broker

import (
	"context"
	"fmt"
	"github.com/AcroManiac/otus-go/homework/calendargrpc/internal/infrastructure/logger"
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

	logger.Info("RabbitMQ broker connected", "host", m.Host)

	return nil
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

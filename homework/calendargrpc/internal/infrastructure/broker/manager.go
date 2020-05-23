package broker

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/ahamtat/otus-go/homework/calendargrpc/internal/infrastructure/logger"
	"github.com/pkg/errors"
	"github.com/streadway/amqp"

	"github.com/cenkalti/backoff/v4"
)

const (
	exchangeName = "otus.homework"
	queueName    = "senderQueue"
	routingKey   = "senderQueue.events"
)

type Manager struct {
	connUrl string
	Conn    *amqp.Connection
	Done    chan error
	wr      io.WriteCloser
	rd      io.ReadCloser
}

func NewManager(protocol, user, password, host string, port int) *Manager {

	// Create manager object
	m := &Manager{
		connUrl: fmt.Sprintf("%s://%s:%s@%s:%d/", protocol, user, password, host, port),
		Conn:    nil,
		Done:    make(chan error),
	}

	// Open connection to broker
	if err := m.connect(); err != nil {
		logger.Error("RabbitMQ connection failed", "error", err)
		return nil
	}

	return m
}

func (m *Manager) connect() (err error) {

	// Open RabbitMQ connection
	m.Conn, err = amqp.Dial(m.connUrl)
	if err != nil {
		return
	}

	return
}

func (m *Manager) ConnectionListener(ctx context.Context) {
	select {
	case <-ctx.Done():
		break
	case connErr := <-m.Conn.NotifyClose(make(chan *amqp.Error)):
		logger.Error("RabbitMQ connection is closed", "error", connErr.Error())
		// Notify clients
		m.Done <- errors.New("connection closed")
	}
}

func (m *Manager) Reconnect() error {

	// Close i/o channels
	_ = m.closeIOChannels()

	// Create reconnect backoff
	be := backoff.NewExponentialBackOff()
	be.MaxElapsedTime = time.Minute
	be.InitialInterval = 1 * time.Second
	be.Multiplier = 2
	be.MaxInterval = 15 * time.Second

	// Do reconnect loop
	boCtx := backoff.WithContext(be, context.Background())
	for {
		boTime := boCtx.NextBackOff()
		if boTime == backoff.Stop {
			return errors.New("backoff timer elapsed")
		}

		select {
		case <-time.After(boTime):
			if err := m.connect(); err != nil {
				logger.Error("couldn't reconnect", "error", err)
				continue
			}
			logger.Info("Reconnect to RabbitMQ succeeded")
			return nil
		}
	}

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

func (m *Manager) closeIOChannels() error {

	if m.rd != nil {
		if err := m.rd.Close(); err != nil {
			//return errors.Wrap(err, "failed closing reader")
		}
		m.rd = nil
	}
	if m.wr != nil {
		if err := m.wr.Close(); err != nil {
			//return errors.Wrap(err, "failed closing writer")
		}
		m.wr = nil
	}
	return nil
}

func (m *Manager) Close() error {

	// Close i/o channels
	_ = m.closeIOChannels()

	// Close connection notify channel
	if m.Done != nil {
		close(m.Done)
	}

	// Close connection
	if m.Conn != nil {
		if err := m.Conn.Close(); err != nil {
			logger.Error("failed closing connection", "error", err)
		}
	}

	return nil
}

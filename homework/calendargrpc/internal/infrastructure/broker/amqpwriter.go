package broker

import (
	"io"

	"github.com/pkg/errors"

	"github.com/streadway/amqp"
)

type AmqpWriter struct {
	ch *amqp.Channel
}

func NewAmqpWriter(ch *amqp.Channel) io.Writer {
	return &AmqpWriter{ch: ch}
}

// Write message to RabbitMQ broker.
// Returns message length on success or error if any
func (w *AmqpWriter) Write(p []byte) (n int, err error) {
	if w.ch == nil {
		return 0, errors.New("no output channel defined")
	}

	// Send message to gateway
	err = w.ch.Publish(
		exchangeName, // exchange
		routingKey,   // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        p,
		})
	if err != nil {
		return 0, errors.Wrap(err, "failed to publish a message")
	}

	return len(p), nil
}

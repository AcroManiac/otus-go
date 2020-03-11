package broker

import (
	"context"
	"io"

	"github.com/AcroManiac/otus-go/homework/calendargrpc/internal/infrastructure/logger"
	"github.com/pkg/errors"

	"github.com/streadway/amqp"
)

type AmqpReader struct {
	ctx  context.Context
	msgs <-chan amqp.Delivery
}

func NewAmqpReader(ctx context.Context, ch *amqp.Channel, que *amqp.Queue) io.Reader {

	// Create consuming channel
	msgs, err := ch.Consume(
		que.Name, // queue
		"",       // consumer
		true,     // auto ack
		false,    // exclusive
		false,    // no local
		false,    // no wait
		nil,      // args
	)
	if err != nil {
		logger.Error("failed to register a consumer", "error", err)
		return nil
	}

	// Create reader object
	return &AmqpReader{
		ctx:  ctx,
		msgs: msgs,
	}
}

// Read one message from RabbitMQ queue. Returns message length in bytes
func (r *AmqpReader) Read(p []byte) (n int, err error) {
	select {
	case <-r.ctx.Done():
		err = errors.New("context cancelled")
	case message, ok := <-r.msgs:
		if ok {
			n = copy(p, message.Body)
		}
	}
	return
}

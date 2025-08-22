package registers

import (
	"fmt"
	"time"

	"github.com/eliabe-portfolio/restaurant-app/internal/connections"
	"github.com/eliabe-portfolio/restaurant-app/internal/constants"
	"github.com/rabbitmq/amqp091-go"
)

const (
	DefaultTimeout = 5 * time.Second
)

type QueueHandler struct {
	channel *amqp091.Channel
}

func New(connections *connections.Provider) *QueueHandler {
	return &QueueHandler{channel: connections.RabbitMQ.Get()}
}

func (h *QueueHandler) declareQueue(queueName string) error {
	if queueName == "" {
		return fmt.Errorf("queue name cannot be empty")
	}

	_, err := h.channel.QueueDeclare(
		queueName,
		true,  // durable
		false, // autoDelete
		false, // exclusive
		false, // noWait
		nil,   // args
	)
	if err != nil {
		return fmt.Errorf("failed to declare queue %q: %w", queueName, err)
	}
	return nil
}

func (h *QueueHandler) DeclareAllQueues() {
	h.declareQueue(string(constants.Queues.ResetPasswordEmail))
}

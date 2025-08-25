package sendresetpasswordemailproducer

import (
	"encoding/json"
	"fmt"

	"github.com/eliabe-portfolio/restaurant-app/internal/connections"
	"github.com/eliabe-portfolio/restaurant-app/internal/constants"
	"github.com/rabbitmq/amqp091-go"
)

type Producer interface {
	Send(message SendPasswordResetEmailMessage) error
}

type SendPasswordResetEmailMessage struct {
	ResetPasswordToken string
	RandomHash         string
}

type producer struct {
	channel *amqp091.Channel
}

func New(connections *connections.Provider) Producer {
	return &producer{channel: connections.RabbitMQ.Get()}
}

func (p *producer) Send(message SendPasswordResetEmailMessage) error {
	body, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("could not marshal message: %w", err)
	}

	err = p.channel.Publish(
		"", // exchange
		constants.Queues.ResetPasswordEmail,
		false, // mandatory
		false, // immediate
		amqp091.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp091.Persistent,
			Body:         body,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish to queue %q: %w", constants.Queues.ResetPasswordEmail, err)
	}

	return nil
}

package consumers

import (
	"context"
	"fmt"
	"log"

	"github.com/eliabe-portfolio/restaurant-app/internal/connections"
	"github.com/eliabe-portfolio/restaurant-app/internal/constants"
	sendresetpasswordemailconsumer "github.com/eliabe-portfolio/restaurant-app/internal/queues/consumers/send-reset-password-email"
	"github.com/eliabe-portfolio/restaurant-app/internal/repositories"
	sendresetpasswordemailcmd "github.com/eliabe-portfolio/restaurant-app/internal/use-cases/notification/email/send-reset-pwd-email"
	"github.com/eliabe-portfolio/restaurant-app/pkg/returns"
	"github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	channel      *amqp091.Channel
	repositories *repositories.Provider
}

func New(connections *connections.Provider, repositories *repositories.Provider) *Consumer {
	return &Consumer{
		channel:      connections.RabbitMQ.Get(),
		repositories: repositories,
	}
}

func (c *Consumer) Start(ctx context.Context) error {
	deliveries, err := c.channel.Consume(
		constants.Queues.ResetPasswordEmail,
		"",    // consumerTag
		false, // autoAck
		false, // exclusive
		false, // noLocal
		false, // noWait
		nil,   // args
	)
	if err != nil {
		return fmt.Errorf("could not start consuming from %q: %w", constants.Queues.ResetPasswordEmail, err)
	}

	go c.handler(ctx, deliveries)

	return nil
}

func (c *Consumer) handler(ctx context.Context, deliveries <-chan amqp091.Delivery) {
	for {
		select {
		case <-ctx.Done():
			log.Printf("Consumer for queue %q stopped: %v", constants.Queues.ResetPasswordEmail, ctx.Err())
			return
		case d, ok := <-deliveries:
			if !ok {
				log.Printf("Channel closed for queue %q", constants.Queues.ResetPasswordEmail)
				return
			}
			if d.RoutingKey == constants.Queues.ResetPasswordEmail {
				consumer := sendresetpasswordemailconsumer.New(c.repositories)
				consumer.Process(ctx, d)
				return
			}
			if d.RoutingKey == constants.Queues.InviteUserEmail {
				consumer := sendresetpasswordemailconsumer.New(c.repositories)
				consumer.Process(ctx, d)
				return
			}
		}
	}
}

type ConsumerProvider interface {
	SendPasswordResetEmail(message sendresetpasswordemailcmd.Params) (returns.Api, error)
}

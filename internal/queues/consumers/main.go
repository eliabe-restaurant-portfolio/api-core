package consumers

import (
	"log"

	"github.com/eliabe-restaurant-portfolio/api-core/internal/adapters"
	"github.com/eliabe-restaurant-portfolio/api-core/internal/connections"
	"github.com/eliabe-restaurant-portfolio/api-core/internal/constants"
	sendresetpasswordemailconsumer "github.com/eliabe-restaurant-portfolio/api-core/internal/queues/consumers/send-reset-password-email"
	"github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	channel  *amqp091.Channel
	adapters *adapters.Adapters
}

func New(connections *connections.Provider, adapters *adapters.Adapters) *Consumer {
	return &Consumer{
		channel:  (*connections).RabbitMQ.Get(),
		adapters: adapters,
	}
}

func (c *Consumer) Start() {
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
		panic(err)
	}

	go c.handler(deliveries)
}

func (c *Consumer) handler(deliveries <-chan amqp091.Delivery) {
	for {
		d, ok := <-deliveries
		if !ok {
			log.Printf("Channel closed for queue %q", constants.Queues.ResetPasswordEmail)
			return
		}
		if d.RoutingKey == constants.Queues.ResetPasswordEmail {
			consumer := sendresetpasswordemailconsumer.New(c.adapters)
			consumer.Process(d)
			return
		}
	}
}

package rabbitmq

import (
	"fmt"
	"log"
	"time"

	"github.com/eliabe-portfolio/restaurant-app/internal/connections/configs"
	"github.com/eliabe-portfolio/restaurant-app/internal/envs"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Connection struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

func (rc *Connection) Get() *amqp.Channel {
	return rc.channel
}

func (rc *Connection) Close() error {
	if err := rc.channel.Close(); err != nil {
		return fmt.Errorf("failed to close RabbitMQ channel: %v", err)
	}
	if err := rc.conn.Close(); err != nil {
		return fmt.Errorf("failed to close RabbitMQ connection: %v", err)
	}
	log.Printf("✅ rabbitmq connection closed.")
	return nil
}

func Connect(conf *configs.Config) (*Connection, error) {
	fmt.Println(conf.RabbitMQ)

	var conn *amqp.Connection
	var err error
	maxRetries := 10
	retryDelay := 5 * time.Second

	for i := 0; i < maxRetries; i++ {
		conn, err = amqp.Dial(conf.RabbitMQ)
		if err == nil {
			break
		}
		log.Printf("WARNING: Could not dial RabbitMQ (attempt %d/%d): %v. Retrying in %s...", i+1, maxRetries, err, retryDelay)
		time.Sleep(retryDelay)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("could not open channel: %v", err)
	}

	if err := ch.Qos(selectPrefetchCount(), 0, false); err != nil {
		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("could not set QoS: %v", err)
	}

	log.Printf("✅ rabbitmq connection established.")
	return &Connection{
		conn:    conn,
		channel: ch,
	}, nil
}

func selectPrefetchCount() int {
	if envs.IsDev() {
		return 1
	}
	return 10
}

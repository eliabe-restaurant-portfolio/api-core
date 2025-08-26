package configs

import (
	"fmt"

	"github.com/eliabe-restaurant-portfolio/api-core/internal/envs"
)

type Config struct {
	ServerName string
	ServerPort string
	Postgres   string
	RabbitMQ   string
}

func New() *Config {
	return &Config{
		ServerName: envs.Get(envs.SERVER_NAME),
		ServerPort: envs.Get(envs.SERVER_PORT),
		Postgres:   buildPostgresPath(),
		RabbitMQ:   buildRabbitMQPath(),
	}
}

func buildPostgresPath() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		envs.Get(envs.POSTGRES_USERNAME),
		envs.Get(envs.POSTGRES_PASSWORD),
		envs.Get(envs.POSTGRES_HOST),
		envs.Get(envs.POSTGRES_PORT),
		envs.Get(envs.POSTGRES_DATABASE),
	)
}

func buildRabbitMQPath() string {
	return fmt.Sprintf(
		"amqp://%s:%s@%s:%s",
		envs.Get(envs.RABBITMQ_USERNAME),
		envs.Get(envs.RABBITMQ_PASSWORD),
		envs.Get(envs.RABBITMQ_HOST),
		envs.Get(envs.RABBITMQ_PORT),
	)
}

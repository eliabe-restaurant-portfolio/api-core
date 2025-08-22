package connections

import (
	"github.com/eliabe-portfolio/restaurant-app/internal/connections/configs"
	"github.com/eliabe-portfolio/restaurant-app/internal/connections/postgres"
	"github.com/eliabe-portfolio/restaurant-app/internal/connections/rabbitmq"
	"gorm.io/gorm"
)

type Provider struct {
	config   *configs.Config
	Postgres postgres.Connection
	RabbitMQ rabbitmq.Connection
}

func New(conf *configs.Config) *Provider {
	return &Provider{config: conf}
}

func (c *Provider) ConnectPostgres() {
	conn, err := postgres.Connect(c.config)
	if err != nil {
		panic(err)
	}

	c.Postgres = *conn
}

func (c *Provider) ConnectRabbitMQ() {
	conn, err := rabbitmq.Connect(c.config)
	if err != nil {
		panic(err)
	}

	c.RabbitMQ = *conn
}

func (c *Provider) ClosePostgres() {
	if err := c.Postgres.Close(); err != nil {
		panic(err)
	}
}

func (c *Provider) GetPostgres() *gorm.DB {
	return c.Postgres.Get()
}

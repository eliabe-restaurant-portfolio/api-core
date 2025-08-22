package queues

import (
	"context"

	"github.com/eliabe-portfolio/restaurant-app/internal/connections"
	"github.com/eliabe-portfolio/restaurant-app/internal/queues/consumers"
	"github.com/eliabe-portfolio/restaurant-app/internal/queues/producers"
	"github.com/eliabe-portfolio/restaurant-app/internal/queues/registers"
	"github.com/eliabe-portfolio/restaurant-app/internal/repositories"
)

type Provider struct {
	Producers producers.Provider
}

func New(connections *connections.Provider, repositories *repositories.Provider) Provider {
	handler := registers.New(connections)
	handler.DeclareAllQueues()

	producer := producers.New(connections)
	consumer := consumers.New(connections, repositories)
	consumer.Start(context.Background())

	return Provider{
		Producers: producer,
	}
}

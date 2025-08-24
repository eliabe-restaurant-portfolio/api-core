package queues

import (
	"github.com/eliabe-portfolio/restaurant-app/internal/connections"
	"github.com/eliabe-portfolio/restaurant-app/internal/queues/producers"
	"github.com/eliabe-portfolio/restaurant-app/internal/queues/registers"
	"github.com/eliabe-portfolio/restaurant-app/internal/repositories"
)

type Provider interface {
	Producers() producers.Provider
}

type provider struct {
	producers producers.Provider
}

func New(connections *connections.Provider, repositories *repositories.Provider) Provider {
	handler := registers.New(connections)
	handler.DeclareAllQueues()

	producer := producers.New(connections)

	return provider{
		producers: producer,
	}
}

func (p provider) Producers() producers.Provider {
	return p.producers
}

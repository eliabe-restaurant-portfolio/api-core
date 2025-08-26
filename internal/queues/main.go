package queues

import (
	"github.com/eliabe-restaurant-portfolio/api-core/internal/connections"
	"github.com/eliabe-restaurant-portfolio/api-core/internal/queues/producers"
	"github.com/eliabe-restaurant-portfolio/api-core/internal/queues/registers"
	"github.com/eliabe-restaurant-portfolio/api-core/internal/repositories"
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

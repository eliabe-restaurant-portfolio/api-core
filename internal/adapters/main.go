package adapters

import (
	"github.com/eliabe-portfolio/restaurant-app/internal/connections"
	"github.com/eliabe-portfolio/restaurant-app/internal/middlewares"
	"github.com/eliabe-portfolio/restaurant-app/internal/queues"
	"github.com/eliabe-portfolio/restaurant-app/internal/queues/producers"

	"github.com/eliabe-portfolio/restaurant-app/internal/repositories"
	uow "github.com/eliabe-portfolio/restaurant-app/internal/unit-of-work"
)

type Adapters interface {
	Repositories() repositories.Provider
	Middlewares() middlewares.Provider
	UnitOfWork() uow.UnitOfWork
	Producers() producers.Provider
}

type adapters struct {
	repositories repositories.Provider
	middlewares  middlewares.Provider
	unitOfWork   uow.UnitOfWork
	producers    producers.Provider
}

func New(connections *connections.Provider) Adapters {
	repositories := repositories.New(connections)

	middlewares := middlewares.New()

	unitOfWork := uow.New(connections)

	queueHandler := queues.New(connections, &repositories)
	producers := queueHandler.Producers()

	return adapters{
		repositories,
		middlewares,
		unitOfWork,
		producers,
	}
}

func (a adapters) Repositories() repositories.Provider {
	return a.repositories
}
func (a adapters) Middlewares() middlewares.Provider {
	return a.middlewares
}
func (a adapters) UnitOfWork() uow.UnitOfWork {
	return a.unitOfWork
}
func (a adapters) Producers() producers.Provider {
	return a.producers
}

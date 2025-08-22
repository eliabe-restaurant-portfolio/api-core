package userhdl

import (
	"github.com/eliabe-portfolio/restaurant-app/internal/adapters"
	"github.com/eliabe-portfolio/restaurant-app/internal/middlewares"
	"github.com/eliabe-portfolio/restaurant-app/internal/queues/producers"
	"github.com/eliabe-portfolio/restaurant-app/internal/repositories"
	uow "github.com/eliabe-portfolio/restaurant-app/internal/unit-of-work"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	repositories repositories.Provider
	uow          uow.UnitOfWork
	producers    producers.Provider
	middlewares  middlewares.Provider
}

func New(adapters adapters.Adapters) *UserHandler {
	return &UserHandler{
		repositories: adapters.Repositories(),
		uow:          adapters.UnitOfWork(),
		producers:    adapters.Producers(),
		middlewares:  adapters.Middlewares(),
	}
}

func (hdl *UserHandler) Register(router *gin.Engine) {
	group := router.Group("/api")

	group.POST("users", hdl.Create)

	group.Use(hdl.middlewares.BearerAuth())
	{
		group.DELETE("users", hdl.Delete)
	}
}

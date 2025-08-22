package authhdl

import (
	"github.com/eliabe-portfolio/restaurant-app/internal/adapters"
	"github.com/eliabe-portfolio/restaurant-app/internal/middlewares"
	"github.com/eliabe-portfolio/restaurant-app/internal/queues/producers"
	"github.com/eliabe-portfolio/restaurant-app/internal/repositories"
	uow "github.com/eliabe-portfolio/restaurant-app/internal/unit-of-work"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	repositories repositories.Provider
	uow          uow.UnitOfWork
	producers    producers.Provider
	middlewares  middlewares.Provider
}

func New(adapters adapters.Adapters) *AuthHandler {
	return &AuthHandler{
		repositories: adapters.Repositories(),
		uow:          adapters.UnitOfWork(),
		producers:    adapters.Producers(),
		middlewares:  adapters.Middlewares(),
	}
}

func (hdl *AuthHandler) Register(router *gin.Engine) {
	group := router.Group("/api/auth")

	group.POST("login", hdl.Login)
	group.POST("activate", hdl.ActivateLogin)
	group.POST("reset", hdl.ResetPassword)

	group.Use(hdl.middlewares.BearerAuth())
	{
		group.POST("change", hdl.ChangePassword)
	}
}

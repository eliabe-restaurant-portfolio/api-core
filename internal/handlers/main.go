package handlers

import (
	"github.com/eliabe-restaurant-portfolio/api-core/internal/adapters"
	healthyhdl "github.com/eliabe-restaurant-portfolio/api-core/internal/handlers/healthy"
	authhdl "github.com/eliabe-restaurant-portfolio/api-core/internal/handlers/password-auth"
	userhdl "github.com/eliabe-restaurant-portfolio/api-core/internal/handlers/user"
	"github.com/gin-gonic/gin"
)

type Handlers interface {
	Register(r *gin.Engine)
}

func New(apt *adapters.Adapters) []Handlers {
	return []Handlers{
		authhdl.New(apt),
		userhdl.New(apt),
		healthyhdl.New(apt),
	}
}

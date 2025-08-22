package handlers

import (
	"github.com/eliabe-portfolio/restaurant-app/internal/adapters"
	authhdl "github.com/eliabe-portfolio/restaurant-app/internal/handlers/password-auth"
	userhdl "github.com/eliabe-portfolio/restaurant-app/internal/handlers/user"
	"github.com/gin-gonic/gin"
)

type Handlers interface {
	Register(r *gin.Engine)
}

func New(apt *adapters.Adapters) []Handlers {
	return []Handlers{
		authhdl.New(*apt),
		userhdl.New(*apt),
	}
}

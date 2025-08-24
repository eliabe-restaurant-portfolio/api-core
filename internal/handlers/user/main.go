package userhdl

import (
	"github.com/eliabe-portfolio/restaurant-app/internal/adapters"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	adapters *adapters.Adapters
}

func New(adapters *adapters.Adapters) *UserHandler {
	return &UserHandler{
		adapters: adapters,
	}
}

func (hdl *UserHandler) Register(router *gin.Engine) {
	group := router.Group("/api")

	group.POST("users", hdl.Create)
}

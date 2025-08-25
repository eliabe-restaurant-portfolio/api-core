package authhdl

import (
	"github.com/eliabe-portfolio/restaurant-app/internal/adapters"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	adapters *adapters.Adapters
}

func New(adapters *adapters.Adapters) *AuthHandler {
	return &AuthHandler{
		adapters: adapters,
	}
}

func (hdl *AuthHandler) Register(router *gin.Engine) {
	group := router.Group("/api/auth")

	group.POST("login", hdl.Login)
	group.POST("activate", hdl.ActivateUser)
	group.POST("reset", hdl.RequestResetPassword)

	group.Use((*hdl.adapters).Middlewares().BearerAuth())
	{
		group.POST("change", hdl.ChangePassword)
	}
}

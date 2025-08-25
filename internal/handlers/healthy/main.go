package healthyhdl

import (
	"github.com/eliabe-portfolio/restaurant-app/internal/adapters"
	"github.com/gin-gonic/gin"
)

type HealthyHandler struct {
}

func New(_ *adapters.Adapters) *HealthyHandler {
	return &HealthyHandler{}
}

func (hdl *HealthyHandler) Register(router *gin.Engine) {
	group := router.Group("/")

	group.GET("", hdl.Healthy)
}

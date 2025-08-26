package healthyhdl

import (
	"github.com/eliabe-restaurant-portfolio/api-core/internal/adapters"
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

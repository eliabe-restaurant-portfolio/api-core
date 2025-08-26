package healthyhdl

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (hdl HealthyHandler) Healthy(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, "Hello, World!")
}

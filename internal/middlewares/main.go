package middlewares

import (
	"github.com/gin-gonic/gin"
)

type Provider interface {
	BearerAuth() gin.HandlerFunc
}

type middlewares struct {
}

func New() Provider {
	return middlewares{}
}

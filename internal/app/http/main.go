package httpserver

import (
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type HttpServer struct {
	router *gin.Engine
}

func New() *HttpServer {
	router := gin.New()

	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "route not found",
		})
	})

	router.NoMethod(func(c *gin.Context) {
		c.JSON(http.StatusMethodNotAllowed, gin.H{
			"error": "method not allowed",
		})
	})

	return &HttpServer{router: router}
}

func (hs *HttpServer) ConfigureTrustedProxies() *HttpServer {
	err := hs.router.SetTrustedProxies(nil)
	if err != nil {
		panic(err)
	}
	return hs
}

func (hs *HttpServer) ConfigureCors() *HttpServer {
	hs.router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	return hs
}

func (hs *HttpServer) ConfigureLogs() *HttpServer {
	hs.router.Use(gin.Logger())
	return hs
}

func (hs *HttpServer) Router() *gin.Engine {
	return hs.router
}

func (hs *HttpServer) ConfigureRecoveryPanic() *HttpServer {
	hs.router.Use(func(c *gin.Context) {
		defer func() {
			if rec := recover(); rec != nil {
				log.Printf("recovered from panic: %v\nstack:\n%s", rec, debug.Stack())

				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"error":   fmt.Sprintf("Erro interno no servidor: %v", rec),
					"message": "Ocorreu um problema inesperado. Tente novamente mais tarde.",
				})
			}
		}()

		c.Next()
	})

	return hs
}

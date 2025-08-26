package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/eliabe-restaurant-portfolio/api-core/internal/adapters"
	httpserver "github.com/eliabe-restaurant-portfolio/api-core/internal/app/http"
	"github.com/eliabe-restaurant-portfolio/api-core/internal/connections"
	"github.com/eliabe-restaurant-portfolio/api-core/internal/connections/configs"
	"github.com/eliabe-restaurant-portfolio/api-core/internal/envs"
	"github.com/eliabe-restaurant-portfolio/api-core/internal/handlers"
	"github.com/eliabe-restaurant-portfolio/api-core/internal/queues/consumers"
	"github.com/gin-gonic/gin"
)

func main() {
	envs.Load()

	conf := configs.New()

	conns := connections.New(conf)
	conns.ConnectPostgres()
	conns.ConnectRabbitMQ()

	adapters := adapters.New(conns)

	consumers := consumers.New(conns, &adapters)
	consumers.Start()

	httpServer := httpserver.New().
		ConfigureTrustedProxies().
		ConfigureRecoveryPanic().
		ConfigureLogs().
		ConfigureCors()

	app := New(httpServer, conf)
	app.StartRoutes(&adapters)
	app.Run()
}

type App struct {
	Server *http.Server
	Router *gin.Engine
}

func New(
	server *httpserver.HttpServer,
	conf *configs.Config,
) *App {
	return &App{
		Router: server.Router(),
		Server: &http.Server{
			Addr:              fmt.Sprintf("0.0.0.0:%s", conf.ServerPort),
			Handler:           server.Router(),
			ReadTimeout:       time.Second * 15,
			ReadHeaderTimeout: time.Second * 15,
			WriteTimeout:      time.Second * 15,
			IdleTimeout:       time.Second * 30,
		},
	}
}

func (app *App) StartRoutes(adapters *adapters.Adapters) {
	for _, registerer := range handlers.New(adapters) {
		registerer.Register(app.Router)
	}
}

func (app *App) Run() error {
	return app.Server.ListenAndServe()
}

func (app *App) Shutdown(ctx context.Context) {
	err := app.Server.Shutdown(ctx)
	if err != nil {
		panic(err)
	}
}

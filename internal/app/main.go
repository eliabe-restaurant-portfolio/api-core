package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/eliabe-portfolio/restaurant-app/internal/adapters"
	httpserver "github.com/eliabe-portfolio/restaurant-app/internal/app/http"
	"github.com/eliabe-portfolio/restaurant-app/internal/connections/configs"
	"github.com/eliabe-portfolio/restaurant-app/internal/handlers"
	"github.com/gin-gonic/gin"
)

type App struct {
	Adapters *adapters.Adapters
	Server   *http.Server
	Router   *gin.Engine
}

func New(
	server *httpserver.HttpServer,
	apt *adapters.Adapters,
	conf *configs.Config,
) *App {
	return &App{
		Adapters: apt,
		Router:   server.Router(),
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

func (app *App) RegisterControllers() {
	for _, registerer := range handlers.New(app.Adapters) {
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

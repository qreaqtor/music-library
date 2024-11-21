package app

import "fmt"

import (
	"context"
	"net"

	"github.com/gorilla/mux"
	"github.com/qreaqtor/music-library/internal/config"
	appserver "github.com/qreaqtor/music-library/pkg/appServer"
	httpserver "github.com/qreaqtor/music-library/pkg/httpServer"
)

type server interface {
	Start() error
	Wait() []error
}

type App struct {
	server server
}

func NewApp(ctx context.Context, cfg *config.Config) *App {
	setupLogger(cfg.Env)

	r := mux.NewRouter()

	appServer := appserver.NewAppServer(
		ctx, 
		httpserver.NewHTTPServer(r), 
		net.JoinHostPort(cfg.Host, fmt.Sprint(cfg.Port)),
	)

	return &App{
		server: appServer,
	}
}

func (a *App) Start() error {
	return a.server.Start()
}

func (a *App) Wait() []error {
	return a.server.Wait()
}

package app

import (
	"context"

	"github.com/qreaqtor/music-library/internal/config"
)

type server interface {
	Start() error
	Wait() []error
}

type App struct {
	server server
}

func NewApp(ctx context.Context, cfg *config.Config) *App {
	return &App{
	}
}

func (a *App) Start() error {
	return a.server.Start()
}

func (a *App) Wait() []error {
	return a.server.Wait()
}

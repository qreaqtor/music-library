package app

import (
	"context"
	"fmt"
	"io"
	"net"

	"github.com/gorilla/mux"
	"github.com/qreaqtor/music-library/internal/api"
	"github.com/qreaqtor/music-library/internal/config"
	"github.com/qreaqtor/music-library/internal/service"
	storage "github.com/qreaqtor/music-library/internal/storage/postgres"

	appserver "github.com/qreaqtor/music-library/pkg/appServer"

	httpserver "github.com/qreaqtor/music-library/pkg/httpServer"
)

type server interface {
	Start() error
	Wait() []error
}

type App struct {
	server server

	cfg *config.Config

	router *mux.Router

	toClose []io.Closer
}

func NewApp(ctx context.Context, cfg *config.Config) *App {
	setupLogger(cfg.Env)

	r := mux.NewRouter().PathPrefix(fmt.Sprintf("/v%d", cfg.Api.Version)).Subrouter()

	appServer := appserver.NewAppServer(
		ctx,
		httpserver.NewHTTPServer(r),
		net.JoinHostPort(cfg.Host, fmt.Sprint(cfg.Port)),
	)

	return &App{
		server:  appServer,
		cfg:     cfg,
		router:  r,
		toClose: make([]io.Closer, 0),
	}
}

func (a *App) Start() error {
	conn, err := getPostgresConn(a.cfg.Postgres)
	if err != nil {
		return err
	}

	st := storage.NewSongsStorage(conn)
	srv := service.NewSongsService(st)
	api := api.NewSongsAPI(srv)
	api.Register(a.router)

	a.toClose = append(a.toClose, conn)

	return a.server.Start()
}

func (a *App) Wait() []error {
	errs := a.server.Wait()

	for _, closer := range a.toClose {
		err := closer.Close()
		if err != nil {
			errs = append(errs, err)
		}
	}

	return errs
}

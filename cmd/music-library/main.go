package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/qreaqtor/music-library/internal/app"
	"github.com/qreaqtor/music-library/internal/config"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalln(err)
	}

	app := app.NewApp(ctx, cfg)

	err = app.Start()
	if err != nil {
		log.Fatalln(err)
	}

	errs := app.Wait()
	if len(errs) != 0 {
		log.Fatalln(errs)
	}
}

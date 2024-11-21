package app

import (
	"log/slog"
	"os"

	"github.com/qreaqtor/music-library/pkg/logging/pretty"
)

const (
	local = "local"
	dev   = "dev"
	prod  = "prod"
)

// Return default slog.Logger if env has unsupported value.
// Posible env values: local, dev, prod.
func setupLogger(env string) {
	var handler slog.Handler

	out := os.Stdout

	switch env {
	case local:
		handler = pretty.NewPrettyHandler(
			out,
			&slog.HandlerOptions{
				Level: slog.LevelDebug,
			})
	case dev:
		handler = slog.NewJSONHandler(
			out,
			&slog.HandlerOptions{
				Level: slog.LevelInfo,
			})
	case prod:
		handler = slog.NewJSONHandler(
			out,
			&slog.HandlerOptions{
				Level: slog.LevelError,
			})
	default:
		return
	}

	slog.SetDefault(slog.New(handler))
}

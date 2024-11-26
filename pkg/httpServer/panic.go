package httpserver

import (
	"log/slog"
	"net/http"

	logmsg "github.com/qreaqtor/music-library/pkg/logging/message"
)

func panic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				http.Error(w, "Internal server error", 500)
				slog.Error(
					"panic",
					"err", err,
					"url", r.URL.Path,
					"method", r.Method,
					"operation", logmsg.ExtractOperationID(r.Context()),
				)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

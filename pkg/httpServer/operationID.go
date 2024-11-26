package httpserver

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	logmsg "github.com/qreaqtor/music-library/pkg/logging/message"
)

func setOperationID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := uuid.New()
		ctx := context.WithValue(r.Context(), logmsg.OperationID, id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

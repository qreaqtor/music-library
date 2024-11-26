package httpserver

import (
	"log/slog"
	"net"
	"net/http"
)

type HTTPServer struct {
	server *http.Server
}

// return http server with added recovery middleware and
func NewHTTPServer(handler http.Handler) *HTTPServer {
	server := &HTTPServer{
		server: &http.Server{
			Handler: handler,
		},
	}

	return server
}

// Added middlewares to http.Server handler in LIFO order.
func (h *HTTPServer) AddMiddlewares(middlewares ...Middleware) {
	if len(middlewares) == 0 {
		return
	}

	cur := h.server.Handler

	for _, middleware := range middlewares {
		cur = middleware(cur)
	}

	h.server.Handler = cur
}

func (h *HTTPServer) Serve(l net.Listener) error {
	h.AddMiddlewares(panic, setOperationID)
	slog.Info("Start http server at " + l.Addr().String())
	return h.server.Serve(l)
}

func (h *HTTPServer) Close() error {
	slog.Info("Stop http server")
	return h.server.Close()
}

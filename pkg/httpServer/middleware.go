package httpserver

import "net/http"

type Middleware func(next http.Handler) http.Handler

package web

import (
	"net/http"
)

type HttpMiddleware func(http.Handler) http.Handler

// ChainMiddleware executes middleware in the given order.
func ChainMiddleware(middlewares ...HttpMiddleware) HttpMiddleware {
	return func(handler http.Handler) http.Handler {
		h := handler
		for i := len(middlewares) - 1; i >= 0; i-- {
			h = middlewares[i](h)
		}

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)
		})
	}
}

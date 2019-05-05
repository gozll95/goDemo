package main

import (
	"net/http"

	"golang.org/x/net/context"
)

func RequestIDMiddleware(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		// More correctly, we'd use a const key of type struct{} and a random ID via
		// crypto/rand.
		ctx := context.WithValue(r.Context(), "app.req.id", "12345")

		h.ServeHTTP(w, r.WithContext(ctx))
	}

	return http.HandlerFunc(fn)
}

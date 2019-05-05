package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// e.g. middleware_test.go
func TestPopulateContext(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/users", nil)
	if err != nil {
		t.Fatal(err)
	}

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if val, ok := r.Context().Value("app.req.id").(string); !ok {
			t.Errorf("app.req.id not in request context: got %q", val)
		}
	})

	rr := httptest.NewRecorder()
	// func RequestIDMiddleware(h http.Handler) http.Handler
	// Stores an "app.req.id" in the request context.
	handler := RequestIDMiddleware(testHandler)
	handler.ServeHTTP(rr, req)
}

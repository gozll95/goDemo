package main

import (
	"testing"
)

func TestGetProjectsHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/users", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	// e.g. func GetUsersHandler(ctx context.Context, w http.ResponseWriter, r *http.Request)
	handler := http.HandlerFunc(GetUsersHandler)

	// Populate the request's context with our test data.
	ctx := req.Context()
	ctx = context.WithValue(ctx, "app.auth.token", "abc123")
	ctx = context.WithValue(ctx, "app.user",
		&YourUser{ID: "qejqjq", Email: "user@example.com"})

	// Add our context to the request: note that WithContext returns a copy of
	// the request, which we must assign.
	req = req.WithContext(ctx)
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

//httptest结合context

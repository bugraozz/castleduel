package health

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestReadinessReportsDependencyFailure(t *testing.T) {
	handler := NewHandler(slog.New(slog.NewTextHandler(httptest.NewRecorder(), nil)),
		Check{Name: "database", Run: func(context.Context) error { return errors.New("unavailable") }},
	)
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/health/ready", nil))

	if response.Code != http.StatusServiceUnavailable || response.Body.String() != "{\"status\":\"not_ready\"}\n" {
		t.Fatalf("response = %d %q", response.Code, response.Body.String())
	}
}

func TestLivenessDoesNotRunDependencyChecks(t *testing.T) {
	called := false
	handler := NewHandler(nil, Check{Name: "database", Run: func(context.Context) error {
		called = true
		return errors.New("unavailable")
	}})
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/health/live", nil))

	if response.Code != http.StatusOK || called {
		t.Fatalf("response = %d, dependency check called = %t", response.Code, called)
	}
}

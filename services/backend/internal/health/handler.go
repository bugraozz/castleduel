package health

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"
)

type Check struct {
	Name string
	Run  func(context.Context) error
}

type response struct {
	Status string `json:"status"`
}

func NewHandler(logger *slog.Logger, checks ...Check) http.Handler {
	if logger == nil {
		logger = slog.Default()
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health/live", func(w http.ResponseWriter, _ *http.Request) {
		writeResponse(logger, w, http.StatusOK, response{Status: "ok"})
	})
	mux.HandleFunc("GET /health/ready", func(w http.ResponseWriter, r *http.Request) {
		for _, check := range checks {
			ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
			err := check.Run(ctx)
			cancel()
			if err != nil {
				logger.WarnContext(r.Context(), "readiness check failed", "check", check.Name, "error", err)
				writeResponse(logger, w, http.StatusServiceUnavailable, response{Status: "not_ready"})
				return
			}
		}
		writeResponse(logger, w, http.StatusOK, response{Status: "ok"})
	})
	return mux
}

func writeResponse(logger *slog.Logger, w http.ResponseWriter, status int, body response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(body); err != nil {
		logger.Error("failed to write health response", "error", err)
	}
}

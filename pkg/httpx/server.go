// Package httpx is the HTTP chassis every service shares: liveness and
// readiness endpoints, graceful shutdown, and an instrumented client.
package httpx

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"sync"
	"time"
)

// Check reports whether one dependency is ready.
type Check func(ctx context.Context) error

// Health serves /healthz (process alive) and /readyz (dependencies ready).
type Health struct {
	mu     sync.RWMutex
	checks map[string]Check
}

// NewHealth returns a Health with no readiness checks.
func NewHealth() *Health { return &Health{checks: map[string]Check{}} }

// AddReadiness registers a named readiness check.
func (h *Health) AddReadiness(name string, c Check) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.checks[name] = c
}

// Mount registers the two endpoints on mux.
func (h *Health) Mount(mux *http.ServeMux) {
	mux.HandleFunc("GET /healthz", h.Live)
	mux.HandleFunc("GET /readyz", h.Ready)
}

// Live answers the liveness probe: the process is up.
func (h *Health) Live(w http.ResponseWriter, _ *http.Request) {
	WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// Ready answers the readiness probe: every registered dependency answers.
func (h *Health) Ready(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()
	h.mu.RLock()
	defer h.mu.RUnlock()
	failed := map[string]string{}
	for name, c := range h.checks {
		if err := c(ctx); err != nil {
			failed[name] = err.Error()
		}
	}
	if len(failed) > 0 {
		WriteJSON(w, http.StatusServiceUnavailable, map[string]any{"status": "unavailable", "failed": failed})
		return
	}
	WriteJSON(w, http.StatusOK, map[string]string{"status": "ready"})
}

// WriteJSON writes v as JSON with status.
func WriteJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

// Run serves srv until ctx ends, then shuts down within 15 seconds.
func Run(ctx context.Context, srv *http.Server) error {
	errc := make(chan error, 1)
	go func() {
		slog.Info("listening", "addr", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errc <- err
		}
		close(errc)
	}()
	select {
	case err := <-errc:
		return err
	case <-ctx.Done():
	}
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	return srv.Shutdown(shutdownCtx)
}

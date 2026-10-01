package httpapi

import (
	"log/slog"
	"net/http"
)

// ready answers GET /ready.
// It pings Redis. A failed ping answers 503. /health does not ping Redis.
func (s *Server) ready(w http.ResponseWriter, r *http.Request) {
	if err := s.svc.Ready(r.Context()); err != nil {
		slog.Error("ready", "err", err)
		writeJSON(w, http.StatusServiceUnavailable, map[string]string{"status": "not ready"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

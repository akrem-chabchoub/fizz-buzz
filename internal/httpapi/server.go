package httpapi

import (
	"log/slog"
	"net/http"
	"time"
)

type Server struct {
	maxLimit     int
	maxStringLen int
}

// New keeps the caps used by GET /v1/fizzbuzz.
func New(maxLimit, maxStringLen int) *Server {
	return &Server{maxLimit: maxLimit, maxStringLen: maxStringLen}
}

func (s *Server) Handler() http.Handler {
	return withLog(s.routes())
}

func withLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ww := &statusWriter{ResponseWriter: w, status: http.StatusOK}
		defer func() {
			if rec := recover(); rec != nil {
				slog.Error("panic", "err", rec)
				if ww.status == http.StatusOK {
					writeError(ww, http.StatusInternalServerError, "internal error")
				}
			}
			logRequest(ww.status, "request",
				"method", r.Method,
				"path", r.URL.Path,
				"status", ww.status,
				"ms", time.Since(start).Milliseconds(),
			)
		}()
		next.ServeHTTP(ww, r)
	})
}

type statusWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

func logRequest(status int, msg string, args ...any) {
	switch {
	case status >= http.StatusInternalServerError:
		slog.Error(msg, args...)
	case status >= http.StatusBadRequest && status < http.StatusInternalServerError:
		slog.Warn(msg, args...)
	default:
		slog.Info(msg, args...)
	}
}

func (s *Server) health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

package httpapi

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/akrem-chabchoub/fizz-buzz/internal/app"
)

type Server struct {
	maxLimit     int
	maxStringLen int
	svc          *app.Service
}

// New keeps the caps used by GET /v1/fizzbuzz.
// svc records successful calls, answers GET /v1/stats and checks the store for GET /ready.
func New(svc *app.Service, maxLimit, maxStringLen int) *Server {
	return &Server{maxLimit: maxLimit, maxStringLen: maxStringLen, svc: svc}
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
				if rec == http.ErrAbortHandler {
					panic(rec)
				}
				slog.Error("panic", "err", rec)
				if !ww.wroteHeader {
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
	status      int
	wroteHeader bool
}

func (w *statusWriter) WriteHeader(code int) {
	if !w.wroteHeader {
		w.status = code
		w.wroteHeader = true
	}
	w.ResponseWriter.WriteHeader(code)
}

func (w *statusWriter) Write(b []byte) (int, error) {
	w.wroteHeader = true
	return w.ResponseWriter.Write(b)
}

func logRequest(status int, msg string, args ...any) {
	switch {
	case status >= http.StatusInternalServerError:
		slog.Error(msg, args...)
	case status >= http.StatusBadRequest:
		slog.Warn(msg, args...)
	default:
		slog.Info(msg, args...)
	}
}

func (s *Server) health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

package httpapi

import "net/http"

func (s *Server) routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", s.health)
	mux.HandleFunc("GET /ready", s.ready)
	mux.HandleFunc("GET /v1/fizzbuzz", s.fizzbuzz)
	mux.HandleFunc("GET /v1/stats", s.stats)
	mux.HandleFunc("GET /openapi.yaml", s.openapi)
	mux.Handle("/swagger/", swaggerUI())
	return mux
}

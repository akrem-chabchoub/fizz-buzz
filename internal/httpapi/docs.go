package httpapi

import (
	"net/http"

	"github.com/swaggest/swgui/v5emb"

	"github.com/akrem-chabchoub/fizz-buzz/api"
)

// openapi serves the API description at GET /openapi.yaml.
func (s *Server) openapi(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/yaml")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(api.Spec)
}

func swaggerUI() http.Handler {
	return v5emb.New("fizzbuzz", "/openapi.yaml", "/swagger/")
}

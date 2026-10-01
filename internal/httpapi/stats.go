package httpapi

import (
	"log/slog"
	"net/http"
)

// stats answers GET /v1/stats.
// No successful calls yet returns {"hits":0}.
// Only a GET fizzbuzz call that answers 200 is counted.
func (s *Server) stats(w http.ResponseWriter, r *http.Request) {
	req, hits, ok, err := s.svc.Top(r.Context())
	if err != nil {
		slog.Error("stats top failed", "err", err)
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}
	if !ok {
		writeJSON(w, http.StatusOK, map[string]uint64{"hits": 0})
		return
	}
	writeJSON(w, http.StatusOK, statsBody{
		Int1:  req.Int1,
		Int2:  req.Int2,
		Limit: req.Limit,
		Str1:  req.Str1,
		Str2:  req.Str2,
		Hits:  hits,
	})
}

type statsBody struct {
	Int1  int    `json:"int1"`
	Int2  int    `json:"int2"`
	Limit int    `json:"limit"`
	Str1  string `json:"str1"`
	Str2  string `json:"str2"`
	Hits  uint64 `json:"hits"`
}

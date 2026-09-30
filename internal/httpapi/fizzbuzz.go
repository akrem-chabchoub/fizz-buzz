package httpapi

import (
	"net/http"

	"github.com/akrem-chabchoub/fizz-buzz/internal/fizzbuzz"
	"github.com/akrem-chabchoub/fizz-buzz/internal/utils"
)

// fizzbuzz answers GET /v1/fizzbuzz.
// It counts from 1 to limit. Multiples of int1 become str1, multiples of int2
// become str2, and multiples of both become the two words stuck together.
//
//	/v1/fizzbuzz?int1=3&int2=5&limit=15&str1=fizz&str2=buzz
//
// returns ["1","2","fizz","4","buzz","fizz","7","8","fizz","buzz","11","fizz","13","14","fizzbuzz"].
// A bad query answers 400, for example int1=0 returns {"error":"int1 must be greater than 0"}.
// An empty word is allowed. When int1 equals int2, every multiple takes the both path.
func (s *Server) fizzbuzz(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	int1, err := utils.ParsePositive(q.Get("int1"), "int1", q.Has("int1"))
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	int2, err := utils.ParsePositive(q.Get("int2"), "int2", q.Has("int2"))
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	limit, err := utils.ParsePositive(q.Get("limit"), "limit", q.Has("limit"))
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if limit > s.maxLimit {
		writeError(w, http.StatusBadRequest, "limit is too big")
		return
	}
	str1, err := utils.ParseWord(q, "str1", s.maxStringLen)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	str2, err := utils.ParseWord(q, "str2", s.maxStringLen)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, fizzbuzz.Generate(int1, int2, limit, str1, str2))
}

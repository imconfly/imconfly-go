package handler

import (
	"net/http"
)

type statusResponseWriter struct {
	http.ResponseWriter
	status int
}

func newSRW(w http.ResponseWriter) *statusResponseWriter {
	return &statusResponseWriter{w, 0}
}

func (s *statusResponseWriter) WriteHeader(code int) {
	s.status = code
	s.ResponseWriter.WriteHeader(code)
}

package server

import (
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type statusResponseWriter struct {
	http.ResponseWriter
	status int
}

func (s *statusResponseWriter) WriteHeader(code int) {
	s.status = code
	s.ResponseWriter.WriteHeader(code)
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			startTime := time.Now()
			sw := &statusResponseWriter{w, 200}
			next.ServeHTTP(sw, r)
			duration := time.Since(startTime)
			if sw.status == http.StatusOK {
				log.Info(sw.status, "\t", duration, "\t", r.RequestURI)
			} else {
				log.Error(sw.status, "\t", duration, "\t", r.RequestURI)
			}
		},
	)
}

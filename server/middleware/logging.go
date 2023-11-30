package middleware

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

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			startTime := time.Now()
			sw := &statusResponseWriter{w, 0}
			next.ServeHTTP(sw, r)
			duration := time.Since(startTime)
			if sw.status == http.StatusOK {
				log.Infof(
					"%d\t%s\t%s",
					sw.status,
					duration,
					r.RequestURI)
			} else {
				log.Errorf(
					"%d %s\t%s\t%s",
					sw.status,
					http.StatusText(sw.status),
					duration,
					r.RequestURI)
			}
		},
	)
}

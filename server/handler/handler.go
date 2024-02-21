package handler

import (
	"errors"
	"fmt"
	"github.com/imconfly/imconfly_go/lib/os_tools"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/imconfly/imconfly_go/core/resolver"
	rsErrors "github.com/imconfly/imconfly_go/core/resolver/resolver_errors"
)

type Handler struct {
	resolver *resolver.Resolver
}

func NewHandler(rs *resolver.Resolver) http.Handler {
	return &Handler{
		resolver: rs,
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	status := http.StatusOK
	var err error
	logName := fmt.Sprintf("Hanler.ServeHTTP(%s)", r.RequestURI)
	log.Debug(logName)

	switch r.RequestURI {
	// for clean logs in case of testing from browser
	case "/favicon.ico":
		http.NotFound(w, r)
		status = http.StatusNotFound
	// simple server health check
	case "/health":
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Write([]byte("Ok"))
	// main action
	default:
		var fileAbsPath os_tools.FileAbsPath
		fileAbsPath, err = h.resolver.Request(r.RequestURI)
		if err != nil {
			status = http.StatusInternalServerError
			var rErr *rsErrors.Error
			if errors.As(err, &rErr) {
				status = rErr.HTTPCode
			}
			http.Error(w, err.Error(), status)
		} else {
			// we want to know status
			srw := newSRW(w)
			http.ServeFile(srw, r, string(fileAbsPath))
			status = srw.status
		}
	}

	// logging

	duration := time.Since(startTime)
	if err != nil {
		log.Errorf(
			"%s: %d %s\t%s\t%s",
			logName,
			status,
			http.StatusText(status),
			duration,
			err.Error())
	} else {
		log.Infof(
			"%s: %d\t%s\t",
			logName,
			status,
			duration)
	}
}

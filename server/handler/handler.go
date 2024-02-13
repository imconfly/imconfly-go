package handler

import (
	"errors"
	"fmt"
	"github.com/imconfly/imconfly_go/core/resolver"
	rsErrors "github.com/imconfly/imconfly_go/core/resolver/resolver_errors"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
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
	logName := fmt.Sprintf("Hanler.ServeHTTP(%s)", r.RequestURI)
	log.Debug(logName)

	fileAbsPath, err := h.resolver.Request(r.RequestURI)
	status := http.StatusOK
	if err != nil {
		status = http.StatusInternalServerError
		var rErr *rsErrors.Error
		if errors.As(err, &rErr) {
			status = rErr.HTTPCode
		}
		http.Error(w, err.Error(), status)
	} else {
		srw := newSRW(w)
		http.ServeFile(srw, r, string(fileAbsPath))
		status = srw.status
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

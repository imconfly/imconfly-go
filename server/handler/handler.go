package handler

import (
	"errors"
	"github.com/imconfly/imconfly_go/core/resolver"
	rsErrors "github.com/imconfly/imconfly_go/core/resolver/resolver_errors"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type Handler struct {
	resolver *resolver.Resolver
	logger   *logrus.Logger
}

func NewHandler(rs *resolver.Resolver, logger *logrus.Logger) http.Handler {
	return &Handler{
		resolver: rs,
		logger:   logger,
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
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
		h.logger.Errorf(
			"%d %s\t%s\t%s\t%s",
			status,
			http.StatusText(status),
			duration,
			r.RequestURI,
			err.Error())
	} else {
		h.logger.Infof(
			"%d\t%s\t%s",
			status,
			duration,
			r.RequestURI)
	}
}

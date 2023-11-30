package handler

import (
	"errors"
	"github.com/imconfly/imconfly_go/core/resolver"
	errors2 "github.com/imconfly/imconfly_go/core/resolver/errors"
	"net/http"
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
	fileAbsPath, err := h.resolver.Request(r.RequestURI)
	if err != nil {
		status := http.StatusInternalServerError
		var rErr *errors2.ResolverError
		if errors.As(err, &rErr) {
			status = rErr.HTTPCode
		}
		http.Error(w, err.Error(), status)
	} else {
		http.ServeFile(w, r, string(fileAbsPath))
	}
}

package handler

import (
	"github.com/imconfly/imconfly_go/core/resolver"
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		http.ServeFile(w, r, string(fileAbsPath))
	}
}

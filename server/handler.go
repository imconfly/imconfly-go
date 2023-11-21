package server

import (
	"github.com/imconfly/imconfly_go/core/resolver"
	"net/http"
)

type Handler struct {
	Resolver *resolver.Resolver
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fileAbsPath, err := h.Resolver.Request(r.RequestURI)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.ServeFile(w, r, string(fileAbsPath))
}

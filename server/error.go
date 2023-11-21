package server

import (
	"fmt"
	"net/http"
)

type Error struct {
	HTTPCode int
	Err      error
}

func (e *Error) Error() string {
	return fmt.Sprintf("[%d: %s] %s", e.HTTPCode, http.StatusText(e.HTTPCode), e.Err.Error())
}

func (e *Error) WriteToResponse(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(e.HTTPCode)
	fmt.Fprintln(w, e.Error())
}

package handler

import (
	"net/http"

	"github.com/a-h/templ"
)

func Render(w http.ResponseWriter, r *http.Request, status int, cmp templ.Component) error {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	if err := cmp.Render(r.Context(), w); err != nil {
		return err
	}

	return nil
}

package server

import (
	"net/http"
	"github.com/go-chi/render"
	"github.com/go-chi/chi"
)

func JSON(w http.ResponseWriter, r *http.Request, v interface{}) {
	render.JSON(w, r, v)
}

func URLParam(r *http.Request, p string) string {
	return chi.URLParam(r, p)
}

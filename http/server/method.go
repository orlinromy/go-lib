package server

import (
	"net/http"
)

func (rtr Router) Get(route string, handler http.HandlerFunc) {
	rtr.Engine.Get(route, handler)
}

func (rtr Router) Patch(route string, handler http.HandlerFunc) {
	rtr.Engine.Patch(route, handler)
}

func (rtr Router) Put(route string, handler http.HandlerFunc) {
	rtr.Engine.Put(route, handler)
}

func (rtr Router) Post(route string, handler http.HandlerFunc) {
	rtr.Engine.Post(route, handler)
}

func (rtr Router) Delete(route string, handler http.HandlerFunc) {
	rtr.Engine.Delete(route, handler)
}

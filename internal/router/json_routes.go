package router

import (
	"net/http"

	"github.com/DilankaHer/sop-in-go/internal/middleware"
	"github.com/go-chi/chi/v5"
)

// JSONRoutes registers handlers that return StandardResponse with Response applied automatically.
type JSONRoutes struct {
	r  chi.Router
	mw middleware.Middleware
}

func NewJSONRoutes(r chi.Router, mw middleware.Middleware) JSONRoutes {
	return JSONRoutes{r: r, mw: mw}
}

func (jr JSONRoutes) Get(pattern string, h func(r *http.Request) middleware.StandardResponse) {
	jr.r.Get(pattern, jr.mw.Response(h))
}

func (jr JSONRoutes) Post(pattern string, h func(r *http.Request) middleware.StandardResponse) {
	jr.r.Post(pattern, jr.mw.Response(h))
}

func (jr JSONRoutes) Put(pattern string, h func(r *http.Request) middleware.StandardResponse) {
	jr.r.Put(pattern, jr.mw.Response(h))
}

func (jr JSONRoutes) Patch(pattern string, h func(r *http.Request) middleware.StandardResponse) {
	jr.r.Patch(pattern, jr.mw.Response(h))
}

func (jr JSONRoutes) Delete(pattern string, h func(r *http.Request) middleware.StandardResponse) {
	jr.r.Delete(pattern, jr.mw.Response(h))
}

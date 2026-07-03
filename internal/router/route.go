package router

import (
	"github.com/DilankaHer/sop-in-go/internal/app"
	"github.com/DilankaHer/sop-in-go/internal/handler"
	"github.com/go-chi/chi/v5"
)

func InitRoutes(app *app.App) *chi.Mux {
	r := chi.NewRouter()

	handler := handler.NewHandler(app)

	r.Use(app.Middleware.Recover)
	r.Use(app.Middleware.AccessLog)

	api := NewJSONRoutes(r, app.Middleware)
	api.Get("/version", handler.ServerHandler.Version)

	return r
}

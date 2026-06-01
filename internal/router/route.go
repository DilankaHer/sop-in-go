package router

import (
	"github.com/DilankaHer/sop-in-go/internal/app"
	"github.com/DilankaHer/sop-in-go/internal/handler"
	"github.com/go-chi/chi/v5"
)

func InitRoutes(app *app.App) *chi.Mux {
	r := chi.NewRouter()

	handler := handler.NewHandler(app)

	r.Get("/version", handler.ServerHandler.Version)
	return r
}

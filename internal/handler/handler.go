package handler

import (
	"github.com/DilankaHer/sop-in-go/internal/app"
	"github.com/DilankaHer/sop-in-go/internal/handler/server"
)

type Handler struct {
	ServerHandler *server.ServerHandler
}

func NewHandler(app *app.App) *Handler {
	serverHandler := server.NewServerHandler(app)
	return &Handler{ServerHandler: serverHandler}
}

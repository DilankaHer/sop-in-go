package server

import (
	"net/http"

	"github.com/DilankaHer/sop-in-go/internal/app"
	"github.com/DilankaHer/sop-in-go/internal/middleware"
)

type ServerHandler struct {
	app *app.App
}

type getVersionResp struct {
	Version string `json:"version"`
}

func NewServerHandler(app *app.App) *ServerHandler {
	return &ServerHandler{
		app: app,
	}
}

func (h *ServerHandler) Version(r *http.Request) (resp middleware.StandardResponse) {
	resp.Data = getVersionResp{Version: "v1.1"}
	resp.Status = http.StatusOK
	resp.Message = http.StatusText(resp.Status)
	return
}

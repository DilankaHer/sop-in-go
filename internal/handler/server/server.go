package server

import (
	"net/http"

	"github.com/DilankaHer/sop-in-go/internal/app"
	"github.com/DilankaHer/sop-in-go/internal/utility"
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

func (h *ServerHandler) Version(w http.ResponseWriter, r *http.Request) {
	utility.WriteJSON(w, http.StatusOK, "success", getVersionResp{Version: "v1.1"})
}

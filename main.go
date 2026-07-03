package main

import (
	"log"
	"net/http"

	"github.com/DilankaHer/sop-in-go/internal/app"
	"github.com/DilankaHer/sop-in-go/internal/router"
)

func main() {
	app, err := app.NewApp()
	if err != nil {
		log.Fatalf("error during startup: \n%s", err.Error())
	}

	r := router.InitRoutes(app)

	app.Logger.Debug("Server running on port", app.Config.Server.Port)
	err = http.ListenAndServe(":"+app.Config.Server.Port, r)
	if err != nil {
		log.Fatalf("error during serving server: %s", err.Error())
	}
}

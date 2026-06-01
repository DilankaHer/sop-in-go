package main

import (
	"fmt"
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

	fmt.Println("Server running on port:", app.Config.Port)
	err = http.ListenAndServe(":"+app.Config.Port, r)
	if err != nil {
		log.Fatalf("error during serving server: %s", err.Error())
	}
}

package main

import (
	"log"

	"github.com/adexcell/delayed-notifier/cmd/app"
)

// @title          Delayed Notifier API
// @version        1.0
// @description    Delayed Notifier
// @host           localhost:8080
// @BasePath       /

func main() {
	app, err := app.New()
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	if err := app.Run(); err != nil {
		log.Fatalf("error: %v", err)
	}
}

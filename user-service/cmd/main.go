package main

import (
	"context"
	"github.com/19parwiz/user/config"
	"github.com/19parwiz/user/internal/app"
	"log"
)

func main() {
	ctx := context.Background()

	cfg, err := config.New()
	if err != nil {
		log.Println("error loading config", err)
		return
	}

	app, err := app.NewApp(ctx, cfg)
	if err != nil {
		log.Println("error initializing App", err)
		return
	}

	err = app.Start()

	if err != nil {
		log.Println("error starting App", err)
		return
	}

}

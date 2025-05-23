package main

import (
	"log"

	"github.com/Amytlisss/GO_old/priyutik/internal/config"

	"github.com/Amytlisss/GO_old/priyutik/internal/app"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	application, err := app.New(cfg)
	if err != nil {
		log.Fatalf("Failed to create app: %v", err)
	}

	if err := application.Run(); err != nil {
		log.Fatalf("Application error: %v", err)
	}
}

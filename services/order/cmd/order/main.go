package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/BeInBloom/big-deal/services/order/internal/app"
	"github.com/BeInBloom/big-deal/services/order/internal/config"
)

func main() {
	if err := run(); err != nil {
		log.Printf("order service failed: %s", err)
		os.Exit(1)
	}

	log.Println("server stopped")
}

func run() error {
	cfg, err := config.LoadFromEnv()
	if err != nil {
		return fmt.Errorf("wrong config: %w", err)
	}

	application, err := app.New(cfg)
	if err != nil {
		return fmt.Errorf("build app: %w", err)
	}

	return application.Run(context.Background())
}

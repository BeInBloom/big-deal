package main

import (
	"context"
	"log"

	"github.com/BeInBloom/big-deal/services/order/internal/app"
	"github.com/BeInBloom/big-deal/services/order/internal/config"
)

func main() {
	cfg, err := config.LoadFromEnv()
	if err != nil {
		log.Fatalf("wrong config: %s", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	app := app.New(cfg)
	if err := app.Run(ctx); err != nil {
		cancel()
		log.Fatalf("server stopped with error: %s", err)
	}

	log.Println("server stopped")
}

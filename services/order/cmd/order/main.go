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
		log.Println("wrong config: %s", err.Error())
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	app := app.New(cfg)
	if err := app.Run(ctx); err != nil {
		log.Println("something wrong")
	}

	log.Println("server stopped")
}

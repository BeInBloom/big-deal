package main

import (
	"context"
	"log"

	"github.com/BeInBloom/big-deal/services/order/internal/app"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	app := app.New()
	if err := app.Run(ctx); err != nil {
		log.Println("something wrong")
	}

	log.Println("server stopped")
}

package main

import (
	"context"
	"log"

	"github.com/google/uuid"

	orders "github.com/BeInBloom/big-deal/generated/go/order-service/openapi/v1"
	"github.com/BeInBloom/big-deal/services/order/internal/app"
)

type Handler struct {
	orders.UnimplementedHandler
}

func (h Handler) OrderServiceCreateOrder(
	ctx context.Context,
	req *orders.CreateOrderRequest,
) (orders.OrderServiceCreateOrderRes, error) {
	return &orders.CreateOrderResponse{
		OrderUUID:  uuid.New(),
		TotalPrice: 1234.56,
	}, nil
}

func (h Handler) OrderServiceGetOrder(
	ctx context.Context,
	params orders.OrderServiceGetOrderParams,
) (orders.OrderServiceGetOrderRes, error) {
	return &orders.Order{
		OrderUUID:  params.OrderUUID,
		UserUUID:   uuid.New(),
		PartUuids:  []uuid.UUID{uuid.New()},
		TotalPrice: 1234.56,
		Status:     orders.OrderStatusPENDINGPAYMENT,
	}, nil
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	app := app.New()
	if err := app.Run(ctx); err != nil {
		log.Fatalln("something wrong")
	}

	log.Println("server stopped")
}

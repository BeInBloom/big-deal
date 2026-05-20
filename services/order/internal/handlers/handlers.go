package handlers

import (
	"context"

	orders "github.com/BeInBloom/big-deal/generated/go/order-service/openapi/v1"
	"github.com/BeInBloom/big-deal/services/order/internal/models"
)

type orderService interface {
	CancelOrder(context.Context, models.OrderId) error
	CreateOrder(context.Context, models.UserId, []models.PartId) (models.PendingOrder, error)
	GetOrder(context.Context, models.OrderId) (models.Order, error)
	PayOrder(context.Context, models.OrderId, models.PaymentMethod) (models.PaidOrder, error)
}

type Handlers struct {
	orderService orderService
}

func New() *Handlers {
	return &Handlers{}
}

func (h *Handlers) OrderServiceCancelOrder(
	ctx context.Context,
	params orders.OrderServiceCancelOrderParams,
) (orders.OrderServiceCancelOrderRes, error) {
	panic("implement me")
}

func (h *Handlers) OrderServiceCreateOrder(
	ctx context.Context,
	req *orders.CreateOrderRequest,
) (orders.OrderServiceCreateOrderRes, error) {
	panic("implement me")
}

func (h *Handlers) OrderServiceGetOrder(
	ctx context.Context,
	params orders.OrderServiceGetOrderParams,
) (orders.OrderServiceGetOrderRes, error) {
	panic("implement me")
}

func (h *Handlers) OrderServicePayOrder(
	ctx context.Context,
	req *orders.PayOrderRequest,
	params orders.OrderServicePayOrderParams,
) (orders.OrderServicePayOrderRes, error) {
	panic("implement me")
}

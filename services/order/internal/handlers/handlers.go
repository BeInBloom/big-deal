package handlers

import (
	"context"

	orders "github.com/BeInBloom/big-deal/generated/go/order-service/openapi/v1"
	"github.com/BeInBloom/big-deal/services/order/internal/models"
	"github.com/google/uuid"
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

func New(orderService orderService) *Handlers {
	return &Handlers{
		orderService: orderService,
	}
}

func (h *Handlers) OrderServiceCancelOrder(
	ctx context.Context,
	params orders.OrderServiceCancelOrderParams,
) (orders.OrderServiceCancelOrderRes, error) {
	if err := h.orderService.CancelOrder(ctx, models.OrderId(params.OrderUUID)); err != nil {
		return mapCancelOrderError(err), nil
	}

	return &orders.OrderServiceCancelOrderNoContent{}, nil
}

func (h *Handlers) OrderServiceCreateOrder(
	ctx context.Context,
	req *orders.CreateOrderRequest,
) (orders.OrderServiceCreateOrderRes, error) {
	order, err := h.orderService.CreateOrder(
		ctx,
		models.UserId(req.UserUUID),
		toPartIds(req.PartUuids),
	)
	if err != nil {
		return mapCreateOrderError(err), nil
	}

	return &orders.CreateOrderResponse{
		OrderUUID:  uuid.UUID(order.Id()),
		TotalPrice: toAPIPrice(order.Price()),
	}, nil
}

func (h *Handlers) OrderServiceGetOrder(
	ctx context.Context,
	params orders.OrderServiceGetOrderParams,
) (orders.OrderServiceGetOrderRes, error) {
	order, err := h.orderService.GetOrder(ctx, models.OrderId(params.OrderUUID))
	if err != nil {
		return mapGetOrderError(err), nil
	}

	return toAPIOrder(order), nil
}

func (h *Handlers) OrderServicePayOrder(
	ctx context.Context,
	req *orders.PayOrderRequest,
	params orders.OrderServicePayOrderParams,
) (orders.OrderServicePayOrderRes, error) {
	order, err := h.orderService.PayOrder(
		ctx,
		models.OrderId(params.OrderUUID),
		toPaymentMethod(req.PaymentMethod),
	)
	if err != nil {
		return mapPayOrderError(err), nil
	}

	return &orders.PayOrderResponse{
		TransactionUUID: uuid.UUID(order.TransactionId()),
	}, nil
}

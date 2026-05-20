package handlers

import (
	"errors"

	orders "github.com/BeInBloom/big-deal/generated/go/order-service/openapi/v1"
	"github.com/BeInBloom/big-deal/services/order/internal/services"
)

func mapCancelOrderError(err error) orders.OrderServiceCancelOrderRes {
	switch {
	case errors.Is(err, services.ErrOrderNotFound):
		res := orders.OrderServiceCancelOrderNotFound(newErrorResponse("ORDER_NOT_FOUND", "order not found"))
		return &res
	case errors.Is(err, services.ErrOrderCannotBeCanceled):
		res := orders.OrderServiceCancelOrderConflict(newErrorResponse("ORDER_CANNOT_BE_CANCELED", "order cannot be canceled"))
		return &res
	default:
		return newCancelOrderInternalServerError()
	}
}

func mapCreateOrderError(err error) orders.OrderServiceCreateOrderRes {
	switch {
	case errors.Is(err, services.ErrOrderPartsNotFound):
		res := orders.OrderServiceCreateOrderBadRequest(newErrorResponse("ORDER_PARTS_NOT_FOUND", "order parts not found"))
		return &res
	default:
		return newCreateOrderInternalServerError()
	}
}

func mapGetOrderError(err error) orders.OrderServiceGetOrderRes {
	switch {
	case errors.Is(err, services.ErrOrderNotFound):
		res := orders.OrderServiceGetOrderNotFound(newErrorResponse("ORDER_NOT_FOUND", "order not found"))
		return &res
	default:
		return newGetOrderInternalServerError()
	}
}

func mapPayOrderError(err error) orders.OrderServicePayOrderRes {
	switch {
	case errors.Is(err, services.ErrOrderNotFound):
		res := orders.OrderServicePayOrderNotFound(newErrorResponse("ORDER_NOT_FOUND", "order not found"))
		return &res
	case errors.Is(err, services.ErrOrderCannotBePaid):
		res := orders.OrderServicePayOrderConflict(newErrorResponse("ORDER_CANNOT_BE_PAID", "order cannot be paid"))
		return &res
	default:
		return newPayOrderInternalServerError()
	}
}

func newErrorResponse(code, message string) orders.ErrorResponse {
	return orders.ErrorResponse{
		Code:    code,
		Message: message,
	}
}

func newInternalServerErrorResponse() orders.ErrorResponse {
	return newErrorResponse("INTERNAL_ERROR", "internal server error")
}

func newCancelOrderInternalServerError() *orders.OrderServiceCancelOrderInternalServerError {
	res := orders.OrderServiceCancelOrderInternalServerError(newInternalServerErrorResponse())
	return &res
}

func newCreateOrderInternalServerError() *orders.OrderServiceCreateOrderInternalServerError {
	res := orders.OrderServiceCreateOrderInternalServerError(newInternalServerErrorResponse())
	return &res
}

func newGetOrderInternalServerError() *orders.OrderServiceGetOrderInternalServerError {
	res := orders.OrderServiceGetOrderInternalServerError(newInternalServerErrorResponse())
	return &res
}

func newPayOrderInternalServerError() *orders.OrderServicePayOrderInternalServerError {
	res := orders.OrderServicePayOrderInternalServerError(newInternalServerErrorResponse())
	return &res
}

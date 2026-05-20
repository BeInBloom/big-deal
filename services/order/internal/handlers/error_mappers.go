package handlers

import orders "github.com/BeInBloom/big-deal/generated/go/order-service/openapi/v1"

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

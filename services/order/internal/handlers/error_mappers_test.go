package handlers

import (
	"testing"

	orders "github.com/BeInBloom/big-deal/generated/go/order-service/openapi/v1"
	"github.com/BeInBloom/big-deal/services/order/internal/services"
	"github.com/stretchr/testify/require"
)

func TestMapCreateOrderErrorReturnsBadRequestForInvalidParts(t *testing.T) {
	tests := []struct {
		name string
		err  error
		code string
	}{
		{name: "required", err: services.ErrOrderPartsRequired, code: "ORDER_PARTS_REQUIRED"},
		{name: "duplicated", err: services.ErrOrderPartsDuplicated, code: "ORDER_PARTS_DUPLICATED"},
		{name: "not found", err: services.ErrOrderPartsNotFound, code: "ORDER_PARTS_NOT_FOUND"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, ok := mapCreateOrderError(tt.err).(*orders.OrderServiceCreateOrderBadRequest)
			require.True(t, ok)
			require.Equal(t, tt.code, res.Code)
		})
	}
}

func TestMapPayOrderErrorReturnsBadRequestForInvalidPaymentMethod(t *testing.T) {
	res, ok := mapPayOrderError(services.ErrInvalidPaymentMethod).(*orders.OrderServicePayOrderBadRequest)

	require.True(t, ok)
	require.Equal(t, "INVALID_PAYMENT_METHOD", res.Code)
}

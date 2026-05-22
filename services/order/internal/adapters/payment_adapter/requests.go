package paymentadapter

import (
	payment "github.com/BeInBloom/big-deal/generated/go/payment/v1"
	"github.com/BeInBloom/big-deal/services/order/internal/models"
	"github.com/google/uuid"
)

func buildReq(
	userId models.UserId,
	orderId models.OrderId,
	method models.PaymentMethod,
) *payment.PayOrderRequest {
	return &payment.PayOrderRequest{
		OrderUuid:     uuid.UUID(orderId).String(),
		UserUuid:      uuid.UUID(userId).String(),
		PaymentMethod: paymentMethodToPaymentMethod(method),
	}
}

package paymentadapter

import (
	"testing"

	payment "github.com/BeInBloom/big-deal/generated/go/payment/v1"
	"github.com/BeInBloom/big-deal/services/order/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestBuildReq(t *testing.T) {
	userUUID := uuid.New()
	orderUUID := uuid.New()

	req := buildReq(
		models.UserId(userUUID),
		models.OrderId(orderUUID),
		models.PaymentMethodCARD,
	)

	require.Equal(t, userUUID.String(), req.GetUserUuid())
	require.Equal(t, orderUUID.String(), req.GetOrderUuid())
	require.Equal(t, payment.PaymentMethod_PAYMENT_METHOD_CARD, req.GetPaymentMethod())
}

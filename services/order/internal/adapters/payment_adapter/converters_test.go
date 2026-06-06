package paymentadapter

import (
	"testing"

	payment "github.com/BeInBloom/big-deal/generated/go/payment/v1"
	"github.com/BeInBloom/big-deal/services/order/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestRawTransactionIdToTransactionId(t *testing.T) {
	transactionUUID := uuid.New()

	transactionId, err := rawTransactionIdToTransactionId(transactionUUID.String())

	require.NoError(t, err)
	require.Equal(t, models.TransactionId(transactionUUID), transactionId)
}

func TestRawTransactionIdToTransactionIdReturnsError(t *testing.T) {
	transactionId, err := rawTransactionIdToTransactionId("not-a-uuid")

	require.Error(t, err)
	require.Zero(t, transactionId)
}

func TestPaymentMethodToPaymentMethod(t *testing.T) {
	tests := []struct {
		name string
		in   models.PaymentMethod
		want payment.PaymentMethod
	}{
		{
			name: "card",
			in:   models.PaymentMethodCARD,
			want: payment.PaymentMethod_PAYMENT_METHOD_CARD,
		},
		{
			name: "sbp",
			in:   models.PaymentMethodSBP,
			want: payment.PaymentMethod_PAYMENT_METHOD_SBP,
		},
		{
			name: "credit card",
			in:   models.PaymentMethodCREDITCARD,
			want: payment.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD,
		},
		{
			name: "investor money",
			in:   models.PaymentMethodINVESTORMONEY,
			want: payment.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY,
		},
		{
			name: "unknown",
			in:   models.PaymentMethodUNKNOWN,
			want: payment.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, paymentMethodToPaymentMethod(tt.in))
		})
	}
}

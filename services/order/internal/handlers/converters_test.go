package handlers

import (
	"testing"

	orders "github.com/BeInBloom/big-deal/generated/go/order-service/openapi/v1"
	"github.com/BeInBloom/big-deal/services/order/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestToAPIOrderPaidOrder(t *testing.T) {
	orderUUID := uuid.New()
	userUUID := uuid.New()
	partUUID := uuid.New()
	transactionUUID := uuid.New()

	order := models.NewPaidOrder(
		models.OrderId(orderUUID),
		models.UserId(userUUID),
		models.Parts{
			{
				Id:          models.PartId(partUUID),
				Description: "engine",
				Price:       1500,
			},
		},
		models.PaymentMethodCARD,
		models.TransactionId(transactionUUID),
	)

	apiOrder := toAPIOrder(order)

	require.Equal(t, orderUUID, apiOrder.OrderUUID)
	require.Equal(t, userUUID, apiOrder.UserUUID)
	require.Equal(t, []uuid.UUID{partUUID}, apiOrder.PartUuids)
	require.Equal(t, 15.0, apiOrder.TotalPrice)
	require.Equal(t, orders.OrderStatusPAID, apiOrder.Status)
	require.Equal(t, transactionUUID, apiOrder.TransactionUUID.Or(uuid.Nil))
	require.Equal(t, orders.PaymentMethodCARD, apiOrder.PaymentMethod.Or(orders.PaymentMethodUNKNOWN))
}

func TestToPartIds(t *testing.T) {
	firstUUID := uuid.New()
	secondUUID := uuid.New()

	ids := toPartIds([]uuid.UUID{firstUUID, secondUUID})

	require.Equal(t, []models.PartId{
		models.PartId(firstUUID),
		models.PartId(secondUUID),
	}, ids)
}

func TestToAPIPartIds(t *testing.T) {
	firstUUID := uuid.New()
	secondUUID := uuid.New()

	ids := toAPIPartIds([]models.PartId{
		models.PartId(firstUUID),
		models.PartId(secondUUID),
	})

	require.Equal(t, []uuid.UUID{firstUUID, secondUUID}, ids)
}

func TestToAPIPrice(t *testing.T) {
	require.Equal(t, 15.0, toAPIPrice(1500))
}

func TestToPaymentMethod(t *testing.T) {
	tests := []struct {
		name string
		in   orders.PaymentMethod
		want models.PaymentMethod
	}{
		{
			name: "card",
			in:   orders.PaymentMethodCARD,
			want: models.PaymentMethodCARD,
		},
		{
			name: "sbp",
			in:   orders.PaymentMethodSBP,
			want: models.PaymentMethodSBP,
		},
		{
			name: "credit card",
			in:   orders.PaymentMethodCREDITCARD,
			want: models.PaymentMethodCREDITCARD,
		},
		{
			name: "investor money",
			in:   orders.PaymentMethodINVESTORMONEY,
			want: models.PaymentMethodINVESTORMONEY,
		},
		{
			name: "unknown",
			in:   orders.PaymentMethodUNKNOWN,
			want: models.PaymentMethodUNKNOWN,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, toPaymentMethod(tt.in))
		})
	}
}

func TestToAPIPaymentMethod(t *testing.T) {
	tests := []struct {
		name string
		in   models.PaymentMethod
		want orders.PaymentMethod
	}{
		{
			name: "card",
			in:   models.PaymentMethodCARD,
			want: orders.PaymentMethodCARD,
		},
		{
			name: "sbp",
			in:   models.PaymentMethodSBP,
			want: orders.PaymentMethodSBP,
		},
		{
			name: "credit card",
			in:   models.PaymentMethodCREDITCARD,
			want: orders.PaymentMethodCREDITCARD,
		},
		{
			name: "investor money",
			in:   models.PaymentMethodINVESTORMONEY,
			want: orders.PaymentMethodINVESTORMONEY,
		},
		{
			name: "unknown",
			in:   models.PaymentMethodUNKNOWN,
			want: orders.PaymentMethodUNKNOWN,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, toAPIPaymentMethod(tt.in))
		})
	}
}

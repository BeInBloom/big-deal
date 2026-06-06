package models

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestOrderSnapshotTryPaidOrder(t *testing.T) {
	orderId := OrderId(uuid.New())
	userId := UserId(uuid.New())
	partId := PartId(uuid.New())
	transactionId := TransactionId(uuid.New())

	snapshot := OrderSnapshot{
		Id:            orderId,
		UserId:        userId,
		Status:        Paid,
		Parts:         newTestParts(partId),
		PaymentMethod: Some(PaymentMethodCARD),
		TransactionId: Some(transactionId),
	}

	order, err := snapshot.TryPaidOrder()
	require.NoError(t, err)
	require.Equal(t, orderId, order.Id())
	require.Equal(t, userId, order.UserId())
	require.Equal(t, Paid, order.Status())
	require.Equal(t, uint(1500), order.Price())
	require.Equal(t, PaymentMethodCARD, order.PaymentMethod())
	require.Equal(t, transactionId, order.TransactionId())
}

func TestOrderSnapshotTryPaidOrderRequiresPaymentData(t *testing.T) {
	orderId := OrderId(uuid.New())
	userId := UserId(uuid.New())
	partId := PartId(uuid.New())

	tests := []struct {
		name     string
		snapshot OrderSnapshot
	}{
		{
			name: "missing payment method",
			snapshot: OrderSnapshot{
				Id:            orderId,
				UserId:        userId,
				Status:        Paid,
				Parts:         newTestParts(partId),
				TransactionId: Some(TransactionId(uuid.New())),
			},
		},
		{
			name: "missing transaction id",
			snapshot: OrderSnapshot{
				Id:            orderId,
				UserId:        userId,
				Status:        Paid,
				Parts:         newTestParts(partId),
				PaymentMethod: Some(PaymentMethodCARD),
			},
		},
		{
			name: "missing all payment data",
			snapshot: OrderSnapshot{
				Id:     orderId,
				UserId: userId,
				Status: Paid,
				Parts:  newTestParts(partId),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			order, err := tt.snapshot.TryPaidOrder()

			require.ErrorIs(t, err, ErrInvalidOrderSnapshot)
			require.Zero(t, order.Id())
		})
	}
}

func TestOrderSnapshotTryCanceledOrder(t *testing.T) {
	orderId := OrderId(uuid.New())
	userId := UserId(uuid.New())
	partId := PartId(uuid.New())

	snapshot := OrderSnapshot{
		Id:     orderId,
		UserId: userId,
		Status: Canceled,
		Parts:  newTestParts(partId),
	}

	order, err := snapshot.TryCanceledOrder()
	require.NoError(t, err)
	require.Equal(t, orderId, order.Id())
	require.Equal(t, userId, order.UserId())
	require.Equal(t, Canceled, order.Status())
	require.Equal(t, uint(1500), order.Price())
}

func TestOrderSnapshotTryCanceledOrderRejectsPaymentData(t *testing.T) {
	orderId := OrderId(uuid.New())
	userId := UserId(uuid.New())
	partId := PartId(uuid.New())
	transactionId := TransactionId(uuid.New())

	tests := []struct {
		name     string
		snapshot OrderSnapshot
	}{
		{
			name: "with payment method",
			snapshot: OrderSnapshot{
				Id:            orderId,
				UserId:        userId,
				Status:        Canceled,
				Parts:         newTestParts(partId),
				PaymentMethod: Some(PaymentMethodCARD),
			},
		},
		{
			name: "with transaction id",
			snapshot: OrderSnapshot{
				Id:            orderId,
				UserId:        userId,
				Status:        Canceled,
				Parts:         newTestParts(partId),
				TransactionId: Some(transactionId),
			},
		},
		{
			name: "with all payment data",
			snapshot: OrderSnapshot{
				Id:            orderId,
				UserId:        userId,
				Status:        Canceled,
				Parts:         newTestParts(partId),
				PaymentMethod: Some(PaymentMethodCARD),
				TransactionId: Some(transactionId),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			order, err := tt.snapshot.TryCanceledOrder()

			require.ErrorIs(t, err, ErrInvalidOrderSnapshot)
			require.Zero(t, order.Id())
		})
	}
}

func newTestParts(partTd PartId) Parts {
	return Parts{
		{
			Id:          partTd,
			Description: "engine",
			Price:       1500,
		},
	}
}

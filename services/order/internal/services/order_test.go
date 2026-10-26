package services

import (
	"context"
	"errors"
	"fmt"
	"math"
	"testing"

	"github.com/BeInBloom/big-deal/services/order/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestOrderServiceCreateOrder(t *testing.T) {
	ctx := context.Background()

	orderId := models.OrderId(uuid.New())
	userId := models.UserId(uuid.New())
	partId := models.PartId(uuid.New())

	parts := newTestParts(partId)
	deps := newOrderServiceTestDeps(t)
	deps.service.newId = func() models.OrderId {
		return orderId
	}

	deps.partService.EXPECT().
		ListParts(ctx, []models.PartId{partId}).
		Return(parts, nil).
		Once()

	deps.repo.EXPECT().
		Create(ctx, models.OrderSnapshot{
			Id:     orderId,
			UserId: userId,
			Status: models.PendingPayment,
			Parts:  parts,
		}).
		Return(nil).
		Once()

	order, err := deps.service.CreateOrder(ctx, userId, []models.PartId{partId})
	require.NoError(t, err)
	require.Equal(t, orderId, order.Id())
	require.Equal(t, models.PendingPayment, order.Status())
	requireOrderPrice(t, order, models.Money(1500))
}

func TestOrderServiceCreateOrderReturnsPartsError(t *testing.T) {
	ctx := context.Background()

	userId := models.UserId(uuid.New())
	partId := models.PartId(uuid.New())
	expectedErr := errors.New("parts unavailable")

	deps := newOrderServiceTestDeps(t)

	deps.partService.EXPECT().
		ListParts(ctx, []models.PartId{partId}).
		Return(nil, expectedErr).
		Once()

	order, err := deps.service.CreateOrder(ctx, userId, []models.PartId{partId})
	require.ErrorIs(t, err, expectedErr)
	require.Zero(t, order.Id())
}

func TestOrderServiceCreateOrderRejectsInvalidPartIds(t *testing.T) {
	partId := models.PartId(uuid.New())

	tests := []struct {
		name    string
		partIds []models.PartId
		wantErr error
	}{
		{name: "empty", partIds: nil, wantErr: ErrOrderPartsRequired},
		{name: "duplicated", partIds: []models.PartId{partId, partId}, wantErr: ErrOrderPartsDuplicated},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			deps := newOrderServiceTestDeps(t)

			order, err := deps.service.CreateOrder(context.Background(), models.UserId(uuid.New()), tt.partIds)
			require.ErrorIs(t, err, tt.wantErr)
			require.Zero(t, order.Id())
		})
	}
}

func TestOrderServiceCreateOrderReturnsNotFoundForMissingPart(t *testing.T) {
	ctx := context.Background()

	foundId := models.PartId(uuid.New())
	missingId := models.PartId(uuid.New())
	partIds := []models.PartId{foundId, missingId}

	deps := newOrderServiceTestDeps(t)

	deps.partService.EXPECT().
		ListParts(ctx, partIds).
		Return(newTestParts(foundId), nil).
		Once()

	order, err := deps.service.CreateOrder(ctx, models.UserId(uuid.New()), partIds)
	require.ErrorIs(t, err, ErrOrderPartsNotFound)
	require.Zero(t, order.Id())
}

func TestOrderServiceCreateOrderKeepsRequestedPartsOrder(t *testing.T) {
	ctx := context.Background()

	orderId := models.OrderId(uuid.New())
	userId := models.UserId(uuid.New())
	firstId := models.PartId(uuid.New())
	secondId := models.PartId(uuid.New())

	first := models.Part{Id: firstId, Description: "engine", Price: 1500}
	second := models.Part{Id: secondId, Description: "wing", Price: 700}

	deps := newOrderServiceTestDeps(t)
	deps.service.newId = func() models.OrderId {
		return orderId
	}

	deps.partService.EXPECT().
		ListParts(ctx, []models.PartId{firstId, secondId}).
		Return(models.Parts{second, first}, nil).
		Once()

	deps.repo.EXPECT().
		Create(ctx, models.OrderSnapshot{
			Id:     orderId,
			UserId: userId,
			Status: models.PendingPayment,
			Parts:  models.Parts{first, second},
		}).
		Return(nil).
		Once()

	order, err := deps.service.CreateOrder(ctx, userId, []models.PartId{firstId, secondId})
	require.NoError(t, err)
	requireOrderPrice(t, order, models.Money(2200))
}

func TestOrderServiceCreateOrderRejectsTotalOverflow(t *testing.T) {
	ctx := context.Background()

	firstId := models.PartId(uuid.New())
	secondId := models.PartId(uuid.New())
	partIds := []models.PartId{firstId, secondId}

	deps := newOrderServiceTestDeps(t)

	deps.partService.EXPECT().
		ListParts(ctx, partIds).
		Return(models.Parts{
			{Id: firstId, Price: math.MaxInt64},
			{Id: secondId, Price: 1},
		}, nil).
		Once()

	order, err := deps.service.CreateOrder(ctx, models.UserId(uuid.New()), partIds)
	require.ErrorIs(t, err, models.ErrInvalidMoney)
	require.Zero(t, order.Id())
}

func TestOrderServicePayOrderRejectsInvalidPaymentMethod(t *testing.T) {
	deps := newOrderServiceTestDeps(t)

	order, err := deps.service.PayOrder(
		context.Background(),
		models.OrderId(uuid.New()),
		models.PaymentMethodUNKNOWN,
	)
	require.ErrorIs(t, err, ErrInvalidPaymentMethod)
	require.Zero(t, order.Id())
}

func TestOrderServicePayOrderSavesPaidOrder(t *testing.T) {
	ctx := context.Background()

	orderId := models.OrderId(uuid.New())
	userId := models.UserId(uuid.New())
	partId := models.PartId(uuid.New())
	transactionId := models.TransactionId(uuid.New())

	parts := newTestParts(partId)
	deps := newOrderServiceTestDeps(t)

	deps.repo.EXPECT().
		Get(ctx, orderId).
		Return(models.OrderSnapshot{
			Id:     orderId,
			UserId: userId,
			Status: models.PendingPayment,
			Parts:  parts,
		}, nil).
		Once()

	deps.paymentService.EXPECT().
		PayOrder(ctx, userId, orderId, models.PaymentMethodCARD).
		Return(transactionId, nil).
		Once()

	deps.repo.EXPECT().
		Update(ctx, models.OrderSnapshot{
			Id:            orderId,
			UserId:        userId,
			Status:        models.Paid,
			Parts:         parts,
			PaymentMethod: models.Some(models.PaymentMethodCARD),
			TransactionId: models.Some(transactionId),
		}, models.PendingPayment).
		Return(nil).
		Once()

	order, err := deps.service.PayOrder(ctx, orderId, models.PaymentMethodCARD)
	require.NoError(t, err)
	require.Equal(t, orderId, order.Id())
	require.Equal(t, models.Paid, order.Status())
	require.Equal(t, models.PaymentMethodCARD, order.PaymentMethod())
	require.Equal(t, transactionId, order.TransactionId())
}

func TestOrderServicePayOrderReturnsCannotBePaidForCanceledOrder(t *testing.T) {
	ctx := context.Background()

	orderId := models.OrderId(uuid.New())
	userId := models.UserId(uuid.New())
	partId := models.PartId(uuid.New())
	parts := newTestParts(partId)

	deps := newOrderServiceTestDeps(t)

	deps.repo.EXPECT().
		Get(ctx, orderId).
		Return(models.OrderSnapshot{
			Id:     orderId,
			UserId: userId,
			Status: models.Canceled,
			Parts:  parts,
		}, nil).
		Once()

	order, err := deps.service.PayOrder(ctx, orderId, models.PaymentMethodCARD)
	require.ErrorIs(t, err, ErrOrderCannotBePaid)
	require.Zero(t, order.Id())
}

func TestOrderServicePayOrderReturnsPaymentError(t *testing.T) {
	ctx := context.Background()

	orderId := models.OrderId(uuid.New())
	userId := models.UserId(uuid.New())
	partId := models.PartId(uuid.New())
	parts := newTestParts(partId)
	expectedErr := errors.New("payment unavailable")

	deps := newOrderServiceTestDeps(t)

	deps.repo.EXPECT().
		Get(ctx, orderId).
		Return(models.OrderSnapshot{
			Id:     orderId,
			UserId: userId,
			Status: models.PendingPayment,
			Parts:  parts,
		}, nil).
		Once()

	deps.paymentService.EXPECT().
		PayOrder(ctx, userId, orderId, models.PaymentMethodCARD).
		Return(models.TransactionId{}, expectedErr).
		Once()

	order, err := deps.service.PayOrder(ctx, orderId, models.PaymentMethodCARD)
	require.ErrorIs(t, err, expectedErr)
	require.Zero(t, order.Id())
}

func TestOrderServicePayOrderReturnsCannotBePaidOnPaymentConflict(t *testing.T) {
	ctx := context.Background()

	orderId := models.OrderId(uuid.New())
	userId := models.UserId(uuid.New())
	parts := newTestParts(models.PartId(uuid.New()))

	deps := newOrderServiceTestDeps(t)

	deps.repo.EXPECT().
		Get(ctx, orderId).
		Return(models.OrderSnapshot{
			Id:     orderId,
			UserId: userId,
			Status: models.PendingPayment,
			Parts:  parts,
		}, nil).
		Once()

	deps.paymentService.EXPECT().
		PayOrder(ctx, userId, orderId, models.PaymentMethodCARD).
		Return(models.TransactionId{}, fmt.Errorf("wrapped: %w", ErrPaymentConflict)).
		Once()

	order, err := deps.service.PayOrder(ctx, orderId, models.PaymentMethodCARD)
	require.ErrorIs(t, err, ErrOrderCannotBePaid)
	require.Zero(t, order.Id())
}

func TestOrderServiceCancelOrderSavesCanceledOrder(t *testing.T) {
	ctx := context.Background()

	orderId := models.OrderId(uuid.New())
	userId := models.UserId(uuid.New())
	partId := models.PartId(uuid.New())
	parts := newTestParts(partId)

	deps := newOrderServiceTestDeps(t)

	deps.repo.EXPECT().
		Get(ctx, orderId).
		Return(models.OrderSnapshot{
			Id:     orderId,
			UserId: userId,
			Status: models.PendingPayment,
			Parts:  parts,
		}, nil).
		Once()

	deps.repo.EXPECT().
		Update(ctx, models.OrderSnapshot{
			Id:     orderId,
			UserId: userId,
			Status: models.Canceled,
			Parts:  parts,
		}, models.PendingPayment).
		Return(nil).
		Once()

	err := deps.service.CancelOrder(ctx, orderId)
	require.NoError(t, err)
}

func TestOrderServiceCancelOrderReturnsCannotBeCanceledForPaidOrder(t *testing.T) {
	ctx := context.Background()

	orderId := models.OrderId(uuid.New())
	userId := models.UserId(uuid.New())
	partId := models.PartId(uuid.New())
	transactionId := models.TransactionId(uuid.New())
	parts := newTestParts(partId)

	deps := newOrderServiceTestDeps(t)

	deps.repo.EXPECT().
		Get(ctx, orderId).
		Return(models.OrderSnapshot{
			Id:            orderId,
			UserId:        userId,
			Status:        models.Paid,
			Parts:         parts,
			PaymentMethod: models.Some(models.PaymentMethodCARD),
			TransactionId: models.Some(transactionId),
		}, nil).
		Once()

	err := deps.service.CancelOrder(ctx, orderId)
	require.ErrorIs(t, err, ErrOrderCannotBeCanceled)
}

func TestOrderServiceGetOrderReturnsOrder(t *testing.T) {
	ctx := context.Background()

	orderId := models.OrderId(uuid.New())
	userId := models.UserId(uuid.New())
	partId := models.PartId(uuid.New())
	parts := newTestParts(partId)

	deps := newOrderServiceTestDeps(t)

	deps.repo.EXPECT().
		Get(ctx, orderId).
		Return(models.OrderSnapshot{
			Id:     orderId,
			UserId: userId,
			Status: models.PendingPayment,
			Parts:  parts,
		}, nil).
		Once()

	order, err := deps.service.GetOrder(ctx, orderId)
	require.NoError(t, err)
	require.Equal(t, orderId, order.Id())
	require.Equal(t, models.PendingPayment, order.Status())
	requireOrderPrice(t, order, models.Money(1500))
}

func TestOrderServiceGetOrderReturnsRepoError(t *testing.T) {
	ctx := context.Background()

	orderId := models.OrderId(uuid.New())
	expectedErr := errors.New("repo unavailable")

	deps := newOrderServiceTestDeps(t)

	deps.repo.EXPECT().
		Get(ctx, orderId).
		Return(models.OrderSnapshot{}, expectedErr).
		Once()

	order, err := deps.service.GetOrder(ctx, orderId)
	require.ErrorIs(t, err, expectedErr)
	require.Nil(t, order)
}

func TestOrderServicePayOrderResolvesStatusConflict(t *testing.T) {
	orderId := models.OrderId(uuid.New())
	userId := models.UserId(uuid.New())
	parts := newTestParts(models.PartId(uuid.New()))
	transactionId := models.TransactionId(uuid.New())

	pending := models.OrderSnapshot{
		Id:     orderId,
		UserId: userId,
		Status: models.PendingPayment,
		Parts:  parts,
	}
	paidWith := func(id models.TransactionId) models.OrderSnapshot {
		return models.OrderSnapshot{
			Id:            orderId,
			UserId:        userId,
			Status:        models.Paid,
			Parts:         parts,
			PaymentMethod: models.Some(models.PaymentMethodCARD),
			TransactionId: models.Some(id),
		}
	}

	tests := []struct {
		name    string
		stored  models.OrderSnapshot
		wantErr error
	}{
		{
			name:   "concurrent pay stored the same transaction",
			stored: paidWith(transactionId),
		},
		{
			name:    "concurrent pay stored another transaction",
			stored:  paidWith(models.TransactionId(uuid.New())),
			wantErr: ErrOrderCannotBePaid,
		},
		{
			name: "concurrent cancel",
			stored: models.OrderSnapshot{
				Id:     orderId,
				UserId: userId,
				Status: models.Canceled,
				Parts:  parts,
			},
			wantErr: ErrOrderCannotBePaid,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			deps := newOrderServiceTestDeps(t)

			deps.repo.EXPECT().Get(ctx, orderId).Return(pending, nil).Once()

			deps.paymentService.EXPECT().
				PayOrder(ctx, userId, orderId, models.PaymentMethodCARD).
				Return(transactionId, nil).
				Once()

			deps.repo.EXPECT().
				Update(ctx, mock.Anything, models.PendingPayment).
				Return(ErrOrderStatusConflict).
				Once()

			deps.repo.EXPECT().Get(ctx, orderId).Return(tt.stored, nil).Once()

			order, err := deps.service.PayOrder(ctx, orderId, models.PaymentMethodCARD)
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				require.Zero(t, order.Id())
				return
			}

			require.NoError(t, err)
			require.Equal(t, orderId, order.Id())
			require.Equal(t, transactionId, order.TransactionId())
		})
	}
}

func TestOrderServiceCancelOrderReturnsCannotBeCanceledOnStatusConflict(t *testing.T) {
	ctx := context.Background()

	orderId := models.OrderId(uuid.New())
	userId := models.UserId(uuid.New())
	partId := models.PartId(uuid.New())
	parts := newTestParts(partId)

	deps := newOrderServiceTestDeps(t)

	deps.repo.EXPECT().
		Get(ctx, orderId).
		Return(models.OrderSnapshot{
			Id:     orderId,
			UserId: userId,
			Status: models.PendingPayment,
			Parts:  parts,
		}, nil).
		Once()

	deps.repo.EXPECT().
		Update(ctx, mock.Anything, models.PendingPayment).
		Return(ErrOrderStatusConflict).
		Once()

	err := deps.service.CancelOrder(ctx, orderId)
	require.ErrorIs(t, err, ErrOrderCannotBeCanceled)
}

type orderServiceTestDeps struct {
	repo           *MockOrderRepo
	partService    *MockPartService
	paymentService *MockPaymentService
	service        *OrderService
}

func newOrderServiceTestDeps(t *testing.T) orderServiceTestDeps {
	t.Helper()

	repo := NewMockOrderRepo(t)
	partService := NewMockPartService(t)
	paymentService := NewMockPaymentService(t)

	return orderServiceTestDeps{
		repo:           repo,
		partService:    partService,
		paymentService: paymentService,
		service:        New(repo, partService, paymentService),
	}
}

func newTestParts(partId models.PartId) models.Parts {
	return models.Parts{
		{
			Id:          partId,
			Description: "engine",
			Price:       1500,
		},
	}
}

func requireOrderPrice(t *testing.T, order models.Order, want models.Money) {
	t.Helper()

	price, err := order.Price()
	require.NoError(t, err)
	require.Equal(t, want, price)
}

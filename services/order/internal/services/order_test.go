package services

import (
	"context"
	"errors"
	"testing"

	"github.com/BeInBloom/big-deal/services/order/internal/models"
	"github.com/google/uuid"
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
		Save(ctx, models.OrderSnapshot{
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
	require.Equal(t, uint(1500), order.Price())
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
		Save(ctx, models.OrderSnapshot{
			Id:            orderId,
			UserId:        userId,
			Status:        models.Paid,
			Parts:         parts,
			PaymentMethod: models.Some(models.PaymentMethodCARD),
			TransactionId: models.Some(transactionId),
		}).
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
		Save(ctx, models.OrderSnapshot{
			Id:     orderId,
			UserId: userId,
			Status: models.Canceled,
			Parts:  parts,
		}).
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
	require.Equal(t, uint(1500), order.Price())
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

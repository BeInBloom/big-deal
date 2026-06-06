package services

import (
	"context"
	"errors"
	"testing"

	"github.com/BeInBloom/big-deal/services/order/internal/models"
	"github.com/google/uuid"
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
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if order.Id() != orderId {
		t.Fatalf("expected order id %v, got %v", orderId, order.Id())
	}

	if order.Status() != models.PendingPayment {
		t.Fatalf("expected status %q, got %q", models.PendingPayment, order.Status())
	}

	if order.Price() != 1500 {
		t.Fatalf("expected price 1500, got %d", order.Price())
	}
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
	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected error %v, got %v", expectedErr, err)
	}

	if order.Id() != (models.OrderId{}) {
		t.Fatalf("expected zero order, got id %v", order.Id())
	}
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
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if order.Id() != orderId {
		t.Fatalf("expected order id %v, got %v", orderId, order.Id())
	}

	if order.Status() != models.Paid {
		t.Fatalf("expected status %q, got %q", models.Paid, order.Status())
	}

	if order.PaymentMethod() != models.PaymentMethodCARD {
		t.Fatalf("expected payment method %v, got %v", models.PaymentMethodCARD, order.PaymentMethod())
	}

	if order.TransactionId() != transactionId {
		t.Fatalf("expected transaction id %v, got %v", transactionId, order.TransactionId())
	}
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
	if !errors.Is(err, ErrOrderCannotBePaid) {
		t.Fatalf("expected error %v, got %v", ErrOrderCannotBePaid, err)
	}

	if order.Id() != (models.OrderId{}) {
		t.Fatalf("expected zero order, got id %v", order.Id())
	}
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
	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected error %v, got %v", expectedErr, err)
	}

	if order.Id() != (models.OrderId{}) {
		t.Fatalf("expected zero order, got id %v", order.Id())
	}
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
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
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
	if !errors.Is(err, ErrOrderCannotBeCanceled) {
		t.Fatalf("expected error %v, got %v", ErrOrderCannotBeCanceled, err)
	}
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
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if order.Id() != orderId {
		t.Fatalf("expected order id %v, got %v", orderId, order.Id())
	}

	if order.Status() != models.PendingPayment {
		t.Fatalf("expected status %q, got %q", models.PendingPayment, order.Status())
	}

	if order.Price() != 1500 {
		t.Fatalf("expected price 1500, got %d", order.Price())
	}
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
	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected error %v, got %v", expectedErr, err)
	}

	if order != nil {
		t.Fatalf("expected nil order, got %v", order)
	}
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

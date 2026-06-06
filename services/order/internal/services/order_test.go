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

	parts := models.Parts{
		{
			Id:          partId,
			Description: "engine",
			Price:       1500,
		},
	}

	repo := NewMockOrderRepo(t)
	partService := NewMockPartService(t)
	paymentService := NewMockPaymentService(t)

	service := New(repo, partService, paymentService)
	service.newId = func() models.OrderId {
		return orderId
	}

	partService.EXPECT().
		ListParts(ctx, []models.PartId{partId}).
		Return(parts, nil).
		Once()

	repo.EXPECT().
		Save(ctx, models.OrderSnapshot{
			Id:     orderId,
			UserId: userId,
			Status: models.PendingPayment,
			Parts:  parts,
		}).
		Return(nil).
		Once()

	order, err := service.CreateOrder(ctx, userId, []models.PartId{partId})
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

	repo := NewMockOrderRepo(t)
	partService := NewMockPartService(t)
	paymentService := NewMockPaymentService(t)

	service := New(repo, partService, paymentService)

	partService.EXPECT().
		ListParts(ctx, []models.PartId{partId}).
		Return(nil, expectedErr).
		Once()

	order, err := service.CreateOrder(ctx, userId, []models.PartId{partId})
	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected error %v, got %v", expectedErr, err)
	}

	if order.Id() != (models.OrderId{}) {
		t.Fatalf("expected zero order, got id %v", order.Id())
	}
}

package services

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"github.com/BeInBloom/big-deal/services/order/internal/models"
)

var (
	ErrOrderNotFound         = errors.New("order not found")
	ErrOrderPartsNotFound    = errors.New("order parts not found")
	ErrOrderCannotBePaid     = errors.New("order cannot be paid")
	ErrOrderCannotBeCanceled = errors.New("order cannot be canceled")
)

type OrderRepo interface {
	Get(ctx context.Context, id models.OrderId) (models.OrderSnapshot, error)
	Save(ctx context.Context, order models.OrderSnapshot) error
}

type PartService interface {
	ListParts(ctx context.Context, ids []models.PartId) (models.Parts, error)
}

type PaymentService interface {
	// TODO: payment adapter must make PayOrder idempotent by orderId.
	PayOrder(
		ctx context.Context,
		userId models.UserId,
		orderId models.OrderId,
		method models.PaymentMethod,
	) (models.TransactionId, error)
}

type OrderService struct {
	repo     OrderRepo
	parts    PartService
	payments PaymentService
	newId    func() models.OrderId
}

func New(
	repo OrderRepo,
	parts PartService,
	payments PaymentService,
) *OrderService {
	return &OrderService{
		repo:     repo,
		parts:    parts,
		payments: payments,
		newId: func() models.OrderId {
			return models.OrderId(uuid.New())
		},
	}
}

func (s *OrderService) CancelOrder(
	ctx context.Context,
	id models.OrderId,
) error {
	snapshot, err := s.repo.Get(ctx, id)
	if err != nil {
		return err
	}

	order, err := snapshot.TryPendingOrder()
	if err != nil {
		return ErrOrderCannotBeCanceled
	}

	return s.repo.Save(ctx, order.Cancel().Snapshot())
}

func (s *OrderService) CreateOrder(
	ctx context.Context,
	userId models.UserId,
	partIds []models.PartId,
) (models.PendingOrder, error) {
	parts, err := s.parts.ListParts(ctx, partIds)
	if err != nil {
		return models.PendingOrder{}, err
	}

	order := models.NewPendingOrder(s.newId(), userId, parts)
	if err := s.repo.Save(ctx, order.Snapshot()); err != nil {
		return models.PendingOrder{}, err
	}

	return order, nil
}

func (s *OrderService) GetOrder(
	ctx context.Context,
	id models.OrderId,
) (models.Order, error) {
	snapshot, err := s.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	return snapshot.TryOrder()
}

func (s *OrderService) PayOrder(
	ctx context.Context,
	id models.OrderId,
	method models.PaymentMethod,
) (models.PaidOrder, error) {
	snapshot, err := s.repo.Get(ctx, id)
	if err != nil {
		return models.PaidOrder{}, err
	}

	order, err := snapshot.TryPendingOrder()
	if err != nil {
		return models.PaidOrder{}, ErrOrderCannotBePaid
	}

	transactionId, err := s.payments.PayOrder(
		ctx, order.UserId(), order.Id(), method)
	if err != nil {
		return models.PaidOrder{}, err
	}

	paidOrder := order.Pay(method, transactionId)
	if err := s.repo.Save(ctx, paidOrder.Snapshot()); err != nil {
		return models.PaidOrder{}, err
	}

	return paidOrder, nil
}

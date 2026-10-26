package services

import (
	"context"
	"errors"

	"github.com/BeInBloom/big-deal/services/order/internal/models"
	"github.com/google/uuid"
)

var (
	ErrOrderNotFound         = errors.New("order not found")
	ErrOrderPartsRequired    = errors.New("order parts are required")
	ErrOrderPartsDuplicated  = errors.New("order parts are duplicated")
	ErrOrderPartsNotFound    = errors.New("order parts not found")
	ErrOrderCannotBePaid     = errors.New("order cannot be paid")
	ErrOrderCannotBeCanceled = errors.New("order cannot be canceled")
	ErrOrderStatusConflict   = errors.New("order status changed concurrently")
	ErrInvalidPaymentMethod  = errors.New("invalid payment method")
)

type OrderRepo interface {
	Get(ctx context.Context, id models.OrderId) (models.OrderSnapshot, error)
	Create(ctx context.Context, order models.OrderSnapshot) error
	// Update persists order only if its stored status is still from,
	// otherwise it returns ErrOrderStatusConflict.
	Update(ctx context.Context, order models.OrderSnapshot, from models.OrderStatus) error
}

type PartService interface {
	ListParts(ctx context.Context, ids []models.PartId) (models.Parts, error)
}

type PaymentService interface {
	// PayOrder must be idempotent by orderId: a retry returns the same transaction.
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

	err = s.repo.Update(ctx, order.Cancel().Snapshot(), order.Status())
	if errors.Is(err, ErrOrderStatusConflict) {
		return ErrOrderCannotBeCanceled
	}

	return err
}

func (s *OrderService) CreateOrder(
	ctx context.Context,
	userId models.UserId,
	partIds []models.PartId,
) (models.PendingOrder, error) {
	if err := validatePartIds(partIds); err != nil {
		return models.PendingOrder{}, err
	}

	found, err := s.parts.ListParts(ctx, partIds)
	if err != nil {
		return models.PendingOrder{}, err
	}

	parts, err := pickRequestedParts(partIds, found)
	if err != nil {
		return models.PendingOrder{}, err
	}

	order := models.NewPendingOrder(s.newId(), userId, parts)
	if err := s.repo.Create(ctx, order.Snapshot()); err != nil {
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
	if !method.Valid() {
		return models.PaidOrder{}, ErrInvalidPaymentMethod
	}

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
	err = s.repo.Update(ctx, paidOrder.Snapshot(), order.Status())
	if errors.Is(err, ErrOrderStatusConflict) {
		return models.PaidOrder{}, ErrOrderCannotBePaid
	}
	if err != nil {
		return models.PaidOrder{}, err
	}

	return paidOrder, nil
}

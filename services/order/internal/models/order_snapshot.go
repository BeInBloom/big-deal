package models

import "errors"

var ErrInvalidOrderSnapshot = errors.New("invalid order snapshot")

type OrderSnapshot struct {
	Id     OrderId
	UserId UserId
	Price  float64
	Status OrderStatus
	Parts  []PartId

	PaymentMethod Option[PaymentMethod]
	TransactionId Option[TransactionId]
}

func (s OrderSnapshot) TryOrder() (Order, error) {
	switch s.Status {
	case PendingPayment:
		return s.TryPendingOrder()
	case Paid:
		return s.TryPaidOrder()
	case Canceled:
		return s.TryCanceledOrder()
	default:
		return nil, ErrInvalidOrderSnapshot
	}
}

func (s OrderSnapshot) TryPendingOrder() (PendingOrderState, error) {
	if s.Status != PendingPayment {
		return PendingOrderState{}, ErrInvalidOrderSnapshot
	}

	return NewPendingOrder(s.Id, s.UserId, s.Price, s.Parts), nil
}

func (s OrderSnapshot) TryPaidOrder() (PaidOrderState, error) {
	if s.Status != Paid {
		return PaidOrderState{}, ErrInvalidOrderSnapshot
	}

	paymentMethod, ok := s.PaymentMethod.Get()
	if !ok {
		return PaidOrderState{}, ErrInvalidOrderSnapshot
	}

	transactionId, ok := s.TransactionId.Get()
	if !ok {
		return PaidOrderState{}, ErrInvalidOrderSnapshot
	}

	return NewPaidOrder(
		s.Id,
		s.UserId,
		s.Price,
		s.Parts,
		paymentMethod,
		transactionId,
	), nil
}

func (s OrderSnapshot) TryCanceledOrder() (CanceledOrderState, error) {
	if s.Status != Canceled {
		return CanceledOrderState{}, ErrInvalidOrderSnapshot
	}

	return NewCanceledOrder(s.Id, s.UserId, s.Price, s.Parts), nil
}

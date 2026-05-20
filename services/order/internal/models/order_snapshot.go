package models

import "errors"

var ErrInvalidOrderSnapshot = errors.New("invalid order snapshot")

type OrderSnapshot struct {
	Id     OrderId
	UserId UserId
	Status OrderStatus
	Parts  Parts

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
	if s.hasPaymentData() {
		return PendingOrderState{}, ErrInvalidOrderSnapshot
	}

	return NewPendingOrder(s.Id, s.UserId, s.Parts), nil
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
		s.Parts,
		paymentMethod,
		transactionId,
	), nil
}

func (s OrderSnapshot) TryCanceledOrder() (CanceledOrderState, error) {
	if s.Status != Canceled {
		return CanceledOrderState{}, ErrInvalidOrderSnapshot
	}
	if s.hasPaymentData() {
		return CanceledOrderState{}, ErrInvalidOrderSnapshot
	}

	return NewCanceledOrder(s.Id, s.UserId, s.Parts), nil
}

func (s OrderSnapshot) hasPaymentData() bool {
	return s.PaymentMethod.IsSet() || s.TransactionId.IsSet()
}

package models

import (
	"github.com/google/uuid"
)

type OrderStatus string

const (
	PendingPayment OrderStatus = "PENDING_PAYMENT"
	Paid           OrderStatus = "PAID"
	Canceled       OrderStatus = "CANCELLED"
)

type PaymentMethod int

const (
	PaymentMethodUNKNOWN PaymentMethod = iota
	PaymentMethodCARD
	PaymentMethodSBP
	PaymentMethodCREDITCARD
	PaymentMethodINVESTORMONEY
)

// Valid reports whether m is a known payment method.
func (m PaymentMethod) Valid() bool {
	switch m {
	case PaymentMethodCARD,
		PaymentMethodSBP,
		PaymentMethodCREDITCARD,
		PaymentMethodINVESTORMONEY:
		return true
	default:
		return false
	}
}

type (
	UserId        uuid.UUID
	OrderId       uuid.UUID
	PartId        uuid.UUID
	TransactionId uuid.UUID
)

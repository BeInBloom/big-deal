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

type (
	UserId        uuid.UUID
	OrderId       uuid.UUID
	PartId        uuid.UUID
	TransactionId uuid.UUID
)

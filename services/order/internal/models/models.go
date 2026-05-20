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

type PaymentMethod string

const (
	PaymentMethodUNKNOWN       PaymentMethod = "UNKNOWN"
	PaymentMethodCARD          PaymentMethod = "CARD"
	PaymentMethodSBP           PaymentMethod = "SBP"
	PaymentMethodCREDITCARD    PaymentMethod = "CREDIT_CARD"
	PaymentMethodINVESTORMONEY PaymentMethod = "INVESTOR_MONEY"
)

type (
	UserId        uuid.UUID
	OrderId       uuid.UUID
	PartId        uuid.UUID
	TransactionId uuid.UUID
)

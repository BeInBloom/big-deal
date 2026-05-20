package models

import "slices"

type CanceledOrder = CanceledOrderState
type PaidOrder = PaidOrderState
type PendingOrder = PendingOrderState

type Order interface {
	Id() OrderId
	UserId() UserId
	Price() float64
	Status() OrderStatus
	Parts() []PartId
	Snapshot() OrderSnapshot
}

type orderData struct {
	id     OrderId
	userId UserId
	price  float64
	parts  []PartId
}

type (
	CanceledOrderState struct {
		orderData
	}

	PaidOrderState struct {
		orderData

		paymentMethod PaymentMethod
		transactionId TransactionId
	}

	PendingOrderState struct {
		orderData
	}
)

func NewPendingOrder(id OrderId, userId UserId, price float64, parts []PartId) PendingOrderState {
	return PendingOrderState{
		orderData: newOrderData(id, userId, price, parts),
	}
}

func NewPaidOrder(
	id OrderId,
	userId UserId,
	price float64,
	parts []PartId,
	paymentMethod PaymentMethod,
	transactionId TransactionId,
) PaidOrderState {
	return PaidOrderState{
		orderData:     newOrderData(id, userId, price, parts),
		paymentMethod: paymentMethod,
		transactionId: transactionId,
	}
}

func NewCanceledOrder(id OrderId, userId UserId, price float64, parts []PartId) CanceledOrderState {
	return CanceledOrderState{
		orderData: newOrderData(id, userId, price, parts),
	}
}

func newOrderData(id OrderId, userId UserId, price float64, parts []PartId) orderData {
	return orderData{
		id:     id,
		userId: userId,
		price:  price,
		parts:  slices.Clone(parts),
	}
}

func (o orderData) Id() OrderId {
	return o.id
}

func (o orderData) UserId() UserId {
	return o.userId
}

func (o orderData) Price() float64 {
	return o.price
}

func (o orderData) Parts() []PartId {
	return slices.Clone(o.parts)
}

func (o CanceledOrderState) Status() OrderStatus {
	return Canceled
}

func (o CanceledOrderState) Snapshot() OrderSnapshot {
	return OrderSnapshot{
		Id:     o.Id(),
		UserId: o.UserId(),
		Price:  o.Price(),
		Status: o.Status(),
		Parts:  o.Parts(),
	}
}

func (o PaidOrderState) Status() OrderStatus {
	return Paid
}

func (o PaidOrderState) PaymentMethod() PaymentMethod {
	return o.paymentMethod
}

func (o PaidOrderState) TransactionId() TransactionId {
	return o.transactionId
}

func (o PaidOrderState) Snapshot() OrderSnapshot {
	return OrderSnapshot{
		Id:            o.Id(),
		UserId:        o.UserId(),
		Price:         o.Price(),
		Status:        o.Status(),
		Parts:         o.Parts(),
		PaymentMethod: Some(o.paymentMethod),
		TransactionId: Some(o.transactionId),
	}
}

func (o PendingOrderState) Status() OrderStatus {
	return PendingPayment
}

func (o PendingOrderState) Pay(method PaymentMethod, transactionId TransactionId) PaidOrderState {
	return NewPaidOrder(o.Id(), o.UserId(), o.Price(), o.Parts(), method, transactionId)
}

func (o PendingOrderState) Cancel() CanceledOrderState {
	return NewCanceledOrder(o.Id(), o.UserId(), o.Price(), o.Parts())
}

func (o PendingOrderState) Snapshot() OrderSnapshot {
	return OrderSnapshot{
		Id:     o.Id(),
		UserId: o.UserId(),
		Price:  o.Price(),
		Status: o.Status(),
		Parts:  o.Parts(),
	}
}

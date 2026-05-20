package models

type (
	CanceledOrder = CanceledOrderState
	PaidOrder     = PaidOrderState
	PendingOrder  = PendingOrderState
)

type Order interface {
	Id() OrderId
	UserId() UserId
	Price() uint
	Status() OrderStatus
	Parts() Parts
	Snapshot() OrderSnapshot
}

type orderData struct {
	id     OrderId
	userId UserId
	parts  Parts
}

func (o orderData) Price() uint {
	return o.parts.Price()
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

func NewPendingOrder(id OrderId, userId UserId, parts Parts) PendingOrderState {
	return PendingOrderState{
		orderData: newOrderData(id, userId, parts),
	}
}

func NewPaidOrder(
	id OrderId,
	userId UserId,
	parts Parts,
	paymentMethod PaymentMethod,
	transactionId TransactionId,
) PaidOrderState {
	return PaidOrderState{
		orderData:     newOrderData(id, userId, parts),
		paymentMethod: paymentMethod,
		transactionId: transactionId,
	}
}

func NewCanceledOrder(id OrderId, userId UserId, parts Parts) CanceledOrderState {
	return CanceledOrderState{
		orderData: newOrderData(id, userId, parts),
	}
}

func newOrderData(id OrderId, userId UserId, parts Parts) orderData {
	return orderData{
		id:     id,
		userId: userId,
		parts:  parts.Clone(),
	}
}

func (o orderData) Id() OrderId {
	return o.id
}

func (o orderData) UserId() UserId {
	return o.userId
}

func (o orderData) Parts() Parts {
	return o.parts.Clone()
}

func (o CanceledOrderState) Status() OrderStatus {
	return Canceled
}

func (o CanceledOrderState) Snapshot() OrderSnapshot {
	return OrderSnapshot{
		Id:     o.Id(),
		UserId: o.UserId(),
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
	return NewPaidOrder(o.Id(), o.UserId(), o.Parts(), method, transactionId)
}

func (o PendingOrderState) Cancel() CanceledOrderState {
	return NewCanceledOrder(o.Id(), o.UserId(), o.Parts())
}

func (o PendingOrderState) Snapshot() OrderSnapshot {
	return OrderSnapshot{
		Id:     o.Id(),
		UserId: o.UserId(),
		Status: o.Status(),
		Parts:  o.Parts(),
	}
}

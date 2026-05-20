package handlers

import (
	orders "github.com/BeInBloom/big-deal/generated/go/order-service/openapi/v1"
	"github.com/BeInBloom/big-deal/services/order/internal/models"
	"github.com/google/uuid"
)

func toAPIOrder(order models.Order) *orders.Order {
	snapshot := order.Snapshot()

	res := &orders.Order{
		OrderUUID:  uuid.UUID(snapshot.Id),
		UserUUID:   uuid.UUID(snapshot.UserId),
		PartUuids:  toAPIPartIds(snapshot.Parts.Ids()),
		TotalPrice: toAPIPrice(order.Price()),
		Status:     orders.OrderStatus(snapshot.Status),
	}

	if transactionId, ok := snapshot.TransactionId.Get(); ok {
		res.TransactionUUID = orders.NewOptNilUUID(uuid.UUID(transactionId))
	}

	if paymentMethod, ok := snapshot.PaymentMethod.Get(); ok {
		res.PaymentMethod = orders.NewOptNilPaymentMethod(orders.PaymentMethod(paymentMethod))
	}

	return res
}

func toPartIds(ids []uuid.UUID) []models.PartId {
	partIds := make([]models.PartId, 0, len(ids))

	for _, id := range ids {
		partIds = append(partIds, models.PartId(id))
	}

	return partIds
}

func toAPIPartIds(ids []models.PartId) []uuid.UUID {
	apiIds := make([]uuid.UUID, 0, len(ids))

	for _, id := range ids {
		apiIds = append(apiIds, uuid.UUID(id))
	}

	return apiIds
}

func toAPIPrice(price uint) float64 {
	return float64(price) / 100
}

package handlers

import (
	"errors"
	"fmt"

	orders "github.com/BeInBloom/big-deal/generated/go/order-service/openapi/v1"
	"github.com/BeInBloom/big-deal/services/order/internal/models"
	"github.com/google/uuid"
)

var (
	errUnknownPaymentMethod = errors.New("unknown payment method")
	errUnknownOrderStatus   = errors.New("unknown order status")
)

func toAPIOrder(order models.Order) (*orders.Order, error) {
	snapshot := order.Snapshot()

	status, err := toAPIOrderStatus(snapshot.Status)
	if err != nil {
		return nil, err
	}

	res := &orders.Order{
		OrderUUID:  uuid.UUID(snapshot.Id),
		UserUUID:   uuid.UUID(snapshot.UserId),
		PartUuids:  toAPIPartIds(snapshot.Parts.Ids()),
		TotalPrice: order.Price().Major(),
		Status:     status,
	}

	if transactionId, ok := snapshot.TransactionId.Get(); ok {
		res.TransactionUUID = orders.NewOptNilUUID(uuid.UUID(transactionId))
	}

	if paymentMethod, ok := snapshot.PaymentMethod.Get(); ok {
		apiPaymentMethod, err := toAPIPaymentMethod(paymentMethod)
		if err != nil {
			return nil, err
		}

		res.PaymentMethod = orders.NewOptNilPaymentMethod(apiPaymentMethod)
	}

	return res, nil
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

func toAPIOrderStatus(status models.OrderStatus) (orders.OrderStatus, error) {
	switch status {
	case models.PendingPayment:
		return orders.OrderStatusPENDINGPAYMENT, nil
	case models.Paid:
		return orders.OrderStatusPAID, nil
	case models.Canceled:
		return orders.OrderStatusCANCELLED, nil
	default:
		return "", fmt.Errorf("%w: %q", errUnknownOrderStatus, status)
	}
}

func toPaymentMethod(paymentMethod orders.PaymentMethod) models.PaymentMethod {
	switch paymentMethod {
	case orders.PaymentMethodCARD:
		return models.PaymentMethodCARD
	case orders.PaymentMethodSBP:
		return models.PaymentMethodSBP
	case orders.PaymentMethodCREDITCARD:
		return models.PaymentMethodCREDITCARD
	case orders.PaymentMethodINVESTORMONEY:
		return models.PaymentMethodINVESTORMONEY
	default:
		return models.PaymentMethodUNKNOWN
	}
}

func toAPIPaymentMethod(paymentMethod models.PaymentMethod) (orders.PaymentMethod, error) {
	switch paymentMethod {
	case models.PaymentMethodCARD:
		return orders.PaymentMethodCARD, nil
	case models.PaymentMethodSBP:
		return orders.PaymentMethodSBP, nil
	case models.PaymentMethodCREDITCARD:
		return orders.PaymentMethodCREDITCARD, nil
	case models.PaymentMethodINVESTORMONEY:
		return orders.PaymentMethodINVESTORMONEY, nil
	default:
		return "", fmt.Errorf("%w: %d", errUnknownPaymentMethod, paymentMethod)
	}
}

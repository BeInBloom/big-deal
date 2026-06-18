package psql

import (
	"fmt"

	"github.com/BeInBloom/big-deal/services/order/internal/models"
	dbenums "github.com/BeInBloom/big-deal/services/order/internal/repo/psql/gen/enums"
	dbmodels "github.com/BeInBloom/big-deal/services/order/internal/repo/psql/gen/models"
	"github.com/aarondl/opt/omit"
	"github.com/aarondl/opt/omitnull"
	"github.com/google/uuid"
)

const maxInt64 = uint64(1<<63 - 1)

func tryIntoSnapshot(order *dbmodels.Order) (models.OrderSnapshot, error) {
	status, err := orderStatusToModel(order.Status)
	if err != nil {
		return models.OrderSnapshot{}, err
	}

	if !order.R.Loaded.OrderParts {
		return models.OrderSnapshot{}, errOrderPartsNotLoaded
	}

	parts, err := orderPartsToModel(order.R.OrderParts)
	if err != nil {
		return models.OrderSnapshot{}, err
	}

	paymentMethod := models.None[models.PaymentMethod]()
	if raw, ok := order.PaymentMethod.Get(); ok {
		method, err := paymentMethodToModel(raw)
		if err != nil {
			return models.OrderSnapshot{}, err
		}
		paymentMethod = models.Some(method)
	}

	transactionId := models.None[models.TransactionId]()
	if id, ok := order.TransactionID.Get(); ok {
		transactionId = models.Some(models.TransactionId(id))
	}

	return models.OrderSnapshot{
		Id:            models.OrderId(order.ID),
		UserId:        models.UserId(order.UserID),
		Status:        status,
		Parts:         parts,
		PaymentMethod: paymentMethod,
		TransactionId: transactionId,
	}, nil
}

func orderStatusToModel(status dbenums.OrderStatus) (models.OrderStatus, error) {
	switch status {
	case dbenums.OrderStatusPENDING_PAYMENT:
		return models.PendingPayment, nil
	case dbenums.OrderStatusPAID:
		return models.Paid, nil
	case dbenums.OrderStatusCANCELLED:
		return models.Canceled, nil
	default:
		return "", fmt.Errorf("%w: %q", errUnknownOrderStatus, status)
	}
}

func paymentMethodToModel(method dbenums.PaymentMethod) (models.PaymentMethod, error) {
	switch method {
	case dbenums.PaymentMethodCARD:
		return models.PaymentMethodCARD, nil
	case dbenums.PaymentMethodSBP:
		return models.PaymentMethodSBP, nil
	case dbenums.PaymentMethodCREDIT_CARD:
		return models.PaymentMethodCREDITCARD, nil
	case dbenums.PaymentMethodINVESTOR_MONEY:
		return models.PaymentMethodINVESTORMONEY, nil
	default:
		return models.PaymentMethodUNKNOWN, fmt.Errorf("%w: %q", errUnknownPaymentMethod, method)
	}
}

func orderPartsToModel(parts dbmodels.OrderPartSlice) (models.Parts, error) {
	result := make(models.Parts, 0, len(parts))

	for _, part := range parts {
		if part.Price < 0 || uint64(part.Price) > uint64(^uint(0)) {
			return nil, fmt.Errorf("%w: %d", errInvalidPartPrice, part.Price)
		}

		result = append(result, models.Part{
			Id:          models.PartId(part.PartID),
			Description: part.Description,
			Price:       uint(part.Price),
		})
	}

	return result, nil
}

func intoOrderSetter(
	order models.OrderSnapshot,
) (*dbmodels.OrderSetter, error) {
	status, err := orderStatusToDB(order.Status)
	if err != nil {
		return nil, err
	}

	paymentMethod, err := nullablePaymentMethod(order.PaymentMethod)
	if err != nil {
		return nil, err
	}

	return &dbmodels.OrderSetter{
		ID:            omit.From(uuid.UUID(order.Id)),
		UserID:        omit.From(uuid.UUID(order.UserId)),
		Status:        omit.From(status),
		PaymentMethod: paymentMethod,
		TransactionID: nullableTransactionID(order.TransactionId),
	}, nil
}

func intoOrderPartSetters(
	order models.OrderSnapshot,
) ([]*dbmodels.OrderPartSetter, error) {
	setters := make([]*dbmodels.OrderPartSetter, 0, len(order.Parts))

	for _, part := range order.Parts {
		if uint64(part.Price) > maxInt64 {
			return nil, fmt.Errorf("%w: %d", errInvalidPartPrice, part.Price)
		}

		setters = append(setters, &dbmodels.OrderPartSetter{
			OrderID:     omit.From(uuid.UUID(order.Id)),
			PartID:      omit.From(uuid.UUID(part.Id)),
			Description: omit.From(part.Description),
			Price:       omit.From(int64(part.Price)),
		})
	}

	return setters, nil
}

func orderStatusToDB(status models.OrderStatus) (dbenums.OrderStatus, error) {
	switch status {
	case models.PendingPayment:
		return dbenums.OrderStatusPENDING_PAYMENT, nil
	case models.Paid:
		return dbenums.OrderStatusPAID, nil
	case models.Canceled:
		return dbenums.OrderStatusCANCELLED, nil
	default:
		return "", fmt.Errorf("%w: %q", errUnknownOrderStatus, status)
	}
}

func paymentMethodToDB(method models.PaymentMethod) (dbenums.PaymentMethod, error) {
	switch method {
	case models.PaymentMethodCARD:
		return dbenums.PaymentMethodCARD, nil
	case models.PaymentMethodSBP:
		return dbenums.PaymentMethodSBP, nil
	case models.PaymentMethodCREDITCARD:
		return dbenums.PaymentMethodCREDIT_CARD, nil
	case models.PaymentMethodINVESTORMONEY:
		return dbenums.PaymentMethodINVESTOR_MONEY, nil
	default:
		return "", fmt.Errorf("%w: %d", errUnknownPaymentMethod, method)
	}
}

func nullablePaymentMethod(
	opt models.Option[models.PaymentMethod],
) (omitnull.Val[dbenums.PaymentMethod], error) {
	value, ok := opt.Get()
	if !ok {
		var out omitnull.Val[dbenums.PaymentMethod]
		out.Null()
		return out, nil
	}

	method, err := paymentMethodToDB(value)
	if err != nil {
		return omitnull.Val[dbenums.PaymentMethod]{}, err
	}

	return omitnull.From(method), nil
}

func nullableTransactionID(
	opt models.Option[models.TransactionId],
) omitnull.Val[uuid.UUID] {
	value, ok := opt.Get()
	if !ok {
		var out omitnull.Val[uuid.UUID]
		out.Null()
		return out
	}

	return omitnull.From(uuid.UUID(value))
}

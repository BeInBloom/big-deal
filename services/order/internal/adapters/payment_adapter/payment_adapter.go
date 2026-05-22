package paymentadapter

import (
	"context"

	payment "github.com/BeInBloom/big-deal/generated/go/payment/v1"
	"github.com/BeInBloom/big-deal/services/order/internal/models"
	"github.com/google/uuid"
)

type PaymentAdapter struct {
	client payment.PaymentServiceClient
}

func New(client payment.PaymentServiceClient) *PaymentAdapter {
	return &PaymentAdapter{client: client}
}

func (a *PaymentAdapter) PayOrder(
	ctx context.Context,
	userId models.UserId,
	orderId models.OrderId,
	method models.PaymentMethod,
) (models.TransactionId, error) {
	req := buildReq(userId, orderId, method)

	res, err := a.client.PayOrder(ctx, req)
	if err != nil {
		return models.TransactionId{}, err
	}

	transactionId, err := rawTransactionIdToTransactionId(res.GetTransactionUuid())
	if err != nil {
		return models.TransactionId{}, err
	}

	return transactionId, nil
}

func buildReq(
	userId models.UserId,
	orderId models.OrderId,
	method models.PaymentMethod,
) *payment.PayOrderRequest {
	return &payment.PayOrderRequest{
		OrderUuid:     uuid.UUID(orderId).String(),
		UserUuid:      uuid.UUID(userId).String(),
		PaymentMethod: paymentMethodToPaymentMethod(method),
	}
}

func rawTransactionIdToTransactionId(rawTransactionId string) (models.TransactionId, error) {
	transactionUUID, err := uuid.Parse(rawTransactionId)
	if err != nil {
		return models.TransactionId{}, err
	}

	return models.TransactionId(transactionUUID), nil
}

func paymentMethodToPaymentMethod(method models.PaymentMethod) payment.PaymentMethod {
	switch method {
	case models.PaymentMethodCARD:
		return payment.PaymentMethod_PAYMENT_METHOD_CARD
	case models.PaymentMethodSBP:
		return payment.PaymentMethod_PAYMENT_METHOD_SBP
	case models.PaymentMethodCREDITCARD:
		return payment.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD
	case models.PaymentMethodINVESTORMONEY:
		return payment.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY

	default:
		return payment.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED
	}
}

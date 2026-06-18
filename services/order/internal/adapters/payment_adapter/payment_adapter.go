package paymentadapter

import (
	"context"
	"time"

	payment "github.com/BeInBloom/big-deal/generated/go/payment/v1"
	"github.com/BeInBloom/big-deal/services/order/internal/models"
)

type PaymentAdapter struct {
	client  payment.PaymentServiceClient
	timeout time.Duration
}

func New(client payment.PaymentServiceClient, timeout time.Duration) *PaymentAdapter {
	return &PaymentAdapter{
		client:  client,
		timeout: timeout,
	}
}

func (a *PaymentAdapter) PayOrder(
	ctx context.Context,
	userId models.UserId,
	orderId models.OrderId,
	method models.PaymentMethod,
) (models.TransactionId, error) {
	ctx, cancel := context.WithTimeout(ctx, a.timeout)
	defer cancel()

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

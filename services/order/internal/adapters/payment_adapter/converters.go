package paymentadapter

import (
	"fmt"

	payment "github.com/BeInBloom/big-deal/generated/go/payment/v1"
	"github.com/BeInBloom/big-deal/services/order/internal/models"
	"github.com/BeInBloom/big-deal/services/order/internal/services"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

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

func payOrderErrorToError(err error) error {
	if status.Code(err) == codes.AlreadyExists {
		return fmt.Errorf("%w: %w", services.ErrPaymentConflict, err)
	}

	return err
}

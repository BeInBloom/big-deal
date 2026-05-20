package builder

import (
	"errors"

	orders "github.com/BeInBloom/big-deal/generated/go/order-service/openapi/v1"
)

var ErrOrderServiceDependenciesNotConfigured = errors.New("order service dependencies are not configured")

func Build() (*orders.Server, error) {
	return nil, ErrOrderServiceDependenciesNotConfigured
}

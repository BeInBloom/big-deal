package builder

import (
	orders "github.com/BeInBloom/big-deal/generated/go/order-service/openapi/v1"
	"github.com/BeInBloom/big-deal/services/order/internal/handlers"
)

func Build() (*orders.Server, error) {
	h := handlers.New()
	return orders.NewServer(h)
}

package builder

import (
	inventory "github.com/BeInBloom/big-deal/generated/go/inventory/v1"
	orders "github.com/BeInBloom/big-deal/generated/go/order-service/openapi/v1"
	payment "github.com/BeInBloom/big-deal/generated/go/payment/v1"
	partsadapter "github.com/BeInBloom/big-deal/services/order/internal/adapters/parts_adapter"
	paymentadapter "github.com/BeInBloom/big-deal/services/order/internal/adapters/payment_adapter"
	"github.com/BeInBloom/big-deal/services/order/internal/handlers"
	"github.com/BeInBloom/big-deal/services/order/internal/repo"
	"github.com/BeInBloom/big-deal/services/order/internal/services"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	inventoryAddr = "localhost:50052"
	paymentAddr   = "localhost:50051"
)

func Build() (*orders.Server, error) {
	invConn, err := grpc.NewClient(
		inventoryAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	pConn, err := grpc.NewClient(
		paymentAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	orderRepo := repo.NewMapRepo()
	parts := partsadapter.New(inventory.NewInventoryServiceClient(invConn))
	payments := paymentadapter.New(payment.NewPaymentServiceClient(pConn))
	orderService := services.New(orderRepo, parts, payments)
	orderHandlers := handlers.New(orderService)

	return orders.NewServer(orderHandlers)
}

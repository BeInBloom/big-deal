package builder

import (
	"context"
	"fmt"
	"net/http"
	"time"

	inventory "github.com/BeInBloom/big-deal/generated/go/inventory/v1"
	orders "github.com/BeInBloom/big-deal/generated/go/order-service/openapi/v1"
	payment "github.com/BeInBloom/big-deal/generated/go/payment/v1"
	partsadapter "github.com/BeInBloom/big-deal/services/order/internal/adapters/parts_adapter"
	paymentadapter "github.com/BeInBloom/big-deal/services/order/internal/adapters/payment_adapter"
	"github.com/BeInBloom/big-deal/services/order/internal/config"
	"github.com/BeInBloom/big-deal/services/order/internal/handlers"
	"github.com/BeInBloom/big-deal/services/order/internal/models"
	maprepo "github.com/BeInBloom/big-deal/services/order/internal/repo/map"
	"github.com/BeInBloom/big-deal/services/order/internal/repo/psql"
	"github.com/BeInBloom/big-deal/services/order/internal/services"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const psqlConnectTimeout = time.Second * 1

type OrderRepo interface {
	Get(ctx context.Context, id models.OrderId) (models.OrderSnapshot, error)
	Save(ctx context.Context, order models.OrderSnapshot) error
}

func Build(cfg config.Config) (*http.Server, error) {
	invConn, err := grpc.NewClient(
		cfg.Inventory().GRPCAddr(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	pConn, err := grpc.NewClient(
		cfg.Payment().GRPCAddr(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	orderRepo, err := orderRepoFactory(cfg.Repo())
	if err != nil {
		return nil, err
	}

	parts := partsadapter.New(
		inventory.NewInventoryServiceClient(invConn),
		cfg.Inventory().Timeout(),
	)
	payments := paymentadapter.New(
		payment.NewPaymentServiceClient(pConn),
		cfg.Payment().Timeout(),
	)
	orderService := services.New(orderRepo, parts, payments)
	orderHandlers := handlers.New(orderService)

	orderServer, err := orders.NewServer(orderHandlers)
	if err != nil {
		return nil, err
	}

	return &http.Server{
		Addr:              cfg.HTTP().Addr(),
		Handler:           orderServer,
		ReadHeaderTimeout: cfg.HTTP().ReadHeaderTimeout(),
	}, nil
}

func orderRepoFactory(cfg config.Repo) (OrderRepo, error) {
	if cfg.IsInMemory() {
		return maprepo.New(), nil
	}

	return psqlRepo(cfg)
}

func psqlRepo(cfg config.Repo) (OrderRepo, error) {
	postgresURL, ok := cfg.PostgresURL()
	if !ok {
		return nil, fmt.Errorf("postgres url is required")
	}

	ctx, cancel := context.WithTimeout(context.Background(), psqlConnectTimeout)
	defer cancel()

	pool, err := pgxpool.New(ctx, postgresURL.String())
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}

	return psql.New(pool), nil
}

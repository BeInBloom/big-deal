package builder

import (
	"context"
	"errors"
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
	maprepo "github.com/BeInBloom/big-deal/services/order/internal/repo/map"
	"github.com/BeInBloom/big-deal/services/order/internal/repo/psql"
	"github.com/BeInBloom/big-deal/services/order/internal/services"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const psqlConnectTimeout = time.Second * 1

func Build(cfg config.Config) (*http.Server, func() error, error) {
	var resources closers

	invConn, err := grpc.NewClient(
		cfg.Inventory().GRPCAddr(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, nil, err
	}
	resources.add(invConn.Close)

	pConn, err := grpc.NewClient(
		cfg.Payment().GRPCAddr(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, nil, errors.Join(err, resources.close())
	}
	resources.add(pConn.Close)

	orderRepo, closeRepo, err := orderRepoFactory(cfg.Repo())
	if err != nil {
		return nil, nil, errors.Join(err, resources.close())
	}
	resources.add(closeRepo)

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
		return nil, nil, errors.Join(err, resources.close())
	}

	return &http.Server{
		Addr:              cfg.HTTP().Addr(),
		Handler:           orderServer,
		ReadHeaderTimeout: cfg.HTTP().ReadHeaderTimeout(),
	}, resources.close, nil
}

func orderRepoFactory(cfg config.Repo) (services.OrderRepo, func() error, error) {
	if cfg.IsInMemory() {
		return maprepo.New(), func() error { return nil }, nil
	}

	return psqlRepo(cfg)
}

func psqlRepo(cfg config.Repo) (services.OrderRepo, func() error, error) {
	postgresURL, ok := cfg.PostgresURL()
	if !ok {
		return nil, nil, fmt.Errorf("postgres url is required")
	}

	ctx, cancel := context.WithTimeout(context.Background(), psqlConnectTimeout)
	defer cancel()

	pool, err := pgxpool.New(ctx, postgresURL.String())
	if err != nil {
		return nil, nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, nil, err
	}

	closePool := func() error {
		pool.Close()
		return nil
	}

	return psql.New(pool), closePool, nil
}

type closers []func() error

func (c *closers) add(closeFn func() error) {
	*c = append(*c, closeFn)
}

func (c closers) close() error {
	errs := make([]error, 0, len(c))

	for i := len(c) - 1; i >= 0; i-- {
		errs = append(errs, c[i]())
	}

	return errors.Join(errs...)
}

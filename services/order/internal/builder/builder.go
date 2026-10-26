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

type builder struct {
	cfg     config.Config
	closers []func() error
}

func Build(cfg config.Config) (*http.Server, func() error, error) {
	b := &builder{cfg: cfg}

	server, err := b.httpServer()
	if err != nil {
		return nil, nil, errors.Join(err, b.close())
	}

	return server, b.close, nil
}

func (b *builder) httpServer() (*http.Server, error) {
	orderService, err := b.orderService()
	if err != nil {
		return nil, err
	}

	handler, err := orders.NewServer(handlers.New(orderService))
	if err != nil {
		return nil, err
	}

	return &http.Server{
		Addr:              b.cfg.HTTP().Addr(),
		Handler:           handler,
		ReadHeaderTimeout: b.cfg.HTTP().ReadHeaderTimeout(),
	}, nil
}

func (b *builder) orderService() (*services.OrderService, error) {
	repo, err := b.orderRepo()
	if err != nil {
		return nil, err
	}

	parts, err := b.partService()
	if err != nil {
		return nil, err
	}

	payments, err := b.paymentService()
	if err != nil {
		return nil, err
	}

	return services.New(repo, parts, payments), nil
}

func (b *builder) orderRepo() (services.OrderRepo, error) {
	if b.cfg.Repo().IsInMemory() {
		return maprepo.New(), nil
	}

	return b.psqlRepo()
}

func (b *builder) psqlRepo() (services.OrderRepo, error) {
	postgresURL, ok := b.cfg.Repo().PostgresURL()
	if !ok {
		return nil, fmt.Errorf("postgres url is required")
	}

	ctx, cancel := context.WithTimeout(context.Background(), psqlConnectTimeout)
	defer cancel()

	pool, err := pgxpool.New(ctx, postgresURL.String())
	if err != nil {
		return nil, err
	}

	b.onClose(func() error {
		pool.Close()
		return nil
	})

	if err := pool.Ping(ctx); err != nil {
		return nil, err
	}

	return psql.New(pool), nil
}

func (b *builder) partService() (services.PartService, error) {
	conn, err := b.grpcConn(b.cfg.Inventory().GRPCAddr())
	if err != nil {
		return nil, err
	}

	return partsadapter.New(
		inventory.NewInventoryServiceClient(conn),
		b.cfg.Inventory().Timeout(),
	), nil
}

func (b *builder) paymentService() (services.PaymentService, error) {
	conn, err := b.grpcConn(b.cfg.Payment().GRPCAddr())
	if err != nil {
		return nil, err
	}

	return paymentadapter.New(
		payment.NewPaymentServiceClient(conn),
		b.cfg.Payment().Timeout(),
	), nil
}

func (b *builder) grpcConn(addr string) (*grpc.ClientConn, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	b.onClose(conn.Close)

	return conn, nil
}

func (b *builder) onClose(closeFn func() error) {
	b.closers = append(b.closers, closeFn)
}

func (b *builder) close() error {
	errs := make([]error, 0, len(b.closers))

	for i := len(b.closers) - 1; i >= 0; i-- {
		errs = append(errs, b.closers[i]())
	}

	return errors.Join(errs...)
}

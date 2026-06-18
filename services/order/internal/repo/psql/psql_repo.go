package psql

import (
	"context"
	"database/sql"
	"errors"

	"github.com/BeInBloom/big-deal/services/order/internal/models"
	dbmodels "github.com/BeInBloom/big-deal/services/order/internal/repo/psql/gen/models"
	"github.com/BeInBloom/big-deal/services/order/internal/services"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stephenafamo/bob"
	bobpgx "github.com/stephenafamo/bob/drivers/pgx"
)

type repo struct {
	pool bobpgx.Pool
}

func New(pool *pgxpool.Pool) *repo {
	return &repo{pool: bobpgx.NewPool(pool)}
}

func (r *repo) Get(
	ctx context.Context,
	id models.OrderId,
) (models.OrderSnapshot, error) {
	order, err := dbmodels.FindOrder(ctx, r.pool, uuid.UUID(id))
	if errors.Is(err, sql.ErrNoRows) {
		return models.OrderSnapshot{}, services.ErrOrderNotFound
	}
	if err != nil {
		return models.OrderSnapshot{}, err
	}

	if err := order.LoadOrderParts(ctx, r.pool); err != nil {
		return models.OrderSnapshot{}, err
	}

	return tryIntoSnapshot(order)
}

func (r *repo) Save(
	ctx context.Context,
	order models.OrderSnapshot,
) error {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel:   pgx.ReadCommitted,
		AccessMode: pgx.ReadWrite,
	})
	if err != nil {
		return err
	}

	if err := r.saveOrder(ctx, tx, order); err != nil {
		_ = tx.Rollback(ctx)
		return err
	}

	if err := r.saveOrderParts(ctx, tx, order); err != nil {
		_ = tx.Rollback(ctx)
		return err
	}

	return tx.Commit(ctx)
}

func (r *repo) saveOrder(
	ctx context.Context,
	exec bob.Executor,
	order models.OrderSnapshot,
) error {
	setter, err := intoOrderSetter(order)
	if err != nil {
		return err
	}

	_, err = dbmodels.Orders.Insert(setter).One(ctx, exec)
	return err
}

func (r *repo) saveOrderParts(
	ctx context.Context,
	exec bob.Executor,
	order models.OrderSnapshot,
) error {
	setter, err := intoOrderPartSetters(order)
	if err != nil {
		return err
	}
	if len(setter) == 0 {
		return nil
	}

	_, err = dbmodels.OrderParts.Insert(
		bob.ToMods(setter...),
	).All(ctx, exec)
	return err
}

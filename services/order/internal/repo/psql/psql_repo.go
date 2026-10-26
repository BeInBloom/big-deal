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
	"github.com/stephenafamo/bob/dialect/psql"
	"github.com/stephenafamo/bob/dialect/psql/um"
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

func (r *repo) Create(
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
		return errors.Join(intoCreateOrderError(err), tx.Rollback(ctx))
	}

	if err := r.saveOrderParts(ctx, tx, order); err != nil {
		return errors.Join(err, tx.Rollback(ctx))
	}

	return tx.Commit(ctx)
}

// Parts are immutable after creation, so only the order row is updated.
func (r *repo) Update(
	ctx context.Context,
	order models.OrderSnapshot,
	from models.OrderStatus,
) error {
	setter, err := intoOrderStateSetter(order)
	if err != nil {
		return err
	}

	fromStatus, err := orderStatusToDB(from)
	if err != nil {
		return err
	}

	affected, err := dbmodels.Orders.Update(
		setter.UpdateMod(),
		um.SetCol(dbmodels.Orders.Columns.UpdatedAt.Name()).To(psql.Raw("now()")),
		dbmodels.UpdateWhere.Orders.ID.EQ(uuid.UUID(order.Id)),
		dbmodels.UpdateWhere.Orders.Status.EQ(fromStatus),
	).Exec(ctx, r.pool)
	if err != nil {
		return err
	}

	if affected == 0 {
		return services.ErrOrderStatusConflict
	}

	return nil
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

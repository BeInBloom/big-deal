package psql

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/BeInBloom/big-deal/services/order/internal/models"
	"github.com/BeInBloom/big-deal/services/order/internal/services"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib" // registers the pgx driver for goose
	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

const (
	postgresImage = "postgres:17-alpine"
	migrationsDir = "../../../migrations"
)

// testDB is a single Postgres container shared by all tests in the package.
// It starts on first use, so tests skipped by -short never need Docker.
var testDB struct {
	once      sync.Once
	pool      *pgxpool.Pool
	container *postgres.PostgresContainer
	err       error
}

func TestMain(m *testing.M) {
	code := m.Run()

	if testDB.pool != nil {
		testDB.pool.Close()
	}

	if testDB.container != nil {
		_ = testcontainers.TerminateContainer(testDB.container)
	}

	os.Exit(code)
}

func TestRepoCreateAndGet(t *testing.T) {
	ctx := context.Background()
	r := newTestRepo(t)

	order := newPendingSnapshot()
	require.NoError(t, r.Create(ctx, order))

	got, err := r.Get(ctx, order.Id)
	require.NoError(t, err)

	require.Equal(t, order.Id, got.Id)
	require.Equal(t, order.UserId, got.UserId)
	require.Equal(t, models.PendingPayment, got.Status)
	require.ElementsMatch(t, order.Parts, got.Parts)
	require.False(t, got.PaymentMethod.IsSet())
	require.False(t, got.TransactionId.IsSet())
}

func TestRepoCreateKeepsLargePrices(t *testing.T) {
	ctx := context.Background()
	r := newTestRepo(t)

	order := newPendingSnapshot()
	order.Parts[0].Price = 1 << 40
	require.NoError(t, r.Create(ctx, order))

	got, err := r.Get(ctx, order.Id)
	require.NoError(t, err)
	require.Equal(t, order.Parts.Price(), got.Parts.Price())
}

func TestRepoCreateReturnsAlreadyExists(t *testing.T) {
	ctx := context.Background()
	r := newTestRepo(t)

	order := newPendingSnapshot()
	require.NoError(t, r.Create(ctx, order))

	err := r.Create(ctx, order)
	require.ErrorIs(t, err, services.ErrOrderAlreadyExists)
}

func TestRepoGetReturnsNotFound(t *testing.T) {
	r := newTestRepo(t)

	_, err := r.Get(context.Background(), models.OrderId(uuid.New()))
	require.ErrorIs(t, err, services.ErrOrderNotFound)
}

func TestRepoUpdatePersistsPayment(t *testing.T) {
	ctx := context.Background()
	r := newTestRepo(t)

	order := newPendingSnapshot()
	require.NoError(t, r.Create(ctx, order))

	paid := newPaidSnapshot(order)
	require.NoError(t, r.Update(ctx, paid, models.PendingPayment))

	got, err := r.Get(ctx, order.Id)
	require.NoError(t, err)
	require.Equal(t, paid.Status, got.Status)
	require.Equal(t, paid.PaymentMethod, got.PaymentMethod)
	require.Equal(t, paid.TransactionId, got.TransactionId)
	require.ElementsMatch(t, order.Parts, got.Parts)
}

func TestRepoUpdateTouchesUpdatedAt(t *testing.T) {
	ctx := context.Background()
	r := newTestRepo(t)

	order := newPendingSnapshot()
	require.NoError(t, r.Create(ctx, order))

	createdAt := readUpdatedAt(t, order.Id)

	canceled := order
	canceled.Status = models.Canceled
	require.NoError(t, r.Update(ctx, canceled, models.PendingPayment))

	require.True(t, readUpdatedAt(t, order.Id).After(createdAt))
}

func TestRepoUpdateReturnsConflictOnUnexpectedStatus(t *testing.T) {
	ctx := context.Background()
	r := newTestRepo(t)

	order := newPendingSnapshot()
	require.NoError(t, r.Create(ctx, order))

	canceled := order
	canceled.Status = models.Canceled
	require.NoError(t, r.Update(ctx, canceled, models.PendingPayment))

	err := r.Update(ctx, newPaidSnapshot(order), models.PendingPayment)
	require.ErrorIs(t, err, services.ErrOrderStatusConflict)

	got, err := r.Get(ctx, order.Id)
	require.NoError(t, err)
	require.Equal(t, models.Canceled, got.Status)
}

func TestRepoUpdateReturnsConflictForMissingOrder(t *testing.T) {
	r := newTestRepo(t)

	err := r.Update(context.Background(), newPaidSnapshot(newPendingSnapshot()), models.PendingPayment)
	require.ErrorIs(t, err, services.ErrOrderStatusConflict)
}

func TestRepoUpdateRejectsPaidOrderWithoutPayment(t *testing.T) {
	ctx := context.Background()
	r := newTestRepo(t)

	order := newPendingSnapshot()
	require.NoError(t, r.Create(ctx, order))

	paid := order
	paid.Status = models.Paid

	// orders_payment_state_check guards the invariant at the database level.
	err := r.Update(ctx, paid, models.PendingPayment)
	require.Error(t, err)
	require.NotErrorIs(t, err, services.ErrOrderStatusConflict)
}

func newTestRepo(t *testing.T) *repo {
	t.Helper()

	if testing.Short() {
		t.Skip("integration test: needs Docker")
	}

	testcontainers.SkipIfProviderIsNotHealthy(t)

	testDB.once.Do(startTestDB)
	require.NoError(t, testDB.err)

	return New(testDB.pool)
}

func startTestDB() {
	ctx := context.Background()

	container, err := postgres.Run(ctx, postgresImage,
		postgres.WithDatabase("orders"),
		postgres.WithUsername("order"),
		postgres.WithPassword("secret"),
		postgres.BasicWaitStrategies(),
	)
	testDB.container = container
	if err != nil {
		testDB.err = err
		return
	}

	dsn, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		testDB.err = err
		return
	}

	if err := migrate(ctx, dsn); err != nil {
		testDB.err = err
		return
	}

	testDB.pool, testDB.err = pgxpool.New(ctx, dsn)
}

func migrate(ctx context.Context, dsn string) error {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return err
	}

	provider, err := goose.NewProvider(goose.DialectPostgres, db, os.DirFS(migrationsDir))
	if err != nil {
		return errors.Join(err, db.Close())
	}

	_, err = provider.Up(ctx)
	return errors.Join(err, db.Close())
}

func readUpdatedAt(t *testing.T, id models.OrderId) time.Time {
	t.Helper()

	var updatedAt time.Time
	err := testDB.pool.QueryRow(
		context.Background(),
		"SELECT updated_at FROM orders WHERE id = $1",
		uuid.UUID(id),
	).Scan(&updatedAt)
	require.NoError(t, err)

	return updatedAt
}

func newPendingSnapshot() models.OrderSnapshot {
	return models.OrderSnapshot{
		Id:     models.OrderId(uuid.New()),
		UserId: models.UserId(uuid.New()),
		Status: models.PendingPayment,
		Parts: models.Parts{
			{
				Id:          models.PartId(uuid.New()),
				Description: "engine",
				Price:       1500,
			},
			{
				Id:          models.PartId(uuid.New()),
				Description: "wing",
				Price:       700,
			},
		},
	}
}

func newPaidSnapshot(order models.OrderSnapshot) models.OrderSnapshot {
	order.Status = models.Paid
	order.PaymentMethod = models.Some(models.PaymentMethodSBP)
	order.TransactionId = models.Some(models.TransactionId(uuid.New()))

	return order
}

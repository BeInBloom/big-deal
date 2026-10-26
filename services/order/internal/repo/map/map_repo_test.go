package maprepo

import (
	"context"
	"testing"

	"github.com/BeInBloom/big-deal/services/order/internal/models"
	"github.com/BeInBloom/big-deal/services/order/internal/services"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestMapRepoCreateAndGet(t *testing.T) {
	ctx := context.Background()

	orderID := models.OrderId(uuid.New())
	userID := models.UserId(uuid.New())
	partID := models.PartId(uuid.New())

	snapshot := models.OrderSnapshot{
		Id:     orderID,
		UserId: userID,
		Status: models.PendingPayment,
		Parts: models.Parts{
			{
				Id:          partID,
				Description: "engine",
				Price:       1500,
			},
		},
	}

	repo := New()

	err := repo.Create(ctx, snapshot)
	require.NoError(t, err)

	got, err := repo.Get(ctx, orderID)
	require.NoError(t, err)

	require.Equal(t, snapshot.Id, got.Id)
	require.Equal(t, snapshot.UserId, got.UserId)
	require.Equal(t, snapshot.Status, got.Status)
	require.Equal(t, snapshot.Parts.Price(), got.Parts.Price())
}

func TestMapRepoGetReturnsNotFound(t *testing.T) {
	ctx := context.Background()

	repo := New()
	orderId := models.OrderId(uuid.New())

	_, err := repo.Get(ctx, orderId)

	require.ErrorIs(t, err, services.ErrOrderNotFound)
}

func TestMapRepoGetReturnsSnapshotClone(t *testing.T) {
	ctx := context.Background()

	orderId := models.OrderId(uuid.New())
	userId := models.UserId(uuid.New())
	partId := models.PartId(uuid.New())

	snapshot := models.OrderSnapshot{
		Id:     orderId,
		UserId: userId,
		Status: models.PendingPayment,
		Parts: models.Parts{
			{
				Id:          partId,
				Description: "engine",
				Price:       1500,
			},
		},
	}

	repo := New()

	err := repo.Create(ctx, snapshot)
	require.NoError(t, err)

	got, err := repo.Get(ctx, orderId)
	require.NoError(t, err)

	got.Parts[0].Price = 9999

	gotAgain, err := repo.Get(ctx, orderId)
	require.NoError(t, err)

	require.Equal(t, uint(1500), gotAgain.Parts[0].Price)
}

func TestMapRepoCreateRejectsExistingOrder(t *testing.T) {
	ctx := context.Background()

	repo := New()
	snapshot := newPendingSnapshot()

	require.NoError(t, repo.Create(ctx, snapshot))
	require.Error(t, repo.Create(ctx, snapshot))
}

func TestMapRepoUpdateChangesOrderWithExpectedStatus(t *testing.T) {
	ctx := context.Background()

	repo := New()
	snapshot := newPendingSnapshot()
	require.NoError(t, repo.Create(ctx, snapshot))

	canceled := snapshot
	canceled.Status = models.Canceled

	err := repo.Update(ctx, canceled, models.PendingPayment)
	require.NoError(t, err)

	got, err := repo.Get(ctx, snapshot.Id)
	require.NoError(t, err)
	require.Equal(t, models.Canceled, got.Status)
}

func TestMapRepoUpdateReturnsConflictOnUnexpectedStatus(t *testing.T) {
	ctx := context.Background()

	repo := New()
	snapshot := newPendingSnapshot()
	require.NoError(t, repo.Create(ctx, snapshot))

	canceled := snapshot
	canceled.Status = models.Canceled
	require.NoError(t, repo.Update(ctx, canceled, models.PendingPayment))

	paid := snapshot
	paid.Status = models.Paid
	paid.PaymentMethod = models.Some(models.PaymentMethodCARD)
	paid.TransactionId = models.Some(models.TransactionId(uuid.New()))

	err := repo.Update(ctx, paid, models.PendingPayment)
	require.ErrorIs(t, err, services.ErrOrderStatusConflict)

	got, err := repo.Get(ctx, snapshot.Id)
	require.NoError(t, err)
	require.Equal(t, models.Canceled, got.Status)
}

func TestMapRepoUpdateReturnsNotFound(t *testing.T) {
	ctx := context.Background()

	repo := New()

	err := repo.Update(ctx, newPendingSnapshot(), models.PendingPayment)
	require.ErrorIs(t, err, services.ErrOrderNotFound)
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
		},
	}
}

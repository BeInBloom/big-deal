package repo

import (
	"context"
	"testing"

	"github.com/BeInBloom/big-deal/services/order/internal/models"
	"github.com/BeInBloom/big-deal/services/order/internal/services"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestMapRepoSaveAndGet(t *testing.T) {
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

	repo := NewMapRepo()

	err := repo.Save(ctx, snapshot)
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

	repo := NewMapRepo()
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

	repo := NewMapRepo()

	err := repo.Save(ctx, snapshot)
	require.NoError(t, err)

	got, err := repo.Get(ctx, orderId)
	require.NoError(t, err)

	got.Parts[0].Price = 9999

	gotAgain, err := repo.Get(ctx, orderId)
	require.NoError(t, err)

	require.Equal(t, uint(1500), gotAgain.Parts[0].Price)
}

package partsadapter

import (
	"testing"

	inventory "github.com/BeInBloom/big-deal/generated/go/inventory/v1"
	"github.com/BeInBloom/big-deal/services/order/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestRawPartToPart(t *testing.T) {
	partUUID := uuid.New()

	part, err := rawPartToPart(&inventory.InventoryPart{
		Uuid:        partUUID.String(),
		Description: "engine",
		Price:       12.345,
	})

	require.NoError(t, err)
	require.Equal(t, models.PartId(partUUID), part.Id)
	require.Equal(t, "engine", part.Description)
	require.Equal(t, models.Money(1235), part.Price)
}

func TestRawPartToPartReturnsError(t *testing.T) {
	part, err := rawPartToPart(&inventory.InventoryPart{
		Uuid: "not-a-uuid",
	})

	require.Error(t, err)
	require.Zero(t, part)
}

func TestRawPartToPartRejectsNegativePrice(t *testing.T) {
	part, err := rawPartToPart(&inventory.InventoryPart{
		Uuid:  uuid.NewString(),
		Price: -1,
	})

	require.ErrorIs(t, err, models.ErrInvalidMoney)
	require.Zero(t, part)
}

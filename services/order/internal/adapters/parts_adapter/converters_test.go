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
	require.Equal(t, uint(1235), part.Price)
}

func TestRawPartToPartReturnsError(t *testing.T) {
	part, err := rawPartToPart(&inventory.InventoryPart{
		Uuid: "not-a-uuid",
	})

	require.Error(t, err)
	require.Zero(t, part)
}

func TestFloatPriceToUintPrice(t *testing.T) {
	tests := []struct {
		name  string
		price float64
		want  uint
	}{
		{
			name:  "whole",
			price: 10,
			want:  1000,
		},
		{
			name:  "cents",
			price: 10.25,
			want:  1025,
		},
		{
			name:  "rounds half up",
			price: 10.235,
			want:  1024,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, floatPriceToUintPrice(tt.price))
		})
	}
}

func TestPartIdsToStrings(t *testing.T) {
	firstUUID := uuid.New()
	secondUUID := uuid.New()

	ids := partIdsToStrings([]models.PartId{
		models.PartId(firstUUID),
		models.PartId(secondUUID),
	})

	require.Equal(t, []string{firstUUID.String(), secondUUID.String()}, ids)
}

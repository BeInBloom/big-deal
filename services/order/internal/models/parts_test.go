package models

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestPartsIds(t *testing.T) {
	firstPartId := PartId(uuid.New())
	secondPartId := PartId(uuid.New())

	parts := Parts{
		{
			Id:          firstPartId,
			Description: "engine",
			Price:       1500,
		},
		{
			Id:          secondPartId,
			Description: "wing",
			Price:       2500,
		},
	}

	ids := parts.Ids()

	require.Equal(t, []PartId{firstPartId, secondPartId}, ids)
}

func TestPartsIdsReturnsEmptySliceForEmptyParts(t *testing.T) {
	ids := Parts{}.Ids()
	require.Empty(t, ids)
}

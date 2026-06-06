package partsadapter

import (
	"testing"

	"github.com/BeInBloom/big-deal/services/order/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestBuildReq(t *testing.T) {
	firstUUID := uuid.New()
	secondUUID := uuid.New()

	req := buildReq([]models.PartId{
		models.PartId(firstUUID),
		models.PartId(secondUUID),
	})

	require.Equal(t, []string{firstUUID.String(), secondUUID.String()}, req.GetFilter().GetUuids())
}

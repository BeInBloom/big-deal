package partsadapter

import (
	inventory "github.com/BeInBloom/big-deal/generated/go/inventory/v1"
	"github.com/BeInBloom/big-deal/services/order/internal/models"
)

func buildReq(ids []models.PartId) *inventory.ListPartsRequest {
	return &inventory.ListPartsRequest{
		Filter: &inventory.InventoryPartsFilter{
			Uuids: partIdsToStrings(ids),
		},
	}
}

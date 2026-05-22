package partsadapter

import (
	"context"
	"math"

	inventory "github.com/BeInBloom/big-deal/generated/go/inventory/v1"
	"github.com/BeInBloom/big-deal/services/order/internal/models"
	"github.com/google/uuid"
)

type PartAdapter struct {
	client inventory.InventoryServiceClient
}

func New(client inventory.InventoryServiceClient) *PartAdapter {
	return &PartAdapter{client: client}
}

func (a *PartAdapter) ListParts(
	ctx context.Context,
	partIds []models.PartId,
) (models.Parts, error) {
	rawParts, err := a.getRawListParts(ctx, partIds)
	if err != nil {
		return nil, err
	}

	parts, err := rawPartsToParts(rawParts)
	if err != nil {
		return nil, err
	}

	return parts, nil
}

func (a *PartAdapter) getRawListParts(
	ctx context.Context,
	parts []models.PartId,
) ([]*inventory.InventoryPart, error) {
	req := buildReq(parts)

	res, err := a.client.ListParts(ctx, req)
	if err != nil {
		return nil, err
	}

	return res.GetParts(), nil
}

func rawPartsToParts(rawParts []*inventory.InventoryPart) (models.Parts, error) {
	parts := make([]models.Part, 0, len(rawParts))

	for _, rawPart := range rawParts {
		part, err := rawPartToPart(rawPart)
		if err != nil {
			return nil, err
		}

		parts = append(parts, part)
	}

	return parts, nil
}

func rawPartToPart(rawPart *inventory.InventoryPart) (models.Part, error) {
	partUUID, err := uuid.Parse(rawPart.GetUuid())
	if err != nil {
		return models.Part{}, err
	}

	price := floatPriceToUintPrice(rawPart.GetPrice())

	return models.Part{
		Id:          models.PartId(partUUID),
		Description: rawPart.GetDescription(),
		Price:       price,
	}, nil
}

func floatPriceToUintPrice(price float64) uint {
	return uint(math.Round(price * 100))
}

func buildReq(ids []models.PartId) *inventory.ListPartsRequest {
	return &inventory.ListPartsRequest{
		Filter: &inventory.InventoryPartsFilter{
			Uuids: partIdsToStrings(ids),
		},
	}
}

func partIdsToStrings(ids []models.PartId) []string {
	stringUids := make([]string, 0, len(ids))

	for _, id := range ids {
		stringUids = append(stringUids, uuid.UUID(id).String())
	}

	return stringUids
}

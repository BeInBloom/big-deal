package partsadapter

import (
	"context"
	"time"

	inventory "github.com/BeInBloom/big-deal/generated/go/inventory/v1"
	"github.com/BeInBloom/big-deal/services/order/internal/models"
)

type PartAdapter struct {
	client  inventory.InventoryServiceClient
	timeout time.Duration
}

func New(client inventory.InventoryServiceClient, timeout time.Duration) *PartAdapter {
	return &PartAdapter{
		client:  client,
		timeout: timeout,
	}
}

func (a *PartAdapter) ListParts(
	ctx context.Context,
	partIds []models.PartId,
) (models.Parts, error) {
	ctx, cancel := context.WithTimeout(ctx, a.timeout)
	defer cancel()

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

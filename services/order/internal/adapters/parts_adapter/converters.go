package partsadapter

import (
	"fmt"

	inventory "github.com/BeInBloom/big-deal/generated/go/inventory/v1"
	"github.com/BeInBloom/big-deal/services/order/internal/models"
	"github.com/google/uuid"
)

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

	price, err := models.MoneyFromMajor(rawPart.GetPrice())
	if err != nil {
		return models.Part{}, fmt.Errorf("part %s price %v: %w", partUUID, rawPart.GetPrice(), err)
	}

	return models.Part{
		Id:          models.PartId(partUUID),
		Description: rawPart.GetDescription(),
		Price:       price,
	}, nil
}

func partIdsToStrings(ids []models.PartId) []string {
	stringUids := make([]string, 0, len(ids))

	for _, id := range ids {
		stringUids = append(stringUids, uuid.UUID(id).String())
	}

	return stringUids
}

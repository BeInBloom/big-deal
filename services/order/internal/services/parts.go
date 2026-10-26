package services

import "github.com/BeInBloom/big-deal/services/order/internal/models"

func validatePartIds(ids []models.PartId) error {
	if len(ids) == 0 {
		return ErrOrderPartsRequired
	}

	seen := make(map[models.PartId]struct{}, len(ids))

	for _, id := range ids {
		if _, ok := seen[id]; ok {
			return ErrOrderPartsDuplicated
		}

		seen[id] = struct{}{}
	}

	return nil
}

func pickRequestedParts(ids []models.PartId, found models.Parts) (models.Parts, error) {
	byId := make(map[models.PartId]models.Part, len(found))

	for _, part := range found {
		byId[part.Id] = part
	}

	parts := make(models.Parts, 0, len(ids))

	for _, id := range ids {
		part, ok := byId[id]
		if !ok {
			return nil, ErrOrderPartsNotFound
		}

		parts = append(parts, part)
	}

	return parts, nil
}

package models

import "slices"

type (
	Parts []Part
)

type Part struct {
	Id          PartId
	Description string
	Price       uint
}

func (p Parts) Clone() Parts {
	return slices.Clone(p)
}

func (p Parts) Price() uint {
	var sum uint

	for _, part := range p {
		sum += part.Price
	}

	return sum
}

func (p Parts) Ids() []PartId {
	ids := make([]PartId, 0, len(p))

	for _, part := range p {
		ids = append(ids, part.Id)
	}

	return ids
}

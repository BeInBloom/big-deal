package models

import "slices"

type Parts []Part

type Part struct {
	Id          PartId
	Description string
	Price       Money
}

func (p Parts) Clone() Parts {
	return slices.Clone(p)
}

func (p Parts) Price() (Money, error) {
	var sum Money

	for _, part := range p {
		var err error
		if sum, err = sum.Add(part.Price); err != nil {
			return 0, err
		}
	}

	return sum, nil
}

func (p Parts) Ids() []PartId {
	ids := make([]PartId, 0, len(p))

	for _, part := range p {
		ids = append(ids, part.Id)
	}

	return ids
}

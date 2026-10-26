package models

import (
	"errors"
	"math"
)

var ErrInvalidMoney = errors.New("invalid money amount")

const minorUnitsPerMajor = 100

// Money is an amount in minor currency units (kopecks).
type Money uint64

// MoneyFromMajor converts an amount in major units (12.34) to Money,
// rounding to the nearest minor unit.
func MoneyFromMajor(amount float64) (Money, error) {
	if math.IsNaN(amount) || math.IsInf(amount, 0) || amount < 0 {
		return 0, ErrInvalidMoney
	}

	minor := math.Round(amount * minorUnitsPerMajor)
	if minor >= math.MaxUint64 {
		return 0, ErrInvalidMoney
	}

	return Money(minor), nil
}

// Major returns the amount in major units (12.34).
func (m Money) Major() float64 {
	return float64(m) / minorUnitsPerMajor
}

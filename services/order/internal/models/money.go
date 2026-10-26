package models

import (
	"errors"
	"math"
)

var ErrInvalidMoney = errors.New("invalid money amount")

const minorUnitsPerMajor = 100

// Storage keeps amounts in a signed bigint.
const maxMoney Money = math.MaxInt64

// Money is an amount in kopecks.
type Money uint64

func MoneyFromMajor(amount float64) (Money, error) {
	if math.IsNaN(amount) || math.IsInf(amount, 0) || amount < 0 {
		return 0, ErrInvalidMoney
	}

	// float64(maxMoney) rounds up to 2^63, so ">=" rejects everything above maxMoney.
	minor := math.Round(amount * minorUnitsPerMajor)
	if minor >= float64(maxMoney) {
		return 0, ErrInvalidMoney
	}

	return Money(minor), nil
}

func (m Money) Add(other Money) (Money, error) {
	if m > maxMoney || other > maxMoney-m {
		return 0, ErrInvalidMoney
	}

	return m + other, nil
}

func (m Money) Major() float64 {
	return float64(m) / minorUnitsPerMajor
}

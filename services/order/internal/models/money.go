package models

import (
	"errors"
	"math"
)

var ErrInvalidMoney = errors.New("invalid money amount")

const minorUnitsPerMajor = 100

// maxMoney keeps amounts within the signed 64-bit range used by storage.
const maxMoney Money = math.MaxInt64

// Money is a non-negative amount in minor currency units (kopecks),
// at most maxMoney.
type Money uint64

// MoneyFromMajor converts an amount in major units (12.34) to Money,
// rounding to the nearest minor unit.
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

// Add returns m + other, or ErrInvalidMoney if the sum exceeds maxMoney.
func (m Money) Add(other Money) (Money, error) {
	if m > maxMoney || other > maxMoney-m {
		return 0, ErrInvalidMoney
	}

	return m + other, nil
}

// Major returns the amount in major units (12.34).
func (m Money) Major() float64 {
	return float64(m) / minorUnitsPerMajor
}

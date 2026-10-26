package models

import (
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMoneyFromMajor(t *testing.T) {
	tests := []struct {
		name   string
		amount float64
		want   Money
	}{
		{name: "zero", amount: 0, want: 0},
		{name: "whole", amount: 15, want: 1500},
		{name: "cents", amount: 12.34, want: 1234},
		{name: "rounds half up", amount: 12.345, want: 1235},
		{name: "rounds down", amount: 12.344, want: 1234},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := MoneyFromMajor(tt.amount)
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestMoneyFromMajorRejectsInvalidAmount(t *testing.T) {
	for _, amount := range []float64{-0.01, math.NaN(), math.Inf(1), math.Inf(-1), 1e20} {
		_, err := MoneyFromMajor(amount)
		require.ErrorIs(t, err, ErrInvalidMoney, "amount %v", amount)
	}
}

func TestMoneyMajor(t *testing.T) {
	require.InEpsilon(t, 12.34, Money(1234).Major(), 1e-9)
	require.Zero(t, Money(0).Major())
}

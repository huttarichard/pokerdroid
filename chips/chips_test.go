package chips

import (
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestChipsMath(t *testing.T) {
	require.Equal(t, NewFromInt(1), Chips(1))
	require.Equal(t, NewFromFloat64(1), Chips(1))
	require.Equal(t, NewFromFloat(1), Chips(1))
	require.Equal(t, NewFromFloat32(1), Chips(1))
	require.Equal(t, NewFromString("1"), Chips(1))
	require.Equal(t, NewFromString("1.0001"), Chips(1.0001))

	require.Equal(t, Chips(1.0001).String(), "1.000100")
	require.Equal(t, Chips(1.0001).StringFixed(2), "1.00")

	require.Equal(t, Chips(1).Add(2), Chips(3))
	require.Equal(t, Chips(1).Sub(2), Chips(-1))
	require.Equal(t, Chips(1).Mul(2), Chips(2))
	require.Equal(t, Chips(1).Div(2), Chips(0.5))
	require.Equal(t, Chips(1).Div(0), Chips(math.Inf(1)))
	require.Equal(t, Chips(1).Equal(1), true)

	require.Equal(t, Chips(1).GreaterThan(1), false)
	require.Equal(t, Chips(1).GreaterThan(0), true)
	require.Equal(t, Chips(1).GreaterThanOrEqual(2), false)
	require.Equal(t, Chips(1).GreaterThanOrEqual(1), true)

	require.Equal(t, Chips(1).LessThan(2), true)
	require.Equal(t, Chips(1).LessThan(1), false)
	require.Equal(t, Chips(1).LessThanOrEqual(1), true)
	require.Equal(t, Chips(1).LessThanOrEqual(0), false)

	require.Equal(t, Chips(-1).Abs(), Chips(1))
	require.Equal(t, Chips(-1).Pow(-1), Chips(-1))
	require.Equal(t, Chips(-1).Float32(), float32(-1))
	require.Equal(t, Chips(-1).Float64(), float64(-1))
}

func TestRound(t *testing.T) {
	require.Equal(t, Chips(1.23456789).Round(2), Chips(1.23))
	require.Equal(t, Chips(1.23456789).Round(3), Chips(1.235))
	require.Equal(t, Chips(1.23456789).Round(4), Chips(1.2346))
	require.Equal(t, Chips(1.23456789).Round(5), Chips(1.23457))
	require.Equal(t, Chips(1.23456789).Round(6), Chips(1.234568))
}

func TestRoundUp(t *testing.T) {
	require.Equal(t, Chips(1.23456789).RoundUp(2), Chips(1.24))
	require.Equal(t, Chips(1.23456789).RoundUp(3), Chips(1.235))
	require.Equal(t, Chips(1.23456789).RoundUp(4), Chips(1.2346))
	require.Equal(t, Chips(1.23456789).RoundUp(5), Chips(1.23457))
	require.Equal(t, Chips(1.23456789).RoundUp(6), Chips(1.234568))
}

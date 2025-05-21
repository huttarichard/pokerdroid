package card

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRangeDist(t *testing.T) {
	t.Run("uniform distribution", func(t *testing.T) {
		r := NewUniformRangeDist()
		require.InDelta(t, float32(1.0), r.Sum(), 0.0001)

		// Check each probability is 1/1326
		expected := float32(1.0 / 1326.0)
		for i := range r {
			require.InDelta(t, expected, r[i], 0.0001)
		}
	})

	t.Run("normalize", func(t *testing.T) {
		r := RangeDist{}
		// Set some arbitrary values
		r[0] = 2.0
		r[1] = 3.0
		r[2] = 5.0

		normalized := r.Normalize()
		require.InDelta(t, float32(1.0), normalized.Sum(), 0.0001)

		// Check relative proportions are maintained
		require.InDelta(t, float32(0.2), normalized[0], 0.0001) // 2/10
		require.InDelta(t, float32(0.3), normalized[1], 0.0001) // 3/10
		require.InDelta(t, float32(0.5), normalized[2], 0.0001) // 5/10
	})
}

func TestRangeIndex(t *testing.T) {
	require.Equal(t, RangeIndex(Cards{CardAC, CardAD}), 0)
	require.Equal(t, RangeIndex(Cards{CardAC, CardAH}), 1)
	require.Equal(t, RangeIndex(Cards{CardAC, CardAS}), 2)
	require.Equal(t, RangeIndex(Cards{CardAD, CardAH}), 3)
	require.Equal(t, RangeIndex(Cards{CardAD, CardAS}), 4)
	require.Equal(t, RangeIndex(Cards{CardAH, CardAS}), 5)

	require.Equal(t, RangeIndex(Cards{Card2C, Card2D}), 1320)
	require.Equal(t, RangeIndex(Cards{Card2C, Card2H}), 1321)
	require.Equal(t, RangeIndex(Cards{Card2C, Card2S}), 1322)
	require.Equal(t, RangeIndex(Cards{Card2D, Card2H}), 1323)
	require.Equal(t, RangeIndex(Cards{Card2D, Card2S}), 1324)
	require.Equal(t, RangeIndex(Cards{Card2H, Card2S}), 1325)

	require.Equal(t, RangeCards(0), Cards{CardAC, CardAD})
	require.Equal(t, RangeCards(1), Cards{CardAC, CardAH})
	require.Equal(t, RangeCards(2), Cards{CardAC, CardAS})
	require.Equal(t, RangeCards(3), Cards{CardAD, CardAH})
	require.Equal(t, RangeCards(4), Cards{CardAD, CardAS})
	require.Equal(t, RangeCards(5), Cards{CardAH, CardAS})

	require.Equal(t, RangeCards(1320), Cards{Card2C, Card2D})
	require.Equal(t, RangeCards(1321), Cards{Card2C, Card2H})
	require.Equal(t, RangeCards(1322), Cards{Card2C, Card2S})
	require.Equal(t, RangeCards(1323), Cards{Card2D, Card2H})
	require.Equal(t, RangeCards(1324), Cards{Card2D, Card2S})
	require.Equal(t, RangeCards(1325), Cards{Card2H, Card2S})
}

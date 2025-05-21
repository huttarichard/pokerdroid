package iso

import (
	"testing"

	"github.com/pokerdroid/poker/card"
	"github.com/stretchr/testify/require"
)

func TestCardToUint(t *testing.T) {
	require.Equal(t, uint8(3), cardToUint(card.Card2C))
	require.Equal(t, uint8(46), cardToUint(card.CardKD))
	require.Equal(t, uint8(49), cardToUint(card.CardAH))

	require.Equal(t, uintToCard(3), card.Card2C)
	require.Equal(t, uintToCard(46), card.CardKD)
	require.Equal(t, uintToCard(49), card.CardAH)
}

func TestSizes(t *testing.T) {
	x := River.Size()
	y := Turn.Size()
	z := Flop.Size()
	w := Preflop.Size()

	require.Equal(t, 123156254, x)
	require.Equal(t, 13960050, y)
	require.Equal(t, 1286792, z)
	require.Equal(t, 169, w)
}

func TestPreflopLookup(t *testing.T) {
	for _, c := range card.Combinations(2) {
		idx := Preflop.Index(c)
		require.Equal(t, preflopLookup[[2]card.Card{c[0], c[1]}], idx)
	}
}

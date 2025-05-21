package mc

import (
	"testing"

	"github.com/pokerdroid/poker/card"
	"github.com/stretchr/testify/require"
)

func TestOpenRange(t *testing.T) {
	a := NewOpenRange()
	c := card.Cards{card.Card2C, card.Card3D}
	require.True(t, a.WeakRange(c))
	c = card.Cards{card.Card2C, card.Card3C}
	require.False(t, a.WeakRange(c))

	c = card.Cards{card.Card3D, card.Card2C}
	require.True(t, a.WeakRange(c))
	c = card.Cards{card.Card3C, card.Card2C}
	require.False(t, a.WeakRange(c))
}

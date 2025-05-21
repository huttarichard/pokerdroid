package kuhndealer

import (
	"testing"

	"github.com/pokerdroid/poker/card"
	"github.com/pokerdroid/poker/tree"
	"github.com/stretchr/testify/require"
)

func TestKuhnSampling(t *testing.T) {
	// r := rand.New(rand.NewSource(0))
	s := &Sample{Cards: card.Cards{card.CardKD, card.CardQD}}

	require.Equal(t, s.Utility(tree.KuhnP1KK, 0), float64(1))
	require.Equal(t, s.Utility(tree.KuhnP1KK, 1), float64(-1))

	require.Equal(t, s.Utility(tree.KuhnP1KBC, 0), float64(2))
	require.Equal(t, s.Utility(tree.KuhnP1KBC, 1), float64(-2))

	require.Equal(t, s.Utility(tree.KuhnP1KBF, 0), float64(-1))
	require.Equal(t, s.Utility(tree.KuhnP1KBF, 1), float64(1))

	require.Equal(t, s.Utility(tree.KuhnP1BC, 0), float64(2))
	require.Equal(t, s.Utility(tree.KuhnP1BC, 1), float64(-2))

	require.Equal(t, s.Utility(tree.KuhnP1BF, 0), float64(1))
	require.Equal(t, s.Utility(tree.KuhnP1BF, 1), float64(-1))

	s = &Sample{Cards: card.Cards{card.CardQD, card.CardKD}}

	require.Equal(t, s.Utility(tree.KuhnP1KK, 0), float64(-1))
	require.Equal(t, s.Utility(tree.KuhnP1KK, 1), float64(1))

	require.Equal(t, s.Utility(tree.KuhnP1KBC, 0), float64(-2))
	require.Equal(t, s.Utility(tree.KuhnP1KBC, 1), float64(2))

	require.Equal(t, s.Utility(tree.KuhnP1KBF, 0), float64(-1))
	require.Equal(t, s.Utility(tree.KuhnP1KBF, 1), float64(1))

	require.Equal(t, s.Utility(tree.KuhnP1BC, 0), float64(-2))
	require.Equal(t, s.Utility(tree.KuhnP1BC, 1), float64(2))

	require.Equal(t, s.Utility(tree.KuhnP1BF, 0), float64(1))
	require.Equal(t, s.Utility(tree.KuhnP1BF, 1), float64(-1))

}

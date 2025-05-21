package card

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRandomCards(t *testing.T) {
	require.Len(t, RandomCards(rand.New(rand.NewSource(0)), 5), 5)
}

func TestCompress(t *testing.T) {
	rc := RandomCards(rand.New(rand.NewSource(0)), 7)
	bb := Compress(rc)
	t.Logf("bytes %d", len(bb))
	cc, err := Decompress(bb)
	require.NoError(t, err)
	require.Equal(t, rc, cc)
}

func TestCardsInCoords(t *testing.T) {
	aces := CardsInCoords(0, 0)
	require.Len(t, aces, 6)
	require.Equal(t, aces[0], Cards{CardAC, CardAD})
	require.Equal(t, aces[1], Cards{CardAC, CardAH})
	require.Equal(t, aces[2], Cards{CardAC, CardAS})
	require.Equal(t, aces[3], Cards{CardAD, CardAH})
	require.Equal(t, aces[4], Cards{CardAD, CardAS})
	require.Equal(t, aces[5], Cards{CardAH, CardAS})

	pairs := CardsInCoords(12, 12)
	require.Equal(t, pairs[0], Cards{Card2C, Card2D})
	require.Equal(t, pairs[1], Cards{Card2C, Card2H})
	require.Equal(t, pairs[2], Cards{Card2C, Card2S})
	require.Equal(t, pairs[3], Cards{Card2D, Card2H})
	require.Equal(t, pairs[4], Cards{Card2D, Card2S})
	require.Equal(t, pairs[5], Cards{Card2H, Card2S})

	sutedAces := CardsInCoords(0, 12)
	require.Len(t, sutedAces, 4)

	require.Equal(t, sutedAces[0], Cards{Card2C, CardAC})
	require.Equal(t, sutedAces[1], Cards{Card2D, CardAD})
	require.Equal(t, sutedAces[2], Cards{Card2H, CardAH})
	require.Equal(t, sutedAces[3], Cards{Card2S, CardAS})

	offsuitedAces := CardsInCoords(12, 0)
	require.Len(t, offsuitedAces, 12)

	// [2c ad]
	require.Equal(t, offsuitedAces[0], Cards{Card2C, CardAD})
	// [2c ah]
	require.Equal(t, offsuitedAces[1], Cards{Card2C, CardAH})
	// [2c as]
	require.Equal(t, offsuitedAces[2], Cards{Card2C, CardAS})
	// [2d ac]
	require.Equal(t, offsuitedAces[3], Cards{Card2D, CardAC})
	// [2d ah]
	require.Equal(t, offsuitedAces[4], Cards{Card2D, CardAH})
	// [2d as]
	require.Equal(t, offsuitedAces[5], Cards{Card2D, CardAS})
	// [2h ac]
	require.Equal(t, offsuitedAces[6], Cards{Card2H, CardAC})
	// [2h ad]
	require.Equal(t, offsuitedAces[7], Cards{Card2H, CardAD})
	// [2h as]
	require.Equal(t, offsuitedAces[8], Cards{Card2H, CardAS})
	// [2s ac]
	require.Equal(t, offsuitedAces[9], Cards{Card2S, CardAC})
	// [2s ad]
	require.Equal(t, offsuitedAces[10], Cards{Card2S, CardAD})
	// [2s ah]
	require.Equal(t, offsuitedAces[11], Cards{Card2S, CardAH})

	eq := CardsInCoordsWithBlockersAt(0, 12, Cards{CardAC})
	require.Equal(t, eq[0], Cards{Card2D, CardAD})
	require.Equal(t, eq[1], Cards{Card2H, CardAH})
	require.Equal(t, eq[2], Cards{Card2S, CardAS})
}

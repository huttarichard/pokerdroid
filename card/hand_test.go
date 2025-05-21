package card

import (
	"testing"

	"github.com/stretchr/testify/require"
)



func TestEqualHands(t *testing.T) {
	h1 := NewHandRank(HandRankStraight, 123)
	h2 := NewHandRank(HandRankStraight, 123)
	require.Equal(t, 2, h1.Compare(h2))
	require.Equal(t, 2, h2.Compare(h1))
}

func TestBetterHandRank(t *testing.T) {
	h1 := NewHandRank(HandRankStraight, 123)
	h2 := NewHandRank(HandRankOnePair, 1234)
	require.Equal(t, 0, h1.Compare(h2))
	require.Equal(t, 1, h2.Compare(h1))
}

func TestEqualHandRank(t *testing.T) {
	h1 := NewHandRank(HandRankStraight, 2)
	h2 := NewHandRank(HandRankStraight, 1)
	require.Equal(t, 0, h1.Compare(h2))
	require.Equal(t, 1, h2.Compare(h1))
}

func TestBinaryEncoding(t *testing.T) {
	rank := HandRank{Kind: HandRankHighCard, Rank: 10}
	xx, err := rank.MarshalBinary()
	require.NoError(t, err)

	rx := &HandRank{}
	rx.UnmarshalBinary(xx)

	require.Equal(t, rank, *rx)
}

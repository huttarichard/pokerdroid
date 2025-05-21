package eval

import (
	"testing"

	"github.com/pokerdroid/poker/card"
	"github.com/stretchr/testify/require"
)

func TestW(t *testing.T) {
	_, err := Eval(card.Card2C, card.Card2D, card.Card2H, card.Card2S, card.Card3C)
	require.NoError(t, err)
}

func TestJudge2(t *testing.T) {
	testCases := []struct {
		desc    string
		hand    []card.Cards
		board   card.Cards
		winners []uint8
	}{
		{
			desc: "0",
			hand: []card.Cards{
				{card.Card2C, card.Card2D},
				{card.CardAC, card.CardAD},
			},
			board:   card.Cards{card.Card3S, card.Card4D, card.Card5H, card.Card8C, card.CardTD},
			winners: []uint8{1},
		},
		{
			desc: "1",
			hand: []card.Cards{
				{card.CardAH, card.CardAS},
				{card.CardAC, card.CardAD},
			},
			board:   card.Cards{card.Card3S, card.Card4D, card.Card5H, card.Card8C, card.CardTD},
			winners: []uint8{0, 1},
		},
		{
			desc: "2",
			hand: []card.Cards{
				{card.CardAH, card.CardAS},
				{card.CardAC, card.CardAD},
				{card.CardAC, card.CardAD},
			},
			board:   card.Cards{card.Card3S, card.Card4D, card.Card5H, card.Card8C, card.CardTD},
			winners: []uint8{0, 1, 2},
		},
		{
			desc: "3",
			hand: []card.Cards{
				{card.CardAH, card.CardAS},
				{card.CardAC, card.CardAD},
				{card.Card3D, card.Card4C},
			},
			board:   card.Cards{card.Card3S, card.Card4D, card.Card5H, card.Card8C, card.CardTD},
			winners: []uint8{2},
		},
		{
			desc: "4",
			hand: []card.Cards{
				{card.CardAH, card.CardAS},
				{card.CardAC, card.CardAD},
				{card.Card3D, card.Card4C},
			},
			board:   card.Cards{card.Card3S, card.Card4D, card.Card5H, card.Card8C, card.CardTD},
			winners: []uint8{2},
		},
		{
			desc: "5",
			hand: []card.Cards{
				{card.CardAH, card.CardAS},
				{card.Card3D, card.Card4C},
				{card.CardAC, card.CardAD},
			},
			board:   card.Cards{card.Card3S, card.Card4D, card.Card5H, card.Card8C, card.CardTD},
			winners: []uint8{1},
		},
		{
			desc: "5",
			hand: []card.Cards{
				{card.CardAH, card.CardAS},
				{card.Card3D, card.Card4C},
				{card.CardAC, card.CardAD},
				{card.Card5D, card.Card5C},
			},
			board:   card.Cards{card.Card3S, card.Card4D, card.Card5H, card.Card8C, card.CardTD},
			winners: []uint8{3},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			winners, err := JudgeBoard(tC.hand, tC.board)
			require.NoError(t, err)

			require.Equal(t, tC.winners, winners)
		})
	}
}

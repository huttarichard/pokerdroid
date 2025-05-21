package mc

import (
	"context"
	"testing"

	"github.com/pokerdroid/poker"
	"github.com/pokerdroid/poker/bot"
	"github.com/pokerdroid/poker/card"
	"github.com/pokerdroid/poker/chips"
	"github.com/pokerdroid/poker/table"
	"github.com/stretchr/testify/require"
)

func TestPolicyWeakRange(t *testing.T) {
	// Create game with initial params
	game, err := table.NewGame(table.GameParams{
		NumPlayers:         2,
		BtnPos:             0,
		SbAmount:           chips.NewFromInt(1),
		BetSizes:           [][]float32{{0.5, 1, 1.5}},
		InitialStacks:      chips.NewList(100, 100),
		TerminalStreet:     table.River,
		MaxActionsPerRound: 4,
	})
	require.NoError(t, err)

	// Create advisor and get action
	pol := NewAdvisor()
	_, err = pol.Advise(context.Background(), poker.VoidLogger{}, bot.State{
		Params:    game.GameParams,
		State:     game.Latest,
		Hole:      []card.Card{card.Card2C, card.Card3D},
		Community: card.Cards{},
	})
	require.NoError(t, err)
}

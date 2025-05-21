package mc

import (
	"testing"

	"github.com/pokerdroid/poker/chips"
	"github.com/pokerdroid/poker/frand"
	"github.com/pokerdroid/poker/table"
	"github.com/stretchr/testify/require"
)

func TestPolicyLegalEVs_Fold(t *testing.T) {
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

	// Make initial actions
	err = game.Action(table.DCall) // SB calls BB
	require.NoError(t, err)

	err = game.Action(table.ActionAmount{
		Action: table.Bet,
		Amount: chips.NewFromFloat(60),
	})
	require.NoError(t, err)

	// Create EVs and get action
	evs := NewEVs(EvsParams{
		Params:      game.GameParams,
		State:       game.Latest,
		Equity:      float32(0.32),
		Passiveness: 0.1,
		Defense:     0.2,
	})
	evs.Normalize()

	act := evs.Choice(frand.NewUnsafeInt(42), 0)
	require.Equal(t, table.DFold, act)
}

func TestPolicyLegalEVs_Call(t *testing.T) {
	// Create game with initial params
	game, err := table.NewGame(table.GameParams{
		NumPlayers:         2,
		BtnPos:             0,
		SbAmount:           chips.NewFromInt(1),
		BetSizes:           [][]float32{{0.25, 0.5, 0.75, 1, 1.5, 2, 4, 6, 8}},
		InitialStacks:      chips.NewList(100, 100),
		TerminalStreet:     table.River,
		MaxActionsPerRound: 8,
	})
	require.NoError(t, err)

	// Create EVs and get action
	evs := NewEVs(EvsParams{
		Params:      game.GameParams,
		State:       game.Latest,
		Equity:      float32(0.63),
		Passiveness: 0.1,
		Defense:     0.2,
	})
	evs.Normalize()

	act := evs.Max()
	require.Equal(t, table.DAllIn, act)
}

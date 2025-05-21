package tree

import (
	"testing"

	"github.com/pokerdroid/poker/chips"
	"github.com/pokerdroid/poker/table"
	"github.com/stretchr/testify/require"
)

func TestExpand(t *testing.T) {
	// Create game params
	p := table.GameParams{
		NumPlayers:         2,
		InitialStacks:      chips.NewList(100, 100),
		SbAmount:           chips.NewFromInt(1),
		BetSizes:           [][]float32{{1}},
		MaxActionsPerRound: 8,
		TerminalStreet:     table.River,
		BtnPos:             0,
		Limp:               true,
	}

	// Create initial game state
	state, err := table.NewGame(p)
	require.NoError(t, err)

	// Create root node
	root := &Root{
		Params: p,
		State:  state.Latest,
	}

	// Test initial expansion
	err = Expand(root, root)
	require.NoError(t, err)
	require.NotNil(t, root.Next)

	// First should be chance node for dealing cards
	chance1, ok := root.Next.(*Chance)
	require.True(t, ok, "First node should be Chance for dealing cards")

	// Expand chance node
	err = Expand(root, chance1)
	require.NoError(t, err)
	require.NotNil(t, chance1.Next)

	// After dealing, SB should act first preflop
	preflopPlayer, ok := chance1.Next.(*Player)
	require.True(t, ok, "Expected Player node after dealing")
	require.Equal(t, uint8(0), preflopPlayer.TurnPos, "SB should act first")

	// Expand preflop player node
	err = Expand(root, preflopPlayer)
	require.NoError(t, err)
	require.NotNil(t, preflopPlayer.Actions)

	// Verify preflop actions
	require.True(t, len(preflopPlayer.Actions.Actions) > 0, "Should have preflop actions")

	// Simulate call action
	callIdx := -1
	for i, act := range preflopPlayer.Actions.Actions {
		if act == table.DCall {
			callIdx = i
			break
		}
	}
	require.NotEqual(t, -1, callIdx, "Call action should exist")

	// Get node after call
	nextNode := preflopPlayer.Actions.Nodes[callIdx]
	require.NotNil(t, nextNode)

	// Should be BB's turn
	bbPlayer, ok := nextNode.(*Player)
	require.True(t, ok, "Expected Player node for BB")
	require.Equal(t, uint8(1), bbPlayer.TurnPos, "BB should act next")

	// Expand BB node
	err = Expand(root, bbPlayer)
	require.NoError(t, err)

	require.NotNil(t, bbPlayer.Actions)

	// Find check action
	checkIdx := -1
	for i, act := range bbPlayer.Actions.Actions {
		if act == table.DCheck {
			checkIdx = i
			break
		}
	}
	require.NotEqual(t, -1, checkIdx, "Check action should exist")

	// Get node after check - should be chance node for flop
	flopChance := bbPlayer.Actions.Nodes[checkIdx]
	require.NotNil(t, flopChance)
	chance2, ok := flopChance.(*Chance)
	require.True(t, ok, "Expected Chance node for flop")

	// Expand flop chance node
	err = Expand(root, chance2)
	require.NoError(t, err)
	require.NotNil(t, chance2.Next)

	// After flop, BB should act first
	flopPlayer, ok := chance2.Next.(*Player)
	require.True(t, ok, "Expected Player node after flop")
	require.Equal(t, uint8(1), flopPlayer.TurnPos, "BB should act first on flop")

	// Test full expansion
	err = Expand(root, flopPlayer)
	require.NoError(t, err)

	// Find all-in action
	allInIdx := -1
	for i, act := range flopPlayer.Actions.Actions {
		if act == table.DAllIn {
			allInIdx = i
			break
		}
	}
	require.NotEqual(t, -1, allInIdx, "All-in action should exist")

	// Get node after all-in
	allInNode := flopPlayer.Actions.Nodes[allInIdx]
	require.NotNil(t, allInNode)

	// Should be SB's turn
	sbPlayer, ok := allInNode.(*Player)
	require.True(t, ok, "Expected Player node for SB")
	require.Equal(t, uint8(0), sbPlayer.TurnPos, "SB should act next")

	// Expand SB node
	err = Expand(root, sbPlayer)
	require.NoError(t, err)
	require.NotNil(t, sbPlayer.Actions)

	// Find call action
	callIdx = -1
	for i, act := range sbPlayer.Actions.Actions {
		if act == table.DCall {
			callIdx = i
			break
		}
	}
	require.NotEqual(t, -1, callIdx, "Call action should exist")

	// Get node after call - should be terminal since both players are all-in
	terminal := sbPlayer.Actions.Nodes[callIdx]
	require.NotNil(t, terminal)
	_, ok = terminal.(*Terminal)
	require.True(t, ok, "Expected Terminal node after call of all-in")
}

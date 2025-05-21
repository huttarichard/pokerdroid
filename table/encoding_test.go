package table

import (
	"testing"

	"github.com/pokerdroid/poker/chips"
	"github.com/stretchr/testify/require"
)

// TestMarshalUnmarshalGame creates a Game with a chain of three states,
// marshals it using MarshalGame, then unmarshals it using UnmarshalGame,
// and verifies that the game parameters and the state chain (and final state)
// match the original.
func TestMarshalUnmarshalGame(t *testing.T) {
	// Create game parameters.
	params := NewGameParams(2, chips.NewFromInt(100))
	params.BtnPos = 0
	params.TerminalStreet = River
	params.BetSizes = BetSizesDeep

	// Create the initial state; modify fields for testing.
	st1 := NewState(params)
	st1.CallAmount = chips.NewFromInt(10)

	// Create a chain of states.
	st2 := st1.Next()
	st2.Street = Flop
	st2.CallAmount = chips.NewFromInt(20)

	st3 := st2.Next()
	st3.Street = Turn
	st3.CallAmount = chips.NewFromInt(30)

	// Marshal the game.
	data, err := MarshalBinary(params, st3)
	if err != nil {
		t.Fatalf("failed to marshal game: %v", err)
	}

	// Unmarshal the game.
	newParams, newState, err := UnmarshalBinary(data)
	if err != nil {
		t.Fatalf("failed to unmarshal game: %v", err)
	}

	// Compare game parameters.
	if params.String() != newParams.String() {
		t.Fatal("game params mismatch after unmarshal")
	}

	// Verify that the state chain length is preserved.
	count := 0
	for cur := newState; cur != nil; cur = cur.Previous {
		count++
	}
	if count != 3 {
		t.Fatalf("expected state chain length 3, got %d", count)
	}

	// Verify that the final state matches.
	if !st3.Equal(newState) {
		t.Fatal("final state does not match after unmarshal")
	}
}

func TestActionMarshalUnmarshal(t *testing.T) {
	// Create game parameters.
	params := NewGameParams(2, chips.NewFromInt(100))
	params.BetSizes = BetSizesDeep
	params.TerminalStreet = River

	// Initialize a new game.
	game, err := NewGame(params)
	require.NoError(t, err, "failed to initialize game")

	// Save the current acting player's position.
	actingPlayer := game.Latest.TurnPos

	// Execute an action: a Raise of 10 chips.
	act := ActionAmount{
		Action: Raise,
		Amount: chips.NewFromInt(10),
	}
	err = game.Action(act)
	require.NoError(t, err, "failed to execute action")

	// Check that the action was recorded in the previous state.
	require.NotNil(t, game.Latest.Previous, "expected previous state holding action details")
	require.Equal(t, Raise, game.Latest.Previous.PSLA[actingPlayer],
		"acting player's last action should be Raise")
	require.Contains(t, game.Latest.Previous.String(), Raise.String(),
		"state String should contain the Raise action")

	// Use MarshalGame to serialize the game (including its parameters and state chain).
	data, err := MarshalBinary(params, game.Latest)
	require.NoError(t, err, "failed to marshal game with state chain")

	// Unmarshal the data back into a new game.
	newParams, newState, err := UnmarshalBinary(data)
	require.NoError(t, err, "failed to unmarshal game data")

	// Verify that the game parameters match.
	require.Equal(t, params.String(), newParams.String(),
		"game parameters mismatch after marshal/unmarshal")

	// Verify that the state chain was preserved.
	count := 0
	for st := newState; st != nil; st = st.Previous {
		count++
	}
	// Assuming at least one action was taken, so there must be at least 2 states.
	require.GreaterOrEqual(t, count, 2, "state chain should consist of at least 2 states after an action")

	// Verify that the acting player's last action is still Raise in the unmarshaled game.
	require.NotNil(t, newState.Previous, "expected previous state in unmarshaled game")
	require.Equal(t, Raise, newState.Previous.PSLA[actingPlayer],
		"acting player's last action should be Raise in unmarshaled game")
	require.Contains(t, newState.Previous.String(), Raise.String(),
		"unmarshaled state String should include the Raise action")
}

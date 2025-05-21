package table

import (
	"testing"

	"github.com/pokerdroid/poker/chips"
	"github.com/stretchr/testify/require"
)

func TestRule_ShiftTurn(t *testing.T) {
	s := &State{
		Players: []Player{
			{Paid: chips.New(2), Status: StatusActive},
			{Paid: chips.New(2), Status: StatusActive},
			{Paid: chips.New(2), Status: StatusFolded},
		},
		TurnPos: 0,
		BSC: struct {
			Amount   chips.Chips `json:"amount"`
			Addition chips.Chips `json:"addition"`
			Action   ActionKind  `json:"action"`
		}{
			Amount: chips.NewFromInt(2),
		},
		PSC: chips.NewList(2, 2, 2),
	}

	err := ShiftTurn(GameParams{
		NumPlayers:     2,
		InitialStacks:  chips.NewList(100, 100),
		TerminalStreet: River,
	}, s)
	require.NoError(t, err)
	require.EqualValues(t, 1, s.TurnPos)
}

func TestRule_MultiplePlayersToAct(t *testing.T) {
	p := GameParams{
		NumPlayers:     3,
		InitialStacks:  chips.NewList(100, 100, 100),
		TerminalStreet: River,
	}
	s := &State{
		Players: []Player{
			{Paid: chips.New(2), Status: StatusActive},
			{Paid: chips.New(2), Status: StatusActive},
			{Paid: chips.Zero, Status: StatusActive},
		},
		Street: Preflop,
		PSC:    chips.NewList(2, 2, 0),
		BSC: struct {
			Amount   chips.Chips `json:"amount"`
			Addition chips.Chips `json:"addition"`
			Action   ActionKind  `json:"action"`
		}{Amount: chips.New(2)},
		PSAC: []uint8{0, 0, 0},
		PSLA: []ActionKind{NoAction, NoAction, NoAction},
	}
	// Expect SHIFT TURN, not SHIFT STREET or FINISH
	ruleType := Rule(p, s)
	require.Equal(t, RuleShiftTurn, ruleType)
}

func TestMove_FinishOnFold(t *testing.T) {
	p := GameParams{
		NumPlayers:    2,
		InitialStacks: chips.NewList(100, 100),
	}

	// One player folded...
	s := &State{
		Players: []Player{
			{Status: StatusFolded, Paid: chips.NewFromInt(10)},
			{Status: StatusActive, Paid: chips.NewFromInt(10)},
		},
		Street: Preflop,
	}
	ns, err := Move(p, s)
	require.NoError(t, err)
	require.Equal(t, Finished, ns.Street)
}

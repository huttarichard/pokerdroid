package table

import (
	"testing"

	"github.com/pokerdroid/poker/chips"
	"github.com/stretchr/testify/require"
)

func TestMakeInitialBets(t *testing.T) {
	p := GameParams{
		NumPlayers:    2,
		InitialStacks: chips.NewList(100, 100),
		SbAmount:      chips.New(1),
	}
	s := &State{
		Players: []Player{
			{Paid: chips.Zero, Status: StatusActive},
			{Paid: chips.Zero, Status: StatusActive},
		},
		TurnPos: 0,
		PSC:     chips.NewList(0, 0),
		PSAC:    []uint8{0, 0},
		PSLA:    []ActionKind{NoAction, NoAction},
	}
	newS, err := MakeInitialBets(p, s)
	require.NoError(t, err)
	require.Equal(t, chips.NewFromInt(1), newS.Players[0].Paid) // SB
	require.Equal(t, chips.NewFromInt(2), newS.Players[1].Paid) // BB
}

func TestMakeAction_Fold(t *testing.T) {
	p := GameParams{
		NumPlayers:         2,
		InitialStacks:      chips.NewList(100, 100),
		MaxActionsPerRound: 4,
		TerminalStreet:     River,
	}
	s := &State{
		Players: []Player{
			{Paid: chips.Zero, Status: StatusActive},
			{Paid: chips.Zero, Status: StatusActive},
		},
		Street:     Preflop,
		TurnPos:    0,
		CallAmount: chips.New(1),
		PSC:        chips.NewList(0, 0),
		PSAC:       []uint8{0, 0},
		PSLA:       []ActionKind{NoAction, NoAction},
	}
	na, err := MakeAction(p, s, ActionAmount{
		Action: Fold,
		Amount: chips.Zero,
	})
	require.NoError(t, err)
	require.Equal(t, StatusFolded, na.Players[0].Status)
}

func TestMakeAction_AllIn(t *testing.T) {
	p := GameParams{
		NumPlayers:         2,
		InitialStacks:      chips.NewList(10, 100),
		MaxActionsPerRound: 4,
		SbAmount:           chips.New(1),
		TerminalStreet:     River,
	}
	s := &State{
		Players: []Player{
			{Paid: chips.Zero, Status: StatusActive}, // short stack
			{Paid: chips.Zero, Status: StatusActive},
		},
		TurnPos:    0,
		CallAmount: chips.New(2),
		PSC:        chips.NewList(0, 0),
		PSAC:       []uint8{0, 0},
		PSLA:       []ActionKind{NoAction, NoAction},
	}
	na, err := MakeAction(p, s, ActionAmount{
		Action: AllIn,
		Amount: chips.NewFromInt(10), // tries to go all in more than they have
	})
	require.NoError(t, err)
	require.Equal(t, chips.NewFromInt(10), na.Players[0].Paid)
	require.Equal(t, StatusAllIn, na.Players[0].Status)
}

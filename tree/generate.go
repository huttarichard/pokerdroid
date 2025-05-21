package tree

import (
	"fmt"

	"github.com/pokerdroid/poker/chips"
	"github.com/pokerdroid/poker/table"
)

type NewFullTreeParams struct {
	BigBlind   chips.Chips  `json:"bb"`
	NumPlayers int          `json:"players"`
	Betting    [][]float32  `json:"betting"`
	MaxActions int          `json:"max_actions"`
	Terminal   table.Street `json:"terminal"` // Changed to table2.Street
	MinBet     bool         `json:"min_bet"`
	Limp       bool         `json:"limp"`
}

func (t NewFullTreeParams) Name() string {
	return fmt.Sprintf(
		"tree_p%d_b%d_ma%d_bb%.0f_t%s",
		t.NumPlayers,
		len(t.Betting),
		t.MaxActions,
		t.BigBlind,
		t.Terminal,
	)
}

func NewFullTree(params NewFullTreeParams) (*Root, error) {
	// Create game params with normalized values (SB = 1)
	gameParams := table.GameParams{
		NumPlayers:         uint8(params.NumPlayers),
		MaxActionsPerRound: uint8(params.MaxActions),
		BtnPos:             0,
		SbAmount:           chips.NewFromInt(1), // Always 1
		BetSizes:           params.Betting,
		InitialStacks:      make(chips.List, params.NumPlayers),
		TerminalStreet:     params.Terminal,
		MinBet:             params.MinBet,
		Limp:               params.Limp,
	}

	// Set initial stacks for all players (in big blinds)
	for i := range gameParams.InitialStacks {
		gameParams.InitialStacks[i] = params.BigBlind
	}

	// Create initial state
	state := table.NewState(gameParams)

	// Make initial bets (SB and BB)
	state, err := table.MakeInitialBets(gameParams, state)
	if err != nil {
		return nil, err
	}

	// Create root node
	root := &Root{
		Params:    gameParams,
		State:     state,
		Iteration: 0,
	}

	// Expand full tree
	err = ExpandFull(root)
	if err != nil {
		return nil, err
	}

	return root, nil
}

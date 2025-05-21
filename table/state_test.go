package table

import (
	"testing"

	"github.com/pokerdroid/poker/chips"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestState(t *testing.T) {
	g, err := NewGame(GameParams{
		NumPlayers: 2,
		BtnPos:     0,
		SbAmount:   chips.NewFromInt(1),
		BetSizes:   [][]float32{{0.5, 1, 1.5}},
		InitialStacks: chips.List{
			chips.NewFromInt(100),
			chips.NewFromInt(100),
		},
		TerminalStreet: River,
	})
	require.NoError(t, err)

	err = g.Action(DCall)
	require.NoError(t, err)

	err = g.Action(DCheck)
	require.NoError(t, err)

	err = g.Action(DCheck)
	require.NoError(t, err)

	err = g.Action(DCheck)
	require.NoError(t, err)

	err = g.Action(DCheck)
	require.NoError(t, err)

	err = g.Action(DCheck)
	require.NoError(t, err)

	err = g.Action(DCheck)
	require.NoError(t, err)

	err = g.Action(DCheck)
	require.NoError(t, err)

	require.Equal(t, g.Latest.Street, Finished)
	require.Equal(t, g.Latest.Finished(), true)
}

func TestGameFlow(t *testing.T) {
	tests := []struct {
		name          string
		setup         func() (*Game, error)
		actions       []DiscreteAction
		expectedState func(*State) bool
	}{
		{
			name: "preflop raise - fold",
			setup: func() (*Game, error) {
				return NewGame(GameParams{
					NumPlayers:         2,
					BtnPos:             0,
					SbAmount:           chips.NewFromInt(1),
					BetSizes:           [][]float32{{0.5, 1, 1.5}},
					InitialStacks:      chips.NewList(100, 100),
					TerminalStreet:     River,
					MaxActionsPerRound: 4,
				})
			},
			actions: []DiscreteAction{
				DiscreteAction(1), // SB raises pot
				DFold,             // BB folds
			},
			expectedState: func(s *State) bool {
				return s.Finished() &&
					s.Players[1].Status == StatusFolded
			},
		},
		{
			name: "check to river",
			setup: func() (*Game, error) {
				return NewGame(GameParams{
					NumPlayers:         2,
					BtnPos:             0,
					SbAmount:           chips.NewFromInt(1),
					BetSizes:           [][]float32{{0.5, 1, 1.5}},
					InitialStacks:      chips.NewList(100, 100),
					TerminalStreet:     River,
					MaxActionsPerRound: 4,
				})
			},
			actions: []DiscreteAction{
				DCall,  // SB calls
				DCheck, // BB checks
				DCheck, // SB checks flop
				DCheck, // BB checks flop
				DCheck, // SB checks turn
				DCheck, // BB checks turn
				DCheck, // SB checks river
				DCheck, // BB checks river
			},
			expectedState: func(s *State) bool {
				return s.Finished() &&
					s.Players[0].Status == StatusActive &&
					s.Players[1].Status == StatusActive
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			game, err := tt.setup()
			require.NoError(t, err)

			for _, action := range tt.actions {
				err = game.Action(action)
				require.NoError(t, err)
			}

			require.True(t, tt.expectedState(game.Latest))
		})
	}
}

func TestPositions(t *testing.T) {
	tests := []struct {
		name        string
		numPlayers  uint8
		btnPos      uint8
		expectedBtn uint8
		expectedSB  uint8
		expectedBB  uint8
	}{
		{
			name:        "heads up",
			numPlayers:  2,
			btnPos:      0,
			expectedBtn: 0,
			expectedSB:  0, // In heads up, BTN is SB
			expectedBB:  1,
		},
		{
			name:        "3 players",
			numPlayers:  3,
			btnPos:      0,
			expectedBtn: 0,
			expectedSB:  1,
			expectedBB:  2,
		},
		{
			name:        "6 players - btn wrap",
			numPlayers:  6,
			btnPos:      5,
			expectedBtn: 5,
			expectedSB:  0,
			expectedBB:  1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gp := NewGameParams(tt.numPlayers, chips.NewFromInt(100))
			gp.BtnPos = tt.btnPos

			game, err := NewGame(gp)
			require.NoError(t, err)

			btn, sb, bb := Positions(game.Latest)
			require.Equal(t, tt.expectedBtn, btn)
			require.Equal(t, tt.expectedSB, sb)
			require.Equal(t, tt.expectedBB, bb)
		})
	}
}

func TestPotAccumulation(t *testing.T) {
	tests := []struct {
		name          string
		setup         func() (*Game, error)
		actions       []DiscreteAction
		expectedPots  Pots
		expectedTotal chips.Chips
	}{
		{
			name: "simple pot - no raises",
			setup: func() (*Game, error) {
				return NewGame(GameParams{
					NumPlayers:         2,
					BtnPos:             0,
					SbAmount:           chips.NewFromInt(1),
					BetSizes:           [][]float32{{0.5, 1, 1.5}},
					InitialStacks:      chips.NewList(100, 100),
					TerminalStreet:     River,
					MaxActionsPerRound: 4,
				})
			},
			actions: []DiscreteAction{
				DCall,  // SB calls
				DCheck, // BB checks
				DCheck, // Flop checks
				DCheck,
				DCheck, // Turn checks
				DCheck,
				DCheck, // River checks
				DCheck,
			},
			expectedTotal: chips.NewFromInt(4), // 2BB total pot
		},
		{
			name: "pot with raises",
			setup: func() (*Game, error) {
				return NewGame(GameParams{
					NumPlayers:         2,
					BtnPos:             0,
					SbAmount:           chips.NewFromInt(1),
					BetSizes:           [][]float32{{0.5, 1, 1.5}},
					InitialStacks:      chips.NewList(100, 100),
					TerminalStreet:     River,
					MaxActionsPerRound: 4,
				})
			},
			actions: []DiscreteAction{
				DiscreteAction(1), // SB raises pot
				DCall,             // BB calls
				DCheck,            // Flop checks
				DCheck,
			},
			expectedTotal: chips.NewFromInt(8), // 4BB total pot after raises
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			game, err := tt.setup()
			require.NoError(t, err)

			for _, action := range tt.actions {
				err = game.Action(action)
				require.NoError(t, err)
			}

			totalPot := game.Latest.Players.PaidSum()
			require.True(t, tt.expectedTotal.Equal(totalPot))
		})
	}
}

func TestStreetProgression(t *testing.T) {
	tests := []struct {
		name           string
		setup          func() (*Game, error)
		actions        []DiscreteAction
		expectedStreet Street
	}{
		{
			name: "progress to flop",
			setup: func() (*Game, error) {
				return NewGame(GameParams{
					NumPlayers:     2,
					BtnPos:         0,
					SbAmount:       chips.NewFromInt(1),
					BetSizes:       [][]float32{{0.5, 1, 1.5}},
					InitialStacks:  chips.NewList(100, 100),
					TerminalStreet: River,
				})
			},
			actions: []DiscreteAction{
				DCall,  // SB calls
				DCheck, // BB checks
			},
			expectedStreet: Flop,
		},
		{
			name: "early finish on fold",
			setup: func() (*Game, error) {
				return NewGame(GameParams{
					NumPlayers:     2,
					BtnPos:         0,
					SbAmount:       chips.NewFromInt(1),
					BetSizes:       [][]float32{{0.5, 1, 1.5}},
					InitialStacks:  chips.NewList(100, 100),
					TerminalStreet: River,
				})
			},
			actions: []DiscreteAction{
				DFold, // SB folds
			},
			expectedStreet: Finished,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			game, err := tt.setup()
			require.NoError(t, err)

			for _, action := range tt.actions {
				err = game.Action(action)
				require.NoError(t, err)
			}

			require.Equal(t, tt.expectedStreet, game.Latest.Street)
		})
	}
}

func TestAllInScenarios(t *testing.T) {
	tests := []struct {
		name          string
		setup         func() (*Game, error)
		actions       []DiscreteAction
		expectedPots  []chips.Chips
		expectedState func(*State) bool
	}{
		{
			name: "simple all-in vs call",
			setup: func() (*Game, error) {
				return NewGame(GameParams{
					NumPlayers:         2,
					BtnPos:             0,
					SbAmount:           chips.NewFromInt(1),
					BetSizes:           [][]float32{{0.5, 1, 1.5}},
					InitialStacks:      chips.NewList(10, 100),
					TerminalStreet:     River,
					MaxActionsPerRound: 4,
				})
			},
			actions: []DiscreteAction{
				DAllIn, // SB all-in (10 chips)
				DCall,  // BB calls
			},
			expectedState: func(s *State) bool {
				return s.Street == Finished &&
					s.Players[0].Status == StatusAllIn &&
					s.Players[1].Status == StatusActive
			},
		},
		{
			name: "all-in vs fold",
			setup: func() (*Game, error) {
				return NewGame(GameParams{
					NumPlayers:         2,
					BtnPos:             0,
					SbAmount:           chips.NewFromInt(1),
					BetSizes:           [][]float32{{0.5, 1, 1.5}},
					InitialStacks:      chips.NewList(10, 100),
					TerminalStreet:     River,
					MaxActionsPerRound: 4,
				})
			},
			actions: []DiscreteAction{
				DAllIn, // SB all-in
				DFold,  // BB folds
			},
			expectedState: func(s *State) bool {
				return s.Street == Finished &&
					s.Players[0].Status == StatusAllIn &&
					s.Players[1].Status == StatusFolded
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			game, err := tt.setup()
			require.NoError(t, err)

			for _, action := range tt.actions {
				err = game.Action(action)
				require.NoError(t, err)
			}

			require.True(t, tt.expectedState(game.Latest))
		})
	}
}

func TestMultiwayGameplay(t *testing.T) {
	// 4-player game, each with 100 chips
	g, err := NewGame(GameParams{
		NumPlayers:     4,
		BtnPos:         0,
		SbAmount:       chips.NewFromInt(1),
		BetSizes:       [][]float32{{0.5, 1, 1.5}},
		InitialStacks:  chips.NewList(100, 100, 100, 100),
		TerminalStreet: River,
	})
	require.NoError(t, err)

	// By default, p0=SB(1), p1=BB(2), p2 next, p3 next...
	// Let p2 call 2, p3 calls 2, p0 calls 1 more, p1 checks
	err = g.Action(DCall) // p2 calls 2
	require.NoError(t, err)
	err = g.Action(DCall) // p3 calls 2
	require.NoError(t, err)
	err = g.Action(DCall) // p0 calls 1 more
	require.NoError(t, err)
	err = g.Action(DCheck) // p1 checks
	require.NoError(t, err)

	// Everyone checks flop
	err = g.Action(DCheck)
	require.NoError(t, err)
	err = g.Action(DCheck)
	require.NoError(t, err)
	err = g.Action(DCheck)
	require.NoError(t, err)
	err = g.Action(DCheck)
	require.NoError(t, err)

	// Everyone checks turn
	for i := 0; i < 4; i++ {
		err = g.Action(DCheck)
		require.NoError(t, err)
	}

	// Everyone checks river
	for i := 0; i < 4; i++ {
		err = g.Action(DCheck)
		require.NoError(t, err)
	}

	require.Equal(t, Finished, g.Latest.Street)
	totalPot := g.Latest.Players.PaidSum()
	// Pot should be 2(BB) + 2(Calls by p2) + 2(Calls by p3) + 2( for p0 total ) = 8 total
	require.Equal(t, chips.NewFromInt(8), totalPot)

	// Positions
	btn, sb, bb := Positions(g.Latest)
	require.EqualValues(t, 0, btn)
	require.EqualValues(t, 1, sb)
	require.EqualValues(t, 2, bb)
}

func Test3PlayersFlopBetRaiseAndTie(t *testing.T) {
	const startChips = 100
	g, err := NewGame(GameParams{
		NumPlayers:         3,
		BtnPos:             0, // p0 is BTN
		SbAmount:           chips.NewFromInt(1),
		BetSizes:           [][]float32{{0.5, 1, 1.5}},
		InitialStacks:      chips.NewList(startChips, startChips, startChips),
		TerminalStreet:     River,
		MaxActionsPerRound: 8,
	})
	require.NoError(t, err)

	// Preflop actions:
	// p0 must act first (BTN in 3-handed is effectively UTG):
	err = g.Action(DCall) // p0 calls BB 2
	require.NoError(t, err)
	err = g.Action(DCall) // p1 calls 1 more
	require.NoError(t, err)
	err = g.Action(DCheck) // p2 (BB) checks
	require.NoError(t, err)

	// Confirm street advanced to flop
	require.Equal(t, Flop, g.Latest.Street)

	// Flop:
	// p2 to act first (since p2 was BB):
	// p2 bets 4
	err = g.Action(DiscreteAction(0.5)) // pot right now is 2+2+2=6, so 0.5 * 6=3 => but must at least be min bet?
	// We force bigger bet: let's do a direct ActionAmount for clarity
	if err != nil { // if pot-based didn't succeed, fallback:
		err = g.Action(ActionAmount{Action: Bet, Amount: chips.NewFromInt(4)})
	}
	require.NoError(t, err)

	// p0 raises to 12
	err = g.Action(ActionAmount{Action: Raise, Amount: chips.NewFromInt(12)})
	require.NoError(t, err)

	// p1 calls 12
	err = g.Action(DCall)
	require.NoError(t, err)

	// p2 must call additional 8 (he bet 4, total 12 => 8 more)
	err = g.Action(DCall)
	require.NoError(t, err)

	// Move on to turn
	require.Equal(t, Turn, g.Latest.Street)

	// Turn: everyone checks
	for i := 0; i < 3; i++ {
		err = g.Action(DCheck)
		require.NoError(t, err)
	}
	// Move on to river
	require.Equal(t, River, g.Latest.Street)

	// River: everyone checks
	for i := 0; i < 3; i++ {
		err = g.Action(DCheck)
		require.NoError(t, err)
	}
	require.Equal(t, Finished, g.Latest.Street)

	// Let's see total pot
	totalPot := g.Latest.Players.PaidSum()

	require.Equal(t, chips.NewFromInt(42), totalPot)

	// Judge scenario: p1 & p2 tie, p0 is worst
	j := testJudge{
		winners: map[uint8]uint8{
			0: 2, // p0 loses
			1: 0, // p1 best
			2: 0, // p2 also best => tie with p1
		},
	}
	winnings := GetWinnings(g.Latest.Players, j)
	// Tied among p1 & p2 => each gets half
	split := chips.NewFromInt(21) // 42 / 2
	expected := chips.NewList(chips.New(0), split, split)
	require.Equal(t, expected, winnings)
}

func TestState_Equal(t *testing.T) {
	baseState := &State{
		Players: []Player{
			{Paid: chips.NewFromInt(10), Status: StatusActive},
			{Paid: chips.NewFromInt(20), Status: StatusAllIn},
			{Paid: chips.NewFromInt(15), Status: StatusFolded},
		},
		Street:       Flop,
		TurnPos:      1,
		BtnPos:       0,
		StreetAction: 2,
		CallAmount:   chips.NewFromInt(20),
		BSC: struct {
			Amount   chips.Chips `json:"amount"`
			Addition chips.Chips `json:"addition"`
			Action   ActionKind  `json:"action"`
		}{
			Amount:   chips.NewFromInt(20),
			Addition: chips.NewFromInt(10),
			Action:   Raise,
		},
		PSC:  chips.NewList(10, 20, 15),
		PSAC: []uint8{1, 2, 1},
		PSLA: []ActionKind{Call, Raise, Call},
		Previous: &State{
			Players: []Player{
				{Paid: chips.NewFromInt(2), Status: StatusActive},
				{Paid: chips.NewFromInt(2), Status: StatusActive},
				{Paid: chips.NewFromInt(2), Status: StatusActive},
			},
			Street:       Preflop,
			TurnPos:      0,
			BtnPos:       0,
			StreetAction: 1,
			CallAmount:   chips.NewFromInt(2),
			BSC: struct {
				Amount   chips.Chips `json:"amount"`
				Addition chips.Chips `json:"addition"`
				Action   ActionKind  `json:"action"`
			}{
				Amount:   chips.NewFromInt(2),
				Addition: chips.Zero,
				Action:   BigBlind,
			},
			PSC:  chips.NewList(2, 2, 2),
			PSAC: []uint8{1, 1, 1},
			PSLA: []ActionKind{SmallBlind, BigBlind},
		},
	}

	tests := []struct {
		name     string
		state1   *State
		state2   *State
		expected bool
	}{
		{
			name:     "identical states",
			state1:   baseState,
			state2:   baseState,
			expected: true,
		},
		{
			name:     "nil states",
			state1:   nil,
			state2:   nil,
			expected: true,
		},
		{
			name:     "one nil state",
			state1:   baseState,
			state2:   nil,
			expected: false,
		},
		{
			name:   "different street",
			state1: baseState,
			state2: func() *State {
				s := *baseState
				s.Street = Turn
				return &s
			}(),
			expected: false,
		},
		{
			name:   "different player status",
			state1: baseState,
			state2: func() *State {
				s := *baseState
				s.Players = make([]Player, len(baseState.Players))
				copy(s.Players, baseState.Players)
				s.Players[0].Status = StatusFolded
				return &s
			}(),
			expected: false,
		},
		{
			name:   "different BSC amount",
			state1: baseState,
			state2: func() *State {
				s := *baseState
				s.BSC.Amount = chips.NewFromInt(30)
				return &s
			}(),
			expected: false,
		},
		{
			name:   "different PSLA",
			state1: baseState,
			state2: func() *State {
				s := *baseState
				s.PSLA = []ActionKind{Fold, Check}
				return &s
			}(),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.state1.Equal(tt.state2)
			require.Equal(t, tt.expected, result)
		})
	}

	// Test deep copy equality
	deepCopy := &State{}
	data, err := baseState.MarshalBinary()
	require.NoError(t, err)
	err = deepCopy.UnmarshalBinary(data)
	require.NoError(t, err)
	require.True(t, baseState.Equal(deepCopy))
}

func TestState_BinaryMarshalUnmarshal(t *testing.T) {
	// Create a complex state with multiple levels
	originalState := &State{
		Players: []Player{
			{Paid: chips.NewFromInt(10), Status: StatusActive},
			{Paid: chips.NewFromInt(20), Status: StatusAllIn},
			{Paid: chips.NewFromInt(15), Status: StatusFolded},
		},
		Street:       Flop,
		TurnPos:      1,
		BtnPos:       0,
		StreetAction: 2,
		CallAmount:   chips.NewFromInt(20),
		BSC: struct {
			Amount   chips.Chips `json:"amount"`
			Addition chips.Chips `json:"addition"`
			Action   ActionKind  `json:"action"`
		}{
			Amount:   chips.NewFromInt(20),
			Addition: chips.NewFromInt(10),
			Action:   Raise,
		},
		PSC:  chips.NewList(10, 20, 15),
		PSAC: []uint8{1, 2, 1},
		PSLA: []ActionKind{Call, Raise, Call},
		Previous: &State{ // Add a previous state
			Players: []Player{
				{Paid: chips.NewFromInt(2), Status: StatusActive},
				{Paid: chips.NewFromInt(2), Status: StatusActive},
				{Paid: chips.NewFromInt(2), Status: StatusActive},
			},
			Street:       Preflop,
			TurnPos:      0,
			BtnPos:       0,
			StreetAction: 1,
			CallAmount:   chips.NewFromInt(2),
			BSC: struct {
				Amount   chips.Chips `json:"amount"`
				Addition chips.Chips `json:"addition"`
				Action   ActionKind  `json:"action"`
			}{
				Amount:   chips.NewFromInt(2),
				Addition: chips.Zero,
				Action:   BigBlind,
			},
			PSC:  chips.NewList(2, 2, 2),
			PSAC: []uint8{1, 1, 1},
			PSLA: []ActionKind{SmallBlind, BigBlind},
		},
	}

	tests := []struct {
		name  string
		state *State
	}{
		{
			name:  "complex state with previous",
			state: originalState,
		},
		{
			name: "state without previous",
			state: &State{
				Players:    []Player{{Paid: chips.NewFromInt(10), Status: StatusActive}},
				Street:     Preflop,
				TurnPos:    0,
				BtnPos:     0,
				CallAmount: chips.NewFromInt(2),
				PSC:        chips.NewList(2),
				PSAC:       []uint8{1},
				PSLA:       []ActionKind{SmallBlind},
				Previous:   nil,
			},
		},
		{
			name: "state with zero values",
			state: &State{
				Players: []Player{{Paid: chips.Zero, Status: StatusActive}},
				Street:  NoStreet,
				PSC:     chips.NewList(0),
				PSAC:    []uint8{0},
				PSLA:    []ActionKind{NoAction},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Marshal the state
			data, err := tt.state.MarshalBinary()
			require.NoError(t, err)
			require.NotEmpty(t, data)

			// Unmarshal into a new state
			newState := &State{}
			err = newState.UnmarshalBinary(data)
			require.NoError(t, err)

			// Compare using Equal method
			require.True(t, tt.state.Equal(newState),
				"states should be equal after marshal/unmarshal")

			// Double check specific fields for debugging
			if !tt.state.Equal(newState) {
				t.Logf("Original state:\n%s", tt.state.String())
				t.Logf("Unmarshaled state:\n%s", newState.String())
			}
		})
	}

	// Test error cases
	t.Run("unmarshal invalid data", func(t *testing.T) {
		invalidState := &State{}
		err := invalidState.UnmarshalBinary([]byte{1, 2, 3}) // Invalid data
		require.Error(t, err)
	})
}

func TestGameParams_BinaryMarshalUnmarshal(t *testing.T) {
	tests := []struct {
		name   string
		params GameParams
	}{
		{
			name: "standard game params",
			params: GameParams{
				NumPlayers:         3,
				MaxActionsPerRound: 9,
				BtnPos:             0,
				SbAmount:           chips.NewFromInt(1),
				BetSizes:           [][]float32{{0.5, 1.0, 1.5}},
				InitialStacks:      chips.NewList(100, 100, 100),
				TerminalStreet:     River,
				DisableV:           false,
			},
		},
		{
			name: "custom game params",
			params: GameParams{
				NumPlayers:         6,
				MaxActionsPerRound: 18,
				BtnPos:             2,
				SbAmount:           chips.NewFromInt(2),
				BetSizes:           [][]float32{{0.25, 0.5, 0.75, 1.0, 1.5, 2.0}},
				InitialStacks:      chips.NewList(200, 200, 200, 200, 200, 200),
				TerminalStreet:     Turn,
				DisableV:           true,
			},
		},
		{
			name: "minimal game params",
			params: GameParams{
				NumPlayers:         2,
				MaxActionsPerRound: 6,
				BtnPos:             0,
				SbAmount:           chips.NewFromInt(1),
				BetSizes:           [][]float32{{1.0}},
				InitialStacks:      chips.NewList(50, 50),
				TerminalStreet:     Flop,
				DisableV:           false,
			},
		},
		{
			name: "zero values",
			params: GameParams{
				NumPlayers:         0,
				MaxActionsPerRound: 0,
				BtnPos:             0,
				SbAmount:           chips.Zero,
				BetSizes:           [][]float32{},
				InitialStacks:      chips.NewList(),
				TerminalStreet:     NoStreet,
				DisableV:           false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Marshal the params
			data, err := tt.params.MarshalBinary()
			require.NoError(t, err)
			require.NotEmpty(t, data)

			// Unmarshal into new params
			var newParams GameParams
			err = newParams.UnmarshalBinary(data)
			require.NoError(t, err)

			// Verify specific fields
			require.Equal(t, tt.params.NumPlayers, newParams.NumPlayers)
			require.Equal(t, tt.params.MaxActionsPerRound, newParams.MaxActionsPerRound)
			require.Equal(t, tt.params.BtnPos, newParams.BtnPos)
			require.True(t, tt.params.SbAmount.Equal(newParams.SbAmount))
			require.Equal(t, tt.params.BetSizes, newParams.BetSizes)
			require.Equal(t, len(tt.params.InitialStacks), len(newParams.InitialStacks))
			for i := range tt.params.InitialStacks {
				require.True(t, tt.params.InitialStacks[i].Equal(newParams.InitialStacks[i]))
			}
			require.Equal(t, tt.params.TerminalStreet, newParams.TerminalStreet)
			require.Equal(t, tt.params.DisableV, newParams.DisableV)
		})
	}

	// Test error cases
	t.Run("unmarshal invalid data", func(t *testing.T) {
		var invalidParams GameParams
		err := invalidParams.UnmarshalBinary([]byte{1, 2, 3}) // Invalid data
		require.Error(t, err)
	})

	t.Run("verify NewGameParams defaults", func(t *testing.T) {
		params := NewGameParams(3, chips.NewFromInt(100))
		data, err := params.MarshalBinary()
		require.NoError(t, err)

		var newParams GameParams
		err = newParams.UnmarshalBinary(data)
		require.NoError(t, err)
	})
}

func TestHistory(t *testing.T) {
	t.Run("empty history for new game", func(t *testing.T) {
		game, err := NewGame(NewGameParams(2, chips.NewFromInt(100)))
		require.NoError(t, err)

		history := game.Latest.History()
		assert.Len(t, history, 2) // Should have SB and BB actions

		// Verify small blind
		assert.Equal(t, uint8(0), history[0].Pos)
		assert.Equal(t, SmallBlind, history[0].Action.Action)
		assert.Equal(t, chips.New(1), history[0].Action.Amount)

		// Verify big blind
		assert.Equal(t, uint8(1), history[1].Pos)
		assert.Equal(t, BigBlind, history[1].Action.Action)
		assert.Equal(t, chips.New(2), history[1].Action.Amount)
	})

	t.Run("all-in sequence", func(t *testing.T) {
		params := NewGameParams(2, chips.NewFromInt(10))
		game, err := NewGame(params)
		require.NoError(t, err)

		// Player 0 goes all-in
		err = game.Action(DAllIn)
		require.NoError(t, err)

		history := game.Latest.History()
		require.Len(t, history, 3) // SB, BB, All-in

		lastAction := history[len(history)-1]
		assert.Equal(t, AllIn, lastAction.Action.Action)
		assert.Equal(t, uint8(0), lastAction.Pos)
	})

	t.Run("fold sequence", func(t *testing.T) {
		params := NewGameParams(2, chips.NewFromInt(100))
		game, err := NewGame(params)
		require.NoError(t, err)

		// Player 0 folds
		err = game.Action(DFold)
		require.NoError(t, err)

		history := game.Latest.History()
		require.Len(t, history, 3) // SB, BB, Fold

		lastAction := history[len(history)-1]
		assert.Equal(t, Fold, lastAction.Action.Action)
		assert.Equal(t, uint8(0), lastAction.Pos)
		assert.Equal(t, chips.Zero, lastAction.Action.Amount)
	})
}

func TestGame(t *testing.T) {
	// Create game params for 2 players, 100BB each
	params := NewGameParams(2, chips.NewFromInt(100))

	// Create new game
	game, err := NewGame(params)
	require.NoError(t, err)

	// After NewGame, blinds should be posted
	// Check initial state
	require.Equal(t, Preflop, game.Latest.Street)
	require.Equal(t, chips.NewFromInt(1), game.Latest.Players[0].Paid) // SB
	require.Equal(t, chips.NewFromInt(2), game.Latest.Players[1].Paid) // BB
	require.Equal(t, uint8(0), game.Latest.BtnPos)                     // Button position

	// First player (SB) to act preflop
	// Let's make a call
	err = game.Action(ActionAmount{
		Action: Call,
		Amount: chips.NewFromInt(1), // Call BB (2-1=1)
	})
	require.NoError(t, err)

	// BB checks
	err = game.Action(ActionAmount{
		Action: Check,
		Amount: chips.Zero,
	})
	require.NoError(t, err)

	// Should be on flop now
	require.Equal(t, Flop, game.Latest.Street)

	// Let's verify the history
	history := game.Latest.History()
	require.Len(t, history, 4) // SB, BB, Call, Check

	// Verify the sequence
	require.Equal(t, SmallBlind, history[0].Action.Action)
	require.Equal(t, BigBlind, history[1].Action.Action)
	require.Equal(t, Call, history[2].Action.Action)
	require.Equal(t, Check, history[3].Action.Action)

	// Check pot size (should be 4BB)
	totalPot := game.Latest.Players[0].Paid.Add(game.Latest.Players[1].Paid)
	require.Equal(t, chips.NewFromInt(4), totalPot)
}

func TestGameRaise(t *testing.T) {
	params := NewGameParams(2, chips.NewFromInt(100))
	game, err := NewGame(params)
	require.NoError(t, err)

	// SB raises to 6
	err = game.Action(ActionAmount{
		Action: Raise,
		Amount: chips.NewFromInt(5),
	})
	require.NoError(t, err)

	// BB calls the raise
	err = game.Action(ActionAmount{
		Action: Call,
		Amount: chips.NewFromInt(4),
	})
	require.NoError(t, err)

	// Verify history
	history := game.Latest.History()
	require.Len(t, history, 4) // SB, BB, Raise, Call

	// Verify each action in sequence
	require.Equal(t, SmallBlind, history[0].Action.Action)
	require.Equal(t, chips.NewFromInt(1), history[0].Action.Amount)

	require.Equal(t, BigBlind, history[1].Action.Action)
	require.Equal(t, chips.NewFromInt(2), history[1].Action.Amount)

	require.Equal(t, Raise, history[2].Action.Action)
	require.Equal(t, chips.NewFromInt(5), history[2].Action.Amount) // Additional amount

	require.Equal(t, Call, history[3].Action.Action)
	require.Equal(t, chips.NewFromInt(4), history[3].Action.Amount)

	// Verify final state
	require.Equal(t, Flop, game.Latest.Street)
	require.Equal(t, chips.NewFromInt(6), game.Latest.Players[0].Paid)
	require.Equal(t, chips.NewFromInt(6), game.Latest.Players[1].Paid)

	totalPot := game.Latest.Players[0].Paid.Add(game.Latest.Players[1].Paid)
	require.Equal(t, chips.NewFromInt(12), totalPot)
}

func TestMappingPath(t *testing.T) {
	prms := NewGameParams(2, chips.NewFromInt(20_000))
	prms.BetSizes = [][]float32{{0.5, 1, 2}}
	prms.SbAmount = chips.NewFromFloat32(50)

	s, err := NewGame(prms)
	require.NoError(t, err)

	err = s.Action(DiscreteAction(1))
	require.NoError(t, err)

	err = s.Action(DCall)
	require.NoError(t, err)

	err = s.Action(DCheck)
	require.NoError(t, err)

	err = s.Action(DiscreteAction(1))
	require.NoError(t, err)

	err = s.Action(DCall)
	require.NoError(t, err)

	err = s.Action(DCheck)
	require.NoError(t, err)

	err = s.Action(DiscreteAction(1))
	require.NoError(t, err)

	err = s.Action(DCall)
	require.NoError(t, err)

	err = s.Action(DCheck)
	require.NoError(t, err)

	err = s.Action(DiscreteAction(2))
	require.NoError(t, err)

	err = s.Action(DAllIn)
	require.NoError(t, err)

	require.Equal(t, s.Latest.Path(prms.SbAmount), "r:n:b4.00:c:n:k:b8.00:c:n:k:b24.00:c:n:k:b144.00:a")
}

func TestState_Size(t *testing.T) {
	tests := []struct {
		name  string
		state *State
	}{
		{
			name:  "empty state",
			state: NewState(NewGameParams(2, chips.NewFromInt(1))),
		},
		{
			name: "state with players and actions",
			state: func() *State {
				params := NewGameParams(3, chips.NewFromInt(100))
				s := NewState(params)
				s, _ = MakeInitialBets(params, s)
				return s
			}(),
		},
		{
			name: "state with maximum players",
			state: func() *State {
				params := NewGameParams(9, chips.NewFromInt(1000))
				s := NewState(params)
				s, _ = MakeInitialBets(params, s)
				return s
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Get the reported size
			reportedSize := tt.state.Size()

			// Marshal the state
			data, err := tt.state.MarshalBinary()
			require.NoError(t, err)

			// Verify the actual marshaled size matches the reported size
			require.Equal(t, reportedSize, uint64(len(data)),
				"marshaled size (%d) does not match reported size (%d)",
				len(data), reportedSize)

			// Verify we can unmarshal back
			newState := &State{}
			err = newState.UnmarshalBinary(data)
			require.NoError(t, err)

			// Verify the unmarshaled state equals the original
			require.True(t, tt.state.Equal(newState),
				"unmarshaled state does not equal original state")

			// Verify the new state reports the same size
			require.Equal(t, reportedSize, newState.Size(),
				"unmarshaled state reports different size")
		})
	}
}

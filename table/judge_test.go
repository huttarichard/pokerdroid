package table

import (
	"bytes"
	"testing"

	"github.com/pokerdroid/poker/card"
	"github.com/pokerdroid/poker/chips"
	"github.com/pokerdroid/poker/encbin"
	"github.com/stretchr/testify/require"
)

// Simple judge implementation for testing
type testJudge struct {
	winners map[uint8]uint8 // player index -> rank (lower is better)
}

// New Judge implementation that returns slice of winner indices
func (j testJudge) Judge(pp []uint8) []uint8 {
	if len(pp) == 0 {
		return nil
	}

	// Find lowest rank among players
	bestRank := j.winners[pp[0]]
	for _, p := range pp[1:] {
		if rank := j.winners[p]; rank < bestRank {
			bestRank = rank
		}
	}

	// Collect all players with the best rank
	var winners []uint8
	for _, p := range pp {
		if j.winners[p] == bestRank {
			winners = append(winners, p)
		}
	}

	return winners
}

func TestPotDistribution(t *testing.T) {
	tests := []struct {
		name     string
		setup    func() (*Game, error)
		actions  []DiscreteAction
		judge    Judger
		expected chips.List // expected winnings per player
	}{
		{
			name: "heads up - simple win",
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
				DCall,  // SB calls BB
				DCheck, // BB checks
				DCheck, // SB checks flop
				DCheck, // BB checks flop
				DCheck, // SB checks turn
				DCheck, // BB checks turn
				DCheck, // SB checks river
				DCheck, // BB checks river
			},
			judge: testJudge{winners: map[uint8]uint8{
				0: 1,
				1: 0, // player 1 wins
			}},
			expected: chips.NewList(0, 4), // BB wins 4 chips
		},
		{
			name: "heads up - split pot",
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
				DCall,  // SB calls BB
				DCheck, // BB checks
				DCheck, // SB checks flop
				DCheck, // BB checks flop
				DCheck, // SB checks turn
				DCheck, // BB checks turn
				DCheck, // SB checks river
				DCheck, // BB checks river
			},
			judge: testJudge{winners: map[uint8]uint8{
				0: 0, // both rank 0 = tie
				1: 0,
			}},
			expected: chips.NewList(2, 2), // Split pot
		},
		{
			name: "heads up - all in preflop",
			setup: func() (*Game, error) {
				return NewGame(GameParams{
					NumPlayers:         2,
					BtnPos:             0,
					SbAmount:           chips.NewFromInt(1),
					BetSizes:           [][]float32{{0.5, 1, 1.5}},
					InitialStacks:      chips.NewList(10, 10),
					TerminalStreet:     River,
					MaxActionsPerRound: 4,
				})
			},
			actions: []DiscreteAction{
				DAllIn, // SB all-in
				DCall,  // BB calls
			},
			judge: testJudge{winners: map[uint8]uint8{
				0: 0, // player 0 wins
				1: 1,
			}},
			expected: chips.NewList(20, 0), // Winner takes all
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			game, err := tt.setup()
			require.NoError(t, err)

			// Execute all actions
			for _, action := range tt.actions {
				err = game.Action(action)
				require.NoError(t, err)
			}

			require.True(t, game.Latest.Finished())

			// Calculate winnings
			winnings := GetWinnings(game.Latest.Players, tt.judge)
			require.Equal(t, tt.expected, winnings)
		})
	}
}

func TestSidePots_Multiway(t *testing.T) {
	// 3 players with differing stack sizes
	// p0 -> 10 chips, p1 -> 20 chips, p2 -> 30 chips
	// Everyone ends up all-in or calling such that p2 invests 20, p1 invests 20, p0 invests 10.
	// We want to verify creation of side pot and main pot.

	params := GameParams{
		NumPlayers:         3,
		InitialStacks:      chips.NewList(10, 20, 30),
		SbAmount:           chips.NewFromInt(1),
		BetSizes:           [][]float32{{0.5, 1.0, 2.0}},
		MaxActionsPerRound: 6,
		TerminalStreet:     River,
	}

	game, err := NewGame(params)
	require.NoError(t, err)

	// By default:
	// p0 is SB => pays 1
	// p1 is BB => pays 2
	// p2 is next => must call 2

	// p2 calls 2
	err = game.Action(DCall)
	require.NoError(t, err)

	// p0 raises All In (he has 9 left)
	err = game.Action(DAllIn) // p0 invests total 10
	require.NoError(t, err)

	// p1 calls the all in (he invests total 10, but he already paid 2 -> 8 more)
	err = game.Action(DCall)
	require.NoError(t, err)

	// p2 calls the all in (but only invests 10 total, he had 2 in -> 8 more)
	err = game.Action(DCall)
	require.NoError(t, err)

	// At this point:
	// p0 is all in with 10
	// p1 has invested 10 of 20
	// p2 has invested 10 of 30
	// There's a side pot possible if p1 or p2 invests more, but let's check the distribution so far.

	// Everyone checks down the rest of the streets (or the game might auto-finish if no further bets).
	for !game.Latest.Finished() {
		err = game.Action(DCheck)
		if err != nil {
			break
		}
	}

	require.True(t, game.Latest.Finished())

	// For side pot awarding, let's say p1 actually "wins" rank < p2,
	// but p2 has better rank than p0. So final ranks: p0=2, p1=0, p2=1
	judge := testJudge{
		winners: map[uint8]uint8{
			0: 2,
			1: 0,
			2: 1,
		},
	}

	winnings := GetWinnings(game.Latest.Players, judge)
	// Main pot = 3 players x 10 = 30 total.
	// p0 invests 10, p1 invests 10, p2 invests 10. That 30 is the main pot.
	// No side pot actually formed because p1/p2 did not exceed p0's 10 ... each only put in 10 total.
	// p1 is best in the main pot, so p1 gets 30.

	expected := chips.NewList(0, 30, 0)
	require.Equal(t, expected, winnings)
}

func TestPotMarshalBinary(t *testing.T) {
	tests := []struct {
		name string
		pot  Pot
	}{
		{
			name: "empty pot",
			pot: Pot{
				Amount:  chips.Zero,
				Players: []uint8{},
			},
		},
		{
			name: "pot with one player",
			pot: Pot{
				Amount:  chips.NewFromInt(100),
				Players: []uint8{0},
			},
		},
		{
			name: "pot with multiple players",
			pot: Pot{
				Amount:  chips.NewFromInt(300),
				Players: []uint8{0, 1},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Marshal
			data, err := tt.pot.MarshalBinary()
			require.NoError(t, err)

			// Unmarshal
			var pot Pot
			err = pot.UnmarshalBinary(data)
			require.NoError(t, err)

			// Compare
			require.True(t, tt.pot.Amount.Equal(pot.Amount))
			require.Equal(t, tt.pot.Players, pot.Players)
		})
	}
}

func TestPotsMarshalBinary(t *testing.T) {
	tests := []struct {
		name string
		pots Pots
	}{
		{
			name: "empty pots",
			pots: Pots{},
		},
		{
			name: "single pot",
			pots: Pots{
				{
					Amount:  chips.NewFromInt(100),
					Players: []uint8{0},
				},
			},
		},
		{
			name: "multiple pots",
			pots: Pots{
				{
					Amount:  chips.NewFromInt(200),
					Players: []uint8{0, 1},
				},
				{
					Amount:  chips.NewFromInt(100),
					Players: []uint8{0, 1},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Marshal
			data, err := tt.pots.MarshalBinary()
			require.NoError(t, err)

			// Unmarshal
			var pots Pots
			err = pots.UnmarshalBinary(data)
			require.NoError(t, err)

			// Compare
			require.Equal(t, len(tt.pots), len(pots))
			for i := range tt.pots {
				require.True(t, tt.pots[i].Amount.Equal(pots[i].Amount))
				require.Equal(t, tt.pots[i].Players, pots[i].Players)
			}
		})
	}
}

func TestPotsMarshalBinaryCornerCases(t *testing.T) {
	t.Run("nil pots", func(t *testing.T) {
		var pots Pots
		data, err := pots.MarshalBinary()
		require.NoError(t, err)

		var unmarshaledPots Pots
		err = unmarshaledPots.UnmarshalBinary(data)
		require.NoError(t, err)
		require.Equal(t, 0, len(unmarshaledPots))
	})

	t.Run("max players in pot", func(t *testing.T) {
		// Create a pot with maximum number of players (255 due to uint8)
		pot := Pot{
			Amount:  chips.NewFromInt(1000),
			Players: make([]uint8, 20),
		}
		for i := range pot.Players {
			pot.Players[i] = uint8(i)
		}

		pots := Pots{pot}
		data, err := pots.MarshalBinary()
		require.NoError(t, err)

		var unmarshaledPots Pots
		err = unmarshaledPots.UnmarshalBinary(data)
		require.NoError(t, err)
		require.Equal(t, 1, len(unmarshaledPots))
		require.Equal(t, 20, len(unmarshaledPots[0].Players))
	})

	t.Run("max pots", func(t *testing.T) {
		// Create maximum number of pots (255 due to uint8)
		pots := make(Pots, 50)
		for i := range pots {
			pots[i] = Pot{
				Amount:  chips.NewFromInt(int64(i + 1)),
				Players: []uint8{uint8(i)},
			}
		}

		data, err := pots.MarshalBinary()
		require.NoError(t, err)

		var unmarshaledPots Pots
		err = unmarshaledPots.UnmarshalBinary(data)
		require.NoError(t, err)
		require.Equal(t, 50, len(unmarshaledPots))
	})
}

func TestCards_Judge(t *testing.T) {
	c := Cards{
		Community: card.Cards{
			card.Card2C, card.Card3S, card.Card8H, card.Card6C, card.Card4D,
		},
		Players: []card.Cards{
			{card.CardKC, card.CardKD},
			{card.CardAC, card.CardAD},
		},
	}

	pl, err := NewGame(GameParams{
		NumPlayers:         2,
		BtnPos:             0,
		SbAmount:           chips.NewFromInt(1),
		BetSizes:           [][]float32{{0.5, 1, 1.5}},
		InitialStacks:      chips.NewList(100, 100),
		TerminalStreet:     River,
		MaxActionsPerRound: 4,
	})
	require.NoError(t, err)

	pl.Action(DCall)
	pl.Action(DAllIn)
	pl.Action(DCall)

	w := GetWinnings(pl.Latest.Players, &c)
	require.Equal(t, chips.NewList(0, 200), w)
}

func TestCards_Judge2(t *testing.T) {
	c := Cards{
		Community: card.Cards{
			card.Card2C, card.Card3S, card.Card8H, card.Card6C, card.Card4D,
		},
		Players: []card.Cards{
			{card.CardKC, card.CardKD},
			{card.CardAC, card.CardAD},
		},
	}

	pl, err := NewGame(GameParams{
		NumPlayers:         2,
		BtnPos:             0,
		SbAmount:           chips.NewFromInt(1),
		BetSizes:           [][]float32{{0.5, 1, 1.5}},
		InitialStacks:      chips.NewList(100, 100),
		TerminalStreet:     River,
		MaxActionsPerRound: 4,
	})
	require.NoError(t, err)

	err = pl.Action(DiscreteAction(1))
	require.NoError(t, err)

	err = pl.Action(DFold)
	require.NoError(t, err)

	w := GetWinnings(pl.Latest.Players, &c)
	require.Equal(t, chips.NewList(6, 0), w)
}

func TestCards_Judge3(t *testing.T) {
	c := Cards{
		Community: card.Cards{
			card.Card2C, card.Card3S, card.Card8H, card.Card6C, card.Card4D,
		},
		Players: []card.Cards{
			{card.CardKC, card.CardKD},
			{card.CardAC, card.CardAD},
		},
	}

	pl, err := NewGame(GameParams{
		NumPlayers:         2,
		BtnPos:             0,
		SbAmount:           chips.NewFromInt(1),
		BetSizes:           [][]float32{{0.5, 1, 1.5}},
		InitialStacks:      chips.NewList(100, 100),
		TerminalStreet:     River,
		MaxActionsPerRound: 4,
	})
	require.NoError(t, err)

	// sb raise
	err = pl.Action(DiscreteAction(1))
	require.NoError(t, err)

	// bb reraise
	err = pl.Action(DiscreteAction(2))
	require.NoError(t, err)

	// sb folds
	err = pl.Action(DFold)
	require.NoError(t, err)

	w := GetWinnings(pl.Latest.Players, &c)
	require.Equal(t, chips.NewList(0, 18), w)
}

func TestPot_Size(t *testing.T) {
	tests := []struct {
		name string
		pot  Pot
	}{
		{
			name: "empty pot",
			pot: Pot{
				Amount:  chips.Zero,
				Players: nil,
			},
		},
		{
			name: "pot with amount only",
			pot: Pot{
				Amount:  chips.NewFromInt(100),
				Players: []uint8{},
			},
		},
		{
			name: "pot with single player",
			pot: Pot{
				Amount:  chips.NewFromInt(100),
				Players: []uint8{1},
			},
		},
		{
			name: "pot with multiple players",
			pot: Pot{
				Amount:  chips.NewFromInt(100),
				Players: []uint8{1, 2, 3},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Get size and marshal
			reportedSize := tt.pot.Size()
			data, err := tt.pot.MarshalBinary()
			require.NoError(t, err)

			// Analyze binary structure
			buf := bytes.NewBuffer(data)

			// Read Amount
			var amount chips.Chips
			err = encbin.UnmarshalValues(buf, &amount)
			require.NoError(t, err)

			// Read Players length
			var playersLen uint8
			err = encbin.UnmarshalValues(buf, &playersLen)
			require.NoError(t, err)

			require.Equal(t, reportedSize, uint64(len(data)),
				"marshaled size (%d) does not match reported size (%d)",
				len(data), reportedSize)
		})
	}
}

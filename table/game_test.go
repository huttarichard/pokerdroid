package table

import (
	"testing"

	"github.com/pokerdroid/poker/chips"
	"github.com/stretchr/testify/require"
)

func TestGameParams_Size(t *testing.T) {
	tests := []struct {
		name   string
		params GameParams
	}{
		{
			name:   "minimal game params",
			params: NewGameParams(2, chips.NewFromInt(1)),
		},
		{
			name: "custom bet sizes",
			params: func() GameParams {
				p := NewGameParams(2, chips.NewFromInt(100))
				p.BetSizes = [][]float32{{0.5, 1.0, 1.5, 2.0}}
				return p
			}(),
		},
		{
			name: "max players",
			params: func() GameParams {
				p := NewGameParams(9, chips.NewFromInt(1000))
				p.BetSizes = [][]float32{{0.25, 0.5, 0.75, 1.0, 1.5, 2.0}}
				return p
			}(),
		},
		{
			name: "empty bet sizes",
			params: func() GameParams {
				p := NewGameParams(2, chips.NewFromInt(100))
				p.BetSizes = [][]float32{}
				return p
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Get the reported size
			reportedSize := tt.params.Size()

			// Marshal the params
			data, err := tt.params.MarshalBinary()
			require.NoError(t, err)

			// Verify the actual marshaled size matches the reported size
			require.Equal(t, reportedSize, uint64(len(data)),
				"marshaled size (%d) does not match reported size (%d)",
				len(data), reportedSize)

			// Verify we can unmarshal back
			newParams := GameParams{}
			err = newParams.UnmarshalBinary(data)
			require.NoError(t, err)

			// Verify the new params reports the same size
			require.Equal(t, reportedSize, newParams.Size(),
				"unmarshaled params reports different size")

			// Debug output if sizes don't match
			if reportedSize != uint64(len(data)) {
				t.Logf("Params details:")
				t.Logf("NumPlayers: %d", tt.params.NumPlayers)
				t.Logf("BetSizes: %d", len(tt.params.BetSizes))
				t.Logf("InitialStacks: %d", len(tt.params.InitialStacks))
			}
		})
	}
}

func TestGameParams_SizeEdgeCases(t *testing.T) {
	tests := []struct {
		name   string
		params GameParams
	}{
		{
			name: "nil slices",
			params: GameParams{
				NumPlayers:    2,
				BetSizes:      nil,
				InitialStacks: nil,
			},
		},
		{
			name: "empty slices",
			params: GameParams{
				NumPlayers:    2,
				BetSizes:      [][]float32{},
				InitialStacks: chips.List{},
			},
		},
		{
			name: "zero values",
			params: GameParams{
				NumPlayers:    0,
				BetSizes:      [][]float32{{0}},
				InitialStacks: chips.List{0},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reportedSize := tt.params.Size()
			data, err := tt.params.MarshalBinary()
			require.NoError(t, err)
			require.Equal(t, reportedSize, uint64(len(data)))
		})
	}
}

func TestGameMarshal(t *testing.T) {
	tests := []struct {
		name   string
		params GameParams
	}{
		{
			name:   "default params",
			params: NewGameParams(2, chips.NewFromInt(100)),
		},
		{
			name: "custom bet sizes",
			params: GameParams{
				NumPlayers:    3,
				BetSizes:      [][]float32{{0.5, 1}, {1, 2}},
				InitialStacks: chips.List{100, 200, 300},
			},
		},
		{
			name: "min bet enabled",
			params: GameParams{
				NumPlayers:    2,
				MinBet:        true,
				BetSizes:      [][]float32{{1}},
				InitialStacks: chips.List{100, 100},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Marshal
			data, err := tt.params.MarshalBinary()
			require.NoError(t, err)

			// Unmarshal
			var unmarshaled GameParams
			err = unmarshaled.UnmarshalBinary(data)
			require.NoError(t, err)

			// Verify fields match
			require.Equal(t, tt.params.NumPlayers, unmarshaled.NumPlayers)
			require.Equal(t, tt.params.MaxActionsPerRound, unmarshaled.MaxActionsPerRound)
			require.Equal(t, tt.params.BtnPos, unmarshaled.BtnPos)
			require.Equal(t, tt.params.SbAmount, unmarshaled.SbAmount)
			require.Equal(t, tt.params.TerminalStreet, unmarshaled.TerminalStreet)
			require.Equal(t, tt.params.DisableV, unmarshaled.DisableV)
			require.Equal(t, tt.params.MinBet, unmarshaled.MinBet)
			require.Equal(t, tt.params.BetSizes, unmarshaled.BetSizes)
			require.Equal(t, tt.params.InitialStacks, unmarshaled.InitialStacks)
		})
	}
}

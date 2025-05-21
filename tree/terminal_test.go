package tree

import (
	"testing"

	"github.com/pokerdroid/poker/chips"
	"github.com/pokerdroid/poker/table"
	"github.com/stretchr/testify/require"
)

func TestTerminal_Basic(t *testing.T) {
	// Create test terminal
	terminal := &Terminal{
		Parent: &Root{},
		Pots:   table.Pots{{Amount: chips.NewFromInt(100)}},
		Players: table.Players{
			{Paid: chips.NewFromInt(50), Status: table.StatusActive},
			{Paid: chips.NewFromInt(50), Status: table.StatusFolded},
		},
	}

	// Test basic properties
	require.Equal(t, NodeKindTerminal, terminal.Kind())
	require.Equal(t, terminal.Parent, terminal.GetParent())
	require.Equal(t, 1, len(terminal.Pots))
	require.Equal(t, 2, len(terminal.Players))
}

func TestTerminal_MarshalBinary(t *testing.T) {
	tests := []struct {
		name     string
		terminal *Terminal
		wantErr  bool
	}{
		{
			name: "basic terminal",
			terminal: &Terminal{
				Pots: table.Pots{{Amount: chips.NewFromInt(100)}},
				Players: table.Players{
					{Paid: chips.NewFromInt(50), Status: table.StatusActive},
				},
			},
			wantErr: false,
		},
		{
			name: "empty pots",
			terminal: &Terminal{
				Pots:    table.Pots{},
				Players: table.Players{{Status: table.StatusActive}},
			},
			wantErr: false,
		},
		{
			name: "empty players",
			terminal: &Terminal{
				Pots:    table.Pots{{Amount: chips.NewFromInt(100)}},
				Players: table.Players{},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Marshal
			data, err := tt.terminal.MarshalBinary()
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, data)

			// Unmarshal into new terminal
			newTerminal := &Terminal{}
			err = newTerminal.UnmarshalBinary(data)
			require.NoError(t, err)

			// Verify fields
			require.Equal(t, tt.terminal.Pots, newTerminal.Pots)
			require.Equal(t, tt.terminal.Players, newTerminal.Players)
		})
	}
}

func TestTerminal_UnmarshalBinary_InvalidKind(t *testing.T) {
	terminal := &Terminal{}
	err := terminal.UnmarshalBinary([]byte{byte(NodeKindPlayer)}) // Wrong node kind
	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid node kind")
}

func TestTerminal_RoundTrip(t *testing.T) {
	// Create a complex terminal node
	original := &Terminal{
		Pots: table.Pots{
			{Amount: chips.NewFromInt(100)},
			{Amount: chips.NewFromInt(50)},
		},
		Players: table.Players{
			{Paid: chips.NewFromInt(75), Status: table.StatusActive},
			{Paid: chips.NewFromInt(50), Status: table.StatusFolded},
			{Paid: chips.NewFromInt(25), Status: table.StatusAllIn},
		},
	}

	// Marshal
	data, err := original.MarshalBinary()
	require.NoError(t, err)

	// Unmarshal
	reconstructed := &Terminal{}
	err = reconstructed.UnmarshalBinary(data)
	require.NoError(t, err)

	// Verify key properties
	require.Equal(t, len(original.Pots), len(reconstructed.Pots))
	require.Equal(t, len(original.Players), len(reconstructed.Players))

	// Verify pots
	for i := range original.Pots {
		require.Equal(t, original.Pots[i].Amount, reconstructed.Pots[i].Amount)
	}

	// Verify players
	for i := range original.Players {
		require.Equal(t, original.Players[i].Paid, reconstructed.Players[i].Paid)
		require.Equal(t, original.Players[i].Status, reconstructed.Players[i].Status)
	}
}

func TestTerminal_Size(t *testing.T) {
	tests := []struct {
		name    string
		setup   func() *Terminal
		players int
		pots    int
	}{
		{
			name: "empty terminal",
			setup: func() *Terminal {
				return &Terminal{
					Players: table.Players{},
					Pots:    table.Pots{},
				}
			},
			players: 0,
			pots:    0,
		},
		{
			name: "single player no pots",
			setup: func() *Terminal {
				return &Terminal{
					Players: table.Players{
						{Status: table.StatusActive, Paid: chips.NewFromInt(100)},
					},
					Pots: table.Pots{},
				}
			},
			players: 1,
			pots:    0,
		},
		{
			name: "multiple players with pots",
			setup: func() *Terminal {
				return &Terminal{
					Players: table.Players{
						{Status: table.StatusActive, Paid: chips.NewFromInt(100)},
						{Status: table.StatusFolded, Paid: chips.NewFromInt(50)},
						{Status: table.StatusAllIn, Paid: chips.NewFromInt(200)},
					},
					Pots: table.Pots{
						{Amount: chips.NewFromInt(150)},
						{Amount: chips.NewFromInt(200)},
					},
				}
			},
			players: 3,
			pots:    2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			terminal := tt.setup()

			reportedSize := terminal.Size()
			data, err := terminal.MarshalBinary()
			require.NoError(t, err)

			require.Equal(t, reportedSize, uint64(len(data)),
				"marshaled size (%d) does not match reported size (%d)",
				len(data), reportedSize)

			// Test round trip
			newTerminal := &Terminal{}
			err = newTerminal.UnmarshalBinary(data)
			require.NoError(t, err)

			// Verify basic properties
			require.Equal(t, len(terminal.Players), len(newTerminal.Players))
			require.Equal(t, len(terminal.Pots), len(newTerminal.Pots))

			// Compare players
			for i, player := range terminal.Players {
				require.Equal(t, player.Status, newTerminal.Players[i].Status)
				require.Equal(t, player.Paid, newTerminal.Players[i].Paid)
			}

			// Compare pots
			for i, pot := range terminal.Pots {
				require.Equal(t, pot.Amount, newTerminal.Pots[i].Amount)
			}
		})
	}
}

func TestTerminal_LastAlive(t *testing.T) {
	tests := []struct {
		name     string
		terminal *Terminal
		playerID uint8
		expected bool
	}{
		{
			name: "single player alive",
			terminal: &Terminal{
				Players: table.Players{
					{Status: table.StatusActive},
					{Status: table.StatusFolded},
					{Status: table.StatusFolded},
				},
			},
			playerID: 0,
			expected: true,
		},
		{
			name: "multiple players alive",
			terminal: &Terminal{
				Players: table.Players{
					{Status: table.StatusActive},
					{Status: table.StatusActive},
					{Status: table.StatusFolded},
				},
			},
			playerID: 0,
			expected: false,
		},
		{
			name: "player folded",
			terminal: &Terminal{
				Players: table.Players{
					{Status: table.StatusFolded},
					{Status: table.StatusActive},
					{Status: table.StatusFolded},
				},
			},
			playerID: 0,
			expected: false,
		},
		{
			name: "all folded except last",
			terminal: &Terminal{
				Players: table.Players{
					{Status: table.StatusFolded},
					{Status: table.StatusFolded},
					{Status: table.StatusActive},
				},
			},
			playerID: 2,
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.terminal.Players.LastAlive(tt.playerID)
			require.Equal(t, tt.expected, result)
		})
	}
}

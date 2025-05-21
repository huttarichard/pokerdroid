package tree

import (
	"testing"

	"github.com/pokerdroid/poker/chips"
	"github.com/pokerdroid/poker/table"
	"github.com/stretchr/testify/require"
)

func TestNewFullTree(t *testing.T) {
	params := NewFullTreeParams{
		BigBlind:   chips.NewFromInt(100),
		NumPlayers: 2,
		Betting:    [][]float32{},
		MaxActions: 2,
		Terminal:   table.Flop,
		MinBet:     false,
		Limp:       true,
	}

	root, err := NewFullTree(params)
	require.NoError(t, err)

	// Expected payouts at terminal nodes using table2.Players
	m := map[string]table.Players{
		"r:n:a:f:t": {
			{Paid: chips.NewFromInt(100), Status: table.StatusAllIn},
			{Paid: chips.NewFromInt(2), Status: table.StatusFolded},
		},
		"r:n:a:c:t": {
			{Paid: chips.NewFromInt(100), Status: table.StatusAllIn},
			{Paid: chips.NewFromInt(100), Status: table.StatusAllIn},
		},
		"r:n:f:t": {
			{Paid: chips.NewFromFloat(1), Status: table.StatusFolded},
			{Paid: chips.NewFromInt(2), Status: table.StatusActive},
		},
		"r:n:c:k:n:k:k:t": {
			{Paid: chips.NewFromInt(2), Status: table.StatusActive},
			{Paid: chips.NewFromInt(2), Status: table.StatusActive},
		},
		"r:n:c:k:n:a:f:t": {
			{Paid: chips.NewFromInt(2), Status: table.StatusFolded},
			{Paid: chips.NewFromInt(100), Status: table.StatusAllIn},
		},
		"r:n:c:k:n:a:c:t": {
			{Paid: chips.NewFromInt(100), Status: table.StatusAllIn},
			{Paid: chips.NewFromInt(100), Status: table.StatusAllIn},
		},
	}

	// Expected pot sizes at terminal nodes
	pots := map[string]float32{
		"r:n:a:f:t":       102,
		"r:n:a:c:t":       200,
		"r:n:f:t":         3,
		"r:n:c:k:n:k:k:t": 4,
		"r:n:c:k:n:a:f:t": 102,
		"r:n:c:k:n:a:c:t": 200,
	}

	leafs := FindLeafNodes(root)
	// require.Equal(t, len(m), len(leafs))

	for _, r := range leafs {
		x := r.(*Terminal)
		path := GetPath(x).String()
		require.Equal(t, m[path], x.Players)
		require.Equal(t, pots[path], x.Pots.Sum().Float32())
	}

	// Test LastAlive status
	for path, players := range m {
		x := findLeafByPath(leafs, path)
		require.NotNil(t, x)

		for i := range players {
			require.Equal(t,
				players[i].Status != table.StatusFolded &&
					isLastActivePlayer(players, uint8(i)),
				x.Players.LastAlive(uint8(i)))
		}
	}
}

// Helper function to find leaf node by path
func findLeafByPath(leafs []Node, path string) *Terminal {
	for _, r := range leafs {
		if x, ok := r.(*Terminal); ok {
			if GetPath(x).String() == path {
				return x
			}
		}
	}
	return nil
}

// Helper function to check if player is last active
func isLastActivePlayer(players table.Players, id uint8) bool {
	if players[id].Status == table.StatusFolded {
		return false
	}
	for i, x := range players {
		if uint8(i) == id {
			continue
		}
		if x.Status != table.StatusFolded {
			return false
		}
	}
	return true
}

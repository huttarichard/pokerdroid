package tree

import (
	"bytes"
	"testing"

	"github.com/pokerdroid/poker/abs"
	"github.com/pokerdroid/poker/chips"
	"github.com/pokerdroid/poker/policy"
	"github.com/pokerdroid/poker/table"
	"github.com/stretchr/testify/require"
)

func createTestTree() *Root {
	// Create terminal node
	terminal := &Terminal{
		Pots: table.Pots{{
			Amount:  chips.NewFromInt(100),
			Players: []uint8{0, 1},
		}},
		Players: table.Players{
			{Paid: chips.NewFromInt(50), Status: table.StatusActive},
			{Paid: chips.NewFromInt(50), Status: table.StatusActive},
		},
	}

	// Create player node with actions
	player := &Player{
		TurnPos: 1,
		State: &table.State{
			Street:  table.Flop,
			TurnPos: 1,
			BtnPos:  0,
		},
	}

	pm := NewPolicies()

	// Add a test policy
	pol := policy.New(2) // 2 actions
	pol.RegretSum = []float64{1.0, -1.0}
	pol.StrategySum = []float64{0.6, 0.4}
	pol.Baseline = []float64{0.1, -0.1}
	pol.BuildStrategy()

	pm.Store(abs.Cluster(1), pol)
	pm.Store(abs.Cluster(2), pol)

	// Create actions for player
	actions := &PlayerActions{
		Parent:   player,
		Actions:  []table.DiscreteAction{table.DCall, table.DFold},
		Nodes:    []Node{terminal, nil},
		Policies: pm,
	}
	player.Actions = actions

	// Create chance node
	chance := &Chance{
		Next:   player,
		Parent: nil,
		State: &table.State{
			Street:  table.Preflop,
			TurnPos: 0,
			BtnPos:  0,
		},
	}

	// Create root node
	root := &Root{
		States: 3,
		Next:   chance,
		Params: table.GameParams{
			NumPlayers:         2,
			MaxActionsPerRound: 2,
			BtnPos:             0,
			SbAmount:           chips.NewFromInt(1),
			BetSizes:           [][]float32{{1.0}},
			InitialStacks:      chips.List{chips.NewFromInt(100), chips.NewFromInt(100)},
			TerminalStreet:     table.River,
		},
		Iteration: 1,
		State: &table.State{
			Street:  table.Preflop,
			TurnPos: 0,
			BtnPos:  0,
		},
	}

	// Set up parent relationships
	terminal.Parent = player
	player.Parent = chance
	chance.Parent = root

	return root
}

func TestNodeMarshalUnmarshal(t *testing.T) {
	original := createTestTree()

	// Marshal the tree
	data, err := original.MarshalBinary()
	require.NoError(t, err)

	// Unmarshal into a new tree
	unmarshaled := &Root{}
	err = unmarshaled.UnmarshalBinary(data)
	require.NoError(t, err)

	require.Equal(t, original.Size(), unmarshaled.Size())

	// Helper function to compare nodes and fix parent relationships
	var compareNodes func(n1, n2 Node, parent Node)
	compareNodes = func(n1, n2 Node, parent Node) {
		if n1 == nil || n2 == nil {
			require.Equal(t, n1, n2)
			return
		}

		// Set parent relationship
		switch x2 := n2.(type) {
		case *Terminal:
			x2.Parent = parent
		case *Player:
			x2.Parent = parent
			if x2.Actions != nil {
				x2.Actions.Parent = x2
			}
		case *Chance:
			x2.Parent = parent
		}

		require.Equal(t, n1.Kind(), n2.Kind())

		switch x1 := n1.(type) {
		case *Root:
			x2 := n2.(*Root)
			require.Equal(t, x1.States, x2.States)
			require.Equal(t, x1.Params, x2.Params)
			require.Equal(t, x1.Iteration, x2.Iteration)
			require.Equal(t, x1.State, x2.State)
			compareNodes(x1.Next, x2.Next, x2)

		case *Chance:
			x2 := n2.(*Chance)
			require.Equal(t, x1.State, x2.State)
			compareNodes(x1.Next, x2.Next, x2)

		case *Player:
			x2 := n2.(*Player)
			require.Equal(t, x1.TurnPos, x2.TurnPos)
			require.Equal(t, x1.State, x2.State)

			if x1.Actions == nil || x2.Actions == nil {
				require.Equal(t, x1.Actions, x2.Actions)
				return
			}

			require.Equal(t, len(x1.Actions.Actions), len(x2.Actions.Actions))
			require.Equal(t, x1.Actions.Actions, x2.Actions.Actions)
			require.True(t, x1.Actions.Policies.Equal(x2.Actions.Policies))

			if x1.Actions.Nodes == nil || x2.Actions.Nodes == nil {
				require.Equal(t, x1.Actions.Nodes, x2.Actions.Nodes)
				return
			}

			require.Equal(t, len(x1.Actions.Nodes), len(x2.Actions.Nodes))
			for i := range x1.Actions.Nodes {
				compareNodes(x1.Actions.Nodes[i], x2.Actions.Nodes[i], x2)
			}

		case *Terminal:
			x2 := n2.(*Terminal)
			require.Equal(t, x1.Pots, x2.Pots)
			require.Equal(t, x1.Players, x2.Players)
		}
	}

	// Compare entire trees and fix parent relationships
	compareNodes(original, unmarshaled, nil)
}

func TestNodeMarshalUnmarshalEdgeCases(t *testing.T) {
	t.Run("empty root", func(t *testing.T) {
		root := &Root{}
		data, err := root.MarshalBinary()
		require.NoError(t, err)

		unmarshaled := &Root{}
		err = unmarshaled.UnmarshalBinary(data)
		require.NoError(t, err)
		require.Equal(t, root.States, unmarshaled.States)
		require.Equal(t, root.Size(), unmarshaled.Size())
	})

	t.Run("nil nodes", func(t *testing.T) {
		root := &Root{
			Next: &Chance{
				Next: nil,
			},
		}
		data, err := root.MarshalBinary()
		require.NoError(t, err)

		unmarshaled := &Root{}
		err = unmarshaled.UnmarshalBinary(data)
		require.NoError(t, err)
		require.NotNil(t, unmarshaled.Next)
		require.Nil(t, unmarshaled.Next.(*Chance).Next)
	})

	t.Run("invalid node kind", func(t *testing.T) {
		data := []byte{255} // Invalid node kind
		root := &Root{}
		err := root.UnmarshalBinary(data)
		require.Error(t, err)
	})
}

func TestNodeMarshalUnmarshalNilCases(t *testing.T) {
	t.Run("nil actions in player", func(t *testing.T) {
		player := &Player{
			TurnPos: 1,
			State: &table.State{
				Street:  table.Flop,
				TurnPos: 1,
			},
			Actions: nil, // Explicitly nil actions
		}

		data, err := player.MarshalBinary()
		require.NoError(t, err)

		unmarshaled := &Player{}
		err = unmarshaled.UnmarshalBinary(data)

		require.NoError(t, err)
		require.Nil(t, unmarshaled.Actions)
	})

	t.Run("nil nodes in actions", func(t *testing.T) {
		pm := NewPolicies()

		// Add a test policy
		pol := policy.New(2) // 2 actions
		pol.Strategy = []float64{0.7, 0.3}
		pol.RegretSum = []float64{1.0, -1.0}
		pol.StrategySum = []float64{0.6, 0.4}
		pol.Baseline = []float64{0.1, -0.1}

		pm.Store(abs.Cluster(1), pol)

		actions := &PlayerActions{
			Actions:  []table.DiscreteAction{table.DCall, table.DFold},
			Nodes:    []Node{nil, nil}, // Explicitly nil nodes
			Policies: pm,
		}

		data, err := actions.MarshalBinary()
		require.NoError(t, err)

		unmarshaled := &PlayerActions{}
		err = unmarshaled.UnmarshalBinary(data)

		require.NoError(t, err)
		require.Equal(t, actions.Nodes, unmarshaled.Nodes)

		require.Equal(t, actions.Size(), unmarshaled.Size())
	})

	t.Run("nil next in chance", func(t *testing.T) {
		chance := &Chance{
			Next: nil,
			State: &table.State{
				Street:  table.Preflop,
				TurnPos: 0,
			},
		}

		data, err := chance.MarshalBinary()
		require.NoError(t, err)

		unmarshaled := &Chance{}
		err = unmarshaled.UnmarshalBinary(data)
		require.NoError(t, err)
		require.Nil(t, unmarshaled.Next)
	})

	t.Run("nil state in nodes", func(t *testing.T) {
		chance := &Chance{State: nil}
		data, err := chance.MarshalBinary()
		require.NoError(t, err)

		unmarshaled := &Chance{}
		err = unmarshaled.UnmarshalBinary(data)
		require.NoError(t, err)
		require.Nil(t, unmarshaled.State)
	})
}

func TestParentRelationships(t *testing.T) {
	// Create original tree
	original := createTestTree()

	// Marshal and unmarshal
	data, err := original.MarshalBinary()
	require.NoError(t, err)

	unmarshaled := &Root{}
	err = unmarshaled.UnmarshalBinary(data)
	require.NoError(t, err)

	// Find all leaf nodes in unmarshaled tree
	leaves := FindLeafNodes(unmarshaled)
	require.NotEmpty(t, leaves, "Should have found leaf nodes")

	// For each leaf node, traverse up to root and verify relationships
	for _, leaf := range leaves {
		current := leaf
		var path []NodeKind

		// Traverse up until we hit root
		for current != nil {
			path = append(path, current.Kind())

			// Special handling for Player nodes to verify Actions parent
			if player, ok := current.(*Player); ok {
				if player.Actions != nil {
					require.Equal(t, player, player.Actions.Parent,
						"Player's Actions should have Player as parent")
				}
			}

			current = current.GetParent()
		}

		// Verify path ends at root
		require.Greater(t, len(path), 1, "Path should have multiple nodes")
		require.Equal(t, NodeKindRoot, path[len(path)-1],
			"Path should end at root node")

		// Verify expected path structure
		// In our test tree: Terminal -> Player -> Chance -> Root
		expectedPath := []NodeKind{
			NodeKindTerminal,
			NodeKindPlayer,
			NodeKindChance,
			NodeKindRoot,
		}
		require.Equal(t, expectedPath, path,
			"Path from leaf to root doesn't match expected structure")
	}
}

func TestNodeWriteRead(t *testing.T) {
	r1 := createTestTree()

	// Write the original tree using WriteBinary
	var buf bytes.Buffer
	err := r1.WriteBinary(&buf)
	require.NoError(t, err, "WriteBinary failed")
	require.Greater(t, buf.Len(), 1, "WriteBinary produced empty buffer")

	mx, err := r1.MarshalBinary()
	require.NoError(t, err)

	require.Equal(t, len(mx), buf.Len())

	r2 := &Root{}
	err = r2.UnmarshalBinary(buf.Bytes())
	require.NoError(t, err)

	// Read the tree back using ReadBinary
	readSeeker := bytes.NewReader(mx)
	r3, err := NewRootFromReadSeeker(readSeeker)
	require.NoError(t, err)

	r1x, err := r1.MarshalBinary()
	require.NoError(t, err)

	r2x, err := r2.MarshalBinary()
	require.NoError(t, err)

	r3x, err := r3.MarshalBinary()
	require.NoError(t, err)

	require.Equal(t, r1x, r2x)
	require.Equal(t, r1x, r3x)
	require.Equal(t, r2x, r3x)
}

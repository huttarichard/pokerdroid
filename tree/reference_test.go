package tree_test

import (
	"bytes"
	"testing"

	"github.com/pokerdroid/poker/chips"
	"github.com/pokerdroid/poker/table"
	"github.com/pokerdroid/poker/tree"
	"github.com/stretchr/testify/require"
)

// TestTerminalReadBinary tests that a Terminal node can be marshaled
// into a buffer and then re-read correctly using ReadBinary.
func TestTerminalReadBinary(t *testing.T) {
	// Construct a Terminal node with some dummy values.
	origTerminal := &tree.Terminal{
		// Assuming table.Pots is defined as an integer type.
		Pots: table.Pots{
			{Amount: chips.NewFromInt(100)},
		},
		// Create two players with dummy values.
		Players: table.Players{
			{Paid: chips.NewFromInt(50), Status: table.StatusActive},
			{Paid: chips.NewFromInt(75), Status: table.StatusFolded},
		},
		Parent: nil,
	}

	// Marshal the Terminal node to binary.
	data, err := origTerminal.MarshalBinary()
	require.NoError(t, err)

	// Create a reader from the byte slice.
	reader := bytes.NewReader(data)

	// Read the binary data back into a fresh Terminal node.
	newTerminal := &tree.Terminal{}
	err = newTerminal.ReadBinary(reader)
	require.NoError(t, err)

	// Verify the re-read values.
	require.Equal(t, origTerminal.Pots, newTerminal.Pots, "Terminal Pots mismatch")
	require.Equal(t, origTerminal.Players, newTerminal.Players, "Terminal Players mismatch")
}

func TestComplexTreeReadBinary(t *testing.T) {
	req := require.New(t)

	// Setup dummy game parameters.
	prms := table.GameParams{
		NumPlayers:         2,
		SbAmount:           chips.NewFromInt(1),
		InitialStacks:      chips.List{chips.NewFromInt(100), chips.NewFromInt(100)},
		BetSizes:           [][]float32{{1.0}},
		MaxActionsPerRound: 10,
		TerminalStreet:     table.River,
	}

	// Create a Root node.
	root, err := tree.NewRoot(prms)
	req.NoError(err)

	// 	Root
	// 	│
	//  (Reference wrapping Chance)
	// 	│
	//   Chance
	// 	│
	//  (Reference wrapping Player)
	// 	│
	//   Player
	// 	│
	//   Actions [with one child node, Terminal]
	// 	│
	//  (Reference wrapping Terminal)
	// 	│
	//   Terminal

	// Create a Chance node; assign a dummy state.
	chanceNode := &tree.Chance{
		Parent: root,
		State:  new(table.State),
	}

	// Create a Terminal node with dummy values.
	terminalNode := &tree.Terminal{
		// Assuming table.Pots is defined as an integer type.
		Pots: table.Pots{
			{Amount: chips.NewFromInt(100)},
		},
		// Create two players with dummy values.
		Players: table.Players{
			{Paid: chips.NewFromInt(50), Status: table.StatusActive},
			{Paid: chips.NewFromInt(75), Status: table.StatusFolded},
		},
		Parent: nil,
	}
	// Terminal node's parent will be set later.

	// Create a Player node.
	playerNode := &tree.Player{
		Parent:  chanceNode,
		TurnPos: 1,
		State:   new(table.State),
	}

	// Create an Actions node for the Player.
	actions := &tree.PlayerActions{
		Parent:   playerNode,
		Actions:  []table.DiscreteAction{table.DCall},
		Nodes:    []tree.Node{terminalNode},
		Policies: tree.NewPolicies(),
	}
	// Bind Actions to the Player.
	playerNode.Actions = actions
	// Now set the Terminal node's parent.
	terminalNode.Parent = playerNode

	// Link the chain of nodes.
	chanceNode.Next = playerNode // Chance node's child is the Player.
	root.Next = chanceNode       // Root node's child is the Chance node.

	// ---------- Marshal the entire tree ----------
	// Note: Under the hood, nodes (other than Root) get written with length prefixes
	// and wrapped in Reference structs upon reading.
	data, err := root.MarshalBinary()
	req.NoError(err)

	require.Equal(t, len(data), int(root.Size()))

	// ---------- Re-read the tree from binary ----------
	// Create a reader from the binary data.
	reader := bytes.NewReader(data)

	// Create a new, empty Root node.
	newRoot := &tree.Root{}
	err = newRoot.ReadBinary(reader)
	req.NoError(err)

	// ---------- Verify the reconstructed tree ----------

	// Verify game parameters and core fields.
	req.Equal(prms, newRoot.Params, "Game parameters mismatch")
	req.Equal(root.Iteration, newRoot.Iteration, "Iteration mismatch")

	// Verify the child chain from the Root.
	req.NotNil(newRoot.Next, "newRoot.Next should not be nil")
	req.Equal(tree.NodeKindChance, newRoot.Next.Kind(), "Expected child of Root to be a Chance node")

	// The child may be wrapped in a Reference. Force expansion if needed.
	var newChance *tree.Chance
	switch n := newRoot.Next.(type) {
	case *tree.Chance:
		newChance = n
	case *tree.Reference:
		ref, err := n.Expand()
		req.NoError(err)
		var ok bool
		newChance, ok = ref.(*tree.Chance)
		req.True(ok, "Expanded node is not of type *Chance")
	default:
		req.Fail("newRoot.Next must be either *Chance or *Reference wrapping a Chance")
	}

	req.NotNil(newChance.Next, "Chance.Next should not be nil")
	req.Equal(tree.NodeKindPlayer, newChance.Next.Kind(), "Expected child of Chance to be a Player node")

	// Expand Player node if wrapped in a Reference.
	var newPlayer *tree.Player
	switch p := newChance.Next.(type) {
	case *tree.Player:
		newPlayer = p
	case *tree.Reference:
		ref, err := p.Expand()
		req.NoError(err)
		var ok bool
		newPlayer, ok = ref.(*tree.Player)
		req.True(ok, "Expanded node is not of type *Player")
	default:
		req.Fail("newChance.Next must be either *Player or *Reference wrapping a Player")
	}

	req.Equal(uint8(1), newPlayer.TurnPos, "Player TurnPos mismatch")
	req.NotNil(newPlayer.Actions, "Player.Actions should not be nil")
	req.Len(newPlayer.Actions.Nodes, 1, "Actions.Nodes length mismatch")

	// Verify that the sole action child is a Terminal node.
	leaf := newPlayer.Actions.Nodes[0]
	req.Equal(tree.NodeKindTerminal, leaf.Kind(), "Expected action node to be a Terminal node")

	var newTerminalRead *tree.Terminal
	switch tnode := leaf.(type) {
	case *tree.Terminal:
		newTerminalRead = tnode
	case *tree.Reference:
		ref, err := tnode.Expand()
		req.NoError(err)
		var ok bool
		newTerminalRead, ok = ref.(*tree.Terminal)
		req.True(ok, "Expanded node is not of type *Terminal")
	default:
		req.Fail("Actions.Nodes[0] must be either *Terminal or *Reference wrapping a Terminal")
	}

	req.Equal(table.Pots(table.Pots{{Amount: 100}}), newTerminalRead.Pots, "Terminal Pots mismatch")
	req.Len(newTerminalRead.Players, 2, "Terminal Players count mismatch")
}

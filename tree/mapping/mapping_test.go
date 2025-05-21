package mapping

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/pokerdroid/poker/chips"
	"github.com/pokerdroid/poker/table"
	"github.com/pokerdroid/poker/tree"
)

// TestMapGameStateToTree_Basic checks the standard flow where an existing action
// in the real game state should match a corresponding DiscreteAction in the tree.
func TestMapGameStateToTree_Basic(t *testing.T) {
	prms := table.NewGameParams(2, chips.NewFromInt(10))
	prms.BetSizes = [][]float32{{1, 2}}
	prms.Limp = true

	s, err := table.NewGame(prms)
	require.NoError(t, err)

	// For demonstration, let's have player 0 Call
	err = s.Action(table.DCall)
	require.NoError(t, err)

	// Build a root node
	root, err := tree.NewRoot(prms)
	require.NoError(t, err)

	// Expand the tree fully (if your tree package has an ExpandFull or similar)
	// or attach minimal nodes required to reflect the preflop scenario.
	err = tree.ExpandFull(root)
	require.NoError(t, err)

	// Now map the real state to the abstract tree.
	playerNode, err := MapGameStateToTree(prms, s.Latest, root)
	require.NoError(t, err)
	require.NotNil(t, playerNode, "Expected to find a valid *tree.Player node")

	// If everything is correct, we should have the next decision at a Player node.
	t.Logf("Found player node with TurnPos=%d", playerNode.TurnPos)
}

// TestMapGameStateToTree_NoRoot tests when we provide a nil root to the
func TestMapGameStateToTree_NoRoot(t *testing.T) {
	prms := table.NewGameParams(2, chips.NewFromInt(10))
	s := table.NewState(prms)
	_, _ = table.MakeInitialBets(prms, s)

	_, err := MapGameStateToTree(prms, s, nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), "root node is nil")
}

// TestMapGameStateToTree_NoDecisionPoint ensures we error out if the mapping
// no longer finds any further decision node in the tree but the real state
// suggests another action should occur.
func TestMapGameStateToTree_NoDecisionPoint(t *testing.T) {
	prms := table.NewGameParams(2, chips.NewFromInt(10))
	s := table.NewState(prms)
	s, err := table.MakeInitialBets(prms, s)
	require.NoError(t, err)

	// The root is valid but we won't expand the tree or link next nodes
	root, err := tree.NewRoot(prms)
	require.NoError(t, err)
	require.NotNil(t, root)

	// The real state has a history that expects a next decision, but the tree has no next nodes
	// so FindDecisionPoint will return nil
	_, err = MapGameStateToTree(prms, s, root)
	require.Error(t, err)
	require.True(t, errors.Is(err, ErrNoDecisionPointFound))
}

func TestMapping(t *testing.T) {
	prms := table.NewGameParams(2, chips.NewFromInt(100))
	prms.BetSizes = [][]float32{{0.5, 1, 2}}
	prms.SbAmount = chips.NewFromFloat32(1)

	x, err := tree.NewRoot(prms)
	require.NoError(t, err)

	err = tree.Expand(x, x)
	require.NoError(t, err)

	err = tree.Expand(x, x.Next)
	require.NoError(t, err)

	err = tree.Expand(x, x.Next.(*tree.Chance))
	require.NoError(t, err)

	err = tree.Expand(x, x.Next.(*tree.Chance).Next)
	require.NoError(t, err)

	n, ok := x.Next.(*tree.Chance).Next.(*tree.Player).GetAction(table.DiscreteAction(1))
	require.True(t, ok)

	err = tree.Expand(x, n)
	require.NoError(t, err)

	require.Equal(t, n.(*tree.Player).State.Path(1), "r:n:b4.00")

	n, ok = n.(*tree.Player).GetAction(table.DCall)
	require.True(t, ok)

	err = tree.Expand(x, n)
	require.NoError(t, err)

	err = tree.Expand(x, n.(*tree.Chance).Next)
	require.NoError(t, err)

	n, ok = n.(*tree.Chance).Next.(*tree.Player).GetAction(table.DiscreteAction(1))
	require.True(t, ok)

	err = tree.Expand(x, n)
	require.NoError(t, err)

	p := n.(*tree.Player)

	mx, err := MapGameStateToTree(prms, p.State, x)
	require.NoError(t, err)

	require.Equal(t, mx.State.Path(1), "r:n:b4.00:c:n:b8.00")
	require.Equal(t, p.State.Path(1), "r:n:b4.00:c:n:b8.00")
	require.Equal(t, tree.GetPath(p).String(), "r:n:b4.00:c:n:b8.00:p")
}

package studiotree

import (
	"errors"
	"fmt"

	absp "github.com/pokerdroid/poker/abs/pack"
	"github.com/pokerdroid/poker/card"
	"github.com/pokerdroid/poker/chips"
	"github.com/pokerdroid/poker/table"
	"github.com/pokerdroid/poker/tree"
)

type Action struct {
	Kind      tree.NodeKind `json:"kind"`
	ActionIdx *int          `json:"actionIdx"`
	Cards     card.Cards    `json:"cards"`
	Solution  *int          `json:"solution"`
}

type PlayerAction struct {
	Player *tree.Player `json:"player"`
	Action int          `json:"action"`
}

type Result struct {
	Kind        tree.NodeKind          `json:"kind"`
	Street      string                 `json:"street"`
	Actions     []table.DiscreteAction `json:"actions"`
	Matrix      Matrix                 `json:"matrix"`
	State       *table.State           `json:"state"`
	Pot         chips.Chips            `json:"pot,omitempty"`
	TreeHistory string                 `json:"tree_history"`
	End         bool                   `json:"end"`
}

type Solutions struct {
	Roots []*tree.Root
	Abs   *absp.Abs
}

type Inspector struct {
	tree     *Tree
	abs      *absp.Abs
	solution *int
	root     []*tree.Root
}

func NewInspector(abs *absp.Abs, roots []*tree.Root) *Inspector {
	return &Inspector{abs: abs, root: roots}
}

func (i *Inspector) Get(actions []Action) (r *Result, err error) {
	var board card.Cards

	if len(actions) == 0 {
		return nil, errors.New("no actions provided")
	}

	first := actions[0]

	// We first check if the solution is valid.
	if first.Solution == nil || *first.Solution >= len(i.root) {
		return nil, errors.New("solution required")
	}

	// We then check if the solution has changed.
	if i.solution == nil || *i.solution != *first.Solution {
		i.tree, err = NewTree(i.root[*first.Solution])
		if err != nil {
			return nil, err
		}
		i.solution = first.Solution
	}

	i.tree.Reset()

	for _, action := range actions[1:] {
		switch action.Kind {
		case tree.NodeKindChance:
			board = append(board, action.Cards...)
			if err := i.tree.Next(); err != nil {
				return nil, err
			}

		case tree.NodeKindPlayer:
			if action.ActionIdx == nil {
				return nil, errors.New("action index required")
			}
			if err := i.tree.Action(*action.ActionIdx); err != nil {
				return nil, err
			}

		default:
			return nil, errors.New("unknown node kind")
		}
	}

	current, state, err := i.tree.Current()
	if err != nil {
		return nil, err
	}

	if state == nil {
		r = &Result{
			Kind:        current.Kind(),
			Street:      table.NoStreet.String(),
			TreeHistory: tree.GetPath(current).String(),
			End:         true,
		}

		return r, nil
	}

	r = &Result{
		Kind:        current.Kind(),
		Street:      state.Street.String(),
		State:       state,
		TreeHistory: tree.GetPath(current).String(),
		Pot:         state.Players.PaidSum(),
	}

	if player, ok := current.(*tree.Player); ok {
		r.Actions = player.Actions.Actions

		matrix := NewMatrixBuilder(player, board, i.abs, i.tree.GetMoves())

		data, err := matrix.Build()
		if err != nil {
			return nil, err
		}
		r.Matrix = data
	}

	return r, nil
}

type Tree struct {
	root    *tree.Root
	fpn     tree.Node
	current tree.Node
	moves   []PlayerAction
}

func NewTree(root *tree.Root) (*Tree, error) {
	t := &Tree{
		root:    root,
		current: root.Next,
	}

	if err := t.Next(); err != nil {
		return nil, err
	}

	t.fpn = t.current
	return t, nil
}

func (t *Tree) Next() error {
	if err := tree.Expand(t.root, t.current); err != nil {
		return err
	}

	switch z := t.current.(type) {
	case *tree.Root:
		t.current = z.Next
	case *tree.Chance:
		t.current = z.Next
	case *tree.Reference:
		t.current = z.Node
		return t.Next()
	case *tree.Terminal:
		return errors.New("terminal node, cannot advance")
	case *tree.Player:
		return errors.New("player node, cannot advance")
	}

	return nil
}

func (t *Tree) Current() (tree.Node, *table.State, error) {
	if t.current == nil {
		return nil, nil, errors.New("current tree is nil")
	}

	if err := tree.Expand(t.root, t.current); err != nil {
		return nil, nil, err
	}

	var state *table.State

	switch z := t.current.(type) {
	case *tree.Root:
		state = z.State
	case *tree.Chance:
		state = z.State
	case *tree.Player:
		state = z.State
	case *tree.Reference:
		t.current = z.Node
		return t.Current()
	}

	return t.current, state, nil
}

func (t *Tree) Action(idx int) error {
	_ = t.Next()

	if err := tree.Expand(t.root, t.current); err != nil {
		return err
	}

	pl, ok := t.current.(*tree.Player)
	if !ok {
		return fmt.Errorf("not a player node: %T", t.current)
	}

	if idx >= len(pl.Actions.Nodes) {
		return errors.New("invalid action index")
	}

	t.moves = append(t.moves, PlayerAction{
		Player: pl,
		Action: idx,
	})

	t.current = pl.Actions.Nodes[idx]
	return nil
}

func (t *Tree) GetMoves() []PlayerAction {
	return t.moves
}

func (t *Tree) Reset() {
	t.current = t.fpn
	t.moves = []PlayerAction{}
}

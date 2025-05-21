package mapping

import (
	"errors"

	"github.com/pokerdroid/poker/chips"
	"github.com/pokerdroid/poker/table"
	"github.com/pokerdroid/poker/tree"
)

// TODO we should accumulate difference on each streat pot vs abstraction
// in next round we add that to deicde which direction to go.

var (
	ErrNoDecisionPointFound  = errors.New("no decision point found")
	ErrNoMatchingActionFound = errors.New("no matching action found")
)

// MapGameStateToTree attempts to map the real game state (s) to a node in the abstracted tree (root).
func MapGameStateToTree(p table.GameParams, s *table.State, rx *tree.Root) (*tree.Player, error) {
	var current tree.Node

	if rx == nil {
		return nil, errors.New("root node is nil")
	}

	// Calculate scaling ratio between training stakes and actual stakes
	if rx.Params.SbAmount.Equal(chips.Zero) {
		return nil, errors.New("training stakes cannot be zero")
	}

	ratio := p.SbAmount.Div(rx.Params.SbAmount)
	current = rx
	var err error

	for _, pa := range s.History() {
		current, err = tree.FindDecisionPoint(current)
		if err != nil {
			return nil, err
		}

		if current == nil {
			return nil, ErrNoDecisionPointFound
		}

		if _, ok := current.(*tree.Chance); ok {
			continue
		}

		if _, ok := current.(*tree.Terminal); ok {
			return nil, errors.New("terminal node found")
		}

		p, ok := current.(*tree.Player)
		if !ok {
			return nil, errors.New("current node is not a player node")
		}

		if pa.Action.Action.IsBlind() {
			continue
		}

		// Scale only the pot and action amount for matching
		pot := pa.State.Players.PaidSum().Div(ratio)
		act := pa.Action
		act.Amount = act.Amount.Div(ratio)

		idx := MatchAction(act, p.Actions.Actions, pot)
		if idx == -1 {
			return nil, ErrNoMatchingActionFound
		}

		next := p.Actions.Nodes[idx]
		current = next
	}

	current, err = tree.FindDecisionPoint(current)
	if err != nil {
		return nil, err
	}

	if p, ok := current.(*tree.Player); ok {
		return p, nil
	}

	return nil, errors.New("no player node found")
}

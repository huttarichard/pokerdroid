package cfr

import (
	"context"
	"errors"

	"github.com/pokerdroid/poker"
	"github.com/pokerdroid/poker/abs"
	"github.com/pokerdroid/poker/bot"
	"github.com/pokerdroid/poker/frand"
	"github.com/pokerdroid/poker/table"
	"github.com/pokerdroid/poker/tree"
	"github.com/pokerdroid/poker/tree/mapping"
)

var (
	ErrNoState  = errors.New("no state")
	ErrNoPolicy = errors.New("no policy")
)

type Simple struct {
	Abs  abs.Mapper
	Tree *tree.Root
	Rand frand.Rand
}

func (a Simple) Advise(ctx context.Context, loggr poker.Logger, state bot.State) (tb table.DiscreteAction, err error) {
	p, err := mapping.MapGameStateToTree(state.Params, state.State, a.Tree)
	if err != nil {
		return 0, err
	}

	loggr.Printf("Table path:       %s\n", state.State.Path(state.Params.SbAmount))
	loggr.Printf("Tree node mapped: %s\n", tree.GetPath(p))

	if p == nil || p.Actions == nil {
		return 0, ErrNoState
	}

	if p.Actions.Policies == nil {
		return 0, ErrNoState
	}

	cards := append(state.Hole, state.Community...)

	cluster := a.Abs.Map(cards)
	loggr.Printf("Cluster: %d\n", cluster)

	pol, ok := p.Actions.Policies.Get(cluster)
	if !ok {
		return 0, ErrNoPolicy
	}

	loggr.Printf("Policy: %s\n", pol.String())

	action, err := SampleState(&SampleParams{
		Actions: p.Actions.Actions,
		Policy:  pol,
		Rng:     a.Rand,
		State:   state,
	})
	if err != nil {
		return 0, err
	}

	loggr.Printf("Chosen action: %s\n", action.String())
	return action, nil
}

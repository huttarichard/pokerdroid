package mc

import (
	"context"

	"github.com/pokerdroid/poker"
	"github.com/pokerdroid/poker/bot"
	"github.com/pokerdroid/poker/equity/omp"
	"github.com/pokerdroid/poker/frand"
	"github.com/pokerdroid/poker/table"
)

type Advisor struct {
	Passiveness float32
	Defense     float32
	OpenRange   OpenRange
	Rand        frand.Rand
	// Getter      equity.Getter
}

func NewAdvisor() *Advisor {
	return &Advisor{
		Defense:     0.4,
		Passiveness: 0.2,
		OpenRange:   NewOpenRange(),
		Rand:        frand.NewHash(),
	}
}

func (a Advisor) Advise(ctx context.Context, loggr poker.Logger, state bot.State) (table.DiscreteAction, error) {
	equity := omp.Equity(
		state.Hole[0],
		state.Hole[1],
		state.Community,
		len(state.State.Players),
	)

	px := state.Params.Clone()

	actions := table.NewDiscreteLegalActions(px, state.State)
	weak := a.OpenRange.WeakRange(state.Hole)

	_, cInEvs := actions[table.DCall]
	_, kInEvs := actions[table.DCheck]

	// Fold Weak ranges
	if weak && cInEvs {
		return table.DFold, nil
	}

	if weak && kInEvs {
		return table.DCheck, nil
	}

	// Preflop open with 1 POT
	if state.State.Street == table.Preflop && state.State.StreetAction == 2 {
		return table.DiscreteAction(1), nil
	}

	evs := NewEVs(EvsParams{
		Params:      px,
		State:       state.State,
		Equity:      float32(equity[0]),
		Passiveness: a.Passiveness,
		Defense:     a.Defense,
	})

	// return evs.Choice(a.Rand, 0), nil
	return evs.Choice(a.Rand, 0), nil
}

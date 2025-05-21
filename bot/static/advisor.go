package static

import (
	"context"

	"github.com/pokerdroid/poker"
	"github.com/pokerdroid/poker/bot"
	"github.com/pokerdroid/poker/frand"
	"github.com/pokerdroid/poker/table"
)

type Advisor struct {
	Rand   frand.Rand
	Params table.GameParams
}

func NewAdvisor(params table.GameParams) *Advisor {
	rnd := frand.NewHash()
	return &Advisor{Rand: rnd, Params: params}
}

func (a Advisor) Advise(ctx context.Context, logger poker.Logger, state *bot.State) (table.DiscreteAction, error) {
	disa := table.NewDiscreteLegalActions(a.Params, state.State)

	keys := make([]table.DiscreteAction, 0, len(disa))
	for k := range disa {
		keys = append(keys, k)
	}

	action := keys[a.Rand.Intn(len(keys))]

	return action, nil
}

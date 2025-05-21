package cfr

import (
	"errors"

	"github.com/pokerdroid/poker/bot"
	"github.com/pokerdroid/poker/float/f64"
	"github.com/pokerdroid/poker/frand"
	"github.com/pokerdroid/poker/policy"
	"github.com/pokerdroid/poker/table"
)

// TODO we can enhance by becoming exploitable and
// considering getting closer to ev baseline suggests
// We could introduce factor of "aggressiveness" to make it
// more or less exploitable and more or less aggressive.

type SampleParams struct {
	State   bot.State
	Actions []table.DiscreteAction
	Policy  *policy.Policy
	Rng     frand.Rand
}

// Sample againts real policy. This is used to sample actions from
// a policy accounting for cases where policy might be bit off.
func SampleState(params *SampleParams) (table.DiscreteAction, error) {
	s := params.State
	strategy := params.Policy.GetAverageStrategy()
	legal := table.NewLegalActions(s.Params, s.State)

	paid := s.State.Players[s.State.TurnPos].Paid
	stack := s.Params.InitialStacks[s.State.TurnPos].Sub(paid)

	// Create new strategy considering only legal actions
	var sum float64
	strat := make([]float64, 0, len(strategy))
	acts := make([]table.DiscreteAction, 0, len(strategy))

	// Copy probabilities only for legal actions
	for i, act := range params.Actions {
		ax, amount := act.GetAction(s.Params, s.State)
		min, ok := legal[ax]
		if !ok {
			continue
		}
		if amount.LessThan(min) {
			continue
		}
		// Sometimes best like 2x POT might exceed the stack size,
		// so we need to consider this case.
		_, ok = legal[table.AllIn]
		if amount.GreaterThanOrEqual(stack) && ok {
			act = table.DAllIn
		}
		sum += strategy[i]
		strat = append(strat, strategy[i])
		acts = append(acts, act)
	}

	// Normalize if sum is not zero
	if sum == 0 {
		// panic("sum is zero")
		return 0, errors.New("sum is zero")
	}

	f64.ScalUnitary(1/sum, strat)
	indx := frand.SampleIndex(params.Rng, strat, 0.00001)
	return acts[indx], nil
}

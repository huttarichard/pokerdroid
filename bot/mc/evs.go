package mc

import (
	"bytes"
	"fmt"

	"github.com/pokerdroid/poker/chips"
	"github.com/pokerdroid/poker/frand"
	"github.com/pokerdroid/poker/table"
)

type Evs map[table.DiscreteAction]chips.Chips

type EvsParams struct {
	Params      table.GameParams
	State       *table.State
	Equity      float32
	Passiveness float32
	Defense     float32
}

func NewEVs(params EvsParams) Evs {
	state := params.State
	equity := params.Equity

	ef := chips.NewFromFloat32(equity)
	efr := chips.NewFromFloat32(1 - equity)

	lgl := table.NewDiscreteLegalActions(params.Params, params.State)
	lge := make(Evs)
	potSize := params.State.Players.PaidSum()

	paid := params.State.Players[params.State.TurnPos].Paid
	stack := params.Params.InitialStacks[params.State.TurnPos].Sub(paid)

	pass := chips.NewFromFloat32(params.Passiveness)

	aggr := Aggressivity(state, state.TurnPos)

	// Define the old range [a1, a2] and the new range [b1, b2]
	a1, a2 := float32(0), float32(2)
	b1, b2 := float32(1)-params.Defense, float32(1)+params.Defense

	// Use the formula to find the scaled value
	aggr = b1 + (b2-b1)*(aggr-a1)/(a2-a1)

	defa := chips.NewFromFloat32(aggr)
	defb := defa.Sub(chips.NewFromFloat32(params.Defense * 2))

	for action, amount := range lgl {
		var ev chips.Chips
		switch action {
		case table.DFold:
			ev = chips.Zero
		case table.DCheck:
			ev = potSize.Mul(ef)
		default:
			ev = potSize.Add(amount).Mul(ef.Sub(pass)).Sub(amount.Mul(efr))
		}

		if !ev.GreaterThanOrEqual(chips.Zero) {
			continue
		}

		riskFactor := amount.Mul(efr).Div(stack)
		ev = ev.Mul(chips.NewFromFloat(1.0).Sub(riskFactor))

		lge[action] = ev

		if aggr == 1 {
			continue
		}

		// Anything that is a bet
		switch action {
		case table.DCheck, table.DFold, table.DCall:
			continue
		}

		if aggr > 1 {
			// we know that opponent is aggressor
			// become more passive
			lge[action] = lge[action].Mul(defb)
		} else {
			lge[action] = lge[action].Mul(defa)
		}
	}

	return lge
}

func (l Evs) String() string {
	buf := bytes.NewBufferString("EVs:\n")
	for k, v := range l {
		buf.WriteString(fmt.Sprintf("\t%s: %s\n", k, v.StringFixed(8)))
	}
	return buf.String()
}

func (les Evs) Choice(r frand.Rand, magnify int) table.DiscreteAction {
	// Generate a random float between 0 and 1
	rx := r.Float64()

	// Select action based on strategy
	cumulative := 0.0
	lasta := table.DFold

	m := make(Evs)
	magnifier := chips.NewFromInt(int64(magnify))
	for k, v := range les {
		m[k] = v.Pow(magnifier)
	}
	m.Normalize()

	for action, prob := range m {
		lasta = action
		cumulative += prob.Float64()

		if rx < cumulative {
			return action
		}
	}
	return lasta
}

func (les Evs) Max() table.DiscreteAction {
	var action table.DiscreteAction
	var maxEv chips.Chips

	for k, v := range les {
		if v.GreaterThanOrEqual(maxEv) {
			action = k
			maxEv = v
		}
	}
	return action
}

func (les Evs) Normalize() {
	maxEV := chips.Zero
	for _, ev := range les {
		abs := ev.Abs()
		if abs.GreaterThan(maxEV) {
			maxEV = abs
		}
	}

	strategy := make(Evs)
	for action, ev := range les {
		strategy[action] = ev.Add(maxEV)
	}

	sum := chips.NewFromFloat(10e-8)
	for _, ev := range strategy {
		sum = sum.Add(ev)
	}

	// Normalize the EVs so they sum up to 1
	for id, ev := range strategy {
		les[id] = ev.Div(sum)
	}
}

func Aggressivity(s *table.State, pos uint8) float32 {
	actions := 0
	aggressive := 0

	for _, record := range s.History() {
		if record.Pos != pos {
			continue
		}

		switch record.Action.Action {
		case table.SmallBlind, table.BigBlind:
			// ignore these
		case table.Check, table.Call:
			actions++
		case table.Bet, table.Raise, table.AllIn:
			aggressive++
			actions++
		}

	}

	// Avoid division by zero
	if actions == 0 {
		return 1.0 // return neutral if no actions recorded
	}

	// Calculate ratio of aggressive to total actions
	aggressivityRatio := float32(aggressive) / float32(actions)

	// Normalize the ratio so that 0.5 becomes 1, 0 becomes 0 and 1 becomes 2
	aggr := 2*(aggressivityRatio-0.5) + 1

	return aggr
}

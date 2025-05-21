package mapping

import (
	"sort"

	"github.com/pokerdroid/poker/chips"
	"github.com/pokerdroid/poker/table"
)

// MatchAction will match ActionAmount to DiscreteAction given a pot size.
// Purpose of this function is to find index of []DiscreteAction
// given. It is useful if some action is specified in ActionAmount
// for example when player plays but not exactly matches discrete
// action space that is available.
//
// Fold, Call, Check, AllIn are matched exactly.
// While bets are matched to closest pot multiplier.
// If no match is found -1 is returned.
func MatchAction(p table.ActionAmount, aa []table.DiscreteAction, pot chips.Chips) int {
	// Handle special actions that don't depend on pot size
	var sameacts = map[table.ActionKind]table.DiscreteAction{
		table.Fold:  table.DFold,
		table.Call:  table.DCall,
		table.Check: table.DCheck,
		table.AllIn: table.DAllIn,
	}

	// First try to match exact actions (fold, call, check, allin)
	for i, a := range aa {
		x, ok := sameacts[p.Action]
		if ok && a == x {
			return i
		}
	}

	// If action is not a bet/raise, and we didn't find a match above, return -1
	if p.Action != table.Bet && p.Action != table.Raise {
		return -1
	}

	// Handle bet/raise actions by comparing to pot multipliers
	var pm []table.DiscreteAction
	for _, a := range aa {
		if a <= 0 { // Positive values are pot multipliers
			continue
		}
		pm = append(pm, a)
	}

	if len(pm) == 0 {
		return -1
	}

	sort.Slice(pm, func(i, j int) bool {
		return pm[i] < pm[j]
	})

	// Convert actual bet to pot multiplier
	am := table.DiscreteAction(p.Amount.Div(pot).Float32())

	// Determine which multiplier to map to using the deterministic
	// Pseudo-Harmonic Mapping approach from Pseudo Harmonic Mapping.md.
	var best table.DiscreteAction

	switch {
	// If only one possible multiplier, must use it
	case len(pm) == 1:
		best = pm[0]

	// If below smallest multiplier, map to smallest
	case am <= pm[0]:
		best = pm[0]

	// If above largest multiplier, map to largest
	case am >= pm[len(pm)-1]:
		best = pm[len(pm)-1]

	default:
		best = PseudoHarmonicMapping(am, pm)
	}

	// Find index of best matching multiplier in original array
	for i, a := range aa {
		if a == best {
			return i
		}
	}

	return -1
}

func PseudoHarmonicMapping(am table.DiscreteAction, pm []table.DiscreteAction) table.DiscreteAction {
	// Find the interval: sorted[i] <= am < sorted[i+1]
	var i int
	for i = 0; i < len(pm)-1; i++ {
		if am < pm[i] || am >= pm[i+1] {
			continue
		}
		// Compute the threshold x* = (A + B + 2AB) / (A + B + 2)
		x := (pm[i] + pm[i+1] + 2.0*pm[i]*pm[i+1]) / (pm[i] + pm[i+1] + 2.0)

		// If am > x*, map to B; otherwise map to A
		if am > x {
			return pm[i+1]
		}
		return pm[i]
	}
	// If we never found it in the loop (edge case), default to last
	return pm[len(pm)-1]
}

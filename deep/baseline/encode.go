package baselinenn

import (
	"github.com/nlpodyssey/spago/mat/float"
	"github.com/pokerdroid/poker/chips"
	"github.com/pokerdroid/poker/equity"
	"github.com/pokerdroid/poker/table"
	"github.com/pokerdroid/poker/tree"
)

type EncodeParams struct {
	Equity equity.Equity
	Params table.GameParams
	Player *tree.Player
}

func Encode[T float.DType](cfg EncodeParams) ([]T, chips.Chips) {
	state := cfg.Player.State

	// Calculate initial stack sum for normalization
	ss := cfg.Params.InitialStacks.Sum()

	// Calculate total features needed
	maxacts := cfg.Params.MaxActionsPerRound
	total, actionh := EncodeFeatures[T](cfg.Params.NumPlayers, maxacts)

	features := make([]T, 0, total)

	// 1. Basic game state features
	// - Small blind amount normalized (1 float)
	features = append(features, T(cfg.Params.SbAmount.Div(ss)))

	// - Current pot size normalized (1 float)
	features = append(features, T(state.Players.PaidSum().Div(ss)))

	// 2. Player position and stack information
	// For each player: [isSB, initialStack/sum, currentStack/sum]
	for i, player := range state.Players {
		stack := T(cfg.Params.InitialStacks[i].Div(ss).Float64())
		paid := T(player.Paid.Div(ss).Float64())
		isSB := T(0)
		if i == int(cfg.Params.BtnPos) {
			isSB = T(1)
		}
		features = append(features, isSB, stack, paid)
	}

	// 3. Action history encoding (4 features per action)
	// - [fold, check, call, raise] for each action
	// Encode for all 4 streets (preflop, flop, turn, river)
	ff := make([]T, actionh)
	actions := tree.ExtractActions(cfg.Player)

	for street := 0; street < 4; street++ {
		offset := street * (int(maxacts) * 4)

		for _, action := range actions {
			if action.State.Street != table.Street(street+1) {
				continue
			}
			act, amnt := action.Action.GetAction(cfg.Params, state)

			switch act {
			case table.Fold:
				ff[offset] = 1
			case table.Check:
				ff[offset+1] = 1
			case table.Call:
				// Normalize call amount by initial stack sum
				ff[offset+2] = T(amnt.Div(ss))
			default: // Raise/Bet/AllIn
				ff[offset+3] = T(amnt.Div(ss))
			}
			offset += 4
		}
	}

	features = append(features, ff...)
	features = append(features, T(cfg.Equity.WinDraw()))

	return features, ss
}

// EncodeFeatures returns the total number of features
func EncodeFeatures[T float.DType](nump uint8, maxacts uint8) (int, int) {
	total := 0
	total += 1             // sb amount
	total += 1             // pot size
	total += int(nump) * 3 // player features (isSB, initStack, currStack)

	// history:  4 streets * max actions * 4 features
	actionh := int(4 * maxacts * 4)
	total += actionh

	total += 2 // buckets

	return total, actionh
}

func EncodeBasline[T float.DType](
	predicted []table.DiscreteAction,
	actual []table.DiscreteAction,
	baselines []float64,
	stuckstum chips.Chips,
) []T {
	if len(actual) != len(baselines) {
		panic("blueprintActs and baselines must be the same length")
	}

	encoded := make([]T, len(predicted))

	for idy, p := range predicted {
		for idx, a := range actual {
			if p == a {
				encoded[idy] = T(baselines[idx] / stuckstum.Float64())
			}
		}
	}

	return encoded
}

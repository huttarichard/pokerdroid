package cfr

import (
	"github.com/pokerdroid/poker/abs"
	"github.com/pokerdroid/poker/card"
	"github.com/pokerdroid/poker/table"
	"github.com/pokerdroid/poker/tree"
)

type ComputeRangeParams struct {
	Actions []tree.Action
	Players uint8
	Board   card.Cards
	Abs     abs.Mapper
}

type ActionRange struct {
	Action tree.Action
	Board  card.Cards
	Ranges []card.RangeDist
}

func ComputeRange(p ComputeRangeParams) ([]card.RangeDist, error) {
	mp := make([]card.RangeDist, p.Players)
	for j := range mp {
		mp[j] = card.NewUniformRangeDist()
	}

	for _, act := range p.Actions {
		var board card.Cards

		switch act.State.Street {
		case table.Preflop:
			board = card.Cards{}
		case table.Flop:
			board = p.Board[:3]
		case table.Turn:
			board = p.Board[:4]
		case table.River:
			board = p.Board[:5]
		}

		var notfound int
		var total int

		for _, cc := range card.Combinations(2) {
			if card.IsAnyMatch(cc, board) {
				continue
			}

			cluster := p.Abs.Map(append(cc, board...))

			total++
			policy, ok := act.Parent.Actions.Policies.Get(cluster)
			if !ok {
				notfound++
				continue
			}

			pol := policy.GetAverageStrategy()
			i := card.RangeIndex(cc)

			mp[act.State.TurnPos][i] *= float64(pol[act.NodeIdx])
		}

		// if float32(notfound)/float32(total) > 0.5 {
		// 	panic("too many missing policies")
		// }
	}

	for i, r := range mp {
		mp[i] = r.Normalize()
	}

	return mp, nil
}

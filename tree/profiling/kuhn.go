package profiling

import (
	"github.com/pokerdroid/poker/card"
	kuhndealer "github.com/pokerdroid/poker/dealer/kuhn"
	"github.com/pokerdroid/poker/tree"
)

func GetKuhnProfile(root *tree.Root) *Profile {
	profile := &Profile{
		BetSizes: root.Params.BetSizes,
		Nodes:    make([]Node, 0),
	}

	tree.MustVisit(root, -1, func(n tree.Node, _ []tree.Node, _ int) bool {
		player, ok := n.(*tree.Player)
		if !ok {
			return true
		}

		for c, cc := range kuhndealer.Clusters {
			policy, ok := player.Actions.Policies.Get(cc)
			if !ok {
				continue
			}

			profile.Nodes = append(profile.Nodes, Node{
				Runes:    tree.GetPath(n),
				Cards:    card.Cards{c},
				Actions:  player.Actions.Actions,
				Strategy: policy.GetAverageStrategy(),
				Policy:   policy,
			})
		}

		return true
	})

	return profile
}

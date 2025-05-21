package profiling

import (
	"fmt"

	"github.com/pokerdroid/poker/abs"
	"github.com/pokerdroid/poker/card"
	"github.com/pokerdroid/poker/tree"
)

type Profile struct {
	BetSizes [][]float32
	Nodes    Nodes
	Board    card.Cards
}

type Params struct {
	PlayerID uint8
	Depth    int
	Board    card.Cards
	Abs      abs.Mapper
	Tree     tree.Node
	BetSizes [][]float32
}

func New(p Params) (ss *Profile, err error) {
	ss = &Profile{
		BetSizes: p.BetSizes,
		Nodes:    make([]Node, 0),
		Board:    p.Board,
	}

	iterator := func(node tree.Node, ccx card.Cards) error {
		var clus abs.Cluster

		allc := append(ccx, p.Board...)

		clus = p.Abs.Map(allc)

		player, ok := node.(*tree.Player)
		if !ok {
			return nil
		}

		policy, ok := player.Actions.Policies.Get(clus)
		if !ok {
			return nil
		}

		if player.TurnPos != p.PlayerID {
			return nil
		}

		strategy := policy.GetAverageStrategy()

		acts := player.Actions.Actions

		if len(acts) != len(strategy) {
			panic(fmt.Sprintf("strategy length mismatch: %d != %d", len(acts), len(strategy)))
		}

		ss.Nodes = append(ss.Nodes, Node{
			Runes:    tree.GetPath(node),
			Cards:    allc,
			Actions:  acts,
			Policy:   policy,
			Strategy: strategy,
			Node:     node,
		})

		return nil
	}

	tree.MustVisit(p.Tree, p.Depth, func(n tree.Node, _ []tree.Node, _ int) bool {
		if err != nil {
			return false
		}

		cc := card.All(p.Board...)
		combos := card.CombinationsFrom(cc, 2)

		for _, ccx := range combos {
			err = iterator(n, ccx)
			if err != nil {
				return false
			}
		}

		return true
	})

	return
}

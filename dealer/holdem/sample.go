package holdemdealer

import (
	"github.com/pokerdroid/poker/abs"
	"github.com/pokerdroid/poker/card"
	"github.com/pokerdroid/poker/dealer"
	"github.com/pokerdroid/poker/eval"
	"github.com/pokerdroid/poker/frand"
	"github.com/pokerdroid/poker/table"
	"github.com/pokerdroid/poker/tree"
)

type Sample struct {
	deck   *deck
	hands  []card.Cards
	getter abs.Mapper
	rng    frand.Rand
	cur    table.Street
	term   table.Street
}

var _ dealer.Sample = &Sample{}

var offsets = []int{0, 2, 5, 6, 7}
var additions = []int{0, 3, 1, 1}

func (c *Sample) Cards(pID uint8) card.Cards {
	return c.hands[pID][:offsets[c.cur]]
}

func (c *Sample) Board() card.Cards {
	return c.hands[0][2:]
}

func (c *Sample) Sample(s table.Street) {
	var ok bool

	for c.cur < s {
		// Get the offset
		of := offsets[c.cur]
		// Get how many elements we need to add from offset
		ad := additions[c.cur]
		// Expand the hand to the new street
		c.hands[0] = c.hands[0][:of+ad]
		// Fill first hand with the board
		ok = fillrnd(c.rng, c.deck, c.hands[0][of:of+ad])
		if !ok {
			panic("fill random failed")
		}
		// Replicate for rest
		for i := 1; i < len(c.hands); i++ {
			c.hands[i] = c.hands[i][:of+ad]
			copy(c.hands[i][of:of+ad], c.hands[0][of:of+ad])
		}
		c.cur++
	}
}

func (c *Sample) Cluster(n dealer.Turner, m abs.Mapper) abs.Cluster {
	return m.Map(c.hands[n.GetTurnPos()])
}

func (c *Sample) Utility(n *tree.Terminal, pID uint8) float64 {
	paid := float64(n.Players[pID].Paid)
	pot := float64(n.Pots.Sum())

	if n.Players[pID].Status == table.StatusFolded {
		return -paid
	}

	if n.Players.LastAlive(pID) {
		return pot - paid
	}

	c.Sample(c.term)

	winners, err := eval.Judge(c.hands)
	if err != nil {
		panic(err)
	}

	switch len(winners) {
	// we have a draw
	case 0:
		panic("no winners")

	case 1:
		if winners[0] == pID {
			return pot - paid
		}
		return -paid
	default:
		// Tie among multiple winners
		// this however needs to involve a multiple pots.
		//
		// for _, wID := range winners {
		// 	if wID == pID {
		// 		builder.WriteString(fmt.Sprintf("paid: %f\n", (pot/float32(len(winners)))-paid))
		// 		// If pID is among the winners, give them an equal fraction of the pot
		// 		return (pot / float32(len(winners))) - paid
		// 	}
		// }
		// If pID is not among 'winners', they lose their contribution
		return 0
	}
}

func cloneSample(dst *Sample, src *Sample) *Sample {
	cloneDeck(dst.deck, src.deck)
	dst.cur = src.cur
	dst.term = src.term
	dst.getter = src.getter
	dst.rng = src.rng

	for i := range dst.hands {
		dst.hands[i] = dst.hands[i][:len(src.hands[i])]
		copy(dst.hands[i], src.hands[i])
	}

	return dst
}

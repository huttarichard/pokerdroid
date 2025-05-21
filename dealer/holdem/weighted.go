package holdemdealer

import (
	"sync"

	"github.com/pokerdroid/poker/abs"
	"github.com/pokerdroid/poker/card"
	"github.com/pokerdroid/poker/dealer"
	"github.com/pokerdroid/poker/frand"
	"github.com/pokerdroid/poker/table"
)

type RangeParams struct {
	NumPlayers uint8
	Clusters   abs.Mapper
	Ranges     []card.RangeDist
	Board      card.Cards
	deck       *deck
}

type RangeSampler struct {
	RangeParams

	pool sync.Pool
}

func NewWeighted(p RangeParams) *RangeSampler {
	s := &RangeSampler{RangeParams: p}
	s.deck = newDeck()
	popcards(s.deck, p.Board)

	s.pool.New = func() interface{} {
		g := &Sample{
			hands:  make([]card.Cards, p.NumPlayers),
			getter: p.Clusters,
			deck:   newDeck(),
			cur:    table.Preflop,
		}
		for i := uint8(0); i < p.NumPlayers; i++ {
			g.hands[i] = make(card.Cards, 7)
		}
		return g
	}

	return s
}

func (c *RangeSampler) Clone() dealer.Dealer {
	return NewWeighted(c.RangeParams)
}

func (c *RangeSampler) Sample(rng frand.Rand) (dealer.Sample, error) {
	g := c.pool.Get().(*Sample)
	g.rng = rng
	cloneDeck(g.deck, c.deck)

	switch len(c.Board) {
	case 0:
		g.cur = table.Preflop
	case 3:
		g.cur = table.Flop
	case 4:
		g.cur = table.Turn
	case 5:
		g.cur = table.River
	default:
		panic("invalid board")
	}

	for i := uint8(0); i < c.NumPlayers; i++ {
		reng := c.Ranges[i]

		for i := range reng {
			if card.IsAnyMatch(card.RangeCards(i), c.Board) {
				reng[i] = 0
			}
		}

		// if reng.Sum() < 1e-4 {
		// 	pretty.Println(reng.Matrix().String())
		// 	panic("failed to sample from range")
		// }

		zeng := reng.Normalize()
		s := zeng.Sample(rng)

		// if card.IsAnyMatch(s, c.Board) {
		// 	pretty.Println(reng.Matrix().String(), s.String(), c.Board.String())
		// 	pretty.Println(card.Coordinates(s))
		// 	pretty.Println(reng)
		// 	panic("failed to sample from range")
		// }

		popcards(g.deck, s)

		g.hands[i][0] = s[0]
		g.hands[i][1] = s[1]
		g.hands[i] = append(g.hands[i][:2], c.Board...)
	}

	return g, nil
}

func (c *RangeSampler) Copy(rng frand.Rand, s dealer.Sample) (dealer.Sample, error) {
	g := c.pool.Get().(*Sample)
	cloneSample(g, s.(*Sample))
	g.rng = rng
	return g, nil
}

func (c *RangeSampler) Put(s dealer.Sample) {
	if s == nil {
		return
	}
	if gs, ok := s.(*Sample); ok {
		c.pool.Put(gs)
	}
}

package holdemdealer

import (
	"errors"
	"sync"

	"github.com/pokerdroid/poker/card"
	"github.com/pokerdroid/poker/dealer"
	"github.com/pokerdroid/poker/frand"
	"github.com/pokerdroid/poker/table"
)

type SamplerParams struct {
	NumPlayers uint8
	Terminal   table.Street
}

type Sampler struct {
	SamplerParams

	pool sync.Pool
}

func New(p SamplerParams) *Sampler {
	if p.Terminal == table.NoStreet {
		p.Terminal = table.River
	}

	s := &Sampler{SamplerParams: p}
	s.pool.New = func() interface{} {
		g := &Sample{
			hands: make([]card.Cards, p.NumPlayers),
			deck:  newDeck(),
			cur:   table.Preflop,
			term:  p.Terminal,
		}
		for i := uint8(0); i < p.NumPlayers; i++ {
			g.hands[i] = make(card.Cards, 0, 7)
		}
		return g
	}

	return s
}

func (c *Sampler) Clone() dealer.Dealer {
	return New(c.SamplerParams)
}

func (c *Sampler) Sample(rng frand.Rand) (dealer.Sample, error) {
	g := c.pool.Get().(*Sample)
	g.rng = rng
	g.deck.reset()
	g.cur = table.Preflop

	var ok bool

	for i := uint8(0); i < c.NumPlayers; i++ {
		g.hands[i] = g.hands[i][:2]

		ok = fillrnd(rng, g.deck, g.hands[i])
		if !ok {
			return nil, errors.New("not enough cards")
		}
	}

	return g, nil
}

func (c *Sampler) Copy(rng frand.Rand, s dealer.Sample) (dealer.Sample, error) {
	g := c.pool.Get().(*Sample)
	cloneSample(g, s.(*Sample))
	g.rng = rng
	return g, nil
}

func (c *Sampler) Put(s dealer.Sample) {
	if s == nil {
		return
	}
	if gs, ok := s.(*Sample); ok {
		c.pool.Put(gs)
	}
}

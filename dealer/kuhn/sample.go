package kuhndealer

import (
	"github.com/pokerdroid/poker/abs"
	"github.com/pokerdroid/poker/card"
	"github.com/pokerdroid/poker/dealer"
	"github.com/pokerdroid/poker/frand"
	"github.com/pokerdroid/poker/table"
	"github.com/pokerdroid/poker/tree"
)

var DealtCards = []int{2, 3, 1, 1}

type Mapping map[card.Card]abs.Cluster

func (m Mapping) Map(c card.Cards) abs.Cluster {
	return m[c[0]]
}

var Clusters = Mapping{
	card.CardJD: 0,
	card.CardQD: 1,
	card.CardKD: 2,
}

var Reverse = map[abs.Cluster]card.Card{
	0: card.CardJD,
	1: card.CardQD,
	2: card.CardKD,
}

var Cards = card.Cards{
	card.CardJD,
	card.CardQD,
	card.CardKD,
}

type Sample struct {
	Cards card.Cards
}

var _ dealer.Sample = &Sample{}

func (c *Sample) Clone() dealer.Sample {
	return &Sample{Cards: append(card.Cards{}, c.Cards...)}
}

func (c *Sample) Sample(s table.Street) {
	// no-op - no chance nodes in kuhn
}

func (c *Sample) Cluster(n dealer.Turner, m abs.Mapper) abs.Cluster {
	return m.Map(card.Cards{c.Cards[n.GetTurnPos()]})
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

	p1 := Clusters[c.Cards[pID]]

	oID := 0
	if pID == 0 {
		oID = 1
	}

	p2 := Clusters[c.Cards[oID]]

	if p1 > p2 {
		return pot - paid
	}
	return -paid
}

type GameHandSample struct {
	Rand frand.Rand
}

func NewGameSampler(r frand.Rand) *GameHandSample {
	return &GameHandSample{Rand: r}
}

func (c *GameHandSample) Clone() dealer.Dealer {
	return &GameHandSample{Rand: c.Rand}
}

func (c *GameHandSample) Copy(rng frand.Rand, s dealer.Sample) (dealer.Sample, error) {
	x := s.(*Sample)
	return &Sample{Cards: x.Cards}, nil
}

func (c *GameHandSample) Sample(r frand.Rand) (dealer.Sample, error) {
	cds := []card.Cards{
		{card.CardJD, card.CardQD},
		{card.CardQD, card.CardJD},

		{card.CardJD, card.CardKD},
		{card.CardKD, card.CardJD},

		{card.CardQD, card.CardKD},
		{card.CardKD, card.CardQD},
	}

	idx := r.Intn(len(cds))
	sample := &Sample{Cards: cds[idx]}
	return sample, nil
}

func (c *GameHandSample) Put(dealer.Sample) {
	// no-op this is from performance reasons
}

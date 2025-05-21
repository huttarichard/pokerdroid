package iso

import (
	"github.com/pokerdroid/poker/card"
)

type Street struct {
	indx  *Indexer
	rx    int
	cards int
}

var (
	River   = &Street{indx: MustNew(2, []uint8{2, 5}), rx: 1, cards: 7}
	Turn    = &Street{indx: MustNew(2, []uint8{2, 4}), rx: 1, cards: 6}
	Flop    = &Street{indx: MustNew(2, []uint8{2, 3}), rx: 1, cards: 5}
	Preflop = &Street{indx: MustNew(1, []uint8{2}), rx: 0, cards: 2}
)

func (hi *Street) Size() int {
	return hi.indx.Size(hi.rx)
}

func (hi *Street) Index(cards card.Cards) uint32 {
	if hi.cards == 2 {
		return preflopLookup[[2]card.Card{cards[0], cards[1]}]
	}
	return uint32(hi.indx.Index(cards))
}

func (hi *Street) Unindex(index uint64) card.Cards {
	return hi.indx.Unindex(hi.rx, index, hi.cards)
}

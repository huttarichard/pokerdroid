package iso

/*
#cgo CFLAGS: -Wno-pointer-bool-conversion
#include "hand_index.h"
*/
import "C"
import (
	"fmt"
	"runtime"

	"github.com/pokerdroid/poker/card"
)

func cardToUint(c card.Card) uint8 {
	rank := uint8(c-1) / 4
	suit := 3 - (uint8(c-1) % 4) // To account for the reversed order of suits in SUIT_TO_CHAR

	return uint8(rank<<2 | suit)
}

func uintToCard(val uint8) card.Card {
	rank := val >> 2
	suit := 3 - (val & 3) // To account for the reversed order of suits in SUIT_TO_CHAR

	return card.Card(1 + rank*4 + suit)
}

type Indexer struct {
	ptr *C.hand_indexer_t
}

func New(rounds int, cardsPerRound []uint8) (*Indexer, error) {
	indexer := &Indexer{ptr: &C.hand_indexer_t{}}

	success := C.hand_indexer_init(
		C.uint_fast32_t(rounds),
		(*C.uint8_t)(&cardsPerRound[0]),
		indexer.ptr,
	)

	if !success {
		return nil, fmt.Errorf("failed to initialize hand indexer")
	}

	runtime.SetFinalizer(indexer, (*Indexer).free)
	return indexer, nil
}

func MustNew(rounds int, cardsPerRound []uint8) *Indexer {
	ind, err := New(rounds, cardsPerRound)
	if err != nil {
		panic(err)
	}
	return ind
}

func (hi *Indexer) Index(cards card.Cards) uint64 {
	cc := make([]uint8, len(cards))
	for i, r := range cards {
		cc[i] = cardToUint(r)
	}
	return uint64(C.hand_index_last(hi.ptr, (*C.uint8_t)(&cc[0])))
}

func (hi *Indexer) Size(round int) int {
	return int(C.hand_indexer_size(hi.ptr, C.uint_fast32_t(round)))
}

func (hi *Indexer) Unindex(round int, index uint64, size int) card.Cards {
	cards := make([]uint8, size)
	C.hand_unindex(hi.ptr, C.uint_fast32_t(round), C.hand_index_t(index), (*C.uint8_t)(&cards[0]))
	cx := cards[:]
	xx := make(card.Cards, len(cx))
	for i, a := range cx {
		xx[i] = uintToCard(a)
	}
	return xx
}

func (hi *Indexer) free() {
	C.hand_indexer_free(hi.ptr)
}

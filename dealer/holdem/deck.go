package holdemdealer

import (
	"slices"

	"github.com/pokerdroid/poker/card"
	"github.com/pokerdroid/poker/frand"
)

var allcards = card.AllCopy()

type deck struct {
	cards card.Cards
}

func newDeck() *deck {
	d := &deck{cards: make(card.Cards, len(allcards))}
	copy(d.cards, allcards)
	return d
}

func cloneDeck(dest, src *deck) *deck {
	dest.cards = dest.cards[:len(src.cards)]
	copy(dest.cards, src.cards)
	return dest
}

func (d *deck) reset() {
	d.cards = d.cards[:52]
	copy(d.cards, allcards)
}

func fillrnd(rng frand.Rand, deck *deck, fill card.Cards) bool {
	var ok bool
	var element card.Card
	for i := 0; i < len(fill); i++ {
		element, deck.cards, ok = poprnd(rng, deck.cards)
		if !ok {
			return false
		}
		fill[i] = element
	}
	return true
}

// poprnd removes and returns a random element from the slice.
// It modifies the original slice and returns the popped element.
// If the slice is empty, it returns the zero value of T and false.
func poprnd[T any](rng frand.Rand, slice []T) (T, []T, bool) {
	if len(slice) == 0 {
		var zero T
		return zero, slice, false
	}

	// Generate random index
	idx := rng.Intn(len(slice))

	// Get the element at the random index
	element := slice[idx]

	// Remove the element by replacing it with the last element and truncating
	lastIdx := len(slice) - 1
	slice[idx] = slice[lastIdx]
	slice = slice[:lastIdx]

	return element, slice, true
}

func popcards(deck *deck, cds card.Cards) bool {
	for _, cd := range cds {
		idx := slices.Index(deck.cards, cd)
		if idx == -1 {
			continue
		}
		lastIdx := len(deck.cards) - 1
		deck.cards[idx] = deck.cards[lastIdx]
		deck.cards = deck.cards[:lastIdx]
	}
	return true
}

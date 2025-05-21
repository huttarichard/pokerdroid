package card

import (
	"encoding/json"
	"strings"

	"github.com/pokerdroid/poker/frand"
)

// Poper is unifiying interface for popping cards from the deck.
type Poper interface {
	Pop() Card
	PopMulti(n int) Cards
}

// Deck is a slice of cards used for dealing cards.
// Not safe for concurrent use.
type Deck struct {
	Cards Cards `json:"cards"`
}

// NewDeck creates new deck from cards.
func NewDeck(cards Cards) *Deck {
	return &Deck{Cards: cards}
}

// MarshalJSON implements the json.Marshaler interface
func (d Deck) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.Cards)
}

// Shuffle shuffles the deck using the provided
// random number generator.
func (d *Deck) Shuffle(r frand.Rand) {
	dest := Cards{}
	perm := r.Perm(len(d.Cards))
	for _, v := range perm {
		dest = append(dest, d.Cards[v])
	}
	d.Cards = dest
}

// Pop removes a card from the deck and returns it.  Pop
// panics if no cards are available.
func (d *Deck) Pop() Card {
	last := len(d.Cards) - 1
	card := d.Cards[last]
	d.Cards = d.Cards[:last]
	return card
}

// PopMulti calls the Pop function on n number of cards.  PopMulti
// panics if n is larger than the number of cards in the deck.
func (d *Deck) PopMulti(n int) Cards {
	if n > len(d.Cards) {
		panic("deck doesn't have enough cards")
	}
	cards := make([]Card, n)
	for i := 0; i < n; i++ {
		cards[i] = d.Pop()
	}
	return cards
}

// Remove will remove cards from deck.
func (d *Deck) Remove(cc ...Card) {
	for _, c := range cc {
		el := -1
		for i, cx := range d.Cards {
			if c != cx {
				continue
			}
			el = i
			break
		}
		if el == -1 {
			continue
		}
		d.Cards = append(d.Cards[:el], d.Cards[el+1:]...)
	}
}

// String implements the fmt.Stringer interface
func (d *Deck) String() string {
	s := []string{}
	for _, c := range d.Cards {
		s = append(s, c.String())
	}
	return strings.Join(s, ",")
}

package card

import (
	"strings"

	"github.com/pokerdroid/poker/frand"
)

// Card single card represented as byte.
type Card uint8

const (
	Card00 Card = iota
	Card2C
	Card2D
	Card2H
	Card2S
	Card3C
	Card3D
	Card3H
	Card3S
	Card4C
	Card4D
	Card4H
	Card4S
	Card5C
	Card5D
	Card5H
	Card5S
	Card6C
	Card6D
	Card6H
	Card6S
	Card7C
	Card7D
	Card7H
	Card7S
	Card8C
	Card8D
	Card8H
	Card8S
	Card9C
	Card9D
	Card9H
	Card9S
	CardTC
	CardTD
	CardTH
	CardTS
	CardJC
	CardJD
	CardJH
	CardJS
	CardQC
	CardQD
	CardQH
	CardQS
	CardKC
	CardKD
	CardKH
	CardKS
	CardAC
	CardAD
	CardAH
	CardAS
)

// New crates new card from rank and suit
func New(r Rank, s Suit) Card {
	return Card(int(r)*4 + int(s))
}

// NewRandom crates new random card.
// Specify omit as to what cards should be left out.
func NewRandom(r frand.Rand, omit ...Card) Card {
	rx := Card(r.Int63n(52) + 1)
	for _, o := range omit {
		if o == rx {
			return NewRandom(r, omit...)
		}
	}
	return rx
}

// Parse return card from it string representation
// Suit can be: s,h,d,c
func Parse(str string) Card {
	return cardsBack[strings.ToLower(str)]
}

// MarshalText implements the encoding.TextMarshaler interface.
// The text format is "4♠".
func (c Card) MarshalText() ([]byte, error) {
	return []byte(c.String()), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
// The card is expected to be in the format "4♠".
func (c *Card) UnmarshalText(text []byte) error {
	*c = Parse(string(text))
	return nil
}

// String convert card to its string representation
func (c Card) String() string {
	return cardsFront[c]
}

// Symbol convert card to its string representation
// Suits will be: ♠,♥,♦,♣
func (c Card) Symbol() string {
	return cardsSymbols[c]
}

// Returns returns rank of the card
func (c Card) Rank() Rank {
	return cards[c].Rank
}

// Suite returns suite of the card
func (c Card) Suite() Suit {
	return cards[c].Suit
}

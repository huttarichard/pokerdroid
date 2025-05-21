package card

// A Suit represents the suit of a card.
type Suit int

const (
	NoSuit Suit = iota
	Clubs
	Diamonds
	Hearts
	Spades
)

var (
	suitsStr = []string{"♠", "♦", "♥", "♣"}
)

// String returns a string in the format "♠"
func (s Suit) String() string {
	return suitsStr[s]
}

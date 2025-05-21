package card

// A Rank represents the rank of a card.
type Rank int

const (
	NoRank Rank = iota
	Two
	Three
	Four
	Five
	Six
	Seven
	Eight
	Nine
	Ten
	Jack
	Queen
	King
	Ace
)

const (
	ranksStr = "23456789TJQKA"
)

var (
	singularNames = []string{"two", "three", "four", "five", "six", "seven", "eight", "nine", "ten", "jack", "queen", "king", "ace"}
	pluralNames   = []string{"twos", "threes", "fours", "fives", "sixes", "sevens", "eights", "nines", "tens", "jacks", "queens", "kings", "aces"}
)

// String returns a string in the format "2"
func (r Rank) String() string {
	return ranksStr[r-1 : r]
}

// SingularName returns the name of the rank in singular form such as "two" for Two.
func (r Rank) SingularName() string {
	return singularNames[r]
}

// PluralName returns the name of the rank in plural form such as "twos" for Two.
func (r Rank) PluralName() string {
	return pluralNames[r]
}

func (r Rank) OneOf(ranks ...Rank) bool {
	for _, rx := range ranks {
		if r == rx {
			return true
		}
	}
	return false
}

// Ranks represents a list of ranks.
type Ranks []Rank

func (r Ranks) Reverse() (x Ranks) {
	for i := len(r) - 1; i != -1; i-- {
		x = append(x, r[i])
	}
	return x
}

// AllRanks returns all ranks.
func AllRanks() Ranks {
	return Ranks{
		Two,
		Three,
		Four,
		Five,
		Six,
		Seven,
		Eight,
		Nine,
		Ten,
		Jack,
		Queen,
		King,
		Ace,
	}
}

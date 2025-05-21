package card

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/pokerdroid/poker/frand"
)

// Cards represent list of cards.
type Cards []Card

// NewCardsFromBytes creates new cards from bytes.
// Note this is not using decompression.
func NewCardsFromBytes(bb []byte) (xx Cards) {
	for _, x := range bb {
		xx = append(xx, Card(x))
	}
	return
}

// NewCardsFromString creates new cards from string.
// Use empty space to split the cards.
// If card is not recognized it will be skipped.
func NewCardsFromString(s string) Cards {
	var bb Cards
	cards := strings.Split(strings.Trim(strings.ToLower(s), " "), " ")
	for _, c := range cards {
		cx := Parse(c)
		if cx == Card00 {
			continue
		}
		bb = append(bb, Parse(c))
	}
	return bb
}

// Bytes returns cards as bytes.
func (cc Cards) Bytes() (xc []byte) {
	for _, x := range cc {
		xc = append(xc, byte(x))
	}
	return
}

// Compress will compress cards into bytes.
// It uses some additional bits that are available
// as cards to not reach limit of single byte.
// This compression can be more significant with
// more cards.
func Compress(cards Cards) []byte {
	var bitSequence uint64

	for _, card := range cards {
		bitSequence = (bitSequence << 6) | uint64(card)
	}

	// Determine the number of bytes based on the number of cards
	var numBytes int
	switch len(cards) {
	case 2:
		numBytes = 2
	case 5:
		numBytes = 4
	case 6:
		numBytes = 5
	case 7:
		numBytes = 6
	default:
		panic("Invalid number of cards")
	}

	// Convert the bit sequence into bytes
	byteSequence := make([]byte, numBytes)
	for i := numBytes - 1; i >= 0; i-- {
		byteSequence[i] = byte(bitSequence & 0xFF)
		bitSequence >>= 8
	}

	return byteSequence
}

// Decompress will decompress bytes into cards.
// See Compress for more details.
func Decompress(bytes []byte) (Cards, error) {
	// Convert the byte sequence into a bit sequence
	var bitSequence uint64
	for _, b := range bytes {
		bitSequence = (bitSequence << 8) | uint64(b)
	}

	// Determine the number of cards based on the number of bytes
	var numCards int
	switch len(bytes) {
	case 2:
		numCards = 2
	case 4:
		numCards = 5
	case 5:
		numCards = 6
	case 6:
		numCards = 7
	default:
		return nil, errors.New("invalid byte sequence length")
	}

	// Extract cards from the bit sequence
	var cards Cards
	for i := 0; i < numCards; i++ {
		card := Card(bitSequence & 0x3F) // Extract the last 6 bits
		cards = append([]Card{card}, cards...)
		bitSequence >>= 6
	}

	return cards, nil
}

// Clone will clone cards into new slice.
func (cc Cards) Clone() Cards {
	return append(Cards{}, cc...)
}

// IsIn returns list of booleans if cards given by
// `c` are in list.
func (cc Cards) IsIn(c Cards) (x []bool) {
	x = make([]bool, len(cc))
	for i, y := range cc {
		found := false
		for _, x := range c {
			if x == y {
				found = true
				break
			}
		}
		x[i] = found
	}
	return
}

// Has will check if one of the cards given is in list.
func (cc Cards) Has(b ...Card) (x bool) {
	for _, a := range cc.IsIn(b) {
		if a {
			return true
		}
	}
	return false
}

// Contains check if all cards are present in list.
func (cc Cards) Contains(b ...Card) (x bool) {
	for _, a := range cc.IsIn(b) {
		if !a {
			return false
		}
	}
	return true
}

// Equals compares two lists of cards.
func (cc Cards) Equals(others Cards) (x bool) {
	if len(cc) != len(others) {
		return false
	}
	for i, c := range cc {
		if c != others[i] {
			return false
		}
	}
	return true
}

// SplitHand will split hand by two cards and board.
// It will also sort cards.
func (cc Cards) SplitHand() (Cards, Cards) {
	var hand, board Cards
	switch len(cc) {
	case 1:
		return cc[:].Clone(), Cards{}
	case 2:
		hand = cc[0:2]
	default:
		hand = cc[0:2]
		board = cc[2:]
	}
	hand, board = hand.Clone(), board.Clone()

	sort.Sort(hand)
	sort.Sort(board)

	return hand, board
}

// Len returns length of cards.
func (a Cards) Len() int { return len(a) }

// Swap implements sorter interface.
func (a Cards) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

// Less implements sorter interface.
func (a Cards) Less(i, j int) bool { return a[i] < a[j] }

// String will retun cards as string with space as delimiter.
func (a Cards) String() string {
	x := []string{}
	for _, a := range a {
		x = append(x, a.String())
	}
	return fmt.Sprintf("[%s]", strings.Join(x, " "))
}

func AllCopy() Cards {
	m := make(Cards, len(allcards))
	copy(m, allcards)
	return m
}

// All returns all cards except those that are specified as argument.
func All(omit ...Card) (xx Cards) {
	xx = make(Cards, 0, len(cards))

	for _, x := range allcards {
		found := false
		for _, y := range omit {
			if x == y {
				found = true
				break
			}
		}
		if !found {
			xx = append(xx, x)
		}
	}
	return
}

// RandomCards will return random `n` cards expect those
// specified by `omit` argument
//
// Specify frand.Rand which can be either rand.Random or
// any other implementation of frand.Rand.
func RandomCards(r frand.Rand, n int, omit ...Card) (xx Cards) {
	xx = make(Cards, 0, n)
	for i := 0; i < n; i++ {
		cds := All(append(omit, xx...)...)
		if len(cds) == 0 {
			return
		}
		ll := int32(len(cds))
		r := r.Int31n(ll)
		xx = append(xx, cds[r])
	}
	return xx
}

// IsIn shorthand for Cards().IsIn()
func IsIn(a Cards, b Cards) (x []bool) {
	return a.IsIn(b)
}

// IsAnyMatch shorthand for Cards().Has()
func IsAnyMatch(a Cards, b Cards) bool {
	return a.Has(b...)
}

// IsNotAnyMatch shorthand for !Cards().IsAnyMatch()
func IsNotAnyMatch(a Cards, b Cards) bool {
	return !IsAnyMatch(a, b)
}

// Combinations generates combinations of cards.
// This is permitions without repetition.
func Combinations(combos int) (x []Cards) {
	for a := range combosCards(All(), combos) {
		x = append(x, a)
	}
	return x
}

// CombinationsLen (where order doesn't matter):
// Will return number of combinations possible.
// Follows equition: n! / r!(n − r)!
func CombinationsLen(deck, n int) int {
	ll := deck
	for i := deck - 1; i > deck-n; i-- {
		ll *= i
	}
	ll /= factorial(n)
	return ll
}

// CombinationsFrom creates combinations from cards and number of
// combinations specified.
func CombinationsFrom(cds Cards, combos int) (x []Cards) {
	for a := range combosCards(cds, combos) {
		x = append(x, a)
	}
	return x
}

// Coordinates will return coordinates of card in matrix.
func Coordinates(hole Cards) (int, int) {
	if len(hole) != 2 {
		panic("invalid number of cards")
	}
	suited := hole[0].Suite() == hole[1].Suite()
	var sorted Cards
	if hole[0].Rank() < hole[1].Rank() {
		sorted = Cards{hole[1], hole[0]}
	} else {
		sorted = Cards{hole[0], hole[1]}
	}

	r0 := 13 - int(sorted[0].Rank())
	r1 := 13 - int(sorted[1].Rank())

	if suited {
		return r0, r1
	} else {
		return r1, r0
	}
}

var mcards = [13][13][]Cards{}

func init() {
	for _, c := range Combinations(2) {
		x, y := Coordinates(c)
		sort.Sort(c)
		mcards[x][y] = append(mcards[x][y], c)
	}
}

func CardsInCoords(x, y int) []Cards {
	return mcards[x][y]
}

func CardsInCoordsWithBlockersAt(x, y int, blockers Cards) []Cards {
	xx := []Cards{}
	for _, cc := range CardsInCoords(x, y) {
		if IsAnyMatch(cc, blockers) {
			continue
		}
		xx = append(xx, cc)
	}
	return xx
}

func ForCoords(f func(x, y int, cards []Cards)) {
	for x := 0; x < 13; x++ {
		for y := 0; y < 13; y++ {
			f(x, y, mcards[x][y])
		}
	}
}

// CombinationsStreamFrom creates combinations from cards same as
// CombinationsFrom expect returns channel instead.
// Channel will close it self once combinations run out.
func CombinationsStreamFrom(cds Cards, combos int) chan Cards {
	return combosCards(cds, combos)
}

func combosCards(iterable Cards, r int) chan Cards {
	ch := make(chan Cards)
	go func() {
		length := len(iterable)
		for comb := range combos(length, r) {
			result := make(Cards, r)
			for i, val := range comb {
				result[i] = iterable[val]
			}
			ch <- result
		}

		close(ch)
	}()
	return ch
}

// combos generates, from two natural numbers n > r,
// all the possible combinations of r indexes taken from 0 to n-1.
// For example if n=3 and r=2, the result will be:
// [0,1], [0,2] and [1,2]
func combos(n, r int) <-chan []int {
	if r > n {
		panic("Invalid arguments")
	}
	ch := make(chan []int)
	go func() {
		result := make([]int, r)
		for i := range result {
			result[i] = i
		}
		temp := make([]int, r)
		copy(temp, result) // avoid overwriting of result
		ch <- temp
		for {
			for i := r - 1; i >= 0; i-- {
				if result[i] < i+n-r {
					result[i]++
					for j := 1; j < r-i; j++ {
						result[i+j] = result[i] + j
					}
					temp := make([]int, r)
					copy(temp, result) // avoid overwriting of result
					ch <- temp
					break
				}
			}
			if result[0] >= n-r {
				break
			}
		}
		close(ch)
	}()
	return ch
}

func factorial(n int) (result int) {
	if n > 0 {
		result = n * factorial(n-1)
		return result
	}
	return 1
}

package card

import (
	"github.com/pokerdroid/poker/frand"
)

type RangeDist [1326]float64

var rcards = map[int]Cards{}
var rindex = map[[2]Card]int{}
var rcoords = map[int][3]int{}

func init() {
	var i int
	ForCoords(func(x, y int, cards []Cards) {
		for z, c := range cards {
			rindex[[2]Card{c[0], c[1]}] = i
			rcards[i] = c
			rcoords[i] = [3]int{x, y, z}
			i++
		}
	})
}

func NewRangeDist(m Matrix) RangeDist {
	r := RangeDist{}

	// Iterate through all possible coordinates and their corresponding cards
	ForCoords(func(x, y int, cards []Cards) {
		// Get the probability from the matrix for these coordinates
		prob := m[x][y]

		// If there are multiple possible card combinations for these coordinates,
		// distribute the probability evenly among them
		if len(cards) > 0 {
			prob = prob / float64(len(cards))
		}

		// Assign the probability to each specific card combination
		for _, c := range cards {
			idx := RangeIndex(c)
			r[idx] = prob
		}
	})

	return r.Normalize()
}

func NewUniformRangeDist() RangeDist {
	r := RangeDist{}
	for i := range r {
		r[i] = 1. / 1326.
	}
	return r
}

func (r RangeDist) Sub(other RangeDist) RangeDist {
	nr := RangeDist{}
	for i := range r {
		nr[i] = r[i] - other[i]
	}
	return nr
}

func (r RangeDist) Sample(rng frand.Rand) Cards {
	var cumProb float64
	x := rng.Float64()

	for i, p := range r {
		cumProb += float64(p)
		if cumProb > x {
			return RangeCards(i)
		}
	}

	// Leave room for floating point error.
	if cumProb < 1.0-1e-6 {
		panic("probability distribution does not sum to 1!")
	}

	return RangeCards(len(r) - 1)
}

func (r RangeDist) Matrix() Matrix {
	m := Matrix{}
	for i := range r {
		x, y, _ := RangeCoords(i)
		m[x][y] += r[i]
	}
	return m.Normalize()
}

func (r RangeDist) Normalize() RangeDist {
	total := r.Sum()
	for i := range r {
		r[i] /= total
	}
	return r
}

func (r RangeDist) Sum() float64 {
	sum := float64(0)
	for i := range r {
		sum += r[i]
	}
	return sum
}

// RangeIndex returns the stable index of a 2-card combination in range [0..1325].
// Cards are automatically ordered so that if c2 < c1, they are swapped first.
func RangeIndex(cc Cards) int {
	return rindex[[2]Card{cc[0], cc[1]}]
}

// RangeCards returns the two cards corresponding to index i in range [0..1325].
func RangeCards(i int) Cards {
	return rcards[i]
}

func RangeCoords(i int) (int, int, int) {
	return rcoords[i][0], rcoords[i][1], rcoords[i][2]
}

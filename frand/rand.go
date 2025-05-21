package frand

import (
	"fmt"
	"math/rand"
)

// Rand is an interface that mimics some of the
// key methods from the math/rand package.
type Rand interface {
	Int63() int64
	Int63n(n int64) int64
	Int31n(n int32) int32
	Intn(n int) int
	Float64() float64
	Float32() float32
	Perm(n int) []int
	Shuffle(n int, swap func(i, j int))
}

// SampleIndex will sample index
func SampleIndex[T float32 | float64](rng Rand, pv []T, eps T) int {
	var cumProb T
	x := rng.Float32()

	for i, p := range pv {
		cumProb += p
		if cumProb > T(x) {
			return i
		}
	}

	if cumProb < 1.0-eps { // Leave room for floating point error.
		panic(fmt.Sprintf("probability distribution does not sum to 1! cumProb: %f, eps: %f, pv: %v,", cumProb, eps, pv))
	}

	return len(pv) - 1
}

func Clone(rng Rand) Rand {
	cl, ok := rng.(interface{ Clone() Rand })
	if ok {
		return cl.Clone()
	}
	// if math/rand
	if x, ok := rng.(*rand.Rand); ok {
		return rand.New(rand.NewSource(x.Int63() + 1))
	}

	panic(fmt.Sprintf("rng does not implement Clone(): %T", rng))
}

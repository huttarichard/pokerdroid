package frand

import "hash/maphash"

// hash is an implementation of the Frand interface using maphash.
type hash struct{}

// NewHash returns a new Frand implementation using maphash.
func NewHash() Rand {
	return &hash{}
}

func (hi *hash) Clone() Rand {
	return &hash{}
}

// Int63 returns a non-negative pseudo-random 63-bit integer as an int64.
func (hi *hash) Int63() int64 {
	return Int63()
}

// Int63n returns, as an int64, a non-negative pseudo-random number in the half-open interval [0,n).
// It panics if n <= 0.
func (hi *hash) Int63n(n int64) int64 {
	return Int63n(n)
}

// Int31n returns, as an int32, a non-negative pseudo-random number in the half-open interval [0,n).
// It panics if n <= 0.
func (hi *hash) Int31n(n int32) int32 {
	return Int31n(n)
}

// Intn returns, as an int, a non-negative pseudo-random number in the half-open interval [0,n).
// It panics if n <= 0.
func (hi *hash) Intn(n int) int {
	return Intn(n)
}

// Float64 returns, as a float64, a pseudo-random number in the half-open interval [0.0,1.0).
func (hi *hash) Float64() float64 {
	return Float64()
}

// Float32 returns, as a float32, a pseudo-random number in the half-open interval [0.0,1.0).
func (hi *hash) Float32() float32 {
	return Float32()
}

// Perm returns, as a slice of n ints, a pseudo-random permutation of the integers
// in the half-open interval [0,n).
func (hi *hash) Perm(n int) []int {
	return Perm(n)
}

// Shuffle pseudo-randomizes the order of elements.
// n is the number of elements. Shuffle panics if n < 0.
// swap swaps the elements with indexes i and j.
func (hi *hash) Shuffle(n int, swap func(i, j int)) {
	Shuffle(n, swap)
}

// Int63 returns a non-negative pseudo-random 63-bit integer as an int64.
func Int63() int64 {
	return int64(hashValue() >> 1)
}

// Int63n returns, as an int64, a non-negative pseudo-random number in the half-open interval [0,n).
// It panics if n <= 0.
func Int63n(n int64) int64 {
	if n <= 0 {
		panic("non-positive argument to Int63n")
	}
	max := int64((1 << 63) - 1 - (1<<63)%uint64(n))
	v := Int63()
	for v > max {
		v = Int63()
	}
	return v % n
}

// Int31n returns, as an int32, a non-negative pseudo-random number in the half-open interval [0,n).
// It panics if n <= 0.
func Int31n(n int32) int32 {
	if n <= 0 {
		panic("non-positive argument to Int31n")
	}
	return int32(Int63n(int64(n)))
}

// Intn returns, as an int, a non-negative pseudo-random number in the half-open interval [0,n).
// It panics if n <= 0.
func Intn(n int) int {
	if n <= 0 {
		panic("non-positive argument to Intn")
	}
	return int(Int31n(int32(n)))
}

// Float64 returns, as a float64, a pseudo-random number in the half-open interval [0.0,1.0).
func Float64() float64 {
	return float64(Int63()) / (1 << 63)
}

// Float32 returns, as a float32, a pseudo-random number in the half-open interval [0.0,1.0).
func Float32() float32 {
	return float32(Float64())
}

// Perm returns, as a slice of n ints, a pseudo-random permutation of the integers
// in the half-open interval [0,n).
func Perm(n int) []int {
	m := make([]int, n)
	for i := 1; i < n; i++ {
		j := Intn(i + 1)
		m[i] = m[j]
		m[j] = i
	}
	return m
}

// Shuffle pseudo-randomizes the order of elements.
// n is the number of elements. Shuffle panics if n < 0.
// swap swaps the elements with indexes i and j.
func Shuffle(n int, swap func(i, j int)) {
	if n < 0 {
		panic("invalid argument to Shuffle")
	}

	// Fisher-Yates shuffle: https://en.wikipedia.org/wiki/Fisher%E2%80%93Yates_shuffle
	// Shuffle really ought not be called with n that doesn't fit in 32 bits.
	// Not only will it take a very long time, but with 2³¹! possible permutations,
	// there's no way that any PRNG can have a big enough internal state to
	// generate even a minuscule percentage of the possible permutations.
	// Nevertheless, the right API signature accepts an int n, so handle it as best we can.
	i := n - 1
	for ; i > 1<<31-1-1; i-- {
		j := int(Int63n(int64(i + 1)))
		swap(i, j)
	}
	for ; i > 0; i-- {
		j := int(Int31n(int32(i + 1)))
		swap(i, j)
	}
}

func hashValue() uint64 {
	return new(maphash.Hash).Sum64()
}

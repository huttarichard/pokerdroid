package frand

import (
	"math/rand/v2"
)

// unsafe is a lock-free implementation
type unsafe struct {
	*rand.Rand
	s Seed
}

func NewUnsafe() Rand {
	s := NewSeedUnix()
	return &unsafe{
		s:    s,
		Rand: rand.New(s.Seed()),
	}
}

func NewUnsafeSeed(seed Seed) Rand {
	return &unsafe{
		s:    seed,
		Rand: rand.New(seed.Seed()),
	}
}

func NewUnsafeInt(seed int64) Rand {
	s := NewSeed(seed)
	return &unsafe{
		s:    s,
		Rand: rand.New(s.Seed()),
	}
}

func (u *unsafe) Clone() Rand {
	return NewUnsafeSeed(NewSeed(u.Int63()))
}

func (u *unsafe) Int31n(n int32) int32 {
	return u.Rand.Int32N(n)
}

func (u *unsafe) Int63n(n int64) int64 {
	return u.Rand.Int64N(n)
}

func (u *unsafe) Int63() int64 {
	return u.Rand.Int64()
}

func (u *unsafe) Intn(n int) int {
	return u.Rand.IntN(n)
}

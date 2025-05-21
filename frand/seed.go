package frand

import (
	"math/rand"
	"time"
)

type Seed int64

func NewSeed(seed int64) Seed {
	return Seed(seed)
}

func NewSeedUnix() Seed {
	return NewSeed(int64(time.Now().UnixNano()))
}

func (s Seed) Seed() rand.Source64 {
	return rand.NewSource(int64(s)).(rand.Source64)
}

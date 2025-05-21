package sampler

import (
	"github.com/pokerdroid/poker/frand"
	"github.com/pokerdroid/poker/policy"
)

// RobustSampler implements cfr.Sampler by sampling a fixed number of actions
// uniformly randomly.
type Robust struct {
	k int
}

func NewRobust(k int) *Robust {
	return &Robust{k: k}
}

func (rs *Robust) Sample(rng frand.Rand, actions int, p *policy.Policy, gi uint64, depth uint8, qs []float64) {
	if actions <= rs.k {
		for i := range qs {
			qs[i] = 1.0
		}
		return
	}

	copy(qs, zeros[:actions])

	for i := 0; i < rs.k; i++ {
		qs[i] = float64(rs.k) / float64(actions)
	}

	rng.Shuffle(actions, func(i, j int) {
		qs[i], qs[j] = qs[j], qs[i]
	})
}

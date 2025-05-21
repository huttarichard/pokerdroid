package sampler

import (
	"github.com/pokerdroid/poker/frand"
	"github.com/pokerdroid/poker/policy"
)

// Combined implements Sampler by sampling one player action
// according to the current strategy or sampling k actions uniformly at random.
type IterCombined struct {
	s1, s2 Sampler
	iter   uint64
}

// If eps is 0, then we sample with strategy
// If eps is 1, then we sample with rand.
func NewIterCombined(s1, s2 Sampler, iter uint64) *IterCombined {
	return &IterCombined{s1: s1, s2: s2, iter: iter}
}

func (os *IterCombined) Sample(rng frand.Rand, actions int, p *policy.Policy, gi uint64, depth uint8, qs []float64) {
	if p.Iteration > os.iter {
		os.s1.Sample(rng, actions, p, gi, depth, qs)
	}
	os.s2.Sample(rng, actions, p, gi, depth, qs)
}

type DepthCombined struct {
	s1, s2 Sampler
	depth  uint8
}

func NewDepthCombined(s1, s2 Sampler, depth uint8) *DepthCombined {
	return &DepthCombined{s1: s1, s2: s2, depth: depth}
}

func (os *DepthCombined) Sample(rng frand.Rand, actions int, p *policy.Policy, gi uint64, depth uint8, qs []float64) {
	if depth < os.depth {
		os.s1.Sample(rng, actions, p, gi, depth, qs)
	} else {
		os.s2.Sample(rng, actions, p, gi, depth, qs)
	}
}

package sampler

import (
	"github.com/pokerdroid/poker/frand"
	"github.com/pokerdroid/poker/policy"
)

type External struct{}

func NewExternal() *External {
	return &External{}
}

func (es *External) Sample(rng frand.Rand, actions int, p *policy.Policy, gi uint64, depth uint8, qs []float64) {
	copy(qs, onces[:actions])
}

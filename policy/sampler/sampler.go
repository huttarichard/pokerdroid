package sampler

import (
	"github.com/pokerdroid/poker/frand"
	"github.com/pokerdroid/poker/policy"
)

// Sampler selects a subset of child nodes to traverse.
type Sampler interface {
	// Sample returns a vector of sampling probabilities for a
	// subset of the N children of Node. Children with
	// p > 0 will be traversed. The returned slice may be reused
	// between calls to sample; a caller must therefore copy the
	// values before the next call to Sample.
	Sample(rng frand.Rand, actions int, p *policy.Policy, gi uint64, depth uint8, qs []float64)
}

var onces = []float64{}
var zeros = []float64{}

func init() {
	for i := 0; i < 100; i++ {
		onces = append(onces, 1.0)
	}

	for i := 0; i < 100; i++ {
		zeros = append(zeros, 0)
	}
}

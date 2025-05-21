package sampler

import (
	"github.com/pokerdroid/poker/frand"
	"github.com/pokerdroid/poker/policy"
)

// SamplerOutcome implements Sampler by sampling one player action
// according to the current strategy.
type Outcome struct {
	eps float64
}

// If eps is 0, then we sample with strategy
// If eps is 1, then we sample with rand.
func NewOutcome(eps float64) *Outcome {
	return &Outcome{eps: eps}
}

func (os *Outcome) Sample(rng frand.Rand, actions int, p *policy.Policy, gi uint64, depth uint8, qs []float64) {
	var selected int

	if rng.Float64() < os.eps {
		selected = rng.Intn(actions)
	} else {
		selected = frand.SampleIndex(rng, p.Strategy, 0.0001)
	}

	q := os.eps * (1.0 / float64(actions))     // Sampled due to exploration.
	q += (1.0 - os.eps) * p.Strategy[selected] // Sampled due to strategy.

	copy(qs, zeros[:actions])
	qs[selected] = q
}

// OutcomeDecay implements Sampler by sampling one player action
// with epsilon decay over iterations.
type OutcomeDecay struct {
	epsMin  float64
	epsMax  float64
	maxIter uint32
}

// NewOutcomeDecay creates a new OutcomeDecay sampler.
// epsMax is initial exploration rate that decays to epsMin over maxIter iterations.
func NewOutcomeDecay(epsMin, epsMax float64, maxIter uint32) *OutcomeDecay {
	return &OutcomeDecay{
		epsMin:  epsMin,
		epsMax:  epsMax,
		maxIter: maxIter,
	}
}

func (os *OutcomeDecay) Sample(rng frand.Rand, actions int, p *policy.Policy, gi uint64, depth uint8, qs []float64) {
	var selected int

	// Calculate current epsilon based on iteration progress
	progress := float64(p.Iteration) / float64(os.maxIter)
	if progress > 1.0 {
		progress = 1.0
	}
	eps := os.epsMax - progress*(os.epsMax-os.epsMin)

	if rng.Float64() < eps {
		selected = rng.Intn(actions)
	} else {
		selected = frand.SampleIndex(rng, p.Strategy, 0.0001)
	}

	q := eps * (1.0 / float64(actions))     // Sampled due to exploration
	q += (1.0 - eps) * p.Strategy[selected] // Sampled due to strategy

	copy(qs, zeros[:actions])
	qs[selected] = q
}

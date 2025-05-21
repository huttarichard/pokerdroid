package sampler

import (
	"github.com/pokerdroid/poker/float/f64"
	"github.com/pokerdroid/poker/frand"
	"github.com/pokerdroid/poker/policy"
)

// MultiOutcome implements cfr.Sampler by sampling at most k player actions
// with probability according to the current strategy.
type MultiOutcome struct {
	k    int
	eps  float64
	p    []float64
	pool *f64.Pool
}

func NewMultiOutcome(k int, explorationEps float64) *MultiOutcome {
	return &MultiOutcome{
		k:    k,
		eps:  explorationEps,
		p:    make([]float64, k),
		pool: f64.NewPool(k),
	}
}

// func (os *MultiOutcomeSampler) Sample(node cfr.GameTreeNode, policy cfr.NodePolicy) []float64 {
func (os *MultiOutcome) Sample(rng frand.Rand, actions int, p *policy.Policy, gi uint64, depth uint8, qs []float64) {
	if actions <= os.k {
		copy(qs, onces[:actions])
		return
	}

	copy(qs, zeros[:actions])

	q := os.pool.Alloc(actions)
	defer os.pool.Free(q)

	copy(q.Slice, p.Strategy)
	f64.AddConst(os.eps/float64(actions), q.Slice)
	f64.ScalUnitary(1.0/f64.Sum(q.Slice), q.Slice) // Renormalize.

	// Compute probability of choosing i if we draw k times.
	qEff := os.pool.Alloc(actions)
	defer os.pool.Free(qEff)

	for j := range q.Slice {
		qEff.Slice[j] = os.ck(q.Slice, j, os.k)
	}

	for i := 0; i < os.k; i++ {
		sampled := frand.SampleIndex(rng, q.Slice, 0.0001)
		qs[sampled] = qEff.Slice[sampled]

		// Remove sampled action from being re-sampled.
		qSample := q.Slice[sampled]
		q.Slice[sampled] = 0
		f64.ScalUnitary(1.0/(1.0-qSample), q.Slice)
	}
}

func (os *MultiOutcome) ck(p []float64, j, k int) float64 {
	if k == 1 {
		return p[j]
	}

	var des float64
	for i := range p {
		if i == j || p[i] <= 0 {
			continue
		}

		z := os.pool.Alloc(len(p))
		copy(z.Slice, p)
		z.Slice[i] = 0
		f64.ScalUnitary(1.0/(1-p[i]), z.Slice)
		des += p[i] * os.ck(z.Slice, j, k-1)
		os.pool.Free(z)
	}

	return p[j] + des
}

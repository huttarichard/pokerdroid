package sampler

import (
	"github.com/pokerdroid/poker/float/f64"
	"github.com/pokerdroid/poker/frand"
	"github.com/pokerdroid/poker/policy"
)

// AvgStrat implements cfr.Sampler by sampling some player actions
// according to the current average strategy strategy.
type AvgStrat struct {
	eps  float64
	tau  float64
	beta float64
}

func NewAvg(eps, tau, beta float64) *AvgStrat {
	return &AvgStrat{
		eps:  eps,
		tau:  tau,
		beta: beta,
	}
}

func (as *AvgStrat) Sample(rng frand.Rand, actions int, p *policy.Policy, gi uint64, depth uint8, qs []float64) {
	x := rng.Float64()
	sSum := f64.Sum(p.StrategySum)
	for i := range qs {
		rho := as.beta + as.tau*p.StrategySum[i]
		rho /= as.beta + sSum

		if rho < as.eps {
			rho = as.eps
		}

		if x < rho {
			qs[i] = f64.Min(rho, 1.0)
		} else {
			qs[i] = 0
		}
	}
}

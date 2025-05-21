package cfr

import (
	"github.com/pokerdroid/poker/dealer"
	"github.com/pokerdroid/poker/frand"
	"github.com/pokerdroid/poker/policy"
)

type Params struct {
	Iterations uint64
	Rng        frand.Rand
	Sampler    dealer.Dealer
}

type Task struct {
	TraversingID uint8
	Sample       dealer.Sample
	Update       *policy.Update
	Rng          frand.Rand
}

type Runner interface {
	Run(Params) (ev float64, up uint64)
}

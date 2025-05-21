package deep

import (
	"github.com/nlpodyssey/spago/initializers"
	"github.com/nlpodyssey/spago/mat"
	"github.com/nlpodyssey/spago/mat/rand"
	"github.com/nlpodyssey/spago/nn"
	"github.com/nlpodyssey/spago/nn/normalization/batchnorm"
	"github.com/pokerdroid/poker/frand"
)

func ApplyXavierUniform(m nn.Model, rng frand.Rand) {
	r := rand.NewLockedRand(uint64(rng.Int63()))

	nn.Apply(m, func(model nn.Model) {
		switch model.(type) {
		case *batchnorm.Model:
			return
		default:
		}
		// Initialize only the linear layers
		nn.ForEachParam(model, func(param *nn.Param) {
			initializers.XavierUniform(param.Value().(mat.Matrix), 1.0, r)
		})
	})
}

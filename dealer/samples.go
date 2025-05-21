package dealer

import (
	"github.com/pokerdroid/poker/abs"
	"github.com/pokerdroid/poker/frand"
	"github.com/pokerdroid/poker/table"
	"github.com/pokerdroid/poker/tree"
)

type Turner interface {
	GetTurnPos() uint8
}

type Sample interface {
	Sample(s table.Street)
	Cluster(n Turner, abs abs.Mapper) abs.Cluster
	Utility(n *tree.Terminal, pID uint8) float64
}

type Dealer interface {
	Sample(rng frand.Rand) (Sample, error)
	Clone() Dealer
	Copy(rng frand.Rand, s Sample) (Sample, error)
	Put(Sample)
}

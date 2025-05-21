package baselinenn

import (
	"math"
	"sort"

	"github.com/nlpodyssey/spago/mat"
	"github.com/pokerdroid/poker/abs"
	absp "github.com/pokerdroid/poker/abs/pack"
	"github.com/pokerdroid/poker/frand"
	"github.com/pokerdroid/poker/tree"
)

type ValidateParams struct {
	Abs    *absp.Abs
	Model  *Model
	Root   *tree.Root
	Rng    frand.Rand
	Depth  int
	Rounds int
}

func Validate(params ValidateParams) (float64, int) {
	var diff float64
	var samples int

	tree.MustVisit(params.Root, params.Depth, func(n tree.Node, children []tree.Node, depth int) bool {
		px, ok := n.(*tree.Player)
		if !ok {
			return true
		}

		var clusters abs.Clusters

		for cl := range px.Actions.Policies.Map {
			clusters = append(clusters, cl)
		}

		sort.Slice(clusters, func(i, j int) bool {
			return clusters[i] < clusters[j]
		})

		cl := clusters[params.Rng.Intn(len(clusters))]
		p := px.Actions.Policies.Map[cl]

		for round := 0; round < params.Rounds; round++ {
			for idx := range px.Actions.Actions {
				equity := params.Abs.Equity(px.State.Street, cl)

				trajs, ss := Encode[float32](EncodeParams{
					Params: params.Root.Params,
					Player: px,
					Equity: equity,
				})

				total, _ := EncodeFeatures[float32](
					params.Root.Params.NumPlayers,
					params.Root.Params.MaxActionsPerRound,
				)

				mix := mat.NewDense[float32](
					mat.WithBacking(trajs),
					mat.WithShape(total),
					mat.WithGrad(false),
				)

				pred := params.Model.Forward(mix)
				baseline := pred[0].Item().F64()

				diff += math.Abs(float64(p.Baseline[idx])/ss.Float64() - baseline)
				samples++
			}
		}

		return true
	})

	return diff / float64(samples), samples
}

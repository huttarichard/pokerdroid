package deep

import (
	"github.com/nlpodyssey/spago/mat"
	"github.com/nlpodyssey/spago/mat/float"
)

type Sample[T float.DType] struct {
	Features []T
	Labels   []T
}

type Samples[T float.DType] []Sample[T]

type Batch struct {
	X []mat.Tensor
	Y []mat.Tensor
}

func NewBatchFromSamples[T float.DType](t Samples[T], size int) (*Batch, Samples[T]) {
	if len(t) == 0 {
		return nil, nil
	}

	totalSize := len(t[:size])
	allX := make([]mat.Tensor, totalSize)
	allY := make([]mat.Tensor, totalSize)

	// Convert trajectories to tensors first
	for i, traj := range t[:size] {
		allX[i] = mat.NewDense[T](
			mat.WithBacking(traj.Features),
			mat.WithShape(len(t[0].Features)),
			mat.WithGrad(true),
		)

		allY[i] = mat.NewDense[T](
			mat.WithBacking(traj.Labels),
			mat.WithShape(len(t[0].Labels)),
			mat.WithGrad(true),
		)
	}

	return &Batch{X: allX, Y: allY}, t[size:]

}

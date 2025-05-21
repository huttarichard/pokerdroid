package deep

import (
	"testing"

	"github.com/nlpodyssey/spago/ag"
	"github.com/nlpodyssey/spago/mat"
	"github.com/nlpodyssey/spago/mat/float"
	"github.com/stretchr/testify/assert"
)

func TestHuberLoss(t *testing.T) {
	t.Run("float32", func(t *testing.T) { testHuberLoss[float32](t, 1.0e-6) })
	t.Run("float64", func(t *testing.T) { testHuberLoss[float64](t, 1.0e-12) })
}

func testHuberLoss[T float.DType](t *testing.T, tol T) {
	// 1) Setup input, target, delta
	x := mat.NewDense[T](mat.WithBacking([]T{0.0, 2.5, 4.0}), mat.WithGrad(true))
	y := mat.NewDense[T](mat.WithBacking([]T{0.0, 1.0, 2.0}))
	delta := 1.0

	// 2) Compute Huber loss with reduceMean = false
	loss := Huber(x, y, delta, false)

	// Sum of the "Huber" terms across 3 elements => ~2.5
	assert.InDelta(t, 2.5, loss.Value().Item().F64(), float64(tol))

	// 3) Backward
	ag.Backward(loss)

	// For this example:
	// d = [0, 1.5, 2.0]; gradient = sign(d)*delta if |d|>delta else d
	// => [0, 1.0, 1.0]
	assert.InDeltaSlice(t, []T{0.0, 1.0, 1.0}, x.Grad().Data(), float64(tol))

	// 4) Test again with reduceMean = true
	x2 := mat.NewDense[T](mat.WithBacking([]T{0.0, 2.5, 4.0}), mat.WithGrad(true))
	y2 := mat.NewDense[T](mat.WithBacking([]T{0.0, 1.0, 2.0}))
	loss2 := Huber(x2, y2, delta, true)

	// The total is 2.5 for 3 elements => 2.5 / 3 = ~0.8333
	assert.InDelta(t, 0.8333333333333333, loss2.Value().Item().F64(), float64(tol))

	ag.Backward(loss2)
	// The gradient is the same shape but divided by 3 => [0, ~0.3333, ~0.3333]
	assert.InDeltaSlice(t, []T{0.0, 0.3333333333333333, 0.3333333333333333}, x2.Grad().Data(), float64(tol))
}

func TestCrossEntropyLoss(t *testing.T) {
	t.Run("float32", func(t *testing.T) { testCrossEntropyLoss[float32](t, 1.0e-6) })
	t.Run("float64", func(t *testing.T) { testCrossEntropyLoss[float64](t, 1.0e-12) })
}

func testCrossEntropyLoss[T float.DType](t *testing.T, tol T) {
	// 1) Setup predicted (logits) and target (true distribution)
	predicted := mat.NewDense[T](mat.WithBacking([]T{0.8, 0.1, 0.1}), mat.WithGrad(true))
	target := mat.NewDense[T](mat.WithBacking([]T{1.0, 0.0, 0.0}))

	// 2) Compute Cross Entropy loss with reduceMean = false
	loss := CrossEntropy(predicted, target, false)

	// -sum(target * log(predicted)) = -(1.0 * log(0.8) + 0.0 * log(0.1) + 0.0 * log(0.1))
	// = -log(0.8) ≈ 0.223
	assert.InDelta(t, 0.2231435513142097, loss.Value().Item().F64(), float64(tol))

	// 3) Backward
	ag.Backward(loss)

	// Gradient with respect to predicted should be -target/predicted
	// => [-1.25, 0, 0]
	assert.InDeltaSlice(t, []T{-1.25, 0, 0}, predicted.Grad().Data(), float64(tol))

	// 4) Test with reduceMean = true
	predicted2 := mat.NewDense[T](mat.WithBacking([]T{0.8, 0.1, 0.1}), mat.WithGrad(true))
	target2 := mat.NewDense[T](mat.WithBacking([]T{1.0, 0.0, 0.0}))
	loss2 := CrossEntropy(predicted2, target2, true)

	// Same value divided by number of elements (3)
	assert.InDelta(t, 0.07438118377140324, loss2.Value().Item().F64(), float64(tol))

	ag.Backward(loss2)
	// Gradient should be divided by number of elements
	assert.InDeltaSlice(t, []T{-0.41666666666666663, 0, 0}, predicted2.Grad().Data(), float64(tol))
}

func TestCrossEntropySeq(t *testing.T) {
	t.Run("float32", func(t *testing.T) { testCrossEntropySeq[float32](t, 1.0e-6) })
	t.Run("float64", func(t *testing.T) { testCrossEntropySeq[float64](t, 1.0e-12) })
}

func testCrossEntropySeq[T float.DType](t *testing.T, tol T) {
	// 1) Setup sequence of predictions and targets
	predicted := []mat.Tensor{
		mat.NewDense[T](mat.WithBacking([]T{0.8, 0.1, 0.1}), mat.WithGrad(true)),
		mat.NewDense[T](mat.WithBacking([]T{0.1, 0.8, 0.1}), mat.WithGrad(true)),
	}
	target := []mat.Tensor{
		mat.NewDense[T](mat.WithBacking([]T{1.0, 0.0, 0.0})),
		mat.NewDense[T](mat.WithBacking([]T{0.0, 1.0, 0.0})),
	}

	// 2) Test sequence loss with reduceMean = false
	loss := CrossEntropySeq(predicted, target, false)

	// Sum of individual losses: -log(0.8) + -log(0.8) ≈ 0.446
	assert.InDelta(t, 0.4462871026284194, loss.Value().Item().F64(), float64(tol))

	// 3) Backward
	ag.Backward(loss)

	// Gradients for first prediction
	assert.InDeltaSlice(t, []T{-1.25, 0, 0}, predicted[0].Grad().Data(), float64(tol))
	// Gradients for second prediction
	assert.InDeltaSlice(t, []T{0, -1.25, 0}, predicted[1].Grad().Data(), float64(tol))

	// 4) Test with reduceMean = true
	predicted2 := []mat.Tensor{
		mat.NewDense[T](mat.WithBacking([]T{0.8, 0.1, 0.1}), mat.WithGrad(true)),
		mat.NewDense[T](mat.WithBacking([]T{0.1, 0.8, 0.1}), mat.WithGrad(true)),
	}
	loss2 := CrossEntropySeq(predicted2, target, true)

	// Average over sequence length (2)
	assert.InDelta(t, 0.2231435513142097, loss2.Value().Item().F64(), float64(tol))

	ag.Backward(loss2)
	// Gradients should be divided by sequence length
	assert.InDeltaSlice(t, []T{-0.625, 0, 0}, predicted2[0].Grad().Data(), float64(tol))
	assert.InDeltaSlice(t, []T{0, -0.625, 0}, predicted2[1].Grad().Data(), float64(tol))
}

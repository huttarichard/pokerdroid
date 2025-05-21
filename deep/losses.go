package deep

import (
	"image/color"

	"github.com/fogleman/gg"
	"github.com/nlpodyssey/spago/ag"
	"github.com/nlpodyssey/spago/mat"
)

// Huber measures the Huber loss between each element in the input x and target y, controlled by
// the threshold (delta). Below the threshold, it behaves like MSE; above it, it becomes linear
// in order to reduce the effect of outliers. If reduceMean is true, it returns the average loss;
// otherwise it returns the sum of the losses.
//
// Huber(d) = { 0.5 * (d^2)              if |d| ≤ δ
//
//	δ * (|d| - 0.5 * δ)      otherwise }
//
// Here, d = x - y.
func Huber(x, y mat.Tensor, delta float64, reduceMean bool) mat.Tensor {
	// 1) Compute d = x - y, then |d|
	d := ag.Sub(x, y)
	absD := ag.Abs(d)

	// 2) Build a scalar tensor from 'delta'
	//    then multiply it by the shape of absD (via OnesLike) to broadcast
	//    the scalar across all elements. This avoids dimension mismatch with Min().
	deltaMat := x.Value().(mat.Matrix).NewScalar(delta)
	deltaVec := ag.ProdScalar(x.Value().(mat.Matrix).OnesLike(), deltaMat)

	// 3) clipped = min(|d|, deltaVec)
	clipped := ag.Min(absD, deltaVec)

	// 4) 0.5 * (clipped)^2
	halfSqr := ag.ProdScalar(ag.Square(clipped), x.Value().(mat.Matrix).NewScalar(0.5))

	// 5) deltaVec * (|d| - clipped)
	linear := ag.Prod(deltaVec, ag.Sub(absD, clipped))

	// 6) Combine
	loss := ag.Add(halfSqr, linear)

	// 7) reduceMean or sum
	if reduceMean {
		return ag.ReduceMean(loss)
	}
	return ag.ReduceSum(loss)
}

// HuberSeq calculates the Huber loss on multiple (predicted, target) pairs.
// It sums the Huber loss across the entire sequence, optionally averaging it
// by the number of elements if reduceMean is true.
func HuberSeq(predicted, target []mat.Tensor, delta float64, reduceMean bool) mat.Tensor {
	if len(predicted) != len(target) {
		panic("losses: predicted and target slices must have the same length")
	}
	if len(predicted) == 0 {
		panic("losses: no tensors provided to HuberSeq")
	}

	// Accumulate the Huber loss across the sequence
	loss := Huber(predicted[0], target[0], delta, false)
	for i := 1; i < len(predicted); i++ {
		loss = ag.Add(loss, Huber(predicted[i], target[i], delta, false))
	}

	// Optionally divide by length to get mean
	if reduceMean {
		return ag.DivScalar(loss, loss.Value().(mat.Matrix).NewScalar(float64(len(predicted))))
	}
	return loss
}

// CrossEntropy computes cross-entropy between two probability distributions.
// The predicted and target tensors must have the same dimensions.
func CrossEntropy(predicted, target mat.Tensor, reduceMean bool) mat.Tensor {
	// Clip predicted values to avoid log(0)
	const epsilon = 1e-7
	eps := predicted.Value().(mat.Matrix).NewScalar(epsilon)
	oneMinusEps := predicted.Value().(mat.Matrix).NewScalar(1.0 - epsilon)

	// Create broadcasted versions using OnesLike
	epsVec := ag.ProdScalar(predicted.Value().(mat.Matrix).OnesLike(), eps)
	oneMinusEpsVec := ag.ProdScalar(predicted.Value().(mat.Matrix).OnesLike(), oneMinusEps)

	// Clip: max(epsilon, min(1-epsilon, predicted))
	clipped := ag.Max(epsVec, ag.Min(oneMinusEpsVec, predicted))

	// Calculate -sum(target * log(predicted))
	logPred := ag.Log(clipped)
	prod := ag.Prod(target, logPred)
	loss := ag.Neg(prod)

	if reduceMean {
		return ag.ReduceMean(loss)
	}
	return ag.ReduceSum(loss)
}

// CrossEntropySeq calculates the cross-entropy loss on multiple (predicted, target) pairs.
func CrossEntropySeq(predicted, target []mat.Tensor, reduceMean bool) mat.Tensor {
	if len(predicted) != len(target) {
		panic("losses: predicted and target slices must have the same length")
	}
	if len(predicted) == 0 {
		panic("losses: no tensors provided to CrossEntropySeq")
	}

	// Accumulate the cross-entropy loss across the sequence
	loss := CrossEntropy(predicted[0], target[0], false)
	for i := 1; i < len(predicted); i++ {
		loss = ag.Add(loss, CrossEntropy(predicted[i], target[i], false))
	}

	if reduceMean {
		return ag.DivScalar(loss, loss.Value().(mat.Matrix).NewScalar(float64(len(predicted))))
	}
	return loss
}

// KLDivergence computes the Kullback-Leibler divergence between two probability distributions.
// The predicted and target tensors must have the same dimensions.
func KLDivergence(predicted, target mat.Tensor, reduceMean bool) mat.Tensor {
	// 1. Cross entropy component with numerical stability
	const epsilon = 1e-7

	eps := predicted.Value().(mat.Matrix).NewScalar(epsilon)
	oneMinusEps := predicted.Value().(mat.Matrix).NewScalar(1.0 - epsilon)

	epsVec := ag.ProdScalar(predicted.Value().(mat.Matrix).OnesLike(), eps)
	oneMinusEpsVec := ag.ProdScalar(predicted.Value().(mat.Matrix).OnesLike(), oneMinusEps)

	clipped := ag.Max(epsVec, ag.Min(oneMinusEpsVec, predicted))

	// 2. KL divergence: KL(P||Q) = sum(P * log(P/Q))
	ratio := ag.Div(target, clipped)
	logRatio := ag.Log(ratio)

	prod := ag.Prod(target, logRatio)

	if reduceMean {
		return ag.ReduceMean(prod)
	}

	return ag.ReduceSum(prod)
}

// KLDivergenceSeq applies KLDivergence across a sequence of predictions and targets
func KLDivergenceSeq(predicted, target []mat.Tensor, reduceMean bool) mat.Tensor {
	if len(predicted) != len(target) {
		panic("losses: predicted and target slices must have the same length")
	}
	if len(predicted) == 0 {
		panic("losses: no tensors provided to RangeLossSeq")
	}

	loss := KLDivergence(predicted[0], target[0], reduceMean)
	for i := 1; i < len(predicted); i++ {
		loss = ag.Add(loss, KLDivergence(predicted[i], target[i], reduceMean))
	}

	return loss
}

// SaveLossPlot draws a simple line chart for the given slice of losses and saves it as filename.
func PlotLoss(losses []float64, filename string) error {
	const (
		width  = 800
		height = 600
		margin = 40
	)

	dc := gg.NewContext(width, height)
	dc.SetColor(color.White)
	dc.Clear()

	minmax := func(vals []float64) (float64, float64) {
		if len(vals) == 0 {
			return 0.0, 1.0
		}
		minV, maxV := vals[0], vals[0]
		for _, v := range vals {
			if v < minV {
				minV = v
			}
			if v > maxV {
				maxV = v
			}
		}
		return minV, maxV
	}

	// 1) Find min/max to scale the plot
	minVal, maxVal := minmax(losses)
	// If your losses are large, you might clamp to avoid a distorted scale (optional).
	// For example:
	//   if maxVal > 100.0 { maxVal = 100.0 } // just an example threshold

	chartW := float64(width - 2*margin)
	chartH := float64(height - 2*margin)

	// 2) Draw axes
	dc.SetColor(color.Black)
	dc.SetLineWidth(2)
	// X-axis
	dc.DrawLine(float64(margin), float64(height-margin), float64(width-margin), float64(height-margin))
	dc.Stroke()
	// Y-axis
	dc.DrawLine(float64(margin), float64(margin), float64(margin), float64(height-margin))
	dc.Stroke()

	// Optional: Add axis labels, ticks, etc.
	dc.SetColor(color.Black)
	dc.DrawStringAnchored("Loss", float64(margin)/2, float64(height)/2, 0.5, 0.5)
	dc.DrawStringAnchored("Iterations", float64(width)/2, float64(height-margin)/1.02, 0.5, 0)

	// 3) Plot the loss line
	dc.SetLineWidth(2)
	dc.SetColor(color.RGBA{R: 255, G: 0, B: 0, A: 255}) // red line

	// yForLoss maps a loss value to the vertical coordinate on the plot
	yForLoss := func(loss float64, margin int, chartH float64, minVal, maxVal float64) float64 {
		// 0% (lowest loss) at bottom, 100% (highest loss) at top
		normalized := (loss - minVal) / (maxVal - minVal + 1e-9)
		// Reverse so higher loss is at the top of chart
		y := float64(margin) + (1.0-normalized)*chartH
		return y
	}

	// Move to the first point
	if len(losses) > 0 {
		dc.MoveTo(
			float64(margin),
			yForLoss(losses[0], margin, chartH, minVal, maxVal),
		)
	}

	for i, val := range losses {
		// Map iteration i to chart X
		x := float64(margin) + float64(i)*(chartW/float64(len(losses)-1))
		// Map val to chart Y (inverted, so bigger loss is higher on the plot)
		y := yForLoss(val, margin, chartH, minVal, maxVal)
		dc.LineTo(x, y)
	}
	dc.Stroke()

	// 4) Save the final image
	return dc.SavePNG(filename)
}

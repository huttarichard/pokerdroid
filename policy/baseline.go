package policy

import "github.com/pokerdroid/poker/float/f64"

type BaselineUpdater func(p *Policy, weight float64, action int, value float64) float64

// UpdateBaseline with an EMA approach that includes sampling weight.
// Essentially: baseline += bf * w * (value - baseline)
func BaselineEMA(bf float64) BaselineUpdater {
	return func(p *Policy, w float64, action int, value float64) float64 {
		wDelta := w * (value - p.Baseline[action])
		step := bf * wDelta
		p.Baseline[action] += step
		return p.Baseline[action]
	}
}

// UpdateBaseline uses an iteration-based scale so bf can ramp up over time.
// This avoids large steps early when sampling weights can be extreme.
// UpdateBaseline with soft-clamping on step size
func BaselineEMAClamp(bf float64, clamp float64) BaselineUpdater {
	return func(p *Policy, w float64, action int, value float64) float64 {
		wDelta := w * (value - p.Baseline[action])
		step := bf * wDelta
		// Soft clamp: limit step magnitude so it won't blow up
		step = f64.ClampS(step, -clamp, clamp)
		p.Baseline[action] += step
		return p.Baseline[action]
	}
}

// UpdateBaseline with a simple iteration-based approach that also includes weight.
func BaselineWithIteration() BaselineUpdater {
	return func(p *Policy, w float64, action int, value float64) float64 {
		iteration := float64(p.Iteration) + 1
		p.Baseline[action] = (p.Baseline[action]*(iteration-1) + w*value) / iteration
		return p.Baseline[action]
	}
}

// BaselineAscended returns a BaselineUpdater that scales its update factor based on the policy iteration.
// It uses the parameters min, max, and eps (here representing an iteration threshold) to control how aggressive the baseline update is over time:
// at iteration 0, the learning rate is min, and after eps iterations, the learning rate reaches max.
// The update performed is: baseline += lr * w * (value - baseline).
func BaselineAscended(min, max, eps float64) BaselineUpdater {
	return func(p *Policy, w float64, action int, value float64) float64 {
		// Compute a linear scaling factor that increases from 0 to 1 as iterations go from 0 to eps
		iter := float64(p.Iteration)
		rate := iter / eps
		if rate > 1 {
			rate = 1
		}
		lr := min + (max-min)*rate // Learning rate transitions linearly from min to max until eps iterations

		wDelta := w * (value - p.Baseline[action])
		step := lr * wDelta
		p.Baseline[action] += step
		return p.Baseline[action]
	}
}

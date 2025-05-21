package abs

import "math"

type Histogram struct {
	Bins   []float32
	Equity float32
}

func (h Histogram) Normalize() Histogram {
	var total float32
	for i := 0; i < len(h.Bins); i++ {
		total += h.Bins[i]
	}
	if total == 0 {
		return h
	}
	for i := 0; i < len(h.Bins); i++ {
		h.Bins[i] = h.Bins[i] / total
	}
	return h
}

func (h Histogram) Increment(f float32) Histogram {
	idx := int(math.Floor(float64(f) * float64(len(h.Bins))))
	if idx >= len(h.Bins) {
		idx = len(h.Bins) - 1
	}
	h.Bins[idx]++
	h.Equity += f
	return h
}

func (h Histogram) Add(other Histogram) Histogram {
	for i := 0; i < len(h.Bins); i++ {
		h.Bins[i] += other.Bins[i]
	}
	h.Equity += other.Equity
	return h
}

func (h Histogram) Div(num float32) Histogram {
	for i := 0; i < len(h.Bins); i++ {
		h.Bins[i] /= num
	}
	h.Equity /= num
	return h
}

func (h Histogram) EMD(other Histogram) float64 {
	// If the number of bins differ, treat as maximum distance
	if len(h.Bins) != len(other.Bins) {
		return 1.0
	}

	// 1. Compute total counts
	var hTotal, oTotal float32
	for i := 0; i < len(h.Bins); i++ {
		hTotal += h.Bins[i]
		oTotal += other.Bins[i]
	}

	// 2. If either histogram is empty (total == 0), EMD is max
	if hTotal == 0 || oTotal == 0 {
		return 1.0
	}

	// 3. In a single pass, accumulate CDF and sum absolute differences
	var hCdf, oCdf float32
	var sum float64

	for i := 0; i < len(h.Bins); i++ {
		hCdf += h.Bins[i] / hTotal
		oCdf += other.Bins[i] / oTotal
		sum += math.Abs(float64(hCdf - oCdf))
	}

	// 4. Normalize by the number of bins => EMD in [0,1]
	return sum / float64(len(h.Bins))
}

func (h Histogram) Distance(other Histogram) float64 {
	a1 := float64(h.Equity - other.Equity)
	a2 := h.EMD(other)
	return math.Sqrt(a1*a1 + a2*a2)
}

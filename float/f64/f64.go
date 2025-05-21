package f64

import (
	"math"
	"sync/atomic"
	"unsafe"
)

// Uniform distribution
func Uniform(n int) []float64 {
	result := make([]float64, n)
	p := 1.0 / float64(n)
	AddConst(p, result)
	return result
}

// Min returns min of float64
func Min(x ...float64) float64 {
	if len(x) == 0 {
		panic("need at least one number")
	}
	f := x[0]
	for _, z := range x {
		if z < f {
			f = z
		}
	}
	return f
}

func ArgMax(x []float64) int {
	if len(x) == 0 {
		panic("need at least one number")
	}

	idx := 0
	best := x[0]
	for i, z := range x {
		if z > best {
			best = z
			idx = i
		}
	}
	return idx
}

// Max returns max of float664
func Max(x ...float64) float64 {
	if len(x) == 0 {
		panic("need at least one number")
	}
	f := x[0]
	for _, z := range x {
		if z > f {
			f = z
		}
	}
	return f
}

func Clamp(x []float64, min, max float64) {
	for i := range x {
		if x[i] < min {
			x[i] = min
		} else if x[i] > max {
			x[i] = max
		}
	}
}

func ClampS(x float64, min, max float64) float64 {
	if x < min {
		return min
	} else if x > max {
		return max
	}
	return x
}

func IsNanInf(x []float64) bool {
	for i := range x {
		if math.IsNaN(x[i]) || math.IsInf(x[i], 0) {
			return true
		}
	}
	return false
}

func Avg(x []float64) float64 {
	return Sum(x) / float64(len(x))
}

func AtomicAdd(val *float64, delta float64) (new float64) {
	for {
		old := *val
		new = old + delta
		if atomic.CompareAndSwapUint64(
			(*uint64)(unsafe.Pointer(val)),
			math.Float64bits(old),
			math.Float64bits(new),
		) {
			break
		}
	}
	return
}

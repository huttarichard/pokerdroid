package f32

import (
	"math"
	"sync/atomic"
	"unsafe"
)

// Uniform distribution
func Uniform(n int) []float32 {
	result := make([]float32, n)
	p := 1.0 / float32(n)
	AddConst(p, result)
	return result
}

// Min returns min of float32
func Min(x ...float32) float32 {
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

func ArgMax(x []float32) int {
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
func Max(x ...float32) float32 {
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

func Clamp(x []float32, min, max float32) {
	for i := range x {
		if x[i] < min {
			x[i] = min
		} else if x[i] > max {
			x[i] = max
		}
	}
}

func ClampS(x float32, min, max float32) float32 {
	if x < min {
		return min
	} else if x > max {
		return max
	}
	return x
}

func IsNanInf(x []float32) bool {
	for i := range x {
		if x[i] != x[i] || (x[i] > math.MaxFloat32 || x[i] < -math.MaxFloat32) {
			return true
		}
	}
	return false
}

func Avg(x []float32) float32 {
	return Sum(x) / float32(len(x))
}

func AddFloat32(val *float32, delta float32) (new float32) {
	for {
		old := *val
		new = old + delta
		if atomic.CompareAndSwapUint32(
			(*uint32)(unsafe.Pointer(val)),
			math.Float32bits(old),
			math.Float32bits(new),
		) {
			break
		}
	}
	return
}

// ScalUnitaryTo is multiply alpha to float32
func ScalUnitaryTo(dst []float32, alpha float32, x []float32) {
	for i, v := range x {
		dst[i] = alpha * v
	}
}

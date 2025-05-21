package f64

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPreventZero(t *testing.T) {
	tests := []struct {
		name string
		x    []float64
		tol  float64
		want []float64
	}{
		{
			name: "empty slice",
			x:    []float64{},
			tol:  1e-20,
			want: []float64{},
		},
		{
			name: "single zero",
			x:    []float64{0},
			tol:  1e-20,
			want: []float64{1e-20},
		},
		{
			name: "multiple zeros",
			x:    []float64{0, 1, 0, 2, 0},
			tol:  1e-20,
			want: []float64{1e-20, 1, 1e-20, 2, 1e-20},
		},
		{
			name: "no zeros",
			x:    []float64{1, 2, 3, 4},
			tol:  1e-20,
			want: []float64{1, 2, 3, 4},
		},
		{
			name: "different tolerance",
			x:    []float64{0, 1, 0},
			tol:  1e-10,
			want: []float64{1e-10, 1, 1e-10},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := make([]float64, len(tt.x))
			copy(x, tt.x)
			PreventZero(x, tt.tol)
			require.Equal(t, tt.want, x)
		})
	}
}

func TestMakePositive(t *testing.T) {
	tests := []struct {
		name string
		x    []float64
		want []float64
	}{
		{
			name: "empty slice",
			x:    []float64{},
			want: []float64{},
		},
		{
			name: "all positive",
			x:    []float64{1, 2, 3, 4},
			want: []float64{1, 2, 3, 4},
		},
		{
			name: "all negative",
			x:    []float64{-1, -2, -3, -4},
			want: []float64{0, 0, 0, 0},
		},
		{
			name: "mixed values",
			x:    []float64{-1, 2, -3, 4, 0},
			want: []float64{0, 2, 0, 4, 0},
		},
		{
			name: "with zeros",
			x:    []float64{0, -1, 0, 1},
			want: []float64{0, 0, 0, 1},
		},
		{
			name: "large slice",
			x:    []float64{-1, -2, -3, -4, -5, -6, -7, -8, -9, -10},
			want: []float64{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := make([]float64, len(tt.x))
			copy(x, tt.x)
			MakePositive(x)
			require.Equal(t, tt.want, x)
		})
	}
}

func TestMinMax(t *testing.T) {
	require.Equal(t, Min(0, 1, 2), float64(0))
	require.Equal(t, Min(1, 1, 2), float64(1))
	require.Equal(t, Min(2, 1), float64(1))
	require.Equal(t, Min(2, -100), float64(-100))
	require.Equal(t, Min(0), float64(0))

	require.Equal(t, Max(0, 1), float64(1))
	require.Equal(t, Max(0, 1, 100), float64(100))
	require.Equal(t, Max(100, -100), float64(100))
	require.Equal(t, Max(0, 1, 100.56565), float64(100.56565))
}

func TestAxpyUnitary(t *testing.T) {
	tests := []struct {
		name  string
		alpha float64
		x     []float64
		y     []float64
		want  []float64
	}{
		{
			name:  "empty slices",
			alpha: 2.0,
			x:     []float64{},
			y:     []float64{},
			want:  []float64{},
		},
		{
			name:  "unit vectors",
			alpha: 2.0,
			x:     []float64{1},
			y:     []float64{1},
			want:  []float64{3}, // 1 + 2*1
		},
		{
			name:  "longer vectors",
			alpha: 2.0,
			x:     []float64{1, 2, 3, 4},
			y:     []float64{4, 5, 6, 7},
			want:  []float64{6, 9, 12, 15}, // y[i] + alpha*x[i]
		},
		{
			name:  "negative alpha",
			alpha: -2.0,
			x:     []float64{1, 2, 3, 4},
			y:     []float64{4, 5, 6, 7},
			want:  []float64{2, 1, 0, -1}, // y[i] - 2*x[i]
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			y := make([]float64, len(tt.y))
			copy(y, tt.y)
			AxpyUnitary(tt.alpha, tt.x, y)
			require.Equal(t, tt.want, y)
		})
	}
}

func TestDotUnitary(t *testing.T) {
	tests := []struct {
		name string
		x    []float64
		y    []float64
		want float64
	}{
		{
			name: "empty slices",
			x:    []float64{},
			y:    []float64{},
			want: 0,
		},
		{
			name: "unit vectors",
			x:    []float64{2},
			y:    []float64{3},
			want: 6,
		},
		{
			name: "longer vectors",
			x:    []float64{1, 2, 3, 4},
			y:    []float64{4, 5, 6, 7},
			want: 60, // 1*4 + 2*5 + 3*6 + 4*7
		},
		{
			name: "with negatives",
			x:    []float64{-1, 2, -3, 4},
			y:    []float64{4, -5, 6, -7},
			want: -60,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DotUnitary(tt.x, tt.y)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestScalUnitary(t *testing.T) {
	tests := []struct {
		name  string
		alpha float64
		x     []float64
		want  []float64
	}{
		{
			name:  "empty slice",
			alpha: 2.0,
			x:     []float64{},
			want:  []float64{},
		},
		{
			name:  "unit vector",
			alpha: 2.0,
			x:     []float64{1},
			want:  []float64{2},
		},
		{
			name:  "longer vector",
			alpha: 2.0,
			x:     []float64{1, 2, 3, 4},
			want:  []float64{2, 4, 6, 8},
		},
		{
			name:  "negative alpha",
			alpha: -2.0,
			x:     []float64{1, 2, 3, 4},
			want:  []float64{-2, -4, -6, -8},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := make([]float64, len(tt.x))
			copy(x, tt.x)
			ScalUnitary(tt.alpha, x)
			require.Equal(t, tt.want, x)
		})
	}
}

func TestAddConst(t *testing.T) {
	tests := []struct {
		name  string
		alpha float64
		x     []float64
		want  []float64
	}{
		{
			name:  "empty slice",
			alpha: 2.0,
			x:     []float64{},
			want:  []float64{},
		},
		{
			name:  "unit vector",
			alpha: 2.0,
			x:     []float64{1},
			want:  []float64{3},
		},
		{
			name:  "longer vector",
			alpha: 2.0,
			x:     []float64{1, 2, 3, 4},
			want:  []float64{3, 4, 5, 6},
		},
		{
			name:  "negative alpha",
			alpha: -2.0,
			x:     []float64{1, 2, 3, 4},
			want:  []float64{-1, 0, 1, 2},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := make([]float64, len(tt.x))
			copy(x, tt.x)
			AddConst(tt.alpha, x)
			require.Equal(t, tt.want, x)
		})
	}
}

func TestSum(t *testing.T) {
	tests := []struct {
		name string
		x    []float64
		want float64
	}{
		{
			name: "empty slice",
			x:    []float64{},
			want: 0,
		},
		{
			name: "unit vector",
			x:    []float64{1},
			want: 1,
		},
		{
			name: "longer vector",
			x:    []float64{1, 2, 3, 4},
			want: 10,
		},
		{
			name: "with negatives",
			x:    []float64{-1, 2, -3, 4},
			want: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Sum(tt.x)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestScalUnitaryToUP(t *testing.T) {
	tests := []struct {
		name      string
		x         []float64
		alphaUp   float64
		alphaDown float64
		want      []float64
	}{
		{
			name:      "empty slice",
			x:         []float64{},
			alphaUp:   2.0,
			alphaDown: 0.5,
			want:      []float64{},
		},
		{
			name:      "all positive",
			x:         []float64{1.0, 2.0, 3.0, 4.0},
			alphaUp:   2.0,
			alphaDown: 0.5,
			want:      []float64{2.0, 4.0, 6.0, 8.0},
		},
		{
			name:      "all negative",
			x:         []float64{-1.0, -2.0, -3.0, -4.0},
			alphaUp:   2.0,
			alphaDown: 0.5,
			want:      []float64{-0.5, -1.0, -1.5, -2.0},
		},
		{
			name:      "mixed values",
			x:         []float64{-1.0, 2.0, -3.0, 4.0},
			alphaUp:   2.0,
			alphaDown: 0.5,
			want:      []float64{-0.5, 4.0, -1.5, 8.0},
		},
		{
			name:      "with zeros",
			x:         []float64{0.0, 1.0, 0.0, -1.0},
			alphaUp:   2.0,
			alphaDown: 0.5,
			want:      []float64{0.0, 2.0, 0.0, -0.5},
		},
		{
			name:      "large slice",
			x:         make([]float64, 100),
			alphaUp:   2.0,
			alphaDown: 0.5,
			want:      make([]float64, 100),
		},
		{
			name:      "negative alphas",
			x:         []float64{1.0, -2.0, 3.0, -4.0},
			alphaUp:   -2.0,
			alphaDown: -0.5,
			want:      []float64{-2.0, 1.0, -6.0, 2.0},
		},
		{
			name:      "small values",
			x:         []float64{0.000001, -0.000001, 0.000002, -0.000002},
			alphaUp:   2.0,
			alphaDown: 0.5,
			want:      []float64{0.000002, -0.0000005, 0.000004, -0.000001},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dst := make([]float64, len(tt.x))
			ScalUnitaryToUP(dst, tt.alphaUp, tt.alphaDown, tt.x)
			require.Equal(t, tt.want, dst)
		})
	}
}

// Native implementations for benchmarking
func nativePreventZero(x []float64, tol float64) {
	for i := range x {
		if x[i] == 0 {
			x[i] = tol
		}
	}
}

func nativeMakePositive(x []float64) {
	for i := range x {
		if x[i] < 0 {
			x[i] = 0
		}
	}
}

func nativeAxpyUnitary(alpha float64, x, y []float64) {
	for i, v := range x {
		y[i] += alpha * v
	}
}

func nativeDotUnitary(x, y []float64) float64 {
	var sum float64
	for i, v := range x {
		sum += y[i] * v
	}
	return sum
}

func nativeScalUnitary(alpha float64, x []float64) {
	for i := range x {
		x[i] *= alpha
	}
}

func nativeScalUnitaryTo(dst []float64, alpha float64, x []float64) {
	for i, v := range x {
		dst[i] = alpha * v
	}
}

func nativeAddConst(alpha float64, x []float64) {
	for i := range x {
		x[i] += alpha
	}
}

func nativeSum(x []float64) float64 {
	var sum float64
	for _, v := range x {
		sum += v
	}
	return sum
}

// Add native implementation for testing
func nativeScalUnitaryToUP(dst []float64, alphaUp, alphaDown float64, x []float64) {
	for i, v := range x {
		if v > 0 {
			dst[i] = alphaUp * v
		} else if v < 0 {
			dst[i] = alphaDown * v
		} else {
			dst[i] = 0
		}
	}
}

// Compare Go vs ASM implementations
func TestGoVsASM(t *testing.T) {
	const n = 1000
	x := make([]float64, n)
	y := make([]float64, n)
	want := make([]float64, n)
	got := make([]float64, n)
	alpha := 2.0

	// Initialize test data with mix of values
	for i := range x {
		x[i] = float64(i - n/3)       // Mix of negative/positive
		y[i] = float64((i * 2) - n/2) // Different mix of negative/positive
	}

	// Test AxpyUnitary
	copy(want, y)
	nativeAxpyUnitary(alpha, x, want)
	copy(got, y)
	AxpyUnitary(alpha, x, got)
	require.Equal(t, want, got)

	// Test DotUnitary
	wantDot := nativeDotUnitary(x, y)
	gotDot := DotUnitary(x, y)
	require.Equal(t, wantDot, gotDot)

	// Test ScalUnitary
	copy(want, x)
	nativeScalUnitary(alpha, want)
	copy(got, x)
	ScalUnitary(alpha, got)
	require.Equal(t, want, got)

	// Test ScalUnitaryTo
	dst1 := make([]float64, n)
	dst2 := make([]float64, n)
	nativeScalUnitaryTo(dst1, alpha, x)
	ScalUnitaryTo(dst2, alpha, x)
	require.Equal(t, dst1, dst2)

	// Test ScalUnitaryToUP
	dst3 := make([]float64, n)
	nativeScalUnitaryToUP(dst3, 2.0, 0.5, x)
	ScalUnitaryToUP(dst2, 2.0, 0.5, x)
	require.Equal(t, dst3, dst2)

	// Test AddConst
	copy(want, x)
	nativeAddConst(alpha, want)
	copy(got, x)
	AddConst(alpha, got)
	require.Equal(t, want, got)

	// Test Sum
	wantSum := nativeSum(x)
	gotSum := Sum(x)
	require.Equal(t, wantSum, gotSum)

	// Test PreventZero with mixed values including zeros
	for i := range x {
		if i%3 == 0 {
			x[i] = 0
		} else if i%3 == 1 {
			x[i] = float64(i)
		} else {
			x[i] = -float64(i)
		}
	}
	copy(want, x)
	nativePreventZero(want, alpha)
	copy(got, x)
	PreventZero(got, alpha)
	require.Equal(t, want, got)

	// Test MakePositive with full range of values
	for i := range x {
		switch i % 4 {
		case 0:
			x[i] = float64(i) // positive
		case 1:
			x[i] = -float64(i) // negative
		case 2:
			x[i] = 0 // zero
		case 3:
			x[i] = -0.000001 // small negative
		}
	}
	copy(want, x)
	nativeMakePositive(want)
	copy(got, x)
	MakePositive(got)
	require.Equal(t, want, got)
}

// Benchmark functions
func BenchmarkPreventZero(b *testing.B) {
	sizes := []int{100000}
	for _, size := range sizes {
		x := make([]float64, size)
		for i := range x {
			if i%3 == 0 {
				x[i] = 0
			} else {
				x[i] = float64(i)
			}
		}

		b.Run(fmt.Sprintf("ASM/%d", size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				xCopy := make([]float64, len(x))
				copy(xCopy, x)
				PreventZero(xCopy, 1e-10)
			}
		})

		b.Run(fmt.Sprintf("Native/%d", size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				xCopy := make([]float64, len(x))
				copy(xCopy, x)
				nativePreventZero(xCopy, 1e-10)
			}
		})
	}
}

func BenchmarkMakePositive(b *testing.B) {
	sizes := []int{100000}
	for _, size := range sizes {
		x := make([]float64, size)
		for i := range x {
			x[i] = float64(i - size/2)
		}

		b.Run(fmt.Sprintf("ASM/%d", size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				xCopy := make([]float64, len(x))
				copy(xCopy, x)
				MakePositive(xCopy)
			}
		})

		b.Run(fmt.Sprintf("Native/%d", size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				xCopy := make([]float64, len(x))
				copy(xCopy, x)
				nativeMakePositive(xCopy)
			}
		})
	}
}

func BenchmarkAxpyUnitary(b *testing.B) {
	sizes := []int{100000}
	for _, size := range sizes {
		x := make([]float64, size)
		y := make([]float64, size)
		for i := range x {
			x[i] = float64(i)
			y[i] = float64(i * 2)
		}

		b.Run(fmt.Sprintf("ASM/%d", size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				yCopy := make([]float64, len(y))
				copy(yCopy, y)
				AxpyUnitary(2.0, x, yCopy)
			}
		})

		b.Run(fmt.Sprintf("Native/%d", size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				yCopy := make([]float64, len(y))
				copy(yCopy, y)
				nativeAxpyUnitary(2.0, x, yCopy)
			}
		})
	}
}

func BenchmarkDotUnitary(b *testing.B) {
	sizes := []int{100000}
	for _, size := range sizes {
		x := make([]float64, size)
		y := make([]float64, size)
		for i := range x {
			x[i] = float64(i)
			y[i] = float64(i * 2)
		}

		b.Run(fmt.Sprintf("ASM/%d", size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				DotUnitary(x, y)
			}
		})

		b.Run(fmt.Sprintf("Native/%d", size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				nativeDotUnitary(x, y)
			}
		})
	}
}

func BenchmarkScalUnitary(b *testing.B) {
	sizes := []int{100000}
	for _, size := range sizes {
		x := make([]float64, size)
		for i := range x {
			x[i] = float64(i)
		}

		b.Run(fmt.Sprintf("ASM/%d", size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				xCopy := make([]float64, len(x))
				copy(xCopy, x)
				ScalUnitary(2.0, xCopy)
			}
		})

		b.Run(fmt.Sprintf("Native/%d", size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				xCopy := make([]float64, len(x))
				copy(xCopy, x)
				nativeScalUnitary(2.0, xCopy)
			}
		})
	}
}

func BenchmarkScalUnitaryTo(b *testing.B) {
	sizes := []int{100000}
	for _, size := range sizes {
		x := make([]float64, size)
		dst := make([]float64, size)
		for i := range x {
			x[i] = float64(i)
		}

		b.Run(fmt.Sprintf("ASM/%d", size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				ScalUnitaryTo(dst, 2.0, x)
			}
		})

		b.Run(fmt.Sprintf("Native/%d", size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				nativeScalUnitaryTo(dst, 2.0, x)
			}
		})
	}
}

func BenchmarkAddConst(b *testing.B) {
	sizes := []int{100000}
	for _, size := range sizes {
		x := make([]float64, size)
		for i := range x {
			x[i] = float64(i)
		}

		b.Run(fmt.Sprintf("ASM/%d", size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				xCopy := make([]float64, len(x))
				copy(xCopy, x)
				AddConst(2.0, xCopy)
			}
		})

		b.Run(fmt.Sprintf("Native/%d", size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				xCopy := make([]float64, len(x))
				copy(xCopy, x)
				nativeAddConst(2.0, xCopy)
			}
		})
	}
}

func BenchmarkSum(b *testing.B) {
	sizes := []int{100000}
	for _, size := range sizes {
		x := make([]float64, size)
		for i := range x {
			x[i] = float64(i)
		}

		b.Run(fmt.Sprintf("ASM/%d", size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				Sum(x)
			}
		})

		b.Run(fmt.Sprintf("Native/%d", size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				nativeSum(x)
			}
		})
	}
}

// Add benchmark
func BenchmarkScalUnitaryToUP(b *testing.B) {
	sizes := []int{100000}
	for _, size := range sizes {
		x := make([]float64, size)
		dst := make([]float64, size)
		for i := range x {
			if i%2 == 0 {
				x[i] = float64(i)
			} else {
				x[i] = -float64(i)
			}
		}

		b.Run(fmt.Sprintf("ASM/%d", size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				ScalUnitaryToUP(dst, 2.0, 0.5, x)
			}
		})

		b.Run(fmt.Sprintf("Native/%d", size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				nativeScalUnitaryToUP(dst, 2.0, 0.5, x)
			}
		})
	}
}

func TestScalUnitaryToUPEquivalence(t *testing.T) {
	// Local implementation of component functions for comparison
	scalUnitaryToGTZero := func(dst []float64, alpha float64, x []float64) {
		for i, v := range x {
			if v > 0 {
				dst[i] = alpha * v
			}
		}
	}

	scalUnitaryToLTZero := func(dst []float64, alpha float64, x []float64) {
		for i, v := range x {
			if v < 0 {
				dst[i] = alpha * v
			}
		}
	}

	tests := []struct {
		name      string
		x         []float64
		alphaUp   float64
		alphaDown float64
	}{
		{
			name:      "empty slice",
			x:         []float64{},
			alphaUp:   2.0,
			alphaDown: 0.5,
		},
		{
			name:      "mixed values",
			x:         []float64{-1.0, 2.0, 0.0, -3.0, 4.0, 0.0, -5.0},
			alphaUp:   2.0,
			alphaDown: 0.5,
		},
		{
			name:      "all positive",
			x:         []float64{1.0, 2.0, 3.0, 4.0},
			alphaUp:   2.0,
			alphaDown: 0.5,
		},
		{
			name:      "all negative",
			x:         []float64{-1.0, -2.0, -3.0, -4.0},
			alphaUp:   2.0,
			alphaDown: 0.5,
		},
		{
			name:      "all zeros",
			x:         []float64{0.0, 0.0, 0.0, 0.0},
			alphaUp:   2.0,
			alphaDown: 0.5,
		},
		{
			name:      "small values",
			x:         []float64{0.000001, -0.000001, 0.0, 0.000002, -0.000002},
			alphaUp:   2.0,
			alphaDown: 0.5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Result using ScalUnitaryToUP
			dstUP := make([]float64, len(tt.x))
			ScalUnitaryToUP(dstUP, tt.alphaUp, tt.alphaDown, tt.x)

			// Result using separate GT/LT operations
			dstSeparate := make([]float64, len(tt.x))
			scalUnitaryToGTZero(dstSeparate, tt.alphaUp, tt.x)
			scalUnitaryToLTZero(dstSeparate, tt.alphaDown, tt.x)

			// Compare results
			require.Equal(t, dstSeparate, dstUP,
				"Results differ for test case %s:\nExpected: %v\nGot: %v",
				tt.name, dstSeparate, dstUP)
		})
	}
}

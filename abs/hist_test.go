package abs

import (
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHist(t *testing.T) {
	t.Run("EMD", func(t *testing.T) {
		tests := []struct {
			name     string
			h1       Histogram
			h2       Histogram
			expected float64
		}{
			{
				name: "identical histograms",
				h1: Histogram{
					Bins:   []float32{1, 2, 3, 4},
					Equity: 0.5,
				},
				h2: Histogram{
					Bins:   []float32{1, 2, 3, 4},
					Equity: 0.5,
				},
				expected: 0.0,
			},
			{
				name: "completely different histograms",
				h1: Histogram{
					Bins:   []float32{10, 0, 0, 0},
					Equity: 0.5,
				},
				h2: Histogram{
					Bins:   []float32{0, 0, 0, 10},
					Equity: 0.5,
				},
				expected: 0.75, // Maximum EMD for 4 bins
			},
			{
				name: "partially different histograms",
				h1: Histogram{
					Bins:   []float32{5, 5, 0, 0},
					Equity: 0.5,
				},
				h2: Histogram{
					Bins:   []float32{0, 0, 5, 5},
					Equity: 0.5,
				},
				expected: 0.5, // Middle EMD value
			},
			{
				name: "different bin counts",
				h1: Histogram{
					Bins:   []float32{1, 2, 3},
					Equity: 0.5,
				},
				h2: Histogram{
					Bins:   []float32{1, 2, 3, 4},
					Equity: 0.5,
				},
				expected: 1.0, // Maximum distance for different bin counts
			},
			{
				name: "empty histogram",
				h1: Histogram{
					Bins:   []float32{0, 0, 0, 0},
					Equity: 0,
				},
				h2: Histogram{
					Bins:   []float32{1, 2, 3, 4},
					Equity: 0.5,
				},
				expected: 1.0, // Maximum distance for empty histogram
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result := tt.h1.EMD(tt.h2)
				if math.Abs(result-tt.expected) > 1e-6 {
					t.Errorf("EMD() = %v, want %v", result, tt.expected)
				}
			})
		}
	})

	t.Run("EMD symmetry", func(t *testing.T) {
		h1 := Histogram{
			Bins:   []float32{1, 3, 5, 7},
			Equity: 0.5,
		}
		h2 := Histogram{
			Bins:   []float32{7, 5, 3, 1},
			Equity: 0.5,
		}

		emd1 := h1.EMD(h2)
		emd2 := h2.EMD(h1)

		if math.Abs(emd1-emd2) > 1e-6 {
			t.Errorf("EMD should be symmetric: h1.EMD(h2)=%v, h2.EMD(h1)=%v", emd1, emd2)
		}
	})
}

func TestHistogram2(t *testing.T) {
	a := Histogram{
		Bins:   []float32{46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 44, 44, 42, 40, 40, 38, 38, 38, 38, 32, 18, 18, 4, 4, 1},
		Equity: 0.9008575677871704,
	}

	b := Histogram{
		Bins:   []float32{46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 43, 43, 43, 43, 43, 35, 23, 11, 11, 11, 11, 11, 11, 11, 9, 6, 6, 6, 6, 3, 3, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		Equity: 0.49706190824508667,
	}

	emd := a.EMD(b)
	require.Equal(t, 0.19133003931492568, emd)

	a2 := Histogram{
		Bins:   []float32{46, 46, 46, 46, 46, 46, 46, 42, 38, 38, 36, 26, 18, 18, 18, 14, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 2, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		Equity: 0.31800219416618347,
	}

	b2 := Histogram{
		Bins:   []float32{46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 40, 40, 40, 40, 37, 28, 28, 28, 22, 22, 22, 22, 14, 14, 14, 14, 14, 14, 14, 14, 14, 14, 14, 14, 14, 14, 14, 14, 11, 11, 8, 3, 0},
		Equity: 0.6221021413803101,
	}

	emd2 := a2.EMD(b2)
	require.Equal(t, 0.1386696982383728, emd2)
}

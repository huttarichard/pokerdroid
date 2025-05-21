package policy

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMarshalBinary(t *testing.T) {
	a := New(2)
	a.Iteration = 10

	a.Baseline = []float64{1, 2}
	a.StrategySum = []float64{5, 6}
	a.RegretSum = []float64{7, 8}

	a.BuildStrategy()

	data, err := a.MarshalBinary()
	require.NoError(t, err)

	px := &Policy{}
	err = px.UnmarshalBinary(data)
	require.NoError(t, err)

	require.Equal(t, a, px)
}

func TestNew(t *testing.T) {
	p := New(3)
	assert.Equal(t, uint64(0), p.Iteration)
	assert.Equal(t, float64(0.0), p.StrategyWeight)
	assert.Equal(t, []float64{1.0 / 3, 1.0 / 3, 1.0 / 3}, p.Strategy)
	assert.Equal(t, []float64{0, 0, 0}, p.RegretSum)
	assert.Equal(t, []float64{0, 0, 0}, p.StrategySum)
	assert.Equal(t, []float64{0, 0, 0}, p.Baseline)
}

func TestPolicy_AddRegrets(t *testing.T) {
	p := New(3)
	regrets := []float64{1.0, 2.0, 3.0}

	// Test with weight 1.0
	p.AddRegrets(1.0, regrets)
	assert.Equal(t, regrets, p.RegretSum)

	// Test with weight 0.5
	p.AddRegrets(0.5, regrets)
	assert.Equal(t, []float64{1.5, 3.0, 4.5}, p.RegretSum)
}

func TestPolicy_AddStrategyWeight(t *testing.T) {
	p := New(3)
	p.AddStrategyWeight(1.0)
	assert.Equal(t, float64(1.0), p.StrategyWeight)
	p.AddStrategyWeight(0.5)
	assert.Equal(t, float64(1.5), p.StrategyWeight)
}

func TestPolicy_GetAverageStrategy(t *testing.T) {
	p := New(2)

	// Test with zero sum
	avg := p.GetAverageStrategy()
	assert.Equal(t, []float64{0.5, 0.5}, avg)

	// Test with non-zero sum
	p.StrategySum = []float64{1.0, 3.0}
	avg = p.GetAverageStrategy()
	assert.Equal(t, []float64{0.25, 0.75}, avg)
}

func TestPolicy_Clone(t *testing.T) {
	p := New(2)
	p.Iteration = 1
	p.StrategyWeight = 1.0
	p.Strategy = []float64{0.3, 0.7}
	p.RegretSum = []float64{1.0, 2.0}
	p.StrategySum = []float64{3.0, 4.0}
	p.Baseline = []float64{0.1, 0.2}

	clone := p.Clone()

	assert.Equal(t, p.Iteration, clone.Iteration)
	assert.Equal(t, p.StrategyWeight, clone.StrategyWeight)
	assert.Equal(t, p.Strategy, clone.Strategy)
	assert.Equal(t, p.RegretSum, clone.RegretSum)
	assert.Equal(t, p.StrategySum, clone.StrategySum)
	assert.Equal(t, p.Baseline, clone.Baseline)

	// Verify deep copy
	clone.Strategy[0] = 0.9
	assert.NotEqual(t, p.Strategy[0], clone.Strategy[0])
}

func TestPolicy_MarshalUnmarshal(t *testing.T) {
	original := New(2)
	original.Iteration = 42
	original.RegretSum = []float64{1.0, 2.0}
	original.StrategySum = []float64{3.0, 4.0}
	original.Baseline = []float64{0.1, 0.2}

	original.BuildStrategy()

	// Marshal
	data, err := original.MarshalBinary()
	require.NoError(t, err)

	// Unmarshal into new policy
	restored := New(2)
	err = restored.UnmarshalBinary(data)
	require.NoError(t, err)

	// Compare all fields
	assert.Equal(t, original.Iteration, restored.Iteration)
	assert.Equal(t, original.Strategy, restored.Strategy)
	assert.Equal(t, original.RegretSum, restored.RegretSum)
	assert.Equal(t, original.StrategySum, restored.StrategySum)
	assert.Equal(t, original.Baseline, restored.Baseline)
}

func TestPolicy_Size(t *testing.T) {
	tests := []struct {
		name    string
		policy  *Policy
		actions int
	}{
		{
			name:    "empty policy",
			policy:  New(0),
			actions: 0,
		},
		{
			name:    "small policy",
			policy:  New(2),
			actions: 2,
		},
		{
			name:    "medium policy",
			policy:  New(5),
			actions: 5,
		},
		{
			name:    "large policy",
			policy:  New(10),
			actions: 10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.policy == nil {
				require.Equal(t, uint64(0), tt.policy.Size())
				return
			}

			expectedSize := uint64(12)                  // fixed fields (4 + 8 + 4)
			expectedSize += uint64(tt.actions) * 8 * 3 // 5 float32 slices

			reportedSize := tt.policy.Size()
			require.Equal(t, expectedSize, reportedSize, "reported size mismatch")

			// Verify against actual marshaled size
			data, err := tt.policy.MarshalBinary()
			require.NoError(t, err)
			require.Equal(t, reportedSize, uint64(len(data)),
				"marshaled size (%d) does not match reported size (%d)",
				len(data), reportedSize)

			// Verify round trip
			newPolicy := New(tt.actions)
			err = newPolicy.UnmarshalBinary(data)
			require.NoError(t, err)
			require.Equal(t, tt.policy.Iteration, newPolicy.Iteration)
			require.InDeltaSlice(t, tt.policy.Strategy, newPolicy.Strategy, 1e-6)
		})
	}
}

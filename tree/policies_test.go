package tree

import (
	"testing"

	"github.com/pokerdroid/poker/abs"
	"github.com/pokerdroid/poker/policy"
	"github.com/stretchr/testify/require"
)

func TestPolicyMap(t *testing.T) {
	t.Run("basic operations", func(t *testing.T) {
		pm := NewPolicies()

		// Test storing and retrieving
		cl := abs.Cluster(1)
		pol := policy.New(3)
		pm.Store(cl, pol)

		retrieved, ok := pm.Get(cl)
		require.True(t, ok)
		require.Equal(t, pol, retrieved)

		// Test non-existent cluster
		_, ok = pm.Get(abs.Cluster(2))
		require.False(t, ok)

		// Test count
		require.Equal(t, uint32(1), pm.Len())
	})

	t.Run("clone", func(t *testing.T) {
		pm := NewPolicies()

		// Add some policies
		cl1 := abs.Cluster(1)
		pol1 := policy.New(3)
		pol1.RegretSum = []float64{0.1, 0.2, 0.7}

		cl2 := abs.Cluster(2)
		pol2 := policy.New(3)
		pol2.RegretSum = []float64{0.3, 0.3, 0.4}

		pm.Store(cl1, pol1)
		pm.Store(cl2, pol2)

		// Clone and verify
		clone := pm.Clone()
		require.Equal(t, pm.Len(), clone.Len())

		// Verify independence
		pol1.RegretSum[0] = 0.9
		clonedPol, _ := clone.Get(cl1)
		require.NotEqual(t, pol1.RegretSum[0], clonedPol.RegretSum[0])
	})

	t.Run("marshal/unmarshal", func(t *testing.T) {
		pm := NewPolicies()

		// Add some policies
		for i := 0; i < 100; i++ {
			pol := policy.New(3)
			pol.RegretSum = []float64{0.1, 0.2, 0.7}
			pol.BuildStrategy()

			pm.Store(abs.Cluster(i), pol)
		}

		// Marshal
		data, err := pm.MarshalBinary()
		require.NoError(t, err)

		// Unmarshal
		pm2 := NewPolicies()
		err = pm2.UnmarshalBinary(data)
		require.NoError(t, err)

		// Verify
		require.Equal(t, pm.Len(), pm2.Len())
		require.Equal(t, pm.Map, pm2.Map)
	})
}

func TestPolicyMap_Size(t *testing.T) {
	tests := []struct {
		name     string
		setup    func() *Policies
		clusters int
		actions  int
	}{
		{
			name: "empty map",
			setup: func() *Policies {
				return NewPolicies()
			},
			clusters: 0,
			actions:  0,
		},
		{
			name: "single policy",
			setup: func() *Policies {
				pm := NewPolicies()
				pm.Store(abs.Cluster(1), policy.New(2))
				return pm
			},
			clusters: 1,
			actions:  2,
		},
		{
			name: "multiple policies",
			setup: func() *Policies {
				pm := NewPolicies()
				pm.Store(abs.Cluster(1), policy.New(2))
				pm.Store(abs.Cluster(2), policy.New(3))
				pm.Store(abs.Cluster(3), policy.New(4))
				return pm
			},
			clusters: 3,
			actions:  0, // varies by policy
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pm := tt.setup()

			reportedSize := pm.Size()
			data, err := pm.MarshalBinary()
			require.NoError(t, err)

			require.Equal(t, reportedSize, uint64(len(data)),
				"marshaled size (%d) does not match reported size (%d)",
				len(data), reportedSize)

			// Verify round trip
			newPm := NewPolicies()
			err = newPm.UnmarshalBinary(data)
			require.NoError(t, err)
			require.Equal(t, pm.Len(), newPm.Len())

			for cl, pol := range pm.Map {
				newPol, ok := newPm.Get(cl)
				require.True(t, ok)
				require.Equal(t, pol.Iteration, newPol.Iteration)
				require.Equal(t, pol.Strategy, newPol.Strategy)
			}
		})
	}
}

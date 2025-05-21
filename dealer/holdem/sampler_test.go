package holdemdealer

import (
	"testing"

	"github.com/pokerdroid/poker/card"
	"github.com/pokerdroid/poker/frand"
	"github.com/pokerdroid/poker/table"
	"github.com/stretchr/testify/require"
)

func TestSamples(t *testing.T) {
	params := SamplerParams{NumPlayers: 2}

	hnd := New(params)
	for i := 0; i < 1000; i++ {
		_, err := hnd.Sample(frand.NewUnsafeInt(0))
		require.NoError(t, err)
	}
}

// Test to verify pooled sampler maintains correctness
func TestPooledSamplerCorrectness(t *testing.T) {
	params := SamplerParams{NumPlayers: 2}

	hnd := New(params)

	r := frand.NewUnsafeInt(0)

	// Generate multiple samples and verify they're unique
	samples := make(map[string]struct{})
	for i := 0; i < 100; i++ {
		sample, err := hnd.Sample(r)
		require.NoError(t, err)
		gs := sample.(*Sample)

		// Convert cards to string for comparison
		key := gs.Cards(0).String() + gs.Cards(1).String()
		_, exists := samples[key]
		require.False(t, exists, "Generated duplicate sample")
		samples[key] = struct{}{}

		hnd.Put(sample)
	}
}

func TestSampleCopy(t *testing.T) {
	params := SamplerParams{NumPlayers: 2}

	hnd := New(params)

	r := frand.NewUnsafeInt(0)

	sample, err := hnd.Sample(r)
	require.NoError(t, err)

	copy, err := hnd.Copy(r, sample)
	require.NoError(t, err)

	require.Equal(t, sample, copy)
}

func TestSampleCopy2(t *testing.T) {
	params := SamplerParams{NumPlayers: 2}

	hnd := New(params)

	r := frand.NewUnsafeInt(0)

	sample, err := hnd.Sample(r)
	require.NoError(t, err)

	copy1, err := hnd.Copy(r, sample)
	require.NoError(t, err)

	copy2, err := hnd.Copy(r, sample)
	require.NoError(t, err)

	require.Equal(t, sample, copy1)
	require.Equal(t, sample, copy2)

	sample.(*Sample).Sample(table.River)
	sample.(*Sample).hands[0][0] = card.Card00

	require.Equal(t, copy1, copy2)
	require.NotEqual(t, sample, copy1)
	require.NotEqual(t, sample, copy2)
}

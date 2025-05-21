package holdemdealer

import (
	"math"
	"testing"

	"github.com/pokerdroid/poker/card"
	"github.com/pokerdroid/poker/frand"
	"github.com/pokerdroid/poker/table"
	"github.com/stretchr/testify/require"
)

func TestWeightedSamplerSample(t *testing.T) {
	// Initialize a reproducible random generator.
	rng := frand.NewUnsafeInt(42)

	// Define the number of players.
	numPlayers := uint8(2)

	// Create a uniform range distribution for each player.
	ranges := make([]card.RangeDist, numPlayers)
	for i := uint8(0); i < numPlayers; i++ {
		ranges[i] = card.NewUniformRangeDist()
	}

	// Create parameters for the weighted sampler.
	params := RangeParams{
		NumPlayers: numPlayers,
		Ranges:     ranges,
	}

	// Instantiate the weighted sampler.
	ws := NewWeighted(params)

	// Get a sample from the weighted sampler.
	sample, err := ws.Sample(rng)
	if err != nil {
		t.Fatalf("Sample() returned error: %v", err)
	}

	s, ok := sample.(*Sample)
	if !ok {
		t.Fatal("Sample is not of type *Sample")
	}

	// Each player's cards slice should contain exactly 7 cards:
	// 2 hole cards (sampled from range) and 5 community cards.
	for i := uint8(0); i < numPlayers; i++ {
		if len(s.Cards(i)) != 2 {
			t.Errorf("Player %d: expected 7 cards, got %d", i, len(s.Cards(i)))
		}
	}

	// Verify that the hole cards (first two cards) for different players do not overlap.
	hole0 := s.Cards(0)
	hole1 := s.Cards(1)
	if card.IsAnyMatch(hole0, hole1) {
		t.Errorf("Hole cards overlap: %v and %v", hole0, hole1)
	}

	t.Logf("Player 0 cards: %v", s.Cards(0))
	t.Logf("Player 1 cards: %v", s.Cards(1))

	// Return the sample to the pool.
	ws.Put(s)
}

// TestWeightedSamplerRangeDistribution verifies that the sampler uses the provided
// range distribution correctly. We use a custom RangeDist with nonzero probability
// only at two indices so that over many samples we obtain the expected frequency ratio.
func TestWeightedSamplerRangeDistribution(t *testing.T) {
	rng := frand.NewUnsafeInt(42)

	// Create a custom RangeDist that has nonzero values at only two indices.
	var customDist card.RangeDist
	// Zero-out the distribution.
	for i := range customDist {
		customDist[i] = 0.0
	}
	// Set only two indices with fixed probabilities.
	customDist[10] = 0.3
	customDist[20] = 0.7

	// For this test, use a single player to avoid card collision issues.
	numPlayers := uint8(1)

	ranges := make([]card.RangeDist, numPlayers)
	for i := uint8(0); i < numPlayers; i++ {
		ranges[i] = customDist
	}

	params := RangeParams{
		NumPlayers: numPlayers,
		Ranges:     ranges,
	}

	ws := NewWeighted(params)

	totalSamples := 10000
	freq := map[int]int{
		10: 0,
		20: 0,
	}

	// Sample many times to assess the frequency of each range index.
	for i := 0; i < totalSamples; i++ {
		sample, err := ws.Sample(rng)
		if err != nil {
			t.Fatalf("Sample() error: %v", err)
		}
		s := sample.(*Sample)
		// The hole cards are the first two cards.
		hole := s.Cards(0)
		idx := card.RangeIndex(hole)
		if idx != 10 && idx != 20 {
			t.Errorf("Unexpected range index: got %d, expected 10 or 20", idx)
		}
		freq[idx]++
		ws.Put(s)
	}

	ratio10 := float64(freq[10]) / float64(totalSamples)
	ratio20 := float64(freq[20]) / float64(totalSamples)

	t.Logf("Frequency of index10: %d (%.3f), index20: %d (%.3f)", freq[10], ratio10, freq[20], ratio20)

	// Allow a tolerance of 5%.
	if math.Abs(ratio10-0.3) > 0.05 {
		t.Errorf("Index 10 frequency out of tolerance: got %.3f, expected ~0.3", ratio10)
	}

	if math.Abs(ratio20-0.7) > 0.05 {
		t.Errorf("Index 20 frequency out of tolerance: got %.3f, expected ~0.7", ratio20)
	}
}

func TestWSampleCopy(t *testing.T) {
	// For this test, use a single player to avoid card collision issues.
	numPlayers := uint8(1)

	ranges := make([]card.RangeDist, numPlayers)
	for i := uint8(0); i < numPlayers; i++ {
		ranges[i] = card.NewUniformRangeDist()
	}

	params := RangeParams{
		NumPlayers: numPlayers,
		Ranges:     ranges,
	}

	hnd := NewWeighted(params)

	r := frand.NewUnsafeInt(0)

	sample, err := hnd.Sample(r)
	require.NoError(t, err)

	copy, err := hnd.Copy(r, sample)
	require.NoError(t, err)

	require.Equal(t, sample, copy)
}

func TestWSampleCopy2(t *testing.T) {
	// For this test, use a single player to avoid card collision issues.
	numPlayers := uint8(1)

	ranges := make([]card.RangeDist, numPlayers)
	for i := uint8(0); i < numPlayers; i++ {
		ranges[i] = card.NewUniformRangeDist()
	}

	params := RangeParams{
		NumPlayers: numPlayers,
		Ranges:     ranges,
	}

	hnd := NewWeighted(params)

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

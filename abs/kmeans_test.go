package abs

import (
	"math"
	"testing"

	"github.com/pokerdroid/poker"
	"github.com/pokerdroid/poker/frand"
)

// TestKmeansGroups verifies that KMeans groups effectively cluster the data.
// We generate two clusters of float64 numbers:
// - One around 0.0 (expected average ~0.45)
// - One around 10.0 (expected average ~10.45)
func TestKmeansGroups(t *testing.T) {
	var data []float64

	// Cluster 1: points around 0.0
	for i := 0; i < 10; i++ {
		data = append(data, float64(i)*0.1)
	}

	// Cluster 2: points around 10.0
	for i := 0; i < 10; i++ {
		data = append(data, 10.0+float64(i)*0.1)
	}

	// Use a deterministic RNG for repeatability.
	rng := frand.NewUnsafeInt(42)

	k := 2
	// Initialize groups randomly from the data.
	initGroups, err := NewRandomGroups(rng, data, k)
	if err != nil {
		t.Fatalf("failed to create random groups: %v", err)
	}

	// Define a distance function for float64: simple absolute difference.
	distance := func(x, y float64) float64 {
		return math.Abs(x - y)
	}

	// Define a recenter function that computes the mean of provided points.
	recenter := func(rng frand.Rand, pts []float64) float64 {
		sum := 0.0
		for _, v := range pts {
			sum += v
		}
		return sum / float64(len(pts))
	}

	// Configure KMeans options.
	opts := KMeansOpts[float64]{
		Clusters:      k,
		MaxIterations: 100,
		MaxDelta:      0.001,
		Logger:        poker.VoidLogger{},
		LogIteration:  0,
		Rng:           rng,
		Groups:        initGroups,
		Recenter:      recenter,
		Distance:      distance,
	}

	// Run KMeans clustering.
	groups, score, err := KMeans(data, opts)
	if err != nil {
		t.Fatalf("KMeans failed: %v", err)
	}

	t.Logf("KMeans score: %f", score)
	if len(groups) != k {
		t.Fatalf("expected %d groups, got %d", k, len(groups))
	}

	// Expect one group center to be near the average of the first cluster (~0.45)
	// and the other group center near the average of the second cluster (~10.45).
	center1 := groups[0].Center
	center2 := groups[1].Center
	t.Logf("Cluster centers: %f, %f", center1, center2)

	const epsilon = 0.5 // tolerance threshold
	cond1 := math.Abs(center1-0.45) < epsilon && math.Abs(center2-10.45) < epsilon
	cond2 := math.Abs(center2-0.45) < epsilon && math.Abs(center1-10.45) < epsilon

	if !cond1 && !cond2 {
		t.Fatalf("cluster centers not as expected: got %f and %f", center1, center2)
	}
}

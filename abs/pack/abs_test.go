package absp

import (
	"math"
	"sort"
	"testing"

	"github.com/pokerdroid/poker/abs"
	"github.com/pokerdroid/poker/abs/river"
	"github.com/pokerdroid/poker/card"
	"github.com/pokerdroid/poker/equity/omp"
	"github.com/pokerdroid/poker/frand"
	"github.com/pokerdroid/poker/table"
	"github.com/stretchr/testify/require"
)

func TestMeasureQuality(t *testing.T) {
	absx := TstNewAbs(t)
	r := frand.NewHash()

	var all float64

	abs.IndexWorkers(10000, func(i int, done uint64) error {
		n := 5 + r.Intn(3)
		var street table.Street
		switch n {
		case 5:
			street = table.Flop
		case 6:
			street = table.Turn
		case 7:
			street = table.River
		}
		c := card.RandomCards(r, n)
		if len(c) > 7 {
			panic(c)
		}
		eq1 := omp.Equity(c[0], c[1], c[2:], 2)
		eq2 := absx.Equity(street, absx.Map(c))

		all += math.Abs(float64(eq1.WinDraw() - eq2.WinDraw()))
		return nil
	})

	t.Logf("avg dist: %f", all/10000)

	require.Less(t, all/10000, 0.02)
}

func TestMeasureQualityRiver(t *testing.T) {
	absx := TstNewAbs(t)
	r := frand.NewHash()

	var dist float64

	abs.IndexWorkers(10000, func(i int, done uint64) error {
		c := card.RandomCards(r, 7)
		eq1 := river.ComputeEquity(c)
		eq2 := absx.Equity(table.River, absx.Map(c))

		dist += math.Abs(float64(eq1.WinDraw() - eq2.WinDraw()))
		return nil
	})

	t.Logf("avg dist: %f", dist/10000)
	require.Less(t, dist/10000, 0.005)
}

func TestStrategicEquivalence(t *testing.T) {
	abs := TstNewAbs(t)

	// Create pairs of hands that should play similarly
	testCases := []struct {
		name   string
		hand1  string
		hand2  string
		board  string
		street table.Street
	}{
		{"suited aces", "ah kh", "as ks", "2c 3d 4h", table.Flop},
		{"pocket pairs", "qh qd", "qs qc", "kh 7d 2s", table.Flop},
		{"nut flush", "ah 5h", "ah 7h", "2h 4h 9h", table.Flop},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			hand1 := card.NewCardsFromString(tc.hand1)
			hand2 := card.NewCardsFromString(tc.hand2)
			board := card.NewCardsFromString(tc.board)

			cards1 := append(hand1, board...)
			cards2 := append(hand2, board...)

			cluster1 := abs.Map(cards1)
			cluster2 := abs.Map(cards2)

			eq1 := abs.Equity(tc.street, cluster1)
			eq2 := abs.Equity(tc.street, cluster2)

			t.Logf("%s: Clusters %d vs %d, Equities %v vs %v",
				tc.name, cluster1, cluster2, eq1, eq2)

			// Check if strategically similar hands map to same/similar clusters
			if cluster1 != cluster2 {
				t.Logf("Warning: Similar hands mapped to different clusters")
			}
		})
	}
}

func TestClusterDistribution(t *testing.T) {
	abs := TstNewAbs(t)
	r := frand.NewHash()

	// Count hands per cluster for each street
	flopClusters := make(map[uint32]int)
	turnClusters := make(map[uint32]int)
	riverClusters := make(map[uint32]int)

	samples := 10000

	for i := 0; i < samples; i++ {
		// Flop
		cards := card.RandomCards(r, 5)
		cluster := abs.Map(cards)
		flopClusters[uint32(cluster)]++

		// Turn
		cards = card.RandomCards(r, 6)
		cluster = abs.Map(cards)
		turnClusters[uint32(cluster)]++

		// River
		cards = card.RandomCards(r, 7)
		cluster = abs.Map(cards)
		riverClusters[uint32(cluster)]++
	}

	t.Logf("Unique clusters - Flop: %d, Turn: %d, River: %d",
		len(flopClusters), len(turnClusters), len(riverClusters))

	// Calculate distribution statistics
	analyzeDistribution(t, "Flop", flopClusters)
	analyzeDistribution(t, "Turn", turnClusters)
	analyzeDistribution(t, "River", riverClusters)
}

func analyzeDistribution(t *testing.T, name string, clusters map[uint32]int) {
	var counts []int
	for _, count := range clusters {
		counts = append(counts, count)
	}

	// Sort counts
	sort.Ints(counts)

	// Calculate statistics
	var sum int
	for _, c := range counts {
		sum += c
	}

	avg := float64(sum) / float64(len(counts))
	max := counts[len(counts)-1]
	min := counts[0]

	t.Logf("%s clusters - Count: %d, Avg: %.2f, Min: %d, Max: %d, Ratio: %.2f",
		name, len(clusters), avg, min, max, float64(max)/float64(min))
}

func TestCardsClusters(t *testing.T) {
	absx := TstNewAbs(t)
	r := frand.NewHash()
	hand := card.RandomCards(r, 2)

	m := make(map[abs.Cluster][]card.Cards)

	for i := 0; i < 10000; i++ {
		bb := card.RandomCards(r, 3)
		if card.IsAnyMatch(bb, hand) {
			continue
		}
		cluster := absx.Map(append(hand, bb...))
		// t.Logf("cluster: %s %d", bb.String(), cluster)

		m[cluster] = append(m[cluster], bb)
	}

	for k, v := range m {
		t.Logf("cluster: %d", k)
		for _, x := range v {
			eq := omp.Equity(hand[0], hand[1], x, 2)
			t.Logf("  %s: %f", card.Cards(append(hand, x...)).String(), eq.WinDraw())
		}
	}
}

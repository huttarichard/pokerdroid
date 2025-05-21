package river

import (
	"bytes"
	"os"
	"sort"
	"testing"

	"github.com/pokerdroid/poker"
	"github.com/pokerdroid/poker/abs"
	"github.com/pokerdroid/poker/encbin"
	"github.com/pokerdroid/poker/equity"
	"github.com/pokerdroid/poker/frand"
)

type Abs struct {
	Equities []equity.Equity
	Clusters [Size]abs.Cluster
}

func NewFromFile(path string) (*Abs, error) {
	abs := new(Abs)

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	err = abs.UnmarshalBinary(data)
	if err != nil {
		return nil, err
	}
	return abs, nil
}

func TstNewAbs(t *testing.T) *Abs {
	abs, err := NewFromFile(os.Getenv("RIVER_ABS_PATH"))
	if err != nil {
		t.Skipf("skipping abs test: %s", err)
	}
	return abs
}

func (a *Abs) Map(cluster uint32) abs.Cluster {
	return a.Clusters[cluster]
}

func (a *Abs) MarshalBinary() ([]byte, error) {
	buf := bytes.NewBuffer(nil)

	// Write bins
	err := encbin.MarshalValues(buf, a.Clusters)
	if err != nil {
		return nil, err
	}

	// Write equity slice
	err = encbin.MarshalSliceLen[equity.Equity, uint32](buf, a.Equities)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (a *Abs) UnmarshalBinary(data []byte) error {
	r := bytes.NewReader(data)

	// Read bins
	err := encbin.UnmarshalValues(r, &a.Clusters)
	if err != nil {
		return err
	}

	// Read equity slice
	a.Equities, err = encbin.UnmarhsalSliceLen[equity.Equity, uint32](r)
	if err != nil {
		return err
	}

	return nil
}

type PartitionOpts struct {
	Clusters       int
	MaxIterations  int
	DeltaThreshold float64
	Logger         poker.Logger
	LogIteration   int
	Rng            frand.Rand
}

// Partition executes the k-means algorithm on the given set of equities and
// partitions it into k clusters.
func Partition(buckets *Buckets, opts PartitionOpts) (*Abs, error) {
	var counter abs.Cluster

	mapping := make(map[equity.Equity]struct{})
	equities := make([]equity.Equity, 0)

	for _, eq := range buckets.db {
		_, ok := mapping[eq]
		if ok {
			continue
		}
		equities = append(equities, eq)
		mapping[eq] = struct{}{}
		counter++
	}

	groups, err := abs.NewRandomGroups(opts.Rng, equities, opts.Clusters)
	if err != nil {
		return nil, err
	}

	gg, _, err := abs.KMeans(equities, abs.KMeansOpts[equity.Equity]{
		Clusters:      opts.Clusters,
		MaxIterations: opts.MaxIterations,
		MaxDelta:      opts.DeltaThreshold,
		Logger:        opts.Logger,
		LogIteration:  opts.LogIteration,
		Rng:           opts.Rng,
		Groups:        groups,

		Recenter: func(rng frand.Rand, e []equity.Equity) equity.Equity {
			var divider float32
			var v1, v2 float32

			for _, c := range e {
				divider++
				v1 += c.Win()
				v2 += c.Tie()
			}

			v1 = v1 / divider
			v2 = v2 / divider

			return equity.NewEquity(v1, v2)
		},

		Distance: func(a, b equity.Equity) float64 {
			return a.Distance(b)
		},
	})
	if err != nil {
		return nil, err
	}

	sort.Slice(gg, func(i, j int) bool {
		return gg[i].Center.WinDraw() < gg[j].Center.WinDraw()
	})

	absx := new(Abs)
	absx.Equities = make([]equity.Equity, 0)

	for _, g := range gg {
		absx.Equities = append(absx.Equities, g.Center)
	}

	err = abs.IndexWorkers(Size, func(c int, done uint64) error {
		indx := gg.Nearest(buckets.db[c], func(a, b equity.Equity) float64 {
			return a.Distance(b)
		})

		absx.Clusters[c] = abs.Cluster(indx)

		if done%uint64(opts.LogIteration) == 0 {
			opts.Logger.Printf("building mapping: %d/%d", done, Size)
		}
		return nil
	})

	return absx, err
}

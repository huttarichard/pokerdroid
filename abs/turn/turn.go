package turn

import (
	"bytes"
	"sort"

	"github.com/pokerdroid/poker"
	"github.com/pokerdroid/poker/abs"
	"github.com/pokerdroid/poker/abs/river"
	"github.com/pokerdroid/poker/card"
	"github.com/pokerdroid/poker/encbin"
	"github.com/pokerdroid/poker/frand"
	"github.com/pokerdroid/poker/iso"
)

const Size = 13_960_050

type Abs struct {
	Bins     uint8
	Equity   []float32
	Clusters [Size]abs.Cluster
}

func (a *Abs) Map(cluster uint32) abs.Cluster {
	return a.Clusters[cluster]
}

func (a *Abs) MarshalBinary() ([]byte, error) {
	buf := bytes.NewBuffer(nil)

	// Write bins
	err := encbin.MarshalValues(buf, a.Bins, a.Clusters)
	if err != nil {
		return nil, err
	}

	// Write equity slice
	err = encbin.MarshalSliceLen[float32, uint32](buf, a.Equity)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (a *Abs) UnmarshalBinary(data []byte) error {
	r := bytes.NewReader(data)

	// Read bins
	err := encbin.UnmarshalValues(r, &a.Bins, &a.Clusters)
	if err != nil {
		return err
	}

	// Read equity slice
	a.Equity, err = encbin.UnmarhsalSliceLen[float32, uint32](r)
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
	Bins           int
}

// Partition executes the k-means algorithm on the given set of equities and
// partitions it into k clusters.
func Partition(hh []abs.Histogram, opts PartitionOpts) (*Abs, error) {
	rng := frand.Clone(opts.Rng)

	groups, err := abs.NewRandomGroups(rng, hh, opts.Clusters)
	if err != nil {
		return nil, err
	}

	data, score, err := abs.KMeans(hh, abs.KMeansOpts[abs.Histogram]{
		Clusters:      opts.Clusters,
		MaxIterations: opts.MaxIterations,
		MaxDelta:      opts.DeltaThreshold,
		Logger:        opts.Logger,
		LogIteration:  opts.LogIteration,
		Rng:           rng,
		Groups:        groups,

		Recenter: func(rng frand.Rand, e []abs.Histogram) abs.Histogram {
			center := abs.Histogram{
				Bins:   make([]float32, opts.Bins),
				Equity: 0,
			}
			for _, h := range e {
				center = center.Add(h)
			}
			return center.Div(float32(len(e)))
		},

		Distance: func(a, b abs.Histogram) float64 {
			return a.Distance(b)
		},
	})
	if err != nil {
		return nil, err
	}

	opts.Logger.Printf("score: %v", score)

	sort.Slice(data, func(i, j int) bool {
		return data[i].Center.Equity < data[j].Center.Equity
	})

	absx := new(Abs)
	absx.Equity = make([]float32, 0)
	absx.Bins = uint8(opts.Bins)

	for _, h := range data {
		absx.Equity = append(absx.Equity, h.Center.Equity)
	}

	err = abs.IndexWorkers(Size, func(c int, done uint64) error {
		indx := data.Nearest(hh[c], func(a, b abs.Histogram) float64 {
			return a.EMD(b)
		})

		absx.Clusters[c] = abs.Cluster(indx)
		return nil
	})

	return absx, err
}

type ComputeOpts struct {
	Buckets      *river.Buckets
	Logger       poker.Logger
	Bins         int
	LogIteration int
}

func Compute(opts ComputeOpts) ([]abs.Histogram, error) {
	hh := make([]abs.Histogram, Size)

	err := abs.IndexWorkers(Size, func(c int, done uint64) error {
		cb := iso.Turn.Unindex(uint64(c))
		h := abs.Histogram{
			Bins:   make([]float32, opts.Bins),
			Equity: 0,
		}

		var counter float32
		for _, cx := range card.All(cb...) {
			cbb := append(cb, cx)
			clus := abs.Cluster(iso.River.Index(cbb))
			eq := opts.Buckets.Get(clus)
			h = h.Increment(eq.WinDraw())
			counter++
		}

		h = h.Normalize()
		h.Equity = h.Equity / counter

		if done%uint64(opts.LogIteration) == 0 {
			opts.Logger.Printf("computing turn histograms: %d/%d", done, Size)
		}

		hh[c] = h
		return nil
	})

	return hh, err
}

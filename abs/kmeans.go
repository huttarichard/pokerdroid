package abs

import (
	"fmt"
	"math"
	"runtime"
	"sync"
	"sync/atomic"

	"github.com/pokerdroid/poker"
	"github.com/pokerdroid/poker/frand"
)

type Distance[H any] func(H, H) float64

type Recenter[H any] func(frand.Rand, []H) H

type Group[H any] struct {
	Center       H
	Observations []H
}

func (c *Group[H]) Append(point H) {
	c.Observations = append(c.Observations, point)
}

type Groups[H any] []Group[H]

func NewRandomGroups[H any](rng frand.Rand, hists []H, k int) (Groups[H], error) {
	var c Groups[H]
	if k == 0 {
		return c, fmt.Errorf("k must be greater than 0")
	}

	if k > len(hists) {
		return nil, fmt.Errorf("the size of the data set must at least equal k")
	}

	for i := 0; i < k; i++ {
		p := hists[rng.Intn(len(hists))]
		c = append(c, Group[H]{Center: p})
	}

	return c, nil
}

// Reset clears all point assignments
func (c Groups[H]) Reset() {
	for i := 0; i < len(c); i++ {
		c[i].Observations = []H{}
	}
}

// Nearest returns the index of the cluster nearest to point
func (c Groups[H]) Nearest(point H, distance Distance[H]) int {
	var ci int
	dist := -1.0
	// Find the nearest cluster for this data point
	for i, cluster := range c {
		d := distance(point, cluster.Center)
		if dist < 0 || d < dist {
			dist = d
			ci = i
		}
	}
	return ci
}

type KMeansOpts[H any] struct {
	Clusters      int
	MaxIterations int
	MaxDelta      float64
	Logger        poker.Logger
	LogIteration  int
	Rng           frand.Rand
	Groups        Groups[H]

	Recenter Recenter[H]
	Distance Distance[H]
}

func KMeans[H any](data []H, opts KMeansOpts[H]) (Groups[H], float64, error) {
	if opts.MaxDelta == 0 {
		opts.MaxDelta = 0.01
	}

	if opts.MaxIterations == 0 {
		opts.MaxIterations = math.MaxInt64
	}

	if opts.Logger == nil {
		opts.Logger = poker.VoidLogger{}
	}

	cpus := runtime.NumCPU()

	chln := len(data)
	points := make([]int, chln)
	changes := 1
	cc := opts.Groups

	for i := 0; changes > 0; i++ {
		opts.Logger.Printf("kmeans cycle: %d", i)

		changes = 0
		cc.Reset()

		var wg sync.WaitGroup
		var mux sync.Mutex

		type eq struct {
			p H
			i int
		}
		ch := make(chan eq)

		for i := 0; i < cpus; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()

				for p := range ch {
					ci := cc.Nearest(p.p, opts.Distance)
					mux.Lock()
					cc[ci].Append(p.p)
					if points[p.i] != ci {
						points[p.i] = ci
						changes++
					}
					mux.Unlock()
				}
			}()
		}

		for p, point := range data {
			if opts.LogIteration > 0 && p%opts.LogIteration == 0 {
				opts.Logger.Printf("kmeans iteration: %d", p)
			}
			ch <- eq{point, p}
		}

		close(ch)
		wg.Wait()

		for ci := 0; ci < len(cc); ci++ {
			if len(cc[ci].Observations) == 0 {
				// During the iterations, if any of the cluster centers has no
				// data points associated with it, assign a random data point
				// to it.
				// Also see: http://user.ceng.metu.edu.tr/~tcan/ceng465_f1314/Schedule/KMeansEmpty.html
				var ri int
				for {
					// find a cluster with at least two data points, otherwise
					// we're just emptying one cluster to fill another
					ri = opts.Rng.Intn(chln)
					if len(cc[points[ri]].Observations) > 1 {
						break
					}
				}
				cc[ci].Append(data[ri])
				points[ri] = ci
				// Ensure that we always see at least one more iteration after
				// randomly assigning a data point to a cluster
				changes = chln
			}
		}

		if changes > 0 {
			opts.Logger.Printf("kmeans recentering, changes %d", changes)

			perw := len(cc) / cpus
			var madech uint32
			for i := 0; i < cpus; i++ {
				wg.Add(1)
				go func(rng frand.Rand) {
					defer wg.Done()

					start := i * perw
					end := start + perw

					if i == cpus-1 {
						end = len(cc)
					}

					for start < end {
						cc[start].Center = opts.Recenter(rng, cc[start].Observations)
						start++

						atomic.AddUint32(&madech, 1)

						if madech%100 == 0 {
							opts.Logger.Printf("kmeans recentering: %d/%d", madech, len(cc))
						}
					}
				}(frand.Clone(opts.Rng))
			}
			wg.Wait()
		}

		if i == opts.MaxIterations || changes < int(float64(chln)*opts.MaxDelta) {
			break
		}
	}

	score := 0.0
	for _, c := range cc {
		for _, o := range c.Observations {
			distance := opts.Distance(c.Center, o)
			score += distance * distance
		}
	}

	opts.Logger.Printf("kmeans score: %f", score)

	// here we need to add best score
	return cc, score, nil
}

func KMedoid[H any](e []H, bins int, distance Distance[H]) H {
	// For k-medoids, we select the actual data point that minimizes
	// the sum of distances to all other points in the cluster
	if len(e) == 0 {
		var v H
		return v
	}

	var best int
	var min float64 = math.MaxFloat64
	var mux sync.Mutex

	IndexWorkers(len(e), func(i int, done uint64) error {
		var dist float64
		for j, other := range e {
			if i != j {
				dist += distance(e[i], other)
			}
		}

		if dist >= min {
			return nil
		}

		mux.Lock()
		defer mux.Unlock()
		min = dist
		best = i

		return nil
	})

	// Return the actual data point as the center
	return e[best]
}

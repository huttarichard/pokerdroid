package cfr

import (
	"context"
	"fmt"
	"math"
	"runtime"
	"sync"
	"time"

	"github.com/pokerdroid/poker"
	"github.com/pokerdroid/poker/abs"
	"github.com/pokerdroid/poker/dealer"
	"github.com/pokerdroid/poker/frand"
	"github.com/pokerdroid/poker/tree"
)

// RunParams defines the parameters to run CFR.
type RunParams struct {
	Workers    int
	Iterations uint64
	BatchSize  uint64
	EpochSize  uint64
	Rng        frand.Rand
	Game       *tree.Root
	Logger     poker.Logger
	Abs        abs.Mapper
	Sampler    dealer.Dealer
	Checkpoint func(it uint64, stop bool)
}

// NewRunParams constructs new RunParams with reasonable defaults.
func NewRunParams(game *tree.Root, sampler dealer.Dealer, abs abs.Mapper) RunParams {
	return RunParams{
		Workers:    runtime.NumCPU(),
		Iterations: math.MaxUint64,
		BatchSize:  20_000 / uint64(runtime.NumCPU()),
		EpochSize:  20_000,
		Rng:        frand.NewUnsafeInt(42),
		Game:       game,
		Logger:     poker.VoidLogger{},
		Abs:        abs,
		Sampler:    sampler,
		Checkpoint: func(it uint64, stop bool) {},
	}
}

// SetBatch sets the batch and epoch size.
func (p *RunParams) SetBatch(bs uint64, x uint64) {
	p.BatchSize = bs
	p.EpochSize = bs * x
}

// SetEpochs sets the total number of epochs to run.
func (p *RunParams) SetEpochs(it uint64) {
	p.Iterations = it * p.EpochSize
}

// Validate ensures all required parameters are provided.
func (p *RunParams) Validate() error {
	if p.BatchSize*uint64(p.Workers) > p.EpochSize {
		return fmt.Errorf("batch size * workers must be less than epoch size")
	}
	if p.Rng == nil {
		return fmt.Errorf("rng is required")
	}
	if p.Game == nil {
		return fmt.Errorf("game is required")
	}
	if p.Workers <= 0 {
		return fmt.Errorf("workers must be greater than 0")
	}
	if p.Game.State == nil {
		return fmt.Errorf("game state is required")
	}
	return nil
}

// Run executes CFR iterations until the total iteration count is reached.
// Each worker updates its own EV and update count into its index slot; these
// are aggregated and reported after each epoch.
func Run(ctx context.Context, c Runner, p RunParams) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(p.Workers + 1)

	epochDone := make(chan struct{}, p.Workers)

	// Per-worker storage for EV and update count.
	var exploits []float64
	var mux sync.Mutex

	evs := make([]float64, p.Workers)
	ups := make([]uint64, p.Workers)
	revs := make([]float64, 0, 10)

	// Worker function using its index.
	worker := func(idx int, rng frand.Rand) {
		defer wg.Done()

		sampler := p.Sampler.Clone()

		for {
			ev, up := c.Run(Params{
				Iterations: p.BatchSize,
				Rng:        rng,
				Sampler:    sampler,
			})

			mux.Lock()
			evs[idx] += ev
			ups[idx] += up
			mux.Unlock()

			select {
			case <-ctx.Done():
				// Stop the worker if the context is done.
				return
			case epochDone <- struct{}{}:
			}

			if p.Game.Iteration > p.Iterations {
				return
			}
		}
	}

	exploit := func(pit uint64, start time.Time, stop bool) {
		// Run exploitation
		ev := Exploit(ctx, ExploitParams{
			Root:       p.Game,
			Iterations: uint64(p.Workers),
			Params:     p.Game.Params,
			Rng:        p.Rng,
			Abs:        p.Abs,
			Sampler:    p.Sampler,
			Workers:    p.Workers,
		})
		if len(exploits) >= 1_000 {
			exploits = exploits[1:]
		}
		exploits = append(exploits, ev)

		// Prune nodes below the threshold.
		// dis := tree.DiscardBelowEpsilon(p.Game, p.PruneT, p.Workers)
		eps := p.Game.Iteration / p.EpochSize

		// Aggregate per-worker EV and update counts.
		sumev := float64(0.0)
		sumup := uint64(0)

		mux.Lock()
		for i := 0; i < p.Workers; i++ {
			sumev += evs[i]
			sumup += ups[i]

			evs[i] = 0
			ups[i] = 0
		}
		mux.Unlock()

		avgev := sumev / float64(p.Workers)

		// Update running average EV
		if len(revs) == 10 {
			// Remove oldest value
			revs = revs[1:]
		}
		revs = append(revs, avgev)

		// Calculate running average
		ravg := float64(0.0)
		for _, ev := range revs {
			ravg += ev
		}
		ravg /= float64(len(revs))

		var exp float64
		for _, ev := range exploits {
			exp += ev
		}
		exp /= float64(len(exploits))

		// Record iteration delta and update stats.
		st := &Stats{Start: start}
		st.It = p.Game.Iteration - pit
		st.TotIt = p.Game.Iteration
		st.States = p.Game.States
		st.Nodes = p.Game.Nodes
		st.Up = sumup
		st.EV = ravg
		st.Exploit = exp
		st.Epoch = eps

		// Report current epoch statistics.
		p.Logger.Printf(st.String())
		p.Checkpoint(eps, stop)
	}

	for i := 0; i < p.Workers; i++ {
		r := frand.Clone(p.Rng)
		go worker(i, r)
	}

	go func() {
		defer wg.Done()

		for {
			pit := p.Game.Iteration
			start := time.Now()

			var counter uint16

		LOOP:
			if p.Game.Iteration > p.Iterations {
				return
			}
			select {
			case <-epochDone:
				counter++
				if counter < uint16(p.Workers) {
					goto LOOP
				}
				counter = 0
				exploit(pit, start, false)

			case <-ctx.Done():
				exploit(pit, start, true)
				return
			}
		}
	}()

	// Wait for epoch to complete
	wg.Wait()
	close(epochDone)
}

package baselinenn

import (
	"context"

	"github.com/nlpodyssey/spago/ag"
	"github.com/nlpodyssey/spago/mat/float"
	"github.com/nlpodyssey/spago/nn"
	"github.com/nlpodyssey/spago/optimizers"
	"github.com/nlpodyssey/spago/optimizers/gradclipper"
	"github.com/pokerdroid/poker"
	absp "github.com/pokerdroid/poker/abs/pack"
	"github.com/pokerdroid/poker/deep"
	"github.com/pokerdroid/poker/frand"
	"github.com/pokerdroid/poker/tree"
	"golang.org/x/sync/errgroup"
)

type TrainerParams struct {
	Workers int

	// Checkpoint
	Checkpoint func(ctx context.Context, stop bool, i int) error

	Roots deep.Roots

	// ML
	BatchSize int
	Model     *Model
	Optimizer *optimizers.Optimizer

	// Common
	Rand   frand.Rand
	Logger poker.Logger
}

type Trainer[T float.DType] struct {
	cfg TrainerParams

	model     *Model
	optimizer *optimizers.Optimizer
}

func NewTrainer[T float.DType](cfg TrainerParams) *Trainer[T] {

	return &Trainer[T]{
		cfg: cfg,

		model:     cfg.Model,
		optimizer: cfg.Optimizer,
	}
}

func (g *Trainer[T]) Train(ctx context.Context) error {
	ch := make(chan deep.Sample[T], g.cfg.Workers)
	defer close(ch)

	gx, ctx := errgroup.WithContext(ctx)

	var run int

	// Start training
	gx.Go(func() error {
		buffer := make(deep.Samples[T], 0, g.cfg.BatchSize)

		// Add gradient clipper
		clipper := &gradclipper.NormClipper{
			MaxNorm:  0.5, // Start conservative, adjust if needed
			NormType: 2.0, // L2 norm (most common choice)
		}

		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case dd := <-ch:
				buffer = append(buffer, dd)
			}

			if len(buffer) < g.cfg.BatchSize {
				continue
			}

			g.cfg.Logger.Printf("training batch %d", run)

			var nb *deep.Batch
			nb, buffer = deep.NewBatchFromSamples(buffer, g.cfg.BatchSize)

			predictions := g.model.Forward(nb.X...)

			// Calculate loss
			loss := deep.HuberSeq(predictions, nb.Y, 0.1, true)
			// loss := losses.MSESeq(predictions, nb.Y, true)

			// Backward pass
			if err := ag.Backward(loss); err != nil {
				return err
			}

			// Clip gradients before optimization
			clipper.ClipGradients(nn.Parameters(g.model))

			// Update weights
			if err := g.optimizer.Optimize(); err != nil {
				return err
			}

			g.cfg.Logger.Printf("loss: %f", loss.Item().F64())

			err := g.cfg.Checkpoint(ctx, false, run)
			if err != nil {
				return err
			}
			run++
		}
	})

	worker := func(rng frand.Rand, abs *absp.Abs, root *tree.Root) error {
		actions, err := tree.SamplePath(root, root.Next, rng)
		if err != nil {
			return err
		}

		for i := range actions {
			acts := actions[:i+1]
			last := acts[len(acts)-1]
			aax := last.Parent.Actions

			for c, pol := range aax.Policies.Map {
				equity := abs.Equity(last.State.Street, c)

				ff, ss := Encode[T](EncodeParams{
					Params: root.Params,
					Player: last.Parent,
					Equity: equity,
				})

				baselines := EncodeBasline[T](
					g.model.Actions,
					aax.Actions,
					pol.Baseline,
					ss,
				)

				traj := deep.Sample[T]{
					Features: ff,
					Labels:   baselines,
				}

				select {
				case <-ctx.Done():
					return ctx.Err()
				case ch <- traj:
				}
			}
		}

		return nil
	}

	for x := 0; x < g.cfg.Workers; x++ {
		rng := frand.Clone(g.cfg.Rand)

		gx.Go(func() error {
			g.cfg.Logger.Printf("launching worker %d", x)

			for {
				select {
				case <-ctx.Done():
					return ctx.Err()
				default:
				}

				rnd := frand.Intn(len(g.cfg.Roots))
				root := g.cfg.Roots[rnd]

				fr, err := root.FileRoot.NewRoot()
				if err != nil {
					return err
				}

				err = worker(rng, root.Abs, fr)
				if err != nil {
					g.cfg.Logger.Printf("error: %v", err)
					continue
				}
			}
		})
	}

	err := gx.Wait()

	g.cfg.Logger.Printf("training done")
	if err != nil {
		return err
	}

	return g.cfg.Checkpoint(ctx, true, run)
}

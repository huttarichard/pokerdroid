package cmddeep

import (
	"bytes"
	"context"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"

	"github.com/nlpodyssey/spago/nn"
	"github.com/nlpodyssey/spago/optimizers"
	"github.com/nlpodyssey/spago/optimizers/adam"
	baseline "github.com/pokerdroid/poker/deep/baseline"
	"github.com/pokerdroid/poker/frand"
	"github.com/spf13/cobra"
)

type baselineArgs struct {
	workers int
	output  string

	batchSize     int
	learningRate  float64
	seed          int64
	checkpointMod int

	model string

	config string
}

var bstf = baselineArgs{}

func init() {
	flags := baselineCMD.Flags()

	flags.Int64Var(&bstf.seed, "seed", 42, "random number generator seed")
	flags.IntVar(&bstf.workers, "workers", runtime.NumCPU(), "worker for each instance")

	// Config
	flags.StringVar(&bstf.config, "config", "./deep/baseline/config.json", "path to trained config")
	cobra.MarkFlagRequired(flags, "config")

	// NN
	flags.IntVar(&bstf.batchSize, "batch-size", 1_000, "training batch size")
	flags.Float64Var(&bstf.learningRate, "learning-rate", 1e-4, "Adam optimizer learning rate")

	flags.StringVar(&bstf.model, "model", "", "model to start training from")

	flags.IntVar(&bstf.checkpointMod, "checkpoint", 10, "checkpoint frequency (iterations)")
	flags.StringVar(&bstf.output, "output", "experiments", "output path")
}

var baselineCMD = &cobra.Command{
	Use:   "baseline",
	Short: "will train baseline model",
	// Args:  cobra.ExactArgs(1),

	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := signal.NotifyContext(cmd.Context(), os.Interrupt)
		defer cancel()

		logger := log.Default()
		logger.Printf("loading config from %s", bstf.config)

		tcfg, err := NewTrainConfig[baseline.ModelParams](bstf.config)
		if err != nil {
			logger.Fatal(err)
		}

		cfgs, err := NewRootsFromConfigs(tcfg.Solutions)
		if err != nil {
			logger.Fatal(err)
		}
		defer cfgs.Close()

		rnd := frand.NewUnsafeInt(bstf.seed)

		model := baseline.NewModel[float32](tcfg.Model)
		model.InitRandom(rnd)

		if bstf.model != "" {
			logger.Printf("loading model from %s", bstf.model)

			f, err := os.Open(bstf.model)
			if err != nil {
				logger.Fatal(err)
			}

			model, err = baseline.NewFromFile[float32](f)
			if err != nil {
				logger.Fatal(err)
			}

			f.Close()
		}

		conf := adam.NewDefaultConfig()
		conf.StepSize = bstf.learningRate

		optimizer := optimizers.New(nn.Parameters(model), adam.New(conf))

		// Configure generator
		cfg := baseline.TrainerParams{
			Workers: bstf.workers,

			// ML
			BatchSize: bstf.batchSize,
			Model:     model,
			Optimizer: optimizer,

			Rand:   rnd,
			Logger: logger,
			Roots:  cfgs,

			// Checkpoint
			Checkpoint: func(ctx context.Context, stop bool, i int) error {
				if i%bstf.checkpointMod != 0 && !stop {
					return nil
				}

				buf := bytes.NewBuffer(nil)
				err := model.Save(buf)
				if err != nil {
					return err
				}

				name := filepath.Join(bstf.output, "baseline.gob")

				err = os.WriteFile(name, buf.Bytes(), 0644)
				if err != nil {
					return err
				}

				// mx, err := baseline.NewFromFile[float32](buf)
				// if err != nil {
				// 	return err
				// }

				// var diff float64
				// var totalSamples int
				// var count int

				// for _, sol := range solutions {
				// 	for _, root := range sol.Roots {
				// 		d, samples := baseline.Validate(baseline.ValidateParams{
				// 			Root:    root.Root,
				// 			Buckets: sol.Buckets,
				// 			Model:   mx,
				// 			Depth:   8,
				// 			Rng:     rnd,
				// 			Rounds:  30,
				// 		})

				// 		diff += d
				// 		totalSamples += samples
				// 		count++

				// 		tree.DiscardReferenceAtDepth(root.Root, 3)
				// 	}
				// }

				// logger.Printf("validation diff: %f, samples: %d", diff/float64(count), totalSamples)

				return nil
			},
		}

		gen := baseline.NewTrainer[float32](cfg)

		err = gen.Train(ctx)
		if err != nil {
			logger.Fatal(err)
		}

		logger.Printf("done")
	},
}

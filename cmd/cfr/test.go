package cmdcfr

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"
	"time"

	absp "github.com/pokerdroid/poker/abs/pack"
	"github.com/pokerdroid/poker/card"
	"github.com/pokerdroid/poker/cfr"
	"github.com/pokerdroid/poker/chips"
	"github.com/pokerdroid/poker/dealer"
	holdemdealer "github.com/pokerdroid/poker/dealer/holdem"
	"github.com/pokerdroid/poker/frand"
	"github.com/pokerdroid/poker/table"
	"github.com/pokerdroid/poker/tree"
	"github.com/pokerdroid/poker/tree/profiling"
	"github.com/spf13/cobra"
)

type testArgs struct {
	abs        string
	runtime    time.Duration
	iterations uint64
	workers    int
}

var ta = testArgs{}

func init() {
	flags := testCMD.Flags()

	flags.StringVar(&ta.abs, "abs", "", "path to the abstraction")
	cobra.MarkFlagRequired(flags, "abs")

	flags.DurationVar(&ta.runtime, "runtime", time.Minute*15, "runtime of the test")
	flags.Uint64Var(&ta.iterations, "it", 100_000, "number of iterations")
	flags.IntVar(&ta.workers, "workers", runtime.NumCPU()*3, "number of workers")
}

var testCMD = &cobra.Command{
	Use:   "test",
	Short: "will test againts simple tree",

	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := signal.NotifyContext(cmd.Context(), os.Interrupt)
		defer cancel()

		logger := log.Default()
		logger.Print("running test")

		abs, err := absp.NewFromFile(ta.abs)
		if err != nil {
			logger.Fatal(err)
		}

		type algo struct {
			name string
			fn   func(game *tree.Root, dealer dealer.Dealer) cfr.Runner
		}

		algos := []algo{
			// {
			// 	name: "mc_outcome_sampler_decay_bu_0.05",
			// 	fn: func(game *tree.Root, dealer dealer.Dealer) cfr.Runner {
			// 		return cfr.NewMC(cfr.MCParams{
			// 			PS: sampler.NewRobust(2),
			// 			TS: sampler.NewRobust(2),
			// 			Tree:     game,
			// 			Discount: policy.CFRD2(2, 0.5, 3),
			// 			Abs:      abs,
			// 			Sampler:  dealer,
			// 			BU:       policy.BaselineEMAClamp(0.05, float32(20)),
			// 		})
			// 	},
			// },

		}

		rng := frand.NewUnsafeInt(42)

		for _, algo := range algos {
			prms := table.NewGameParams(uint8(2), chips.NewFromInt(400))
			prms.BetSizes = table.BetSizesDeep
			prms.SbAmount = chips.NewFromFloat(1)
			prms.TerminalStreet = table.River
			prms.MaxActionsPerRound = uint8(5)

			game, err := tree.NewRoot(prms)
			if err != nil {
				logger.Fatal(err)
			}

			err = tree.ExpandFull(game)
			if err != nil {
				logger.Fatal(err)
			}

			logger.Printf("running %s", algo.name)

			sampler := holdemdealer.New(holdemdealer.SamplerParams{
				NumPlayers: game.Params.NumPlayers,
				Terminal:   table.River,
			})

			rprms := cfr.NewRunParams(game, sampler, abs)
			rprms.Workers = ta.workers
			rprms.SetBatch(10_000, uint64(ta.workers))

			wctx, cancel := context.WithTimeout(ctx, ta.runtime)
			cfr.Run(wctx, algo.fn(game, sampler), rprms)
			cancel()

			exploit := cfr.Exploit(ctx, cfr.ExploitParams{
				Sampler:    sampler,
				Abs:        abs,
				Rng:        rng,
				Root:       game,
				Params:     game.Params,
				Iterations: ta.iterations,
				Workers:    tf.workers,
			})

			logger.Printf("exploit %s: %f", algo.name, exploit)

			pf, err := profiling.New(profiling.Params{
				PlayerID: 0,
				Abs:      abs,
				Depth:    3,
				Board:    card.Cards{},
				Tree:     game,
				BetSizes: prms.BetSizes,
			})
			if err != nil {
				logger.Fatal(err)
			}

			err = profiling.SaveVisualization(pf, profiling.SaveAllDir(fmt.Sprintf("experiments/%s", algo.name)))
			if err != nil {
				logger.Fatal(err)
			}
		}

	},
}

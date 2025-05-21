package cmdbench

import (
	"log"
	"os"
	"os/signal"
	"runtime"

	absp "github.com/pokerdroid/poker/abs/pack"
	"github.com/pokerdroid/poker/bot"
	"github.com/pokerdroid/poker/bot/cfr"
	"github.com/pokerdroid/poker/bot/mc"
	"github.com/pokerdroid/poker/frand"
	"github.com/pokerdroid/poker/tree"
	"github.com/spf13/cobra"
)

type benchMcCFRArgs struct {
	tree string
	abs  string

	rounds  uint64
	workers int
}

var pmccfr = benchMcCFRArgs{}

func init() {
	flags := mcCFR_CMD.Flags()

	flags.StringVar(&pmccfr.tree, "tree", "", "path to the tree")
	flags.StringVar(&pmccfr.abs, "abs", "", "path to the abstraction")
	flags.Uint64Var(&pmccfr.rounds, "rounds", 100_000, "how many rounds to run")
	flags.IntVar(&pmccfr.workers, "workers", runtime.NumCPU(), "worker for each instance")
}

var mcCFR_CMD = &cobra.Command{
	Use:   "mc-cfr",
	Short: "will bench cfr againts mc",

	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := signal.NotifyContext(cmd.Context(), os.Interrupt)
		defer cancel()

		logger := log.Default()
		logger.Print("loading tree")

		game, err := tree.NewFromFile(pmccfr.tree)
		if err != nil {
			log.Fatal(err)
		}

		logger.Printf("nodes: %+v", game.States)
		logger.Printf("full: %+v", game.Full)
		logger.Printf("iteration: %+v", game.Iteration)
		logger.Printf("%s", game.Params.String())

		logger.Print("loading abstraction")

		abs, err := absp.NewFromFile(pmccfr.abs)
		if err != nil {
			log.Fatal(err)
		}

		rng := frand.NewHash()

		advisors := []bot.Advisor{
			mc.NewAdvisor(),

			cfr.Simple{
				Abs:  abs,
				Tree: game,
				Rand: rng,
			},
		}

		for id, p := range advisors {
			logger.Printf("advisor %d: %T", id, p)
		}

		logger.Print("starting benchmark")

		bot.Benchmark(ctx, bot.BenchmarkParams{
			Advisors: advisors,
			Logger:   logger,
			Workers:  pmccfr.workers,
			Rounds:   int(pmccfr.rounds),
			Rand:     rng,
			Params:   game.Params,
		})

	},
}

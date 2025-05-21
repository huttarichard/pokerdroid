package cmdbench

import (
	"log"
	"os"
	"os/signal"
	"runtime"
	"time"

	absp "github.com/pokerdroid/poker/abs/pack"
	"github.com/pokerdroid/poker/abs/river"
	"github.com/pokerdroid/poker/bot"
	"github.com/pokerdroid/poker/bot/cfr"
	"github.com/pokerdroid/poker/frand"
	"github.com/pokerdroid/poker/slumbot"
	"github.com/pokerdroid/poker/tree"
	"github.com/spf13/cobra"
)

type benchSlumbotCFRArgs struct {
	tree    string
	abs     string
	river   string
	rounds  uint64
	workers int
	search  bool
}

var pscfr = benchSlumbotCFRArgs{}

func init() {
	flags := slumbotCFR_CMD.Flags()

	flags.StringVar(&pscfr.tree, "tree", "", "path to the tree")
	flags.StringVar(&pscfr.abs, "abs", "", "path to the abstraction")

	flags.BoolVar(&pscfr.search, "search", false, "use search")
	flags.StringVar(&pscfr.river, "river", "", "path to the river abstraction")

	flags.Uint64Var(&pscfr.rounds, "rounds", 100_000, "how many rounds to run")

	flags.IntVar(&pscfr.workers, "workers", runtime.NumCPU()*8, "worker for each instance")
}

var slumbotCFR_CMD = &cobra.Command{
	Use:   "slumbot-cfr",
	Short: "will bench cfr againts mc",

	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := signal.NotifyContext(cmd.Context(), os.Interrupt)
		defer cancel()

		logger := log.Default()
		logger.Print("loading tree")

		game, err := tree.NewFromFile(pscfr.tree)
		if err != nil {
			log.Fatal(err)
		}

		logger.Printf("nodes: %+v", game.States)
		logger.Printf("full: %+v", game.Full)
		logger.Printf("iteration: %+v", game.Iteration)
		logger.Printf("%s", game.Params.String())

		logger.Print("loading abstraction")

		abs, err := absp.NewFromFile(pscfr.abs)
		if err != nil {
			log.Fatal(err)
		}

		rng := frand.NewHash()

		logger.Print("starting benchmark")

		workers := min(100, pscfr.workers)

		var adv bot.Advisor = cfr.Simple{
			Abs:  abs,
			Rand: rng,
			Tree: game,
		}

		if pscfr.search {
			var rabs *river.Abs

			if pscfr.river != "" {
				rabs, err = river.NewFromFile(pscfr.river)
				if err != nil {
					log.Fatal(err)
				}
			}

			adv = cfr.SearchAdvisor{
				Abs:         abs,
				Root:        game,
				Rand:        rng,
				MaxDuration: time.Second * 8,
				RiverAbs:    rabs,
			}
		}

		slumbot.Benchmark(ctx, slumbot.BenchmarkParams{
			Advisor:  adv,
			Username: "brownass",
			Password: "brownass",
			Logger:   logger,
			Workers:  workers,
			Rounds:   int(pscfr.rounds),
		})

	},
}

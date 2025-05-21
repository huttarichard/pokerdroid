package cmdcfr

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"runtime/pprof"

	"github.com/pokerdroid/poker/abs"
	absp "github.com/pokerdroid/poker/abs/pack"
	"github.com/pokerdroid/poker/card"
	"github.com/pokerdroid/poker/cfr"
	"github.com/pokerdroid/poker/chips"
	holdemdealer "github.com/pokerdroid/poker/dealer/holdem"
	"github.com/pokerdroid/poker/policy"
	"github.com/pokerdroid/poker/policy/sampler"
	"github.com/pokerdroid/poker/table"
	"github.com/pokerdroid/poker/tree"
	"github.com/pokerdroid/poker/tree/profiling"
	"github.com/spf13/cobra"
)

type traingArgs struct {
	abs string

	batch   uint64
	workers int
	save    int
	output  string
	tree    string

	depth      int
	players    int
	maxactions int
	limp       bool
	minBet     bool

	cpupprof string
	memprof  string
}

var tf = traingArgs{}

func init() {
	flags := trainCMD.Flags()
	flags.StringVar(&tf.abs, "abs", "", "path to the abstraction")
	cobra.MarkFlagRequired(flags, "abs")

	batch := uint64(200000)
	workers := runtime.NumCPU() * 4

	flags.IntVar(&tf.workers, "workers", workers, "worker for each instance")
	flags.Uint64Var(&tf.batch, "batch", batch, "what is the batch before reporting")

	flags.IntVar(&tf.save, "save", 200, "how many epochs to save")
	flags.StringVar(&tf.output, "output", "experiments", "output path")

	// Its either tree
	flags.StringVar(&tf.tree, "tree", "", "tree to load - continue training")

	// Or we generate a new tree
	flags.IntVar(&tf.depth, "depth", 100, "effective stack of players (default 100bb)")
	flags.IntVar(&tf.players, "players", 2, "number of players")
	flags.IntVar(&tf.maxactions, "maxactions", 12, "max actions per round")

	flags.BoolVar(&tf.limp, "limp", false, "use limp")
	flags.BoolVar(&tf.minBet, "minbet", false, "use min bet")

	flags.StringVar(&tf.cpupprof, "cpuprof", "", "cpu profile path")
	flags.StringVar(&tf.memprof, "memprof", "", "memory profile path")

}

var trainCMD = &cobra.Command{
	Use:   "train",
	Short: "will compute strategy",

	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := signal.NotifyContext(cmd.Context(), os.Interrupt)
		defer cancel()

		// flags := cmd.Flags()
		logger := log.Default()

		err := os.MkdirAll(tf.output, 0755)
		if err != nil {
			logger.Fatal(err)
		}

		logger.Printf("loading abstraction")

		var abs abs.Mapper

		abs = absp.NewIso()

		if tf.abs != "" {
			abs, err = absp.NewFromFile(tf.abs)
		}
		if err != nil {
			logger.Fatal(err)
		}

		var game *tree.Root

		if tf.tree != "" {
			logger.Printf("loading tree")
			game, err = tree.NewFromFile(tf.tree)
		} else {
			prms := table.NewGameParams(uint8(tf.players), chips.NewFromInt(int64(tf.depth)*2))
			prms.Limp = tf.limp
			prms.SbAmount = chips.NewFromFloat(1)
			prms.TerminalStreet = table.River
			prms.MaxActionsPerRound = uint8(tf.maxactions)
			prms.MinBet = tf.minBet
			prms.DisableV = true
			prms.SetBetSizes()
			game, err = tree.NewRoot(prms)

			if gm, ok := abs.(*absp.Abs); ok {
				game.AbsID = gm.UID
			}
		}

		if err != nil {
			logger.Fatal(err)
		}

		logger.Printf("experiment: %s", args[0])
		logger.Printf("abs: %s", game.AbsID.String())
		logger.Printf("%s", game.Params.String())

		dealer := holdemdealer.New(holdemdealer.SamplerParams{
			NumPlayers: game.Params.NumPlayers,
			Terminal:   table.River,
		})

		algo := cfr.NewMC(cfr.MCParams{
			PS: sampler.NewOutcome(0.4),
			TS: sampler.NewOutcome(0.2),

			Tree:     game,
			Discount: policy.CFRD(1.5, 0.5, 2),
			Abs:      abs,
			Sampler:  dealer,

			// Baseline EMA alpha
			// BU: policy.BaselineEMAClamp(0.01, float64(tf.depth)),
			BU: policy.BaselineEMA(0.01),

			// Prune: -500_000_000,
			Prune: 0,
		})

		// algo := cfr.NewSimpleMC(cfr.SimpleMCParams{
		// 	Tree:     game,
		// 	Discount: policy.CFRD3(1.5, 0.5, 2),
		// 	Abs:      abs,
		// 	BU:       policy.BaselineEMAClamp(0.01, float64(tf.depth)),
		// })

		rprms := cfr.NewRunParams(game, dealer, abs)
		rprms.Logger = logger
		rprms.Workers = tf.workers
		rprms.SetBatch(tf.batch, uint64(tf.workers))

		// CPU profile
		if tf.cpupprof != "" {
			var cpuf *os.File
			pth := filepath.Join(tf.output, tf.cpupprof)
			logger.Printf("cpu profile: %s", pth)
			cpuf, err = os.Create(pth)
			if err != nil {
				logger.Fatal(err)
			}

			pprof.StartCPUProfile(cpuf)

			defer func() {
				pprof.StopCPUProfile()
				cpuf.Close()
			}()
		}

		// Memory profile
		if tf.memprof != "" {
			var memf *os.File
			pth := filepath.Join(tf.output, tf.memprof)
			logger.Printf("memory profile: %s", pth)
			memf, err = os.Create(pth)
			if err != nil {
				logger.Fatal(err)
			}

			defer func() {
				// Write memory profile before exit
				if err := pprof.WriteHeapProfile(memf); err != nil {
					logger.Printf("could not write memory profile: %v", err)
				}
				memf.Close()
			}()
		}

		// debug.SetGCPercent(-1)

		rprms.Checkpoint = func(epoch uint64, stop bool) {
			runtime.GC()

			if epoch%uint64(tf.save) != 0 && !stop {
				return
			}

			// Create a new directory for the epoch
			output := filepath.Join(tf.output, fmt.Sprintf("epoch_%d", epoch))
			err = os.MkdirAll(output, 0755)
			if err != nil {
				logger.Fatal(err)
			}

			// Save the tree
			logger.Printf("saving policies")

			tree := filepath.Join(tf.output, "tree.bin")

			if err := os.Remove(tree); err != nil && !os.IsNotExist(err) {
				logger.Fatal(err)
			}

			// Create new tree.bin file
			f, err := os.Create(tree)
			if err != nil {
				logger.Fatal(err)
			}

			// Write tree directly to file
			err = game.WriteBinary(f)
			if err != nil {
				logger.Fatal(err)
			}

			f.Close()

			// Save the profile
			logger.Printf("building profile")

			buf := new(bytes.Buffer)

			st, err := profiling.New(profiling.Params{
				PlayerID: 0,
				Tree:     rprms.Game,
				Depth:    3,
				Board:    card.Cards{},
				BetSizes: rprms.Game.Params.BetSizes,
				Abs:      abs,
			})

			if err != nil {
				logger.Fatal(err)
			}

			err = profiling.SaveVisualization(st, profiling.SaveAllDir(output))
			if err != nil {
				logger.Fatal(err)
			}

			buf.WriteString("\n Total Info States: ")
			buf.WriteString(fmt.Sprintf("%d", rprms.Game.States))

			buf.WriteString("\n Total Iterations: ")
			buf.WriteString(fmt.Sprintf("%d", rprms.Game.Iteration))

			// buf.WriteString("\n Exploitability: ")
			// buf.WriteString(fmt.Sprintf("%f", exploit))

			buf.WriteString("\n====================\n")

			buf.WriteString("\nAverage policy: \n")
			st.WriteTable(profiling.FormatAveragePolicy, buf)

			buf.WriteString("\nCurrent policy: \n")
			st.WriteTable(profiling.FormatCurrentPolicy, buf)

			buf.WriteString("\nRegrets summary\n")
			st.WriteTable(profiling.FormatRegretSummary, buf)

			buf.WriteString("\nStrategy summary\n")
			st.WriteTable(profiling.FormatStrategySummary, buf)

			buf.WriteString("\nBaseline\n")
			st.WriteTable(profiling.FormatBaseline, buf)

			base := filepath.Join(output, "profile.txt")

			err = os.WriteFile(base, buf.Bytes(), 0644)

			if err != nil {
				logger.Fatalf("failed to create profile: %s", err)
			}
		}

		logger.Printf("starting trainer")
		cfr.Run(ctx, algo, rprms)

		logger.Printf("done")
	},
}

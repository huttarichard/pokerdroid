package cmdcfr

import (
	"fmt"
	"log"
	"math"
	"strings"

	"github.com/pokerdroid/poker/tree"
	"github.com/spf13/cobra"
)

type analyzeArgs struct {
	tree string
}

var an = analyzeArgs{}

func init() {
	flags := analyzeCMD.Flags()

	flags.StringVar(&an.tree, "tree", "", "path to the tree")
	cobra.MarkFlagRequired(flags, "tree")
}

var analyzeCMD = &cobra.Command{
	Use:   "analyze",
	Short: "will analyze the tree",

	Run: func(cmd *cobra.Command, args []string) {
		logger := log.Default()
		logger.Print("loading tree")

		game, err := tree.NewFromFile(an.tree)
		if err != nil {
			log.Fatal(err)
		}

		logger.Print("running analysis")

		type s struct {
			count        int
			terminal     int
			policies     int
			iterationMin float64
			iterationMax float64
			regretMin    float64
			regretMax    float64
			baselineMin  float64
			baselineMax  float64
		}

		var stats s

		countP := func(x *tree.Player) {
			if x.Actions == nil {
				return
			}

			if x.Actions.Policies.Len() == 0 {
				return
			}

			stats.policies += int(x.Actions.Policies.Len())

			for _, pol := range x.Actions.Policies.Map {
				stats.iterationMin = math.Min(stats.iterationMin, float64(pol.Iteration))
				stats.iterationMax = math.Max(stats.iterationMax, float64(pol.Iteration))

				for _, r := range pol.RegretSum {
					stats.regretMin = math.Min(stats.regretMin, float64(r))
					stats.regretMax = math.Max(stats.regretMax, float64(r))
				}

				for _, b := range pol.Baseline {
					stats.baselineMin = math.Min(stats.baselineMin, float64(b))
					stats.baselineMax = math.Max(stats.baselineMax, float64(b))
				}
			}
		}

		type counter struct {
			count int
			sum   float64
		}
		itavgdepth := make(map[int]counter)

		tree.Visit(game, -1, func(n tree.Node, children []tree.Node, depth int) bool {
			stats.count++

			switch x := n.(type) {
			case *tree.Terminal:
				stats.terminal++
				return true
			case *tree.Player:
				countP(x)

				if _, ok := itavgdepth[depth]; !ok {
					itavgdepth[depth] = counter{
						count: 0,
						sum:   0,
					}
				}

				if x == nil || x.Actions == nil || x.Actions.Policies == nil {
					return true
				}

				var sum float64
				var count int

				for _, p := range x.Actions.Policies.Map {
					sum += float64(p.Iteration)
					count++
				}

				itavgdepth[depth] = counter{
					count: itavgdepth[depth].count + count,
					sum:   itavgdepth[depth].sum + sum,
				}
			}

			return true
		})

		var ss strings.Builder

		ss.WriteString(fmt.Sprintf("game: %s\n", game.Params.String()))

		ss.WriteString(fmt.Sprintf("total nodes: %d\n", stats.count))
		ss.WriteString(fmt.Sprintf("terminal nodes: %d\n", stats.terminal))
		ss.WriteString(fmt.Sprintf("policies: %d\n", stats.policies))
		ss.WriteString(fmt.Sprintf("iterations total: %d\n", game.Iteration))
		ss.WriteString(fmt.Sprintf("iterations min: %f\n", stats.iterationMin))
		ss.WriteString(fmt.Sprintf("iterations max: %f\n", stats.iterationMax))
		ss.WriteString(fmt.Sprintf("regret min: %f\n", stats.regretMin))
		ss.WriteString(fmt.Sprintf("regret max: %f\n", stats.regretMax))
		ss.WriteString(fmt.Sprintf("baseline min: %f\n", stats.baselineMin))
		ss.WriteString(fmt.Sprintf("baseline max: %f\n", stats.baselineMax))

		for i := 0; i < 100; i++ {
			dp, ok := itavgdepth[i]
			if !ok {
				continue
			}
			ss.WriteString(fmt.Sprintf("depth %d: %f\n", i, dp.sum/float64(dp.count)))
		}

		logger.Print(ss.String())
	},
}

package cmdclus

import (
	"log"
	"os"

	"github.com/pokerdroid/poker/abs/river"
	"github.com/pokerdroid/poker/abs/turn"
	"github.com/pokerdroid/poker/frand"
	"github.com/spf13/cobra"
)

func init() {
	flags := turnCMD.Flags()

	flags.String("output", "turn.bin", "path to turn abstraction")
	flags.Int("clusters", 10000, "number of clusters")
	flags.Int("bins", 20, "number of bins")
	flags.Int("maxiter", 5000, "max iterations")

	flags.String("equities", "equities.bin", "path to equities buckets")
	cobra.MarkFlagRequired(flags, "equities")
}

var turnCMD = &cobra.Command{
	Use:   "turn",
	Short: "build turn abstraction",
	Run: func(cmd *cobra.Command, args []string) {
		flags := cmd.Flags()
		logger := log.Default()

		output, err := flags.GetString("output")
		if err != nil {
			panic(err)
		}

		clusters, err := flags.GetInt("clusters")
		if err != nil {
			logger.Fatal(err)
		}

		buckets, err := flags.GetString("equities")
		if err != nil {
			logger.Fatal(err)
		}

		maxiter, err := flags.GetInt("maxiter")
		if err != nil {
			logger.Fatal(err)
		}

		buks, err := river.NewBucketsFromFile(buckets, logger)
		if err != nil {
			logger.Fatal(err)
		}

		bins, err := flags.GetInt("bins")
		if err != nil {
			logger.Fatal(err)
		}

		logger.Printf("computing turn histograms")

		hh, err := turn.Compute(turn.ComputeOpts{
			Buckets:      buks,
			Logger:       logger,
			Bins:         bins,
			LogIteration: 100_000,
		})
		if err != nil {
			logger.Fatal(err)
		}

		logger.Printf("partitioning turn histograms")

		abs, err := turn.Partition(hh, turn.PartitionOpts{
			Clusters:       clusters,
			Logger:         logger,
			Rng:            frand.NewUnsafeInt(42),
			MaxIterations:  maxiter,
			LogIteration:   1_000_000,
			DeltaThreshold: 0.0000001,
			Bins:           bins,
		})
		if err != nil {
			logger.Fatal(err)
		}

		bb, err := abs.MarshalBinary()
		if err != nil {
			logger.Fatal(err)
		}

		err = os.WriteFile(output, bb, 0644)
		if err != nil {
			logger.Fatal(err)
		}
	},
}

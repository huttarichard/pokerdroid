package cmdclus

import (
	"log"
	"os"

	"github.com/pokerdroid/poker/abs/river"
	"github.com/pokerdroid/poker/frand"
	"github.com/spf13/cobra"
)

func init() {
	flags := riverCMD.Flags()

	flags.String("output", "river.bin", "path to river abstraction")
	flags.Int("clusters", 10000, "number of clusters")

	flags.String("equities", "equities.bin", "path to equities buckets")
	cobra.MarkFlagRequired(flags, "equities")
}

var riverCMD = &cobra.Command{
	Use:   "river",
	Short: "build river abstraction",
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

		buks, err := river.NewBucketsFromFile(buckets, logger)
		if err != nil {
			logger.Fatal(err)
		}

		logger.Printf("computing river abstraction")

		hh, err := river.Partition(buks, river.PartitionOpts{
			Clusters:       clusters,
			Logger:         logger,
			Rng:            frand.NewUnsafeInt(42),
			MaxIterations:  5000,
			LogIteration:   10_000,
			DeltaThreshold: 0.0000001,
		})
		if err != nil {
			logger.Fatal(err)
		}

		bb, err := hh.MarshalBinary()
		if err != nil {
			logger.Fatal(err)
		}

		err = os.WriteFile(output, bb, 0644)
		if err != nil {
			logger.Fatal(err)
		}
	},
}

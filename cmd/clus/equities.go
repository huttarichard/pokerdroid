package cmdclus

import (
	"log"
	"os"

	"github.com/pokerdroid/poker/abs/river"
	"github.com/spf13/cobra"
)

func init() {
	flags := equitiesCMD.Flags()
	flags.String("output", "equities.bin", "path to buckets")
}

var equitiesCMD = &cobra.Command{
	Use:   "equities",
	Short: "build equities buckets",
	Run: func(cmd *cobra.Command, args []string) {
		flags := cmd.Flags()
		logger := log.Default()

		output, err := flags.GetString("output")
		if err != nil {
			panic(err)
		}

		buckets, err := river.ComputeBuckets(logger)
		if err != nil {
			panic(err)
		}

		b, err := buckets.MarshalBinary()
		if err != nil {
			panic(err)
		}

		err = os.WriteFile(output, b, 0644)
		if err != nil {
			panic(err)
		}

		logger.Printf("%s written", output)

	},
}

package cmdclus

import (
	"log"
	"os"

	"github.com/google/uuid"
	"github.com/pokerdroid/poker/abs/flop"
	absp "github.com/pokerdroid/poker/abs/pack"
	"github.com/pokerdroid/poker/abs/river"
	"github.com/pokerdroid/poker/abs/turn"
	"github.com/spf13/cobra"
)

func init() {
	flags := packCMD.Flags()

	flags.String("output", "abs.bin", "path to turn abstraction")

	flags.String("flop", "flop.bin", "path to flop abstraction")
	cobra.MarkFlagRequired(flags, "flop")

	flags.String("turn", "turn.bin", "path to turn abstraction")
	cobra.MarkFlagRequired(flags, "turn")

	flags.String("river", "river.bin", "path to river abstraction")
	cobra.MarkFlagRequired(flags, "river")
}

var packCMD = &cobra.Command{
	Use:   "pack",
	Short: "pack turn abstraction",
	Run: func(cmd *cobra.Command, args []string) {
		flags := cmd.Flags()
		logger := log.Default()

		output, err := flags.GetString("output")
		if err != nil {
			panic(err)
		}

		f, err := flags.GetString("flop")
		if err != nil {
			logger.Fatal(err)
		}

		t, err := flags.GetString("turn")
		if err != nil {
			logger.Fatal(err)
		}

		r, err := flags.GetString("river")
		if err != nil {
			logger.Fatal(err)
		}

		logger.Printf("packing abstraction")

		abx := absp.Abs{
			UID:   uuid.New(),
			Flop:  &flop.Abs{},
			Turn:  &turn.Abs{},
			River: &river.Abs{},
		}

		// read flop abstraction

		ff, err := os.ReadFile(f)
		if err != nil {
			logger.Fatal(err)
		}

		err = abx.Flop.UnmarshalBinary(ff)
		if err != nil {
			logger.Fatal(err)
		}

		// read turn abstraction

		tt, err := os.ReadFile(t)
		if err != nil {
			logger.Fatal(err)
		}

		err = abx.Turn.UnmarshalBinary(tt)
		if err != nil {
			logger.Fatal(err)
		}

		// read river abstraction

		rr, err := os.ReadFile(r)
		if err != nil {
			logger.Fatal(err)
		}

		err = abx.River.UnmarshalBinary(rr)
		if err != nil {
			logger.Fatal(err)
		}

		// write packed abstraction

		bb, err := abx.MarshalBinary()
		if err != nil {
			logger.Fatal(err)
		}

		err = os.WriteFile(output, bb, 0644)
		if err != nil {
			logger.Fatal(err)
		}
	},
}

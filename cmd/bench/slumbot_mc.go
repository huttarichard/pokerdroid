package cmdbench

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"github.com/pokerdroid/poker/bot/mc"
	"github.com/pokerdroid/poker/slumbot"
	"github.com/spf13/cobra"
)

var mcSlumbotCMD = &cobra.Command{
	Use:   "slumbot-mc",
	Short: "will bench given mc againts slumbot",

	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGABRT)
		defer cancel()

		adv := mc.NewAdvisor()

		params := slumbot.BenchmarkParams{
			Advisor:  adv,
			Username: "brownass",
			Password: "brownass",
			Rounds:   100_000,
			Workers:  100,
			Logger:   log.Default(),
		}

		slumbot.Benchmark(ctx, params)

	},
}

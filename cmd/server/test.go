package cmdserver

import (
	"log"

	"github.com/pokerdroid/poker/bot"
	"github.com/pokerdroid/poker/card"
	"github.com/pokerdroid/poker/chips"
	"github.com/pokerdroid/poker/table"
	"github.com/spf13/cobra"
)

type testArgs struct {
	endpoint string
}

var tft = testArgs{}

func init() {
	flags := testCMD.Flags()
	flags.StringVar(&tft.endpoint, "endpoint", "http://localhost:8080/advise", "server address")
}

var testCMD = &cobra.Command{
	Use:   "test",
	Short: "Test server connection with sample request",
	Run: func(cmd *cobra.Command, args []string) {
		logger := log.Default()

		p := table.NewGameParams(2, chips.NewFromInt(100))

		game, err := table.NewGame(p)
		if err != nil {
			logger.Fatal("Failed to create game:", err)
		}

		state := bot.State{
			Params:    p,
			State:     game.Latest,
			Hole:      card.NewCardsFromString("as kh"),
			Community: card.Cards{},
		}

		// Send request to server
		action, err := bot.Request(tft.endpoint, state)
		if err != nil {
			logger.Fatal("Request failed:", err)
		}

		logger.Printf("Received action: %s (%.2f)", action.String(), action)
	},
}

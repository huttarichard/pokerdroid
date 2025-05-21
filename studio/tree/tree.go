package studiotree

import (
	"github.com/pokerdroid/poker"
	absp "github.com/pokerdroid/poker/abs/pack"
	"github.com/pokerdroid/poker/bot"
	"github.com/pokerdroid/poker/bot/cfr"
	"github.com/pokerdroid/poker/bot/mc"
	"github.com/pokerdroid/poker/card"
	"github.com/pokerdroid/poker/chips"
	"github.com/pokerdroid/poker/frand"
	"github.com/pokerdroid/poker/table"
	"github.com/pokerdroid/poker/tree"
	webview "github.com/webview/webview_go"
)

type BindParams struct {
	Roots   []*tree.Root
	Abs     *absp.Abs
	WebView webview.WebView
	Logger  poker.Logger
}

func Bind(p BindParams) error {
	inspector := NewInspector(p.Abs, p.Roots)

	p.WebView.Bind("rpc_tree_solutions", func() (response []*tree.Root, err error) {
		return p.Roots, nil
	})

	p.WebView.Bind("rpc_tree_get_state", func(actions []Action) (response *Result, err error) {
		response, err = inspector.Get(actions)
		return
	})

	agents := make(map[string]bot.Advisor)

	agents["simple"] = cfr.Advisor{
		Roots:   p.Roots,
		Rand:    frand.NewHash(),
		Abs:     p.Abs,
		Advisor: cfr.AdvisorSimple,
	}

	agents["search"] = cfr.Advisor{
		Roots:   p.Roots,
		Rand:    frand.NewHash(),
		Abs:     p.Abs,
		Advisor: cfr.AdvisorWithSearch,
	}

	gm := NewGameManager(p.Logger)

	type NewGameInput struct {
		Players int
		Stack   chips.Chips
		SB      chips.Chips
	}

	p.WebView.Bind("rpc_game_new_mc", func(params NewGameInput) (response int64, err error) {
		response, err = gm.New(NewGameParams{
			Players: []bot.Advisor{
				NewUserPlayer(),
				mc.NewAdvisor(),
			},
			StackSize: params.Stack,
			SB:        params.SB,
		})
		return response, err
	})

	type GameStateOuput struct {
		Params      table.GameParams           `json:"params"`
		State       *table.State               `json:"state"`
		Hole        card.Cards                 `json:"hole"`
		Community   card.Cards                 `json:"community"`
		Legal       table.DiscreteLegalActions `json:"legal"`
		Pot         chips.Chips                `json:"pot"`
		Round       int                        `json:"round"`
		LastWinners []int                      `json:"winner"`
	}

	p.WebView.Bind("rpc_game_get_state", func(gameID int64) (response GameStateOuput, err error) {
		game, err := gm.Get(gameID)
		if err != nil {
			return GameStateOuput{}, err
		}

		st, err := gm.State(gameID)
		if err != nil {
			return GameStateOuput{}, err
		}

		lg := table.NewDiscreteLegalActions(st.Params, st.State)

		output := GameStateOuput{
			Params:      st.Params,
			State:       st.State,
			Hole:        st.Hole,
			Community:   st.Community,
			Legal:       lg,
			Pot:         st.State.Players.PaidSum(),
			Round:       game.Round,
			LastWinners: game.LastWinners,
		}

		return output, nil
	})

	p.WebView.Bind("rpc_game_action", func(gameID int64, action table.DiscreteAction) error {
		return gm.Action(gameID, action)
	})

	return nil
}

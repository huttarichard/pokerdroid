package cfr

import (
	"context"
	"errors"
	"time"

	"github.com/pokerdroid/poker"
	"github.com/pokerdroid/poker/abs"
	"github.com/pokerdroid/poker/bot"
	"github.com/pokerdroid/poker/chips"
	"github.com/pokerdroid/poker/frand"
	"github.com/pokerdroid/poker/table"
	"github.com/pokerdroid/poker/tree"
)

type AdvisorFn func(s *Advisor, root *tree.Root) bot.Advisor

type Advisor struct {
	Roots   []*tree.Root
	Rand    frand.Rand
	Logger  poker.Logger
	Abs     abs.Mapper
	Advisor AdvisorFn
}

func AdvisorSimple(s *Advisor, root *tree.Root) bot.Advisor {
	return &Simple{
		Rand: s.Rand,
		Tree: root,
		Abs:  s.Abs,
	}
}

func AdvisorWithSearch(s *Advisor, root *tree.Root) bot.Advisor {
	return &SearchAdvisor{
		Abs:         s.Abs,
		Rand:        s.Rand,
		MaxDuration: time.Second * 7,
		Root:        root,
	}
}
func (a Advisor) Advise(ctx context.Context, loggr poker.Logger, state bot.State) (tb table.DiscreteAction, err error) {
	root := tree.FindClosestRoot(a.Roots, state.Params, state.State.TurnPos)

	if root == nil {
		err = errors.New("no solution found")
		return
	}

	bbs := chips.NewListAlloc(len(state.Params.InitialStacks))
	bbz := chips.NewListAlloc(len(state.Params.InitialStacks))

	for i, bb := range state.Params.InitialStacks {
		bbx := state.Params.SbAmount.Mul(2)
		bbz[i] = root.Params.InitialStacks[i].Div(2)
		bbs[i] = bb.Div(bbx).Sub(bbz[i]).Abs()
	}

	a.Logger.Printf("Found: %v distance: %v", bbz, bbs)

	advx := a.Advisor(&a, root)
	return advx.Advise(ctx, loggr, state)
}

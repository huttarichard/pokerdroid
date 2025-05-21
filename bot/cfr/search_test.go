package cfr

import (
	"context"
	"runtime"
	"testing"
	"time"

	"github.com/kr/pretty"
	"github.com/pokerdroid/poker"
	absp "github.com/pokerdroid/poker/abs/pack"
	"github.com/pokerdroid/poker/abs/river"
	"github.com/pokerdroid/poker/card"
	"github.com/pokerdroid/poker/cfr"
	"github.com/pokerdroid/poker/frand"
	"github.com/pokerdroid/poker/table"
	"github.com/pokerdroid/poker/tree"
	"github.com/stretchr/testify/require"
)

func TestSearch(t *testing.T) {
	t.Logf("loading root")
	blueprint := tree.TstRootFromEnv(t, "SMALL_TREE")
	riv := river.TstNewAbs(t)

	rng := frand.NewHash()

	// r1 := blueprint.Next.(*tree.Chance)
	// r2 := r1.Next.(*tree.Player)
	// r3 := r2.Actions.Nodes[2].(*tree.Player)
	// r4 := r3.Actions.Nodes[1].(*tree.Chance)
	// r5 := r4.Next.(*tree.Player)
	// r6 := r5.Actions.Nodes[1].(*tree.Player)
	// r7 := r6.Actions.Nodes[1].(*tree.Chance)
	// r8 := r7.Next.(*tree.Player)
	// r9 := r8.Actions.Nodes[1].(*tree.Player)
	// r10 := r9.Actions.Nodes[1].(*tree.Chance)
	// r11 := r10.Next.(*tree.Player)

	t.Logf("loading abs")
	abs := absp.TstNewAbs(t)

	game, err := table.NewGame(blueprint.Params)
	require.NoError(t, err)

	require.NoError(t, game.Action(table.DCall))
	require.NoError(t, game.Action(table.DCheck))

	require.NoError(t, game.Action(table.DCheck))
	require.NoError(t, game.Action(table.DCheck))

	require.NoError(t, game.Action(table.DCheck))
	require.NoError(t, game.Action(table.DCheck))

	board := card.Cards{
		card.Card2C,
		card.Card3C,
		card.Card8S,
		card.Card5C,
		card.Card2D,
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*16)
	defer cancel()

	pretty.Println(game.Latest.String())

	rx, err := Search(ctx, &SearchParams{
		Abs:       abs,
		Logger:    &poker.TestingLogger{T: t},
		State:     game.Latest,
		Board:     board,
		Rng:       rng,
		BatchSize: 1000,
		EpochSize: 10,
		Tree:      blueprint,
		Workers:   runtime.NumCPU(),
		Params:    blueprint.Params,
		RiverAbs:  riv,
	})
	require.NoError(t, err)

	p, ok := rx.Player.Actions.Policies.Get(rx.Abs.Map(append(card.Cards{card.CardAS, card.CardAH}, board...)))
	require.True(t, ok)

	str := p.GetAverageStrategy()
	for i, act := range rx.Player.Actions.Actions {
		t.Logf("Action: %-20s | Baseline: %.4f | Strategy: %.4f%% \n", act.String(), p.Baseline[i], str[i]*100)
	}

	exploit := cfr.Exploit(context.Background(), cfr.ExploitParams{
		Abs:        rx.Abs,
		Sampler:    rx.Sampler,
		Rng:        rng,
		Root:       rx.Tree,
		Iterations: 100_000,
		Params:     rx.Tree.Params,
		Workers:    runtime.NumCPU(),
	})

	t.Logf("Exploitability: %f", exploit)
}

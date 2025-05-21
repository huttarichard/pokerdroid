package cfr

import (
	"context"
	"os"
	"testing"

	"github.com/pokerdroid/poker"
	"github.com/pokerdroid/poker/abs"
	absp "github.com/pokerdroid/poker/abs/pack"
	"github.com/pokerdroid/poker/bot"
	"github.com/pokerdroid/poker/card"
	"github.com/pokerdroid/poker/chips"
	"github.com/pokerdroid/poker/frand"
	"github.com/pokerdroid/poker/policy"
	"github.com/pokerdroid/poker/table"
	"github.com/pokerdroid/poker/tree"
	"github.com/stretchr/testify/require"
)

// mockGetter satisfies cluster.Getter for testing
type mockGetter struct{}

func (m *mockGetter) Map(_ card.Cards) abs.Cluster {
	return 0
}

// TestAdvisor_BasicScenario checks basic usage of the CFR Advisor.
// Ensures that an action is returned for a minimal scenario.
func TestAdvisor_BasicScenario(t *testing.T) {
	prms := table.NewGameParams(2, chips.NewFromInt(10))
	prms.BetSizes = [][]float32{{1, 2}}
	prms.MaxActionsPerRound = 2
	prms.TerminalStreet = table.Preflop

	s := table.NewState(prms)

	s, err := table.MakeInitialBets(prms, s)
	require.NoError(t, err)

	// Build a minimal tree root
	root, err := tree.NewRoot(prms)
	require.NoError(t, err)
	require.NotNil(t, root)

	err = tree.ExpandFull(root)
	require.NoError(t, err)

	p := root.Next.(*tree.Chance).Next.(*tree.Player)

	pol := policy.New(len(p.Actions.Actions))

	// Set all in
	idx := p.Actions.GetIdx(table.DiscreteAction(1.0))
	for i := range pol.Strategy {
		pol.StrategySum[i] = 0
	}
	pol.StrategySum[idx] = 1

	p.Actions.Policies.Store(abs.Cluster(0), pol)

	a := &Simple{
		Tree: root,
		Abs:  &mockGetter{},
		Rand: frand.NewHash(),
	}

	st := bot.State{
		Params:    prms,
		State:     s,
		Hole:      card.NewCardsFromString("as ks"),
		Community: card.Cards{},
	}

	act, err := a.Advise(context.Background(), &poker.TestingLogger{T: t}, st)
	require.NoError(t, err)
	t.Logf("Action returned: %v", act)

	require.Equal(t, table.DiscreteAction(1.0), act)
}

func TestFileAdvisor(t *testing.T) {
	abs := absp.TstNewAbs(t)

	es := os.Getenv("TRAINED_PATH")

	if es == "" {
		t.Skip("TRAINED_PATH is not set")
	}

	roots, err := tree.NewFileRootsFromDir(es)
	require.NoError(t, err)
	require.NotNil(t, roots)

	advs := &Advisor{
		Roots:  roots.Roots(),
		Logger: &poker.TestingLogger{T: t},
		Rand:   frand.NewHash(),
		Abs:    abs,
	}

	// Hole: [ad 3s]
	// Community: []
	// === TABLE =============================
	// Street: preflop
	// SB Amount: 0.02
	// BB Amount: 0.04
	// Street Action Count: 3
	// Call Amount: 0.00

	// Positions:
	// BTN: 1, SB: 1, BB: 0
	// Next to act: 0

	// Players:
	// P0: Active (Initial: 5.90, Paid: 0.04, Stack: 5.86) [TO ACT]
	// P1: Active (Initial: 5.07, Paid: 0.04, Stack: 5.03)

	// Street Commitment:
	// 	P0: 0.04
	// 	P1: 0.04

	// Street Action Count:
	// 	P0: 1 actions
	// 	P1: 2 actions

	// Legal Actions:
	// 	check: 0.00
	// 	bet: 0.04
	// 	allin: 5.86

	// Action History:
	// 	P1 sb 0.02
	// 	P0 bb 0.04
	// 	P1 call 0.02
	// ======================================

	// Tree node mapped: r:n:c:p
	// Cluster: 101
	// Policy: Policy:
	//   Iteration: 0
	//   StrategyWeight: 0.000000
	//   Strategy: [0.4685, 0.0006, 0.1253, 0.0203, 0.3854]
	//   RegretSum: [221014.8016, 255.9426, 59110.4523, 9551.4752, 181822.5560]
	//   StrategySum: [5641666.9212, 469901.1465, 784084.6025, 1327459.9937, 1869781.6685]
	//   Baseline: [-0.7505, -1.9221, 7.1889, -5.0859, -42.4914]

	// - Act: allin ($5.86) P: 55.90
	// - Act: check ($0.00) P: 4.66
	// - Act: bet ($0.04) P: 7.77
	// - Act: bet ($0.08) P: 13.15
	// - Act: bet ($0.16) P: 18.53

	// Chosen action: All In
	// 2025/02/08 19:53:57 action advised from agent: *cfr.FileAdvisor: All In

	prms := table.GameParams{
		NumPlayers:         2,
		MaxActionsPerRound: 5,
		TerminalStreet:     table.River,
		BtnPos:             1,
		InitialStacks:      chips.List{5.90, 5.07},
		SbAmount:           chips.NewFromFloat(0.02),
		BetSizes:           [][]float32{{0.5, 1, 2}},
		DisableV:           false,
	}

	st, err := table.NewGame(prms)
	require.NoError(t, err)

	err = st.Action(table.DCall)
	require.NoError(t, err)

	// t.Logf("State: %v", table.Debug(st.Latest, prms))

	ds, err := advs.Advise(context.Background(), &poker.TestingLogger{T: t}, bot.State{
		Params:    prms,
		State:     st.Latest,
		Hole:      card.NewCardsFromString("ad 3s"),
		Community: card.Cards{},
	})
	require.NoError(t, err)

	t.Logf("Action: %v", ds)
}

func TestFileAdvisor2(t *testing.T) {
	abs := absp.TstNewAbs(t)

	es := os.Getenv("TRAINED_PATH")

	if es == "" {
		t.Skip("TRAINED_PATH is not set")
	}

	roots, err := tree.NewFileRootsFromDir(es)
	require.NoError(t, err)
	require.NotNil(t, roots)

	advs := &Advisor{
		Abs:    abs,
		Roots:  roots.Roots(),
		Logger: &poker.TestingLogger{T: t},
		Rand:   frand.NewHash(),
	}

	prms := table.GameParams{
		NumPlayers:         2,
		MaxActionsPerRound: 5,
		TerminalStreet:     table.River,
		BtnPos:             1,
		InitialStacks:      chips.List{50, 50},
		SbAmount:           chips.NewFromFloat(1),
		BetSizes:           [][]float32{{0.5, 1, 2}},
		DisableV:           false,
	}

	st, err := table.NewGame(prms)
	require.NoError(t, err)

	err = st.Action(table.DCall)
	require.NoError(t, err)

	// t.Logf("State: %v", table.Debug(st.Latest, prms))

	a := map[table.DiscreteAction]float64{}
	var total float64

	state := bot.State{
		Params:    prms,
		State:     st.Latest,
		Hole:      card.NewCardsFromString("kd 9d"),
		Community: card.Cards{},
	}

	_, err = advs.Advise(context.Background(), &poker.TestingLogger{T: t}, state)
	require.NoError(t, err)

	advs.Logger = poker.VoidLogger{}

	for i := 0; i < 1000; i++ {
		ds, err := advs.Advise(context.Background(), &poker.TestingLogger{T: t}, state)
		require.NoError(t, err)

		if _, ok := a[ds]; !ok {
			a[ds] = 0
		}
		a[ds]++
		total++
	}

	for k, v := range a {
		a[k] = (v / total) * 100
	}

	t.Logf("Total: %v", a)
}

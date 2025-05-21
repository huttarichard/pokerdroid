package cfr

import (
	"context"
	"errors"
	"math"
	"runtime"
	"time"

	"github.com/pokerdroid/poker"
	"github.com/pokerdroid/poker/abs"
	absp "github.com/pokerdroid/poker/abs/pack"
	"github.com/pokerdroid/poker/abs/river"
	"github.com/pokerdroid/poker/bot"
	"github.com/pokerdroid/poker/card"
	"github.com/pokerdroid/poker/cfr"
	"github.com/pokerdroid/poker/dealer"
	holdemdealer "github.com/pokerdroid/poker/dealer/holdem"
	"github.com/pokerdroid/poker/frand"
	"github.com/pokerdroid/poker/iso"
	"github.com/pokerdroid/poker/policy"
	"github.com/pokerdroid/poker/table"
	"github.com/pokerdroid/poker/tree"
	"github.com/pokerdroid/poker/tree/mapping"
)

type SearchParams struct {
	Tree *tree.Root
	// Abs mapper
	Abs abs.Mapper
	// LargeAbs
	RiverAbs *river.Abs
	// Params
	Params table.GameParams
	// State we are searching from
	State *table.State
	// Board cards
	Board card.Cards
	// Rng
	Rng frand.Rand
	// Logger
	Logger poker.Logger
	// Batch size
	BatchSize uint64
	// Epoch size
	EpochSize uint64
	// Workers
	Workers int
}

type SearchResult struct {
	Tree    *tree.Root
	Player  *tree.Player
	Origin  *tree.Player
	Sampler dealer.Dealer
	Abs     abs.Mapper
}

func Search(ctx context.Context, params *SearchParams) (*SearchResult, error) {
	if params.Logger == nil {
		params.Logger = poker.VoidLogger{}
	}

	p, err := mapping.MapGameStateToTree(params.Params, params.State, params.Tree)
	if err != nil {
		return nil, err
	}

	result := &SearchResult{
		Tree:   params.Tree,
		Player: p,
		Origin: p,
		Abs:    params.Abs,
	}
	switch params.State.Street {
	case table.Preflop, table.Flop, table.Turn:
		return result, nil
	case table.River:
	default:
		return nil, errors.New("invalid street")
	}

	tp := params.Params.Clone()
	tp.BetSizes = table.BetSizesDeep
	tp.MaxActionsPerRound = 3
	tp.DisableV = true

	root := &tree.Root{
		Iteration: 0,
		States:    0,
		Nodes:     0,
		Params:    tp,
		Next:      nil,
		State:     params.State,
		Full:      false,
	}
	result.Tree = root

	err = tree.ExpandFull(root)
	if err != nil {
		return nil, err
	}

	logf := params.Logger.Printf
	logf("\n================================")
	logf("Starting search\n")

	actions := tree.ExtractActions(p)

	ranges, err := cfr.ComputeRange(cfr.ComputeRangeParams{
		Actions: actions,
		Players: root.Params.NumPlayers,
		Board:   params.Board,
		Abs:     params.Abs,
	})
	if err != nil {
		return nil, err
	}

	dealer := holdemdealer.NewWeighted(holdemdealer.RangeParams{
		NumPlayers: root.Params.NumPlayers,
		Board:      params.Board,
		Clusters:   params.Abs,
		Ranges:     ranges,
	})
	result.Sampler = dealer

	var absx abs.Mapper

	if params.RiverAbs != nil {
		absx = absp.AbsFn(func(cds card.Cards) abs.Cluster {
			return params.RiverAbs.Map(iso.River.Index(cds))
		})
	} else {
		absx = absp.NewIso()
	}

	result.Abs = absx

	opts := cfr.SimpleMCParams{
		Tree:     root,
		Abs:      absx,
		Discount: policy.CFRP,
		BU:       policy.BaselineEMA(0.25),
	}

	mc := cfr.NewSimpleMC(opts)

	rp := cfr.NewRunParams(root, dealer, absx)
	rp.SetBatch(params.BatchSize, params.EpochSize)
	rp.Iterations = math.MaxUint64
	rp.Logger = params.Logger
	rp.Rng = params.Rng

	cfr.Run(ctx, mc, rp)

	// No root, nor rollout should be here
	switch n := root.Next.(type) {
	case *tree.Player:
		result.Player = n
		return result, nil
	case *tree.Chance:
		result.Player = n.Next.(*tree.Player)
		return result, nil
	}

	return result, errors.New("invalid output node")
}

type SearchAdvisor struct {
	Abs         abs.Mapper
	RiverAbs    *river.Abs
	Rand        frand.Rand
	MaxDuration time.Duration
	Root        *tree.Root
}

func (a SearchAdvisor) Advise(ctx context.Context, loggr poker.Logger, state bot.State) (tb table.DiscreteAction, err error) {
	logf := loggr.Printf

	ctx, cancel := context.WithTimeout(ctx, a.MaxDuration)
	defer cancel()

	if state.State.Previous == nil {
		return 0, errors.New("state history is empty")
	}

	p, err := Search(ctx, &SearchParams{
		Abs:       a.Abs,
		State:     state.State,
		Board:     state.Community,
		Rng:       a.Rand,
		Logger:    loggr,
		EpochSize: 1,
		BatchSize: 500,
		Tree:      a.Root,
		Workers:   runtime.NumCPU(),
		Params:    state.Params,
		RiverAbs:  a.RiverAbs,
	})
	if err != nil {
		return 0, err
	}
	logf("Table Path:  %s\n", state.State.Path(state.Params.SbAmount))
	logf("Tree origin: %s\n", tree.GetPath(p.Origin))

	if p.Player == nil {
		return 0, ErrNoState
	}

	cards := append(state.Hole, state.Community...)
	cluster := p.Abs.Map(cards)

	pol, ok := p.Player.Actions.Policies.Get(cluster)
	if !ok {
		return 0, ErrNoPolicy
	}

	str := pol.GetAverageStrategy()

	logf("Iteration: %d\n", pol.Iteration)
	for i, act := range p.Player.Actions.Actions {
		logf("Action: %-20s | Baseline: %.4f | Strategy: %.4f%% \n", act.String(), pol.Baseline[i], str[i]*100)
	}

	action, err := SampleState(&SampleParams{
		Actions: p.Player.Actions.Actions,
		Policy:  pol,
		Rng:     a.Rand,
		State:   state,
	})
	if err != nil {
		return 0, err
	}

	logf("Chosen action: %s\n", action.String())
	return action, nil
}

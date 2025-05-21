package bot

import (
	"context"
	"math"
	"sync"

	"github.com/pokerdroid/poker"
	"github.com/pokerdroid/poker/card"
	"github.com/pokerdroid/poker/chips"
	"github.com/pokerdroid/poker/frand"
	"github.com/pokerdroid/poker/table"
)

type BenchmarkParams struct {
	Params   table.GameParams
	Advisors []Advisor
	Logger   poker.Logger
	Workers  int
	Rounds   int
	Rand     frand.Rand

	BeforeAction func(state State)
	OnAction     func(state State, action table.DiscreteAction)
}

func Benchmark(ctx context.Context, opts BenchmarkParams) chips.List {
	var wg sync.WaitGroup
	var mux sync.Mutex

	sumTotal := chips.NewListAlloc(opts.Params.NumPlayers)
	roundsTotal := int64(0)

	if opts.Logger == nil {
		opts.Logger = poker.VoidLogger{}
	}

	run := func(_ int, rounds int) {
		defer wg.Done()

	Loop:
		for r := 0; r <= rounds; r++ {
			deck := card.NewDeck(card.All())
			deck.Shuffle(opts.Rand)

			prms := opts.Params.Clone()
			prms.BtnPos = uint8(r % int(prms.NumPlayers))

			round, err := table.NewGame(prms)
			if err != nil {
				opts.Logger.Printf("error: %s", err)
				continue
			}

			var v []card.Cards
			for i := 0; i < int(prms.NumPlayers); i++ {
				v = append(v, deck.PopMulti(2))
			}

			var c card.Cards
			var s table.Street = table.Preflop

			for !round.Latest.Finished() {
				if ctx.Err() != nil {
					opts.Logger.Printf("context error: %s", ctx.Err())
					break Loop
				}

				advisor := opts.Advisors[round.Latest.TurnPos]

				state := State{
					Params:    prms,
					State:     round.Latest,
					Hole:      v[round.Latest.TurnPos],
					Community: c,
				}

				if opts.BeforeAction != nil {
					opts.BeforeAction(state)
				}

				action, err := advisor.Advise(ctx, opts.Logger, state)
				if err != nil {
					opts.Logger.Printf("error: advisor: %s", err)
					continue Loop
				}

				err = round.Action(action)
				if err != nil {
					opts.Logger.Printf("error: table: %s", err)
					continue Loop
				}

				if opts.OnAction != nil {
					opts.OnAction(State{
						Params:    prms,
						State:     round.Latest,
						Hole:      v[round.Latest.TurnPos],
						Community: c,
					}, action)
				}

				// Shift street
				if round.Latest.Street != s {
					switch round.Latest.Street {
					case table.Flop:
						c = append(c, deck.PopMulti(3)...)
					case table.Turn:
						c = append(c, deck.Pop())
					case table.River:
						c = append(c, deck.Pop())
					}

					s = round.Latest.Street
				}

			}

			winnings := table.GetWinnings(round.Latest.Players, &table.Cards{
				Community: c,
				Players:   v,
			})

			mux.Lock()
			for i, p := range round.Latest.Players {
				sumTotal[i] = sumTotal[i].Add(winnings[i].Sub(p.Paid))
			}
			roundsTotal += 1
			mux.Unlock()
		}
	}

	perW := float64(opts.Rounds) / float64(opts.Workers)
	perW = math.Ceil(perW)

	opts.Logger.Printf("running %d workers", opts.Workers)

	for x := 0; x < opts.Workers; x++ {
		wg.Add(1)
		go run(x, int(perW))
	}

	wg.Wait()

	opts.Logger.Printf("benchmark results")
	opts.Logger.Printf("total rounds: %d", roundsTotal)

	for p, sum := range sumTotal {
		result := sum.
			Div(chips.NewFromInt(roundsTotal)).
			Div(opts.Params.SbAmount.Mul(2))

		opts.Logger.Printf("player %d ======", p)
		opts.Logger.Printf("total winnings: %s", sum.StringFixed(2))
		opts.Logger.Printf("bb/hand: %s", result.StringFixed(4))

		sumTotal[p] = result
	}

	return sumTotal
}

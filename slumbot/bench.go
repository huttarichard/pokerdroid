package slumbot

import (
	"context"
	"fmt"
	"math"
	"os"
	"sync"
	"time"

	"github.com/pokerdroid/poker"
	"github.com/pokerdroid/poker/bot"
	"github.com/pokerdroid/poker/chips"
	"github.com/pokerdroid/poker/table"
)

func Run(ctx context.Context, tk string, a bot.Advisor) (ch chips.Chips, err error) {
	s := new(poker.SBLogger)
	var reward chips.Chips

	r, err := NewRound(ctx, tk)
	if err != nil {
		return chips.Zero, err
	}

	defer func() {

		// Log final state and results
		s.Printf("\n\n\n")
		s.Printf(fmt.Sprintf("Player Hand: %v\n", r.PHand.String()))
		if r.SHand != nil {
			s.Printf(fmt.Sprintf("Slumbot Hand: %v\n", r.SHand.String()))
		}
		s.Printf(fmt.Sprintf("Final Board: %v\n", r.Board.String()))
		s.Printf(fmt.Sprintf("Player Winnings: %s\n", reward.StringFixed(2)))

		if r.s != nil {
			s.Printf(fmt.Sprintf("Player Total Paid: %s\n", r.s.Players[0].Paid.StringFixed(2)))
			s.Printf(fmt.Sprintf("Slumbot Total Paid: %s\n", r.s.Players[1].Paid.StringFixed(2)))
			s.Printf(fmt.Sprintf("Pot Size: %s\n", r.s.Players[0].Paid.Add(r.s.Players[1].Paid).StringFixed(2)))
			s.Printf("\n\n")
			s.Printf(table.Debug(r.s, r.Params))
		}

		if err != nil {
			s.Printf("Error: %s\n", err)
		}

		name := fmt.Sprintf("experiments/slumbot/%d_%d-round.log", time.Now().UnixNano(), int(reward))
		os.WriteFile(name, s.Bytes(), 0644)
	}()

	for !r.s.Finished() {
		da, err := a.Advise(ctx, s, bot.State{
			State:     r.State(),
			Hole:      r.PHand,
			Community: r.Board,
			Params:    r.Params,
		})
		if err != nil {
			return chips.Zero, err
		}

		err = r.Action(da)
		if err != nil {
			return chips.Zero, err
		}
		s.Printf("\n\n\n")
	}

	reward = r.s.Players[0].Paid.Mul(-1).Add(r.Winnings[0])

	return reward, nil
}

type BenchmarkParams struct {
	Advisor  bot.Advisor
	Username string
	Password string
	Rounds   int
	Workers  int
	Logger   poker.Logger
}

func Benchmark(ctx context.Context, opts BenchmarkParams) chips.Chips {
	var wg sync.WaitGroup
	var mux sync.Mutex

	if opts.Logger == nil {
		opts.Logger = poker.VoidLogger{}
	}

	sumTotal := chips.Zero
	roundsTotal := 0

	run := func(_ int, rounds int) {
		defer wg.Done()

		token, err := Login(opts.Username, opts.Password)
		if err != nil {
			opts.Logger.Printf("error: %s", err)
			return
		}

		total := 0
		sum := chips.Zero

		for r := 0; r <= rounds; r++ {
			if ctx.Err() != nil {
				break
			}
			nctx, cancel := context.WithTimeout(ctx, time.Second*30)

			dd, err := Run(nctx, token, opts.Advisor)
			if err != nil {
				cancel()
				opts.Logger.Printf("error: %s", err)
				continue
			}
			total++

			opts.Logger.Printf(
				"winnings: %-6s (%.2f bb/hand, %d rounds)",
				dd.StringFixed(0),
				sum.Div(chips.New(total)).Div(chips.New(BigBlind)),
				total,
			)
			sum = sum.Add(dd)
			cancel()
		}
		mux.Lock()
		sumTotal = sumTotal.Add(sum)
		roundsTotal += total
		mux.Unlock()
	}

	perW := float64(opts.Rounds) / float64(opts.Workers)
	perW = math.Ceil(perW)

	opts.Logger.Printf("rounds: %d", opts.Rounds)
	opts.Logger.Printf("workers: %d", opts.Workers)

	for x := 0; x < opts.Workers; x++ {
		wg.Add(1)
		go run(x, int(perW))
	}

	wg.Wait()

	result := sumTotal.Div(chips.New(roundsTotal)).Div(chips.New(BigBlind))

	opts.Logger.Printf("total winnings: %s", sumTotal.StringFixed(2))
	opts.Logger.Printf("total rounds: %d", roundsTotal)
	opts.Logger.Printf("bb/hand: %s", result.StringFixed(4))

	return result
}

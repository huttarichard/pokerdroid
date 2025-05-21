package table

import (
	"fmt"

	"github.com/pokerdroid/poker/chips"
)

func MakeAction(p GameParams, r *State, aa Actioner) (*State, error) {
	action, amount := aa.GetAction(p, r)

	if !p.DisableV {
		if r.Street > p.TerminalStreet {
			return nil, fmt.Errorf("cannot perform action on street: %s", r.Street)
		}

		la := NewLegalActions(p, r)
		err := la.Validate(p, r, ActionAmount{
			Action: action,
			Amount: amount,
		})
		if err != nil {
			return nil, err
		}
	}

	np := r.Players[r.TurnPos]
	stack := p.InitialStacks[r.TurnPos].Sub(np.Paid)

	if action == AllIn || amount.GreaterThanOrEqual(stack) {
		amount = stack
		np.Status = StatusAllIn
	}

	if action == Fold {
		np.Status = StatusFolded
	}

	np.Paid = np.Paid.Add(amount)

	state := r.Next()
	state.Players[r.TurnPos] = np

	state.StreetAction++

	state.PSC[r.TurnPos] = state.PSC[r.TurnPos].Add(amount)
	state.PSAC[r.TurnPos]++
	state.PSLA[r.TurnPos] = action

	if amount.GreaterThan(r.BSC.Amount) {
		state.BSC.Amount = amount
		state.BSC.Addition = amount.Sub(r.BSC.Amount)
		state.BSC.Action = action
	}

	if action == Bet || action == Raise {
		state.BetAction++
	}

	return state, nil
}

func MakeInitialBets(p GameParams, r *State) (state *State, err error) {
	state, err = MakeAction(p, r, ActionAmount{
		Action: SmallBlind,
		Amount: p.SbAmount,
	})
	if err != nil {
		return
	}
	err = ShiftTurn(p, state)
	if err != nil {
		return
	}
	state, err = MakeAction(p, state, ActionAmount{
		Action: BigBlind,
		Amount: p.SbAmount.Mul(chips.NewFromInt(2)),
	})
	if err != nil {
		return
	}
	err = ShiftTurn(p, state)
	if err != nil {
		return
	}
	return state, nil
}

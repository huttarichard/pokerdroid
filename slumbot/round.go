package slumbot

import (
	"context"
	"errors"

	"github.com/pokerdroid/poker/card"
	"github.com/pokerdroid/poker/chips"
	"github.com/pokerdroid/poker/table"
)

type Round struct {
	s        *table.State
	Deck     *card.Deck
	PHand    card.Cards
	SHand    card.Cards
	Board    card.Cards
	Token    string
	Params   table.GameParams
	Response HandResponse
	Winnings chips.List
	ctx      context.Context
}

func NewRound(ctx context.Context, token string) (*Round, error) {
	nh, err := NewHand(ctx, token)
	if err != nil {
		return nil, err
	}

	// client pos 0 == slumbot is sb
	// client pos 1 == slumbot is bb
	// client btnPos 0 == player is bb
	// client btnPos 1 == player is sb

	btnPos := uint8(1)
	if nh.ClientPos == 1 {
		btnPos = 0
	}

	prms := table.NewGameParams(2, chips.NewFromFloat32(StackSize))
	prms.SbAmount = SmallBlind
	prms.BtnPos = btnPos
	prms.MaxActionsPerRound = 12
	prms.BetSizes = table.BetSizesDeep

	s := table.NewState(prms)

	hole := nh.GetHoleCards()

	s, err = table.MakeInitialBets(prms, s)
	if err != nil {
		return nil, err
	}

	r := &Round{
		s:        s,
		Deck:     card.NewDeck(card.All(hole...)),
		Token:    token,
		Response: nh,
		ctx:      ctx,
		PHand:    nh.GetHoleCards(),
		Board:    nh.GetBoardCards(),
		Params:   prms,
	}

	for _, a := range nh.Actions {
		if s.TurnPos != a.Player {
			panic("player mismatch")
		}

		r.s, err = table.MakeAction(prms, s, a)
		if err != nil {
			return nil, err
		}

		r.s, err = Move(r.Params, r.s, r)
		if err != nil {
			return nil, err
		}
	}

	return r, err
}

func (r *Round) Action(action table.Actioner) error {
	paid := r.s.PSC[0]
	act, amount := action.GetAction(r.Params, r.s)
	encoded := EncodeAction(act, amount.Add(paid))

	var err error
	r.Response, err = Act(r.ctx, r.Token, encoded, r.Response)
	if err != nil {
		return err
	}
	r.Deck.Remove(r.Response.GetBoardCards()...)

	firstAct := r.Response.Actions[0]
	if r.s.TurnPos != firstAct.Player {
		panic("player mismatch")
	}

	r.s, err = table.MakeAction(r.Params, r.s, firstAct)
	if err != nil {
		return err
	}

	r.s, err = Move(r.Params, r.s, r)
	if err != nil {
		return err
	}

	for _, a := range r.Response.Actions[1:] {
		if r.s.TurnPos != a.Player {
			panic("player mismatch")
		}
		r.s, err = table.MakeAction(r.Params, r.s, a.ActionAmount)
		if err != nil {
			return err
		}

		r.s = r.s.Next()

		r.s, err = Move(r.Params, r.s, r)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *Round) State() *table.State {
	return r.s
}

func Move(p table.GameParams, r *table.State, rx *Round) (state *table.State, err error) {
	state = r.Next()

	switch table.Rule(p, r) {
	case table.RuleShiftTurn:
		return state, table.ShiftTurn(p, state)

	case table.RuleShiftStreet:
		board := rx.Response.GetBoardCards()
		rx.Board = board

		return state, table.ShiftStreet(state)

	case table.RuleShiftStreetUntilEnd:
		board := rx.Response.GetBoardCards()
		rx.Board = board

		err := table.ShiftStreet(state)
		if err != nil {
			return nil, err
		}
		return Move(p, state, rx)

	case table.RuleFinish:
		if rx.Response.BotHoleCards != nil {
			rx.SHand = rx.Response.GetBotHoleCards()
			rx.Deck.Remove(rx.SHand...)
		}

		// player is always r.Params.BtnPos
		//  btnPos 0 == player is bb
		//  btnPos 1 == player is sb

		winner, losser := 0, 1
		if rx.Response.WonPot <= 0 {
			winner, losser = 1, 0
		}

		rx.Winnings = chips.NewListAlloc(2)
		rx.Winnings[winner] = chips.NewFromInt(int64(rx.Response.WonPot)).Abs()
		rx.Winnings[losser] = chips.Zero

		state.Street = table.Finished

		return state, nil

	default:
		return nil, errors.New("unknown rule")
	}
}

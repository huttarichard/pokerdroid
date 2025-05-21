package table

import (
	"errors"
	"fmt"
	"strings"

	"github.com/pokerdroid/poker/chips"
)

type RuleKind uint8

const (
	RuleShiftTurn RuleKind = iota
	RuleShiftStreet
	RuleShiftStreetUntilEnd
	RuleFinish
)

func (r RuleKind) String() string {
	switch r {
	case RuleShiftTurn:
		return "ShiftTurn"
	case RuleShiftStreet:
		return "ShiftStreet"
	case RuleShiftStreetUntilEnd:
		return "ShiftStreetUntilEnd"
	case RuleFinish:
		return "Finish"
	}
	return "unknown"
}

func (r RuleKind) MarshalJSON() ([]byte, error) {
	return []byte(`"` + r.String() + `"`), nil
}

func (r *RuleKind) UnmarshalJSON(data []byte) error {
	s := string(data)
	s = strings.Trim(s, `"`)

	switch s {
	case "ShiftTurn":
		*r = RuleShiftTurn
	case "ShiftStreet":
		*r = RuleShiftStreet
	case "ShiftStreetUntilEnd":
		*r = RuleShiftStreetUntilEnd
	case "Finish":
		*r = RuleFinish
	default:
		return fmt.Errorf("unknown RuleKind: %s", s)
	}
	return nil
}

func Rule(p GameParams, r *State) RuleKind {
	h2a := HaveToAct(r)

	// Count number of actions on street
	s := 0
	for _, x := range r.PSAC {
		s += int(x)
	}

	// If there is more than one player that needs to act
	// or if there is one player that needs to act and
	// there is an action on the street, this cover one
	// player going all in (h2a == 1) have to act but
	// it was on that given street when player went all in s > 0.

	// todo is s -2?
	if h2a > 1 || (h2a == 1 && s > 0) {
		return RuleShiftTurn
	}

	if r.Players.Len(IsWaitingAskPlayer) > 1 {
		if r.Street >= p.TerminalStreet {
			return RuleFinish
		}
		return RuleShiftStreet
	}

	if r.Players.Len(IsActivePlayer) <= 1 {
		return RuleFinish
	}

	if r.Street < p.TerminalStreet {
		return RuleShiftStreetUntilEnd
	}

	return RuleFinish
}

func HaveToAct(s *State) (hta int) {
	active := s.Players.Indicies(IsActivePlayer)
	if len(active) == 1 {
		return 0
	}

	psc := s.PSC.Max()

	for pi, p := range s.Players {
		if p.Status != StatusActive {
			continue
		}

		// Need to commit more than maxBet
		// if s.PSC[pi].LessThan(s.BSC.Amount) ||
		if s.PSC[pi].LessThan(psc) ||
			// Need to act
			s.PSAC[pi] == 0 ||
			// Need to act on preflop if you are big blind
			// If you are small blind you have to call so that is
			// covered by the previous check
			(s.Street == Preflop && s.PSAC[pi] == 1 && s.PSLA[pi] == BigBlind) {
			hta++
			continue
		}
	}
	return hta
}

func Move(p GameParams, r *State) (state *State, err error) {
	state = r.Next()

	switch Rule(p, r) {
	case RuleShiftTurn:
		return state, ShiftTurn(p, state)

	case RuleShiftStreet:
		return state, ShiftStreet(state)

	case RuleShiftStreetUntilEnd:
		err := ShiftStreet(state)
		if err != nil {
			return nil, err
		}
		return Move(p, state)

	case RuleFinish:
		state.Street = Finished
		return state, nil

	default:
		return nil, errors.New("unknown rule")
	}
}

func ShiftTurn(p GameParams, r *State) error {
	next := r.Players.FindWaitingPos(int(r.TurnPos + 1))
	if next == -1 {
		return errors.New("could not find next player pos")
	}
	r.TurnPos = uint8(next)

	// Call amount is max commitment
	// minus current player commitment
	cab := r.PSC.Max().Sub(r.PSC[r.TurnPos])
	if cab.LessThan(chips.Zero) {
		cab = chips.Zero
	}

	paid := r.Players[r.TurnPos].Paid
	stack := p.InitialStacks[r.TurnPos].Sub(paid)

	r.CallAmount = chips.Min(cab, stack)
	return nil
}

func ShiftStreet(r *State) error {
	r.Street++

	if r.Street == Finished {
		return nil
	}

	ShiftTurnStreetStart(r)
	r.StreetAction = 0
	r.BetAction = 0

	r.PSC = chips.NewListAlloc(len(r.Players))
	r.PSAC = make([]uint8, len(r.Players))
	r.PSLA = make([]ActionKind, len(r.Players))
	//
	r.BSC.Addition = chips.Zero
	r.BSC.Amount = chips.Zero
	r.BSC.Action = NoAction

	r.CallAmount = chips.Zero
	return nil
}

func ShiftTurnStreetStart(r *State) {
	btn, sb, bb := Positions(r)

	if len(r.Players) == 2 {
		if r.Street == Preflop {
			r.TurnPos = sb
		} else {
			r.TurnPos = bb
		}
		return
	}

	nis := bb
	if r.Street == Preflop {
		nis = btn
	}

	np := r.Players.FindWaitingPos(int(nis + 1))
	if np == -1 {
		return
	}

	r.TurnPos = uint8(np)
}

func Positions(s *State) (uint8, uint8, uint8) {
	pLen := uint8(len(s.Players))
	btn := s.BtnPos % pLen
	var sb uint8
	if pLen == 2 {
		sb = btn
	} else {
		sb = btn + 1
	}
	sb = sb % pLen
	bb := (sb + 1) % pLen
	return uint8(btn), uint8(sb), uint8(bb)
}

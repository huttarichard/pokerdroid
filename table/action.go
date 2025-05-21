package table

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/pokerdroid/poker/chips"
)

// Actioner is an interface for actions.
// Both DiscreteAction and ActionAmount implements it.
type Actioner interface {
	fmt.Stringer
	GetAction(GameParams, *State) (ActionKind, chips.Chips)
}

// ActionKind represents action kind.
type ActionKind uint8

const (
	NoAction ActionKind = iota
	Bet
	SmallBlind
	BigBlind
	Fold
	Check
	Call
	Raise
	AllIn
)

// NewActionFromString returns ActionKind from string.
func NewActionFromString(str string) (ActionKind, error) {
	switch str {
	case "fold":
		return Fold, nil
	case "check":
		return Check, nil
	case "call":
		return Call, nil
	case "raise":
		return Raise, nil
	case "bet":
		return Bet, nil
	case "sb":
		return SmallBlind, nil
	case "bb":
		return BigBlind, nil
	case "allin":
		return AllIn, nil
	default:
		return 0, errors.New("invalid action")
	}
}

// String returns string representation of ActionKind.
func (a ActionKind) String() string {
	switch a {
	case Bet:
		return "bet"
	case SmallBlind:
		return "sb"
	case BigBlind:
		return "bb"
	case Fold:
		return "fold"
	case Check:
		return "check"
	case Call:
		return "call"
	case Raise:
		return "raise"
	case AllIn:
		return "allin"
	case NoAction:
		return "no action"
	default:
		return "unknown"
	}
}

// Short similar to String() but shorter.
func (a ActionKind) Short() string {
	switch a {
	case Bet:
		return "b"
	case SmallBlind:
		return "sb"
	case BigBlind:
		return "bb"
	case Fold:
		return "f"
	case Check:
		return "k"
	case Call:
		return "c"
	case Raise:
		return "b"
	case AllIn:
		return "a"
	default:
		return "u"
	}
}

// IsBet if bet or blinds.
func (a ActionKind) IsBet() bool {
	return a == Bet || a == SmallBlind || a == BigBlind
}

// IsBlind if small or big blind.
func (a ActionKind) IsBlind() bool {
	return a == SmallBlind || a == BigBlind
}

// IsRaise if raise or all-in.
func (a ActionKind) IsRaise() bool {
	return a == Raise || a == AllIn
}

// ActionAmount is an ActionKind with chips.
type ActionAmount struct {
	Action ActionKind  `json:"action"`
	Amount chips.Chips `json:"amount"`
}

var _ Actioner = ActionAmount{}

// String returns string representation of ActionAmount.
// Chips are fixed for 2 decimal places.
func (a ActionAmount) String() string {
	return fmt.Sprintf("%s %s", a.Action, a.Amount.StringFixed(2))
}

// GetAction implements Actioner interface.
func (a ActionAmount) GetAction(p GameParams, r *State) (ActionKind, chips.Chips) {
	switch a.Action {
	case Raise, Bet:
		return a.Action, a.Amount.Round(2)
	case Call:
		return a.Action, r.CallAmount
	case Check:
		return a.Action, chips.Zero
	case Fold:
		return a.Action, chips.Zero
	}
	return a.Action, a.Amount
}

// DiscreteAction is an action with discrete amount.
// This is useful for abstraction and working in
// discrete space. Everythign above 0 is multiple of pot.
//
// - DiscreteAction(-4) = All In
// - DiscreteAction(-3) = Fold
// - DiscreteAction(-2) = Call
// - DiscreteAction(-1) = Check
// - DiscreteAction(0)  = No Action
// - DiscreteAction(>0) = Raise/Bet
type DiscreteAction float32

var _ Actioner = DiscreteAction(0)

const (
	DAllIn    DiscreteAction = -4
	DFold     DiscreteAction = -3
	DCall     DiscreteAction = -2
	DCheck    DiscreteAction = -1
	DNoAction DiscreteAction = 0
)

// String returns string representation of DiscreteAction.
func (da DiscreteAction) String() string {
	switch da {
	case DFold:
		return "Fold"
	case DCheck:
		return "Check"
	case DCall:
		return "Call"
	case DAllIn:
		return "All In"
	default:
		return fmt.Sprintf("Raise %.2f POT", da)
	}
}

// GetAction implements Actioner interface.
func (a DiscreteAction) GetAction(p GameParams, r *State) (ActionKind, chips.Chips) {
	switch a {
	case DFold:
		return Fold, chips.Zero
	case DCheck:
		return Check, chips.Zero
	case DCall:
		return Call, r.CallAmount
	case DAllIn:
		paid := r.Players[r.TurnPos].Paid
		stack := p.InitialStacks[r.TurnPos]
		return AllIn, stack.Sub(paid)
	}

	pt := r.Players.PaidSum()
	amount := pt.Mul(chips.NewFromFloat32(float32(a)))

	if r.CallAmount.Equal(chips.Zero) {
		return Bet, amount.Round(2)
	}

	return Raise, amount.Round(2)
}

func (da DiscreteAction) Hash() int {
	return int(da)
}

// IsRaise if raise/bet or all-in.
func (da DiscreteAction) IsRaise() bool {
	return da > 0 || da == DAllIn
}

// Short similar to String() but shorter.
func (da DiscreteAction) Short() string {
	switch da {
	case DFold:
		return "f"
	case DCheck:
		return "k"
	case DCall:
		return "c"
	case DAllIn:
		return "a"
	default:
		return fmt.Sprintf("b%.2f", da)
	}
}

// LegalActions represents legal actions for a player.
// It is a map of ActionKind to chips.Chips.
// Chip amount is minimum of chips needed to perform action.
// Maximum in NLH is player stack.
type LegalActions map[ActionKind]chips.Chips

// NewLegalActions will create LegalActions for a given state.
func NewLegalActions(p GameParams, r *State) LegalActions {
	actions := make(LegalActions)

	if r.Street > p.TerminalStreet {
		return actions
	}

	nextPlayer := r.Players[r.TurnPos]
	stack := p.InitialStacks[r.TurnPos].Sub(nextPlayer.Paid)

	callAmount := r.CallAmount

	max := p.MaxActionsPerRound
	if r.Street == Preflop {
		max += 2
	}

	reachedMax := r.StreetAction+1 >= max

	if callAmount.Equal(chips.Zero) && !reachedMax {
		actions[Check] = chips.Zero
		bb := p.SbAmount.Mul(chips.NewFromInt(2))
		if stack.GreaterThan(bb) {
			actions[Bet] = bb
		}
		actions[AllIn] = stack

		return actions
	}

	if callAmount.Equal(chips.Zero) && reachedMax {
		actions[Check] = chips.Zero
		return actions
	}

	actions[Fold] = chips.Zero

	if reachedMax {
		actions[Call] = chips.Min(callAmount, stack)
		return actions
	}

	if stack.LessThan(callAmount) {
		actions[Call] = stack
		return actions
	}

	minRaise := callAmount

	if r.BSC.Amount.GreaterThan(chips.Zero) {
		if r.BSC.Action.IsBet() {
			minRaise = minRaise.Add(r.BSC.Amount)
		}
		if r.BSC.Action.IsRaise() {
			minRaise = r.BSC.Addition.Add(r.BSC.Amount)
		}
	}

	rest := chips.Zero
	for i, s := range p.InitialStacks {
		if uint8(i) == r.TurnPos {
			continue
		}
		s := s.Sub(r.Players[i].Paid)
		if s.GreaterThan(rest) {
			rest = s
		}
	}

	restaz := rest.GreaterThan(chips.Zero)

	if minRaise.LessThan(stack) && restaz {
		actions[Raise] = minRaise
	}

	actions[Call] = callAmount

	if callAmount.LessThan(stack) && restaz {
		actions[AllIn] = stack
	}

	return actions
}

// String returns string representation of LegalActions.
func (l LegalActions) String() string {
	buf := bytes.NewBufferString("Legal Actions:\n")
	for k, v := range l {
		buf.WriteString(fmt.Sprintf("\t%s: %s\n", k, v.StringFixed(2)))
	}
	return buf.String()
}

// Validate will validate action for a given player.
func (la LegalActions) Validate(p GameParams, r *State, ax ActionAmount) error {
	nps := r.TurnPos
	stack := p.InitialStacks[r.TurnPos].Sub(r.Players[r.TurnPos].Paid)

	a := ax.Action
	amount := ax.Amount

	if _, ok := la[a]; !ok && a != BigBlind && a != SmallBlind {
		acts := make([]string, 0)
		for act := range la {
			acts = append(acts, act.String())
		}
		return fmt.Errorf("%d: illegal action %s for actions %s", nps, a, strings.Join(acts, ","))
	}

	if (a == Check || a == Fold) && !amount.Equal(chips.Zero) {
		return fmt.Errorf("%d: illegal check/fold amount", nps)
	}

	if a == Call && !la[Call].Equal(amount) {
		return fmt.Errorf("%d: illegal call amount: %s - must be: %s", nps, amount.String(), la[Call].String())
	}

	if a == Raise && amount.LessThan(la[Raise]) {
		return fmt.Errorf("%d: illegal raise amount: %s < %s", nps, amount.String(), la[Raise].String())
	}

	if a == Bet && amount.LessThan(la[Bet]) {
		return fmt.Errorf("%d: illegal bet amount: %s < %s", nps, amount.String(), la[Bet].String())
	}

	if a == AllIn && !amount.Equal(stack) {
		return fmt.Errorf("%d: illegal all-in amount: %s != %s", nps, amount.String(), stack.String())
	}

	if amount.GreaterThan(stack) {
		return fmt.Errorf("%d: illegal %s amount: %s > %s", nps, a, amount.String(), stack.String())
	}

	return nil
}

// DiscreteLegalActions represents legal discrete actions for a player.
// It is similar in nature to LegalActions except chips are excat amounts
// to be used for action.
type DiscreteLegalActions map[DiscreteAction]chips.Chips

func (d DiscreteLegalActions) MarshalJSON() ([]byte, error) {
	// Convert map to array of [action, amount] pairs
	pairs := make([][2]float32, 0, len(d))
	for act, amt := range d {
		pairs = append(pairs, [2]float32{float32(act), amt.Float32()})
	}

	// Sort by action value
	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i][0] < pairs[j][0]
	})

	return json.Marshal(pairs)

}

// List returns list of actions.
func (d DiscreteLegalActions) List() []DiscreteAction {
	actions := make([]DiscreteAction, 0, len(d))
	for a := range d {
		actions = append(actions, a)
	}
	sort.Slice(actions, func(i, j int) bool {
		return actions[i] < actions[j]
	})
	return actions
}

// NewDiscreteLegalActions creates new DiscreteLegalActions for a given state.
// This is using NewLegalActions and converting it to discrete space.
func NewDiscreteLegalActions(p GameParams, r *State) DiscreteLegalActions {
	actions := make(DiscreteLegalActions)
	legalActions := NewLegalActions(p, r)

	if _, ok := legalActions[Fold]; ok {
		actions[DFold] = chips.Zero
	}

	if _, ok := legalActions[Check]; ok {
		actions[DCheck] = chips.Zero
	}

	openAct := r.Street == Preflop && r.StreetAction == 2
	amount, ok := legalActions[Call]

	if openAct && p.Limp && ok {
		actions[DCall] = amount
	}

	if !openAct && ok {
		actions[DCall] = amount
	}

	var minRaise chips.Chips

	if _, ok := legalActions[Raise]; ok {
		minRaise = legalActions[Raise]
	}

	if _, ok := legalActions[Bet]; ok {
		minRaise = legalActions[Bet]
	}

	if _, ok := legalActions[AllIn]; ok {
		actions[DAllIn] = legalActions[AllIn]
	}

	np := r.Players[r.TurnPos]
	stack := p.InitialStacks[r.TurnPos].Sub(np.Paid)

	if len(p.BetSizes) == 0 {
		return actions
	}

	pot := r.Players.PaidSum()
	betSizes := p.BetSizes[len(p.BetSizes)-1]

	if r.BetAction < uint8(len(p.BetSizes)) {
		betSizes = p.BetSizes[r.BetAction]
	}

	if !minRaise.Equal(chips.Zero) {
		// Add minraise as discrete action.
		if p.MinBet {
			actions[DiscreteAction(minRaise.Div(pot))] = minRaise
		}

		for _, amount := range betSizes {
			na := pot.Mul(chips.NewFromFloat32(amount))
			if na.LessThan(minRaise) || na.GreaterThan(stack) {
				continue
			}
			actions[DiscreteAction(amount)] = na
		}
	}

	return actions
}

// Equal compares two legal discrete action maps.
func (l DiscreteLegalActions) Equal(b DiscreteLegalActions) bool {
	if len(l) != len(b) {
		return false
	}
	for k, v := range l {
		if y, ok := b[k]; !ok || !v.Equal(y) {
			return false
		}
	}
	return true
}

// String returns string representation of DiscreteLegalActions.
func (l DiscreteLegalActions) String() string {
	buf := bytes.NewBufferString("Discrete Legal Actions:\n")
	for k, v := range l {
		buf.WriteString(fmt.Sprintf("\t%s: %s\n", k, v.StringFixed(2)))
	}
	return buf.String()
}

// FromActionToDiscrete will convert ActionKind and chips into DiscreteAction.
func FromActionToDiscrete(p GameParams, r *State, ax ActionAmount) DiscreteAction {
	action := ax.Action
	amount := ax.Amount

	if action.IsBlind() {
		return DNoAction
	}

	switch action {
	case Fold:
		return DFold
	case Check:
		return DCheck
	case Call:
		return DCall
	case AllIn:
		return DAllIn
	}

	if !action.IsRaise() && !action.IsBet() {
		return DNoAction
	}

	dla := NewDiscreteLegalActions(p, r)

	type actionDiff struct {
		Action DiscreteAction
		Diff   chips.Chips
	}

	closest := make([]actionDiff, 0)

	for k, ch := range dla {
		closest = append(closest, actionDiff{
			Action: k,
			Diff:   ch.Sub(amount).Abs(),
		})
	}

	if _, ok := dla[DAllIn]; ok && len(closest) == 0 {
		return DAllIn
	}

	sort.Slice(closest, func(i, j int) bool {
		return closest[i].Diff.GreaterThan(closest[j].Diff)
	})

	return closest[0].Action
}

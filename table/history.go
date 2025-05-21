package table

import (
	"fmt"
	"strings"

	"github.com/pokerdroid/poker/chips"
)

// PlayerAction represents a single action taken by a player
type PlayerAction struct {
	Pos    uint8 
	Street Street
	Action ActionAmount
	State  *State
}

type History []PlayerAction

func (h History) String() string {
	var str strings.Builder
	var currentStreet Street
	var firstAction bool = true

	var y int
	for z, action := range h {
		if action.Action.Action.IsBlind() {
			y++
			continue
		}

		i := z - y

		// Check for street change
		if i == 0 {
			currentStreet = action.Street
		} else if action.Street != currentStreet {
			str.WriteString("/")
			currentStreet = action.Street
			firstAction = true
		}

		// Add action separator within same street
		if !firstAction {
			str.WriteString(":")
		}
		firstAction = false

		// Convert action to rune format
		var rune string
		switch action.Action.Action {
		case Fold:
			rune = "f"
		case Check:
			rune = "k"
		case Call:
			rune = "c"
		case AllIn:
			rune = "a"
		case Bet, Raise:
			rune = fmt.Sprintf("b%.2f", action.Action.Amount.Float64())
		}

		str.WriteString(rune)
	}

	return str.String()
}

func (s *State) Path(sb chips.Chips) string {
	var h strings.Builder
	h.WriteString("r")
	var st Street

	for _, x := range s.History() {
		if x.Action.Action.IsBlind() {
			continue
		}
		if x.State.Street != st {
			st = x.State.Street
			h.WriteString(":n")
		}
		psc := x.State.PSC[x.Pos]
		psc = psc.Add(x.Action.Amount).Div(sb)
		isAllin := x.Action.Action == AllIn || x.State.Players[x.Pos].Status == StatusAllIn

		if x.Action.Action.IsRaise() || x.Action.Action.IsBet() {
			if isAllin {
				h.WriteString(":a")
			} else {
				h.WriteString(fmt.Sprintf(":b%.2f", psc.Float64()))
			}
		} else if x.Action.Action == Call {
			h.WriteString(":c")
		} else if x.Action.Action == Check {
			h.WriteString(":k")
		}
	}

	return h.String()
}

// History traverses the Previous state chain and reconstructs the action history
func (s *State) History() History {
	var actions History

	cur := s
	for cur.Previous == nil {
		return actions
	}

	for cur.Previous != nil {
		prev := cur.Previous

		// Skip if street changed to prevent transition artifacts
		if cur.Street != prev.Street {
			cur = prev
			continue
		}

		for i := range cur.PSAC {
			if cur.PSAC[i] == prev.PSAC[i] {
				continue
			}

			// Skip if not a real action
			if cur.PSLA[i] == NoAction {
				continue
			}

			actions = append(actions, PlayerAction{
				Pos:    uint8(i),
				Street: cur.Street,
				Action: ActionAmount{
					Action: cur.PSLA[i],
					Amount: cur.PSC[i].Sub(prev.PSC[i]),
				},
				State: prev,
			})
		}

		cur = prev
	}

	// Reverse to get chronological order
	for i := 0; i < len(actions)/2; i++ {
		j := len(actions) - 1 - i
		actions[i], actions[j] = actions[j], actions[i]
	}

	return actions
}

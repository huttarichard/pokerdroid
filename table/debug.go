package table

import (
	"fmt"
	"strings"
)

func Debug(s *State, params GameParams) string {
	var b strings.Builder

	b.WriteString("=== TABLE =============================\n")
	// Basic game info
	b.WriteString(fmt.Sprintf("Street: %s\n", s.Street))
	b.WriteString(fmt.Sprintf("SB Amount: %s\n", params.SbAmount.StringFixed(2)))
	b.WriteString(fmt.Sprintf("BB Amount: %s\n", params.SbAmount.Mul(2).StringFixed(2)))
	b.WriteString(fmt.Sprintf("Street Action Count: %d\n", s.StreetAction))
	b.WriteString(fmt.Sprintf("Call Amount: %s\n", s.CallAmount.StringFixed(2)))

	// Get positions
	btn, sb, bb := Positions(s)

	// Position info
	b.WriteString("\nPositions:\n")
	b.WriteString(fmt.Sprintf("BTN: %d, SB: %d, BB: %d\n", btn, sb, bb))
	b.WriteString(fmt.Sprintf("Next to act: %d\n", s.TurnPos))

	// Player info
	b.WriteString("\nPlayers:\n")
	for i, player := range s.Players {
		b.WriteString(fmt.Sprintf("P%d: ", i))

		// Status
		switch player.Status {
		case StatusActive:
			b.WriteString("Active")
		case StatusFolded:
			b.WriteString("Folded")
		case StatusAllIn:
			b.WriteString("All-in")
		}

		// Stack and total paid
		initStack := params.InitialStacks[i]
		b.WriteString(fmt.Sprintf(" (Initial: %s", initStack.StringFixed(2)))
		b.WriteString(fmt.Sprintf(", Paid: %s", player.Paid.StringFixed(2)))
		b.WriteString(fmt.Sprintf(", Stack: %s", initStack.Sub(player.Paid).StringFixed(2)))
		b.WriteString(")")

		if i == int(s.TurnPos) {
			b.WriteString(" [TO ACT]")
		}
		b.WriteString("\n")
	}

	// Street commitment info
	if len(s.PSC) > 0 {
		b.WriteString("\nStreet Commitment:\n")
		for i, psc := range s.PSC {
			b.WriteString(fmt.Sprintf("\tP%d: %s\n", i, psc.StringFixed(2)))
		}
	}

	if len(s.PSAC) > 0 {
		b.WriteString("\nStreet Action Count:\n")
		for i, psac := range s.PSAC {
			b.WriteString(fmt.Sprintf("\tP%d: %d actions\n", i, psac))
		}
	}

	// Current betting info
	// Usually not needed, but can be useful

	// if !s.BSC.Amount.Equal(chips.Zero) {
	// 	b.WriteString("\nBetting Info:\n")
	// 	b.WriteString(fmt.Sprintf("Amount: %s\n", s.BSC.Amount.StringFixed(2)))
	// 	b.WriteString(fmt.Sprintf("Addition: %s\n", s.BSC.Addition.StringFixed(2)))
	// 	b.WriteString(fmt.Sprintf("Action: %s\n", s.BSC.Action))
	// 	if !s.CallAmount.Equal(chips.Zero) {
	// 		b.WriteString(fmt.Sprintf("Call Amount: %s\n", s.CallAmount.StringFixed(2)))
	// 	}
	// }

	// Legal actions for current player
	b.WriteString("\nLegal Actions:\n")
	la := NewLegalActions(params, s)
	for act, amount := range la {
		b.WriteString(fmt.Sprintf("\t%s: %s\n", act, amount.StringFixed(2)))
	}

	// Action history
	history := s.History()
	if len(history) > 0 {
		b.WriteString("\nAction History:\n")
		for _, act := range history {
			b.WriteString(fmt.Sprintf("\tP%d %s\n", act.Pos, act.Action.String()))
		}
	}

	b.WriteString("======================================\n")

	return b.String()
}

func DebugPrint(s *State, params GameParams) {
	fmt.Println(Debug(s, params))
}

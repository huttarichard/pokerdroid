package tree

import (
	"fmt"
	"strings"

	"github.com/pokerdroid/poker/chips"
	"github.com/pokerdroid/poker/table"
)

type Rune string

const (
	RuneUnknown Rune = "u"

	// Nodes
	RuneRoot      Rune = "r"
	RuneChance    Rune = "n"
	RuneTerminal  Rune = "t"
	RunePlayer    Rune = "p"
	RuneReference Rune = "f"
	RuneRollout   Rune = "l"

	// Actions
	RuneBet   Rune = "b"
	RuneFold  Rune = "f"
	RuneCall  Rune = "c"
	RuneCheck Rune = "k"
	RuneAllIn Rune = "a"
)

func NewRuneFromAction(t table.DiscreteAction) Rune {
	switch t {
	case table.DNoAction:
		return RuneUnknown
	case table.DFold:
		return RuneFold
	case table.DCall:
		return RuneCall
	case table.DCheck:
		return RuneCheck
	case table.DAllIn:
		return RuneAllIn
	default:
		return RuneBet
	}
}

func NewRuneFromNode(node Node) Rune {
	switch node.Kind() {
	case NodeKindRoot:
		return RuneRoot
	case NodeKindChance:
		return RuneChance
	case NodeKindPlayer:
		return RunePlayer
	case NodeKindTerminal:
		return RuneTerminal
	case NodeKindRollout:
		return RuneRollout
	default:
		return RuneUnknown
	}
}

func (r Rune) String() string {
	return string(r)
}

func (r Rune) WithAmount(f chips.Chips) Rune {
	return Rune(fmt.Sprintf("%s%.2f", r, f))
}

type Runes []Rune

func (r Runes) String() string {
	rx := make([]string, len(r))
	for i, rr := range r {
		rx[i] = rr.String()
	}
	return strings.Join(rx, ":")
}

func FindLeafRunes(n Node) (nn []Runes) {
	for _, r := range FindLeafNodes(n) {
		nn = append(nn, GetPath(r))
	}
	return
}

func GetPath(n Node) (b Runes) {
	if n == nil {
		return nil
	}

	// Start with the current node.
	cur := n
	b = append(b, NewRuneFromNode(cur))

	// Climb up the parent chain.
Loop:
	for {
		parent := cur.GetParent()
		if parent == nil {
			break
		}

		if cur == parent {
			panic("nodes are self-referencing")
		}

		var node Rune

		switch px := parent.(type) {
		case *Player:
			idx, found := px.GetActionIdx(cur)
			if !found {
				cur = parent
				continue Loop
			}

			ract := px.Actions.Actions[idx]
			act := NewRuneFromAction(ract)
			if act == RuneBet {
				pt := potmul(px.State, ract)
				act = act.WithAmount(pt)
			}
			node = act
		default:
			node = NewRuneFromNode(parent)
		}

		// Append the parent's node (root, chance, player, or terminal).
		b = append(b, node)
		cur = parent
	}

	// Reverse to get the path from root to leaf.
	for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]
	}
	return b
}

func potmul(st *table.State, action table.DiscreteAction) chips.Chips {
	if st == nil || st.PSC == nil {
		return chips.Zero
	}
	pt := st.Players.PaidSum()
	psc := st.PSC[st.TurnPos]
	return psc.Add(pt.Mul(chips.New(float32(action))))
}

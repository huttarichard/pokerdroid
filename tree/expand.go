package tree

import (
	"sync/atomic"

	"github.com/pokerdroid/poker/table"
)

func Expand(e *Root, n Node) error {
	switch x := n.(type) {
	case *Root:
		return ExpandRoot(e)
	case *Chance:
		return ExpandChance(e, x)
	case *Player:
		return ExpandPlayer(e, x)
	case *Reference:
		_, err := x.Expand()
		return err
	case *Terminal:
		return nil
	}
	return nil
}

func MustExpand(e *Root, n Node) {
	err := Expand(e, n)
	if err != nil {
		panic(err)
	}
}

func ExpandRoot(e *Root) error {
	if e.Next != nil {
		return nil
	}

	sa := e.State.StreetAction
	if e.State.Street == table.Preflop {
		sa -= 2 // Preflop has 2 actions (SB and BB)
	}

	if sa > 0 {
		player := &Player{
			Parent:  e,
			TurnPos: e.State.TurnPos,
			State:   e.State,
		}
		e.Next = player
		atomic.AddUint32(&e.Nodes, 1)
		return nil
	}

	chance := &Chance{
		Parent: e,
		State:  e.State,
	}
	e.Next = chance
	atomic.AddUint32(&e.Nodes, 1)
	return nil
}

func ExpandChance(e *Root, n *Chance) error {
	if n.Next != nil {
		return nil
	}
	player := &Player{
		Parent:  n,
		TurnPos: n.State.TurnPos,
		State:   n.State,
	}
	n.Next = player
	atomic.AddUint32(&e.Nodes, 1)
	return nil
}

func ExpandPlayer(e *Root, n *Player) (err error) {
	if n.Actions != nil {
		return ExpandPlayerActions(e, n.Actions)
	}

	legal := table.NewDiscreteLegalActions(e.Params, n.State)

	actions := &PlayerActions{
		Parent:   n,
		Actions:  make([]table.DiscreteAction, len(legal)),
		Nodes:    make([]Node, len(legal)),
		Policies: NewStoreBacking(),
	}

	for i, a := range legal.List() {
		actions.Nodes[i], err = MakeNode(e, n, a)
		if err != nil {
			return err
		}
		actions.Actions[i] = a
		atomic.AddUint32(&e.Nodes, 1)
	}

	n.Actions = actions
	return nil
}

func ExpandPlayerActions(e *Root, n *PlayerActions) (err error) {
	for i, a := range n.Actions {
		if n.Nodes[i] != nil {
			continue
		}

		n.Nodes[i], err = MakeNode(e, n.Parent, a)
		if err != nil {
			return err
		}

		atomic.AddUint32(&e.Nodes, 1)
		i++
	}

	return nil
}

func ExpandFull(root *Root) error {
	stack := []Node{root}

	for len(stack) > 0 {
		// Pop item from stack
		current := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		// Expand current node
		err := Expand(root, current)
		if err != nil {
			return err
		}

		// Add children to stack based on node type
		switch n := current.(type) {
		case *Root:
			if n.Next != nil {
				stack = append(stack, n.Next)
			}

		case *Chance:
			if n.Next != nil {
				stack = append(stack, n.Next)
			}

		case *Player:
			if n.Actions != nil {
				stack = append(stack, n.Actions.Nodes...)
			}
		}
	}

	root.Full = true

	return nil
}

package tree

import (
	"errors"
	"fmt"

	"github.com/pokerdroid/poker/frand"
	"github.com/pokerdroid/poker/table"
)

// VisitFn is a callback used during tree traversal.
type VisitFn func(n Node, children []Node, depth int) bool

// Visit starts the traversal with a current depth of zero.
// The parameter maxDepth indicates the maximum level to traverse.
// If maxDepth is -1 the traversal never stops due to depth.
func Visit(n Node, maxDepth int, cb VisitFn) error {
	return visit(n, 0, maxDepth, cb)
}

func MustVisit(n Node, maxDepth int, cb VisitFn) {
	if err := Visit(n, maxDepth, cb); err != nil {
		panic(err)
	}
}

// visit recursively traverses the tree, increasing the current depth at each level.
func visit(n Node, curDepth, maxDepth int, cb VisitFn) error {
	if n == nil {
		return nil
	}

	// If a maximum depth is specified (i.e. non -1) and we've exceeded it, stop the recursion.
	if maxDepth != -1 && curDepth > maxDepth {
		return nil
	}

	var childs []Node
	switch x := n.(type) {
	case *Root:
		if x == nil {
			return nil
		}
		nodes := []Node{}
		if x.Next != nil {
			nodes = []Node{x.Next}
		}
		if cb(x, nodes, curDepth) {
			childs = append(childs, nodes...)
		}

	case *Chance:
		if x == nil {
			return nil
		}
		nodes := []Node{}
		if x.Next != nil {
			nodes = []Node{x.Next}
		}
		if cb(x, nodes, curDepth) {
			childs = append(childs, nodes...)
		}

	case *Terminal:
		if x == nil {
			return nil
		}
		cb(x, []Node{}, curDepth)

	case *Player:
		if x == nil {
			return nil
		}
		nodes := x.Children()
		if cb(x, nodes, curDepth) {
			childs = append(childs, nodes...)
		}

	case *Reference:
		if x == nil {
			return nil
		}
		r, err := x.Expand()
		if err != nil {
			return err
		}
		if cb(x, []Node{r}, curDepth) {
			childs = append(childs, r)
		}

	default:
		return fmt.Errorf("unknown node kind: %T", n)
	}

	// Recurse on each child while incrementing the current depth.
	for _, ch := range childs {
		if err := visit(ch, curDepth+1, maxDepth, cb); err != nil {
			return err
		}
	}

	return nil
}

func FindLeafNodes(n Node) (nn []Node) {
	MustVisit(n, -1, func(n Node, children []Node, depth int) bool {
		if len(children) == 0 {
			nn = append(nn, n)
		}
		return true
	})
	return
}

func CountNodes(n Node) (t uint32) {
	var pn uint32
	MustVisit(n, -1, func(n Node, children []Node, depth int) bool {
		pn += 1
		return true
	})
	return pn
}

func CountStates(n Node) (t uint32) {
	var pn uint32
	MustVisit(n, -1, func(n Node, children []Node, depth int) bool {
		p, ok := n.(*Player)
		if !ok {
			return true
		}

		if p.Actions == nil {
			return true
		}

		pn += uint32(p.Actions.Policies.Len())
		return true
	})
	return pn
}

func FindDecisionPoint(current Node) (Node, error) {
	if r, ok := current.(*Root); ok {
		return FindDecisionPoint(r.Next)
	}
	if c, ok := current.(*Chance); ok {
		return FindDecisionPoint(c.Next)
	}
	if c, ok := current.(*Reference); ok {
		r, err := c.Expand()
		if err != nil {
			return nil, err
		}
		return FindDecisionPoint(r)
	}
	// This is gonna be other player nodes such
	// as Terminal or Player nodes
	return current, nil
}

type Action struct {
	Action  table.DiscreteAction
	Parent  *Player
	Node    Node
	NodeIdx int
	State   *table.State
}

func ExtractActions(n Node) []Action {
	if n == nil {
		return nil
	}

	actions := make([]Action, 0, 16)
	cur := n

	for {
		parent := cur.GetParent()
		if parent == nil {
			break
		}

		if cur == parent {
			panic("nodes are self-referencing")
		}

		if px, ok := parent.(*Player); ok {
			idx, found := px.GetActionIdx(cur)
			if !found {
				cur = parent
				continue
			}

			ract := px.Actions.Actions[idx]

			actions = append(actions, Action{
				Action:  ract,
				Parent:  px,
				Node:    cur,
				NodeIdx: idx,
				State:   px.State,
			})
		}

		cur = parent
	}

	// Reverse actions to get chronological order
	for i, j := 0, len(actions)-1; i < j; i, j = i+1, j-1 {
		actions[i], actions[j] = actions[j], actions[i]
	}

	return actions
}

// SamplePath traverses the tree by sampling actions until reaching a terminal node.
// Returns the terminal node and the sequence of actions taken to reach it.
func SamplePath(r *Root, n Node, rng frand.Rand) (actions []Action, err error) {
	current := n
	actions = make([]Action, 0)

	for {
		err = Expand(r, current)
		if err != nil {
			return nil, err
		}

		switch x := current.(type) {
		case *Root:
			if x.Next == nil {
				return nil, errors.New("root has no next node")
			}
			current = x.Next

		case *Chance:
			if x.Next == nil {
				return nil, errors.New("chance has no next node")
			}
			current = x.Next

		case *Player:
			if x.Actions == nil || len(x.Actions.Actions) == 0 {
				return nil, errors.New("player has no actions")
			}

			// Sample an action index based on probabilities
			idx := rng.Intn(len(x.Actions.Actions))

			current = x.Actions.Nodes[idx]

			// Record the action taken
			actions = append(actions, Action{
				Action:  x.Actions.Actions[idx],
				Parent:  x,
				Node:    current,
				NodeIdx: idx,
				State:   x.State,
			})

		case *Terminal:
			return actions, nil

		case *Reference:
			// Expand reference node
			expanded, err := x.Expand()
			if err != nil {
				return nil, fmt.Errorf("failed to expand reference: %w", err)
			}
			current = expanded

		default:
			return nil, fmt.Errorf("unknown node type: %T", x)
		}
	}
}

const (
	boxVertical         = "│"
	boxHorizontal       = "─"
	boxVerticalAndRight = "├"
	boxUpAndRight       = "└"
)

// PrintTree prints a visual tree structure
// PrintTree prints a visual tree structure
func PrintTree(n Node, maxDepth int) {
	printTree(n, "", true, maxDepth, make(map[Node]bool))
}

func printTree(n Node, prefix string, isLast bool, maxDepth int, visited map[Node]bool) {
	if n == nil {
		return
	}

	if maxDepth != -1 && len(prefix)/4 > maxDepth {
		return
	}

	// Detect cycles
	if visited[n] {
		fmt.Printf("%s%s─ CYCLE: %T(%p)\n", prefix, boxUpAndRight, n, n)
		return
	}
	visited[n] = true

	// Print current node
	fmt.Printf("%s%s─ %T(%p)\n", prefix,
		map[bool]string{true: boxUpAndRight, false: boxVerticalAndRight}[isLast],
		n, n)

	// Child prefix
	childPrefix := prefix
	if isLast {
		childPrefix += "    "
	} else {
		childPrefix += boxVertical + "   "
	}

	// Print parent info
	parent := n.GetParent()
	if parent != nil {
		fmt.Printf("%s%s─ parent: %T(%p)\n", childPrefix, boxVerticalAndRight, parent, parent)
	}

	// Print actions for Player nodes
	if x, ok := n.(*Player); ok && x.Actions != nil {
		fmt.Printf("%s%s─ actions: %v\n", childPrefix, boxVerticalAndRight, x.Actions.Actions)
	}

	// Get children
	var children []Node
	switch x := n.(type) {
	case *Root:
		if x.Next != nil {
			children = []Node{x.Next}
		}
	case *Chance:
		if x.Next != nil {
			children = []Node{x.Next}
		}
	case *Player:
		if x.Actions != nil {
			children = x.Actions.Nodes
		}

	case *Reference:
		expanded, err := x.Expand()
		if err != nil {
			fmt.Printf("%s%s─ ERROR: %v\n", childPrefix, boxVerticalAndRight, err)
			return
		}
		children = []Node{expanded}
	}

	// Print children
	for i, child := range children {
		printTree(child, childPrefix, i == len(children)-1, maxDepth, visited)
	}
}

func GetState(n Node) *table.State {
	switch n := n.(type) {
	case *Chance:
		return n.State
	case *Root:
		return n.State
	case *Player:
		return n.State
	}
	return nil
}

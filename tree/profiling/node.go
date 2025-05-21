package profiling

import (
	"slices"

	"github.com/pokerdroid/poker/card"
	"github.com/pokerdroid/poker/float/f64"
	"github.com/pokerdroid/poker/policy"
	"github.com/pokerdroid/poker/table"
	"github.com/pokerdroid/poker/tree"
)

type Node struct {
	Runes    tree.Runes
	Cards    card.Cards
	Actions  []table.DiscreteAction
	Strategy []float64
	Policy   *policy.Policy
	Node     tree.Node
}

type Nodes []Node

type Matrix [13][13][]float64

type ActionMatrix struct {
	Actions []table.DiscreteAction
	Matrix  map[tree.Node]Matrix
}

func (n Nodes) Dist() ActionMatrix {
	matrix := make(map[tree.Node]Matrix)
	actions := n.Actions()

	for _, node := range n {
		x, y := card.Coordinates(node.Cards[:2])

		ax, ok := matrix[node.Node]
		if !ok {
			ax = Matrix{}
		}

		if ax[x][y] == nil {
			ax[x][y] = make([]float64, len(actions))
		}

		for i, a := range node.Actions {
			index := slices.Index(actions, a)
			if index == -1 {
				continue
			}
			ax[x][y][index] += node.Strategy[i]
		}

		matrix[node.Node] = ax
	}

	for n, nn := range matrix {
		for k, v := range nn {
			for k2, v2 := range v {
				f64.ScalUnitary(1.0/f64.Sum(v2), matrix[n][k][k2])
			}
		}
	}

	return ActionMatrix{Actions: actions, Matrix: matrix}
}

func (n Nodes) Actions() []table.DiscreteAction {
	var actions []table.DiscreteAction

	for _, node := range n {
		for _, a := range node.Actions {
			if slices.Contains(actions, a) {
				continue
			}
			actions = append(actions, a)
		}
	}

	slices.SortFunc(actions, func(a, b table.DiscreteAction) int {
		// Helper function to get sort priority by action type
		getPriority := func(act table.DiscreteAction) int {
			switch act {
			case table.DFold:
				return 1 // Fold first
			case table.DCheck:
				return 2 // Check second
			case table.DCall:
				return 3 // Call third
			case table.DAllIn:
				return 999 // All-in last
			default:
				// Positive values are raises, sort them by size
				if act > 0 {
					return 10 + int(act*100) // Between call and all-in, ordered by size
				}
				return 500 // Unknown actions
			}
		}

		// Compare priorities
		prioA := getPriority(a)
		prioB := getPriority(b)
		if prioA < prioB {
			return -1
		} else if prioA > prioB {
			return 1
		}
		return 0
	})

	return actions
}

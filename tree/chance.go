package tree

import (
	"github.com/pokerdroid/poker/table"
)

type Chance struct {
	Next   Node
	Parent Node
	State  *table.State
}

func (ch *Chance) Kind() NodeKind {
	return NodeKindChance
}

func (ch *Chance) GetParent() Node {
	return ch.Parent
}

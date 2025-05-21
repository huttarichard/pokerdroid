package tree

import (
	"github.com/pokerdroid/poker/table"
)

type Terminal struct {
	Parent  Node
	Pots    table.Pots
	Players table.Players
}

func (ch *Terminal) Kind() NodeKind {
	return NodeKindTerminal
}

func (ch *Terminal) GetParent() Node {
	return ch.Parent
}

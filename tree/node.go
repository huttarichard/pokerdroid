package tree

import (
	"encoding"
	"io"

	"github.com/pokerdroid/poker/abs"
	"github.com/pokerdroid/poker/policy"
)

type NodeKind uint8

const (
	NodeKindRoot NodeKind = iota
	NodeKindChance
	NodeKindTerminal
	NodeKindPlayer
	NodeKindRollout
)

// String returns string representation of NodeKind
func (k NodeKind) String() string {
	switch k {
	case NodeKindRoot:
		return "root"
	case NodeKindChance:
		return "chance"
	case NodeKindTerminal:
		return "terminal"
	case NodeKindPlayer:
		return "player"
	case NodeKindRollout:
		return "rollout"
	default:
		return "unknown"
	}
}

type Node interface {
	encoding.BinaryMarshaler
	encoding.BinaryUnmarshaler

	ReadBinary(r io.ReadSeeker) error
	WriteBinary(w io.Writer) error

	Kind() NodeKind
	GetParent() Node
	Size() uint64
}

type DecisionPoint interface {
	Node
	Len() int
	GetNode(i int) Node
	GetTurnPos() uint8
	Acquire(tree *Root, cluster abs.Cluster) *policy.Policy
	Get(cluster abs.Cluster) (*policy.Policy, bool)
	IsNil(i int) bool
	GetActionIdx(n Node) (int, bool)
}

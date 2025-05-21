package tree

import (
	"errors"
	"io"
	"sync"
)

// Reference is a lazy node loader that holds a file offset and length.
// It implements Node so that when its methods are called, it loads its full content.
type Reference struct {
	Parent  Node   // Parent node (if any)
	Offset  int64  // File offset where this node's data is stored
	Length  uint64 // Length in bytes of the node data
	Pointer io.ReadSeeker
	Node    Node
	Type    NodeKind

	mu sync.Mutex
}

// load reads the node data from the underlying reader if it has not been loaded yet.
func (r *Reference) Expand() (node Node, err error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.Node != nil {
		return r.Node, nil
	}

	// Save current position to restore later.
	cur, err := r.Pointer.Seek(0, io.SeekCurrent)
	if err != nil {
		return nil, err
	}
	defer r.Pointer.Seek(cur, io.SeekStart) // restore

	// Seek to the start of the node data.
	if _, err = r.Pointer.Seek(r.Offset, io.SeekStart); err != nil {
		return nil, err
	}

	// Create empty node of correct type
	switch r.Type {
	case NodeKindRoot:
		r.Node = &Root{}
	case NodeKindChance:
		r.Node = &Chance{Parent: r.Parent}
	case NodeKindPlayer:
		r.Node = &Player{Parent: r.Parent}
	case NodeKindTerminal:
		r.Node = &Terminal{Parent: r.Parent}
	default:
		return nil, errors.New("unknown node kind at reference")
	}

	// Unmarshal using your existing UnmarshalNodeBinary (which creates a full node).
	err = r.Node.ReadBinary(r.Pointer)
	if err != nil {
		return nil, err
	}

	return r.Node, nil
}

func (r *Reference) MustExpand() Node {
	n, err := r.Expand()
	if err != nil {
		panic(err)
	}
	return n
}

// Kind implements Node by lazily retrieving the node kind.
func (r *Reference) Kind() NodeKind {
	return r.Type
}

// GetParent returns the parent of the fully loaded node.
func (r *Reference) GetParent() Node {
	return r.Parent
}

func (r *Reference) Size() uint64 {
	// For a reference, we return the known Length field
	// since it represents the exact size of the referenced node
	return r.Length
}

// Replace replaces the node with a new node.
// TODO test this.
// func (r *Reference) Replace() error {
// 	n, err := r.Expand()
// 	if err != nil {
// 		return err
// 	}

// 	switch x := r.Parent.(type) {
// 	case *Root:
// 		x.Next = n
// 	case *Player:
// 		idx, ok := x.GetActionIdx(r.Node)
// 		if !ok {
// 			return errors.New("player action index not found")
// 		}
// 		x.Actions.Nodes[idx] = n
// 	case *Chance:
// 		x.Next = n
// 	default:
// 		return errors.New("unknown parent node type")
// 	}
// 	r = nil
// 	return nil
// }

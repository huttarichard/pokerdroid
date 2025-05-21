package tree

import (
	"errors"
	"sync/atomic"

	"github.com/pokerdroid/poker/abs"
	"github.com/pokerdroid/poker/policy"
	"github.com/pokerdroid/poker/table"
)

var NewStoreBacking = NewPolicies

type Player struct {
	Parent  Node
	TurnPos uint8
	State   *table.State
	Actions *PlayerActions
}

var _ Node = &Player{}
var _ DecisionPoint = &Player{}

func (ch *Player) Kind() NodeKind {
	return NodeKindPlayer
}

func (ch *Player) GetParent() Node {
	return ch.Parent
}

func (ch *Player) Len() int {
	if ch.Actions == nil {
		return 0
	}
	return len(ch.Actions.Actions)
}

func (ch *Player) GetTurnPos() uint8 {
	return ch.TurnPos
}

func (ch *Player) Acquire(r *Root, c abs.Cluster) *policy.Policy {
	return ch.Actions.Acquire(r, c)
}

func (ch *Player) Get(c abs.Cluster) (*policy.Policy, bool) {
	return ch.Actions.Get(c)
}

func (ch *Player) GetNode(i int) Node {
	return ch.Actions.Nodes[i]
}

func (ch *Player) IsNil(i int) bool {
	return ch.Actions == nil || ch.Actions.Nodes[i] == nil
}

func (ch *Player) Children() []Node {
	if ch.Actions == nil {
		return []Node{}
	}
	return ch.Actions.Nodes
}

func (p *Player) GetActionIdx(n Node) (int, bool) {
	for i, x := range p.Actions.Nodes {
		var realnode Node
		switch x := x.(type) {
		case *Reference:
			realnode = x.Node
		default:
			realnode = x
		}
		if realnode == n {
			return i, true
		}
	}
	return -1, false
}

func (p *Player) GetAction(a table.DiscreteAction) (Node, bool) {
	for i, x := range p.Actions.Actions {
		if x == a {
			return p.Actions.Nodes[i], true
		}
	}
	return nil, false
}

type PlayerActions struct {
	Parent   *Player
	Actions  []table.DiscreteAction
	Nodes    []Node
	Policies *Policies
}

func (p *PlayerActions) GetIdx(action table.DiscreteAction) int {
	for i, a := range p.Actions {
		if a == action {
			return i
		}
	}
	return -1
}

func (p *PlayerActions) Acquire(r *Root, c abs.Cluster) *policy.Policy {
	px, ok := p.Policies.Acquire(c, len(p.Actions))
	if ok {
		return px
	}
	atomic.AddUint32(&r.States, 1)
	return px
}

func (p *PlayerActions) Get(c abs.Cluster) (*policy.Policy, bool) {
	return p.Policies.Get(c)
}

func (p *PlayerActions) Validate() error {
	if len(p.Nodes) != len(p.Actions) {
		return errors.New("actions and nodes length mismatch")
	}

	if p.Nodes == nil {
		return errors.New("nodes is nil")
	}

	if p.Policies == nil {
		return errors.New("policies is nil")
	}

	if p.Actions == nil {
		return errors.New("actions is nil")
	}

	return nil
}

func (p *PlayerActions) AddAction(r *Root, action table.DiscreteAction) (int, error) {
	if p.Actions == nil {
		return -1, errors.New("player actions not initialized")
	}

	if r == nil {
		return -1, errors.New("root node is nil")
	}

	// Check if action already exists
	for i, a := range p.Actions {
		if a == action {
			return i, nil // Action already exists, return its index
		}
	}

	if p.Policies.Len() != 0 {
		return -1, errors.New("policies are initialized already")
	}

	// Find the right position to insert the new action (keep sorted)
	insertIdx := 0
	for i, a := range p.Actions {
		if action > a {
			insertIdx = i + 1
		}
	}

	// Create new node for the action
	newNode, err := MakeNode(r, p.Parent, action)
	if err != nil {
		return -1, err
	}

	// Insert the new action and node at the right position
	p.Actions = append(
		p.Actions[:insertIdx],
		append(
			[]table.DiscreteAction{action},
			p.Actions[insertIdx:]...,
		)...,
	)

	p.Nodes = append(
		p.Nodes[:insertIdx],
		append(
			[]Node{newNode},
			p.Nodes[insertIdx:]...,
		)...,
	)

	// Update node count in root
	atomic.AddUint32(&r.Nodes, 1)

	// Return the index of the newly added action
	return insertIdx, nil
}

func MakeNode(e *Root, n *Player, a table.DiscreteAction) (Node, error) {
	o, err := table.MakeAction(e.Params, n.State, a)
	if err != nil {
		return nil, err
	}

	s := o.Next()

	var nd Node

	switch table.Rule(e.Params, s) {
	case table.RuleShiftTurn:
		err = table.ShiftTurn(e.Params, s)
		if err != nil {
			return nil, err
		}

		nd = &Player{
			Parent:  n,
			TurnPos: s.TurnPos,
			State:   s,
		}

	case table.RuleShiftStreet:
		err = table.ShiftStreet(s)
		if err != nil {
			return nil, err
		}

		nd = &Chance{
			Parent: n,
			State:  s,
		}

	case table.RuleShiftStreetUntilEnd:
		nd = &Terminal{
			Parent:  n,
			Pots:    table.GetPots(s.Players),
			Players: s.Players,
		}

	case table.RuleFinish:
		nd = &Terminal{
			Parent:  n,
			Pots:    table.GetPots(s.Players),
			Players: s.Players,
		}

	default:
		return nil, errors.New("unknown rule")
	}

	return nd, nil
}

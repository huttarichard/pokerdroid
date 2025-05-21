package tree

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"sort"

	"github.com/pokerdroid/poker/abs"
	"github.com/pokerdroid/poker/encbin"
	"github.com/pokerdroid/poker/policy"
	"github.com/pokerdroid/poker/table"
)

func UnmarshalNodeBinary(data []byte, parent Node) (node Node, err error) {
	k := NodeKind(data[0])

	switch k {
	case NodeKindChance:
		n := new(Chance)
		n.Parent = parent
		node = n
	case NodeKindPlayer:
		n := new(Player)
		n.Parent = parent
		node = n
	case NodeKindTerminal:
		n := new(Terminal)
		n.Parent = parent
		node = n
	case NodeKindRoot:
		n := new(Root)
		node = n
	default:
		return nil, errors.New("unknown type")
	}

	return node, node.UnmarshalBinary(data)
}

func MarshalNodeBinaryWithLen(buf io.Writer, node Node) error {
	return encbin.MarshalWithLen[uint64](buf, node)
}

func UnmarshalNodeBinaryWithLen(buf io.Reader, parent Node) (node Node, err error) {
	var length uint64
	err = encbin.UnmarshalValues(buf, &length)
	if err != nil {
		return nil, err
	}

	// Check for nil marker before trying to make slice
	if length == ^uint64(0) { // maxValue for uint64
		return nil, nil
	}

	if length == 0 {
		return nil, nil
	}

	data := make([]byte, length)
	_, err = io.ReadFull(buf, data)
	if err != nil {
		return nil, err
	}

	node, err = UnmarshalNodeBinary(data, parent)
	if err != nil {
		return nil, err
	}

	return node, nil
}

// UnmarshalNodeBinaryWithRef reads the length header and returns a lazy-loading Reference node.
// Unlike traditional UnmarshalNodeBinary it does not load the full node immediately.
func UnmarshalNodeBinaryWithRef(rs io.ReadSeeker, parent Node) (Node, error) {
	// Read the length of the node data (stored as uint64).
	var length uint64
	if err := encbin.UnmarshalValues(rs, &length); err != nil {
		return nil, err
	}

	// Check for nil marker or zero length.
	if length == ^uint64(0) || length == 0 {
		return nil, nil
	}

	// Capture the current offset (this is where the node blob begins).
	offset, err := rs.Seek(0, io.SeekCurrent)
	if err != nil {
		return nil, err
	}

	// Read first byte to determine node kind
	data := make([]byte, 1)
	if _, err := rs.Read(data); err != nil {
		return nil, err
	}
	k := NodeKind(data[0])

	// Seek back to start of node data
	if _, err := rs.Seek(offset, io.SeekStart); err != nil {
		return nil, err
	}

	// Optionally, you can read one byte (peek) to determine the node kind
	// before creating the reference. Here we simply create the reference.
	ref := &Reference{
		Parent:  parent,
		Offset:  offset,
		Length:  length,
		Pointer: rs,
		Type:    k,
	}

	// Advance the read pointer beyond this node's blob so that subsequent nodes can be processed.
	if _, err := rs.Seek(int64(length), io.SeekCurrent); err != nil {
		return nil, err
	}

	return ref, nil
}

func (p *PlayerActions) Size() uint64 {
	size := uint64(0)

	// Policies with length prefix
	size += 4 // Just the length prefix for nil
	if p.Policies != nil {
		size += p.Policies.Size()
	}

	// Number of actions/nodes
	size += 1 // uint8 for nodes count

	// For each action/node pair
	for i := range p.Actions {
		size += 4 // DiscreteAction (float32)
		size += 8 // uint64 length prefix

		// Node with length prefix
		if p.Nodes[i] != nil {
			size += p.Nodes[i].Size()
		}
	}

	return size
}

func (p PlayerActions) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)

	err := p.Validate()
	if err != nil {
		return nil, err
	}

	err = encbin.MarshalWithLen[uint32](buf, p.Policies)
	if err != nil {
		return nil, err
	}

	err = encbin.MarshalValues(buf, uint8(len(p.Actions)))
	if err != nil {
		return nil, err
	}

	for i := range p.Actions {
		err = encbin.MarshalValues(buf, p.Actions[i])
		if err != nil {
			return nil, err
		}
		err = MarshalNodeBinaryWithLen(buf, p.Nodes[i])
		if err != nil {
			return nil, err
		}
	}

	return buf.Bytes(), nil
}

func (p *PlayerActions) WriteBinary(w io.Writer) error {
	err := p.Validate()
	if err != nil {
		return err
	}

	err = encbin.MarshalWithLen[uint32](w, p.Policies)
	if err != nil {
		return err
	}

	err = encbin.MarshalValues(w, uint8(len(p.Actions)))
	if err != nil {
		return err
	}

	for i := range p.Actions {
		err = encbin.MarshalValues(w, p.Actions[i])
		if err != nil {
			return err
		}
		// Use WriteWithLen for direct streaming
		err = encbin.WriteWithLen(w, p.Nodes[i])
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *PlayerActions) UnmarshalBinary(data []byte) error {
	r := bytes.NewReader(data)

	pol := NewStoreBacking()
	ok, err := encbin.UnmarshalWithLenNil[uint32](r, pol)
	if err != nil {
		return err
	}
	if ok {
		p.Policies = pol
	}

	var nodes uint8
	err = encbin.UnmarshalValues(r, &nodes)
	if err != nil {
		return err
	}

	p.Actions = make([]table.DiscreteAction, nodes)
	p.Nodes = make([]Node, nodes)

	for i := 0; i < int(nodes); i++ {
		// Unmarshal action  using encbin utility
		err = encbin.UnmarshalValues(r, &p.Actions[i])
		if err != nil {
			return err
		}

		p.Nodes[i], err = UnmarshalNodeBinaryWithLen(r, p.Parent)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *PlayerActions) ReadBinary(rs io.ReadSeeker) error {
	pol := NewStoreBacking()
	ok, err := encbin.UnmarshalWithLenNil[uint32](rs, pol)
	if err != nil {
		return err
	}
	if ok {
		p.Policies = pol
	}

	var nodes uint8
	err = encbin.UnmarshalValues(rs, &nodes)
	if err != nil {
		return err
	}

	p.Actions = make([]table.DiscreteAction, nodes)
	p.Nodes = make([]Node, nodes)

	for i := 0; i < int(nodes); i++ {
		// Unmarshal action  using encbin utility
		err = encbin.UnmarshalValues(rs, &p.Actions[i])
		if err != nil {
			return err
		}

		p.Nodes[i], err = UnmarshalNodeBinaryWithRef(rs, p.Parent)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Chance) Size() uint64 {
	size := uint64(1) // NodeKind byte

	// State with length prefix
	size += 2 // Just the length prefix for nil
	if c.State != nil {
		size += c.State.Size()
	}

	// Next node with length prefix
	size += 8 // Just the length prefix for nil
	if c.Next != nil {
		size += c.Next.Size()
	}

	return size
}

func (c Chance) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)
	buf.WriteByte(byte(NodeKindChance))

	err := encbin.MarshalWithLen[uint16](buf, c.State)
	if err != nil {
		return nil, err
	}

	err = MarshalNodeBinaryWithLen(buf, c.Next)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (c *Chance) UnmarshalBinary(data []byte) error {
	if data[0] != byte(NodeKindChance) {
		return errors.New("invalid node kind")
	}

	r := bytes.NewReader(data[1:])

	state := new(table.State)
	ok, err := encbin.UnmarshalWithLenNil[uint16](r, state)
	if err != nil {
		return err
	}
	if ok {
		c.State = state
	}

	c.Next, err = UnmarshalNodeBinaryWithLen(r, c)
	if err != nil {
		return err
	}

	return nil
}

// WriteBinary writes the Chance node directly to an io.Writer.
func (c *Chance) WriteBinary(w io.Writer) error {
	_, err := w.Write([]byte{byte(NodeKindChance)})
	if err != nil {
		return err
	}

	// Write State with length prefix
	err = encbin.MarshalWithLen[uint16](w, c.State)
	if err != nil {
		return err
	}

	// Use WriteWithLen for the next node to avoid deep marshaling
	err = encbin.WriteWithLen(w, c.Next)
	if err != nil {
		return err
	}

	return nil
}

// ReadBinary implements the Node interface for Chance, reading from a ReadSeeker.
func (c *Chance) ReadBinary(rs io.ReadSeeker) error {
	var b [1]byte
	// Read the node kind byte
	if _, err := rs.Read(b[:]); err != nil {
		return err
	}
	if b[0] != byte(NodeKindChance) {
		return errors.New("Chance.ReadBinary: invalid node kind")
	}

	// Read the State field using encbin.UnmarshalWithLenNil.
	state := new(table.State)
	ok, err := encbin.UnmarshalWithLenNil[uint16](rs, state)
	if err != nil {
		return err
	}
	if ok {
		c.State = state
	}

	// Read the Next node via the helper which reads a length header.
	c.Next, err = UnmarshalNodeBinaryWithRef(rs, c)
	if err != nil {
		return err
	}

	return nil
}

func (p *Player) Size() uint64 {
	size := uint64(1) // NodeKind byte
	size += 1         // TurnPos (uint8)

	// State with length prefix
	if p.State != nil {
		size += p.State.Size() + 2 // uint16 length prefix
	} else {
		size += 2 // Just the length prefix for nil
	}

	// Actions with length prefix
	size += 8 // Just the length prefix for nil
	if p.Actions != nil {
		size += p.Actions.Size()
	}

	return size
}

func (p *Player) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)
	buf.WriteByte(byte(NodeKindPlayer))

	err := encbin.MarshalValues(buf, p.TurnPos)
	if err != nil {
		return nil, err
	}

	err = encbin.MarshalWithLen[uint16](buf, p.State)
	if err != nil {
		return nil, err
	}

	err = encbin.MarshalWithLen[uint64](buf, p.Actions)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// WriteBinary writes the Player node directly to an io.Writer.
func (p *Player) WriteBinary(w io.Writer) error {
	_, err := w.Write([]byte{byte(NodeKindPlayer)})
	if err != nil {
		return err
	}

	err = encbin.MarshalValues(w, p.TurnPos)
	if err != nil {
		return err
	}

	err = encbin.MarshalWithLen[uint16](w, p.State)
	if err != nil {
		return err
	}

	// Use WriteWithLen for Actions to avoid deep marshaling
	err = encbin.WriteWithLen(w, p.Actions)
	if err != nil {
		return err
	}

	return nil
}

func (p *Player) UnmarshalBinary(data []byte) error {
	if data[0] != byte(NodeKindPlayer) {
		return errors.New("invalid node kind")
	}

	r := bytes.NewReader(data[1:])

	err := encbin.UnmarshalValues(r, &p.TurnPos)
	if err != nil {
		return err
	}

	state := new(table.State)
	ok, err := encbin.UnmarshalWithLenNil[uint16](r, state)
	if err != nil {
		return err
	}
	if ok {
		p.State = state
	}

	actions := new(PlayerActions)
	actions.Parent = p
	ok, err = encbin.UnmarshalWithLenNil[uint64](r, actions)
	if err != nil {
		return err
	}
	if !ok {
		return nil
	}

	p.Actions = actions
	return nil
}

// ReadBinary implements the Node interface for Player, reading directly from the ReadSeeker.
func (p *Player) ReadBinary(rs io.ReadSeeker) error {
	var b [1]byte
	// Read the node kind byte.
	if _, err := rs.Read(b[:]); err != nil {
		return err
	}
	if b[0] != byte(NodeKindPlayer) {
		return errors.New("Player.ReadBinary: invalid node kind")
	}

	// Read TurnPos.
	if err := encbin.UnmarshalValues(rs, &p.TurnPos); err != nil {
		return err
	}

	// Read the State field.
	state := new(table.State)
	ok, err := encbin.UnmarshalWithLenNil[uint16](rs, state)
	if err != nil {
		return err
	}
	if ok {
		p.State = state
	}

	var length uint64
	if err := encbin.UnmarshalValues(rs, &length); err != nil {
		return err
	}

	if length == ^uint64(0) || length == 0 {
		return nil
	}

	// Read the Actions field.
	actions := new(PlayerActions)
	actions.Parent = p

	err = actions.ReadBinary(rs)
	if err != nil {
		return err
	}

	p.Actions = actions

	return nil
}

// Implement Size() for PolicyMap
func (p *Policies) Size() uint64 {
	p.mux.RLock()
	defer p.mux.RUnlock()

	size := uint64(8) // Initial length (uint64)

	// For each cluster/policy pair
	for _, pol := range p.Map {
		// Each entry in the map is marshaled as:
		size += 4 // Cluster (uint32)
		size += 2 // Policy length prefix (uint16)
		if pol != nil {
			size += pol.Size() // Policy data
		}
	}

	return size
}

func (p *Policies) MarshalBinary() ([]byte, error) {
	type clupol struct {
		cl  abs.Cluster
		pol *policy.Policy
	}

	// Snapshot the map under read lock.
	p.mux.RLock()
	entries := make([]clupol, 0, len(p.Map))
	for cl, pol := range p.Map {
		entries = append(entries, clupol{cl: cl, pol: pol})
	}
	p.mux.RUnlock()

	// Ensure deterministic order.
	sort.Slice(entries, func(i, j int) bool { return entries[i].cl < entries[j].cl })

	// Create buffer and write number of entries.
	buf := new(bytes.Buffer)
	if err := encbin.MarshalValues(buf, uint64(len(entries))); err != nil {
		return nil, err
	}

	// Marshal entries sequentially.
	for _, e := range entries {
		if err := encbin.MarshalValues(buf, e.cl); err != nil {
			return nil, err
		}
		if err := encbin.MarshalWithLen[uint16](buf, e.pol); err != nil {
			return nil, err
		}
	}

	return buf.Bytes(), nil
}

// WriteBinary writes the Policies map directly to an io.Writer.
func (p *Policies) WriteBinary(w io.Writer) error {
	type clupol struct {
		cl  abs.Cluster
		pol *policy.Policy
	}

	// Snapshot the map under read lock.
	p.mux.RLock()
	entries := make([]clupol, 0, len(p.Map))
	for cl, pol := range p.Map {
		entries = append(entries, clupol{cl: cl, pol: pol})
	}
	p.mux.RUnlock()

	// Ensure deterministic order.
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].cl < entries[j].cl
	})

	// Number of entries
	if err := encbin.MarshalValues(w, uint64(len(entries))); err != nil {
		return err
	}

	// Write each entry sequentially.
	for _, e := range entries {
		if err := encbin.MarshalValues(w, e.cl); err != nil {
			return err
		}
		if err := encbin.MarshalWithLen[uint16](w, e.pol); err != nil {
			return err
		}
	}

	return nil
}

func (p *Policies) UnmarshalBinary(data []byte) error {
	r := bytes.NewReader(data)

	var length uint64

	err := encbin.UnmarshalValues(r, &length)
	if err != nil {
		return err
	}

	for x := uint64(0); x < length; x++ {
		var cluster abs.Cluster

		err := encbin.UnmarshalValues(r, &cluster)
		if err != nil {
			return err
		}

		policy := &policy.Policy{}

		err = encbin.UnmarshalWithLen[uint16](r, policy)
		if err != nil {
			return err
		}

		p.Map[cluster] = policy
	}
	return nil
}

// MarshalBinary returns the binary serialization of the underlying node.
// It forces expansion if the node has not yet been loaded.
func (r *Reference) MarshalBinary() ([]byte, error) {
	n, err := r.Expand()
	if err != nil {
		return nil, err
	}
	return n.MarshalBinary()
}

// WriteBinary writes the underlying node's data to the writer.
// It forces expansion if the reference hasn't been loaded.
func (r *Reference) WriteBinary(w io.Writer) error {
	return errors.New("reference should not be written")
}

// UnmarshalBinary loads the underlying node from the given data and sets it in the reference.
// This method forces the reference into the "expanded" state.
func (r *Reference) UnmarshalBinary(data []byte) error {
	return errors.New("unmarshaling into reference")
}

// ReadBinary implements the Node interface for Reference.
// It triggers expansion so that subsequent calls use the fully loaded node.
func (r *Reference) ReadBinary(rs io.ReadSeeker) error {
	return errors.New("unmarshaling into reference")
}

func (r Root) Size() uint64 {
	size := uint64(1) // NodeKind byte

	// Fixed fields
	size += 16                  // AbsID (uuid)
	size += 4                   // States (uint32)
	size += 4                   // Nodes (uint32)
	size += 8                   // Iteration (uint64)
	size += r.Params.Size() + 8 // GameParams + length prefix

	// State with length prefix
	size += 2 // Just the length prefix for nil
	if r.State != nil {
		size += r.State.Size()
	}

	// Next node with length prefix
	size += 8 // Just the length prefix for nil
	if r.Next != nil {
		size += r.Next.Size()
	}

	return size
}

func (r Root) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)
	buf.WriteByte(byte(NodeKindRoot))

	err := encbin.MarshalValues(
		buf,
		r.AbsID,
		r.States,
		r.Nodes,
		r.Iteration,
	)
	if err != nil {
		return nil, err
	}

	err = encbin.MarshalWithLen[uint64](buf, r.Params)
	if err != nil {
		return nil, err
	}

	err = encbin.MarshalWithLen[uint16](buf, r.State)
	if err != nil {
		return nil, err
	}

	err = MarshalNodeBinaryWithLen(buf, r.Next)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// WriteBinary writes the Root node directly to an io.Writer.
func (r *Root) WriteBinary(w io.Writer) error {
	_, err := w.Write([]byte{byte(NodeKindRoot)})
	if err != nil {
		return err
	}

	err = encbin.MarshalValues(
		w,
		r.AbsID,
		r.States,
		r.Nodes,
		r.Iteration,
	)
	if err != nil {
		return err
	}

	err = encbin.MarshalWithLen[uint64](w, r.Params)
	if err != nil {
		return err
	}

	err = encbin.MarshalWithLen[uint16](w, r.State)
	if err != nil {
		return err
	}

	// Use WriteWithLen for the next node to avoid deep marshaling
	err = encbin.WriteWithLen(w, r.Next)
	if err != nil {
		return err
	}

	return nil
}

func (r *Root) UnmarshalBinary(data []byte) error {
	if data[0] != byte(NodeKindRoot) {
		return errors.New("invalid node kind")
	}

	buf := bytes.NewReader(data[1:])
	var err error

	err = encbin.UnmarshalValues(
		buf,
		&r.AbsID,
		&r.States,
		&r.Nodes,
		&r.Iteration,
	)
	if err != nil {
		return err
	}

	err = encbin.UnmarshalWithLen[uint64](buf, &r.Params)
	if err != nil {
		return err
	}

	state := new(table.State)
	ok, err := encbin.UnmarshalWithLenNil[uint16](buf, state)
	if err != nil {
		return err
	}
	if ok {
		r.State = state
	}

	r.Next, err = UnmarshalNodeBinaryWithLen(buf, r)
	if err != nil {
		return err
	}

	return nil
}

func (r *Root) ReadBinary(rs io.ReadSeeker) error {
	var b [1]byte
	// Read the node kind byte.
	if _, err := rs.Read(b[:]); err != nil {
		return err
	}
	if b[0] != byte(NodeKindRoot) {
		return errors.New("Root.ReadBinary: invalid node kind")
	}

	err := encbin.UnmarshalValues(
		rs,
		&r.AbsID,
		&r.States,
		&r.Nodes,
		&r.Iteration,
	)
	if err != nil {
		return err
	}

	err = encbin.UnmarshalWithLen[uint64](rs, &r.Params)
	if err != nil {
		return err
	}

	state := new(table.State)
	ok, err := encbin.UnmarshalWithLenNil[uint16](rs, state)
	if err != nil {
		return err
	}
	if ok {
		r.State = state
	}

	r.Next, err = UnmarshalNodeBinaryWithRef(rs, r)
	if err != nil {
		return err
	}

	return nil
}

func (t *Terminal) Size() uint64 {
	size := uint64(1) // NodeKind byte

	// Pots with length prefix
	size += 2 // uint16 length prefix
	size += t.Pots.Size()

	// Players slice with length prefix
	size += 1 // uint8 length prefix
	for _, player := range t.Players {
		size += player.Size()
	}

	return size
}

func (r *Terminal) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)
	buf.WriteByte(byte(NodeKindTerminal))

	err := encbin.MarshalWithLen[uint16](buf, r.Pots)
	if err != nil {
		return nil, err
	}

	err = encbin.MarshalSliceLen[table.Player, uint8](buf, r.Players)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// WriteBinary writes the Terminal node directly to an io.Writer.
func (t *Terminal) WriteBinary(w io.Writer) error {
	_, err := w.Write([]byte{byte(NodeKindTerminal)})
	if err != nil {
		return err
	}

	// Write Pots with length prefix
	err = encbin.MarshalWithLen[uint16](w, t.Pots)
	if err != nil {
		return err
	}

	// Write Players slice with length prefix
	err = encbin.MarshalSliceLen[table.Player, uint8](w, t.Players)
	if err != nil {
		return err
	}

	return nil
}

func (r *Terminal) UnmarshalBinary(data []byte) error {
	if data[0] != byte(NodeKindTerminal) {
		return errors.New("invalid node kind")
	}

	buf := bytes.NewReader(data[1:])

	err := encbin.UnmarshalWithLen[uint16](buf, &r.Pots)
	if err != nil {
		return err
	}

	players, err := encbin.UnmarhsalSliceLen[table.Player, uint8](buf)
	if err != nil {
		return err
	}
	r.Players = players

	return nil
}

// ReadBinary implements the Node interface for Terminal, reading its content from the ReadSeeker.
func (t *Terminal) ReadBinary(rs io.ReadSeeker) error {
	var b [1]byte
	// Read the node kind byte.
	if _, err := rs.Read(b[:]); err != nil {
		return err
	}
	if b[0] != byte(NodeKindTerminal) {
		return errors.New("Terminal.ReadBinary: invalid node kind")
	}

	// Read the Pots field.
	if err := encbin.UnmarshalWithLen[uint16](rs, &t.Pots); err != nil {
		return err
	}

	// Read the Players slice.
	players, err := encbin.UnmarhsalSliceLen[table.Player, uint8](rs)
	if err != nil {
		return err
	}
	t.Players = players

	return nil
}

// MarshalJSON implements json.Marshaler
func (k NodeKind) MarshalJSON() ([]byte, error) {
	var s string
	switch k {
	case NodeKindRoot:
		s = "root"
	case NodeKindChance:
		s = "chance"
	case NodeKindPlayer:
		s = "player"
	case NodeKindTerminal:
		s = "terminal"
	case NodeKindRollout:
		s = "rollout"
	default:
		return nil, errors.New("unknown node kind")
	}
	return json.Marshal(s)
}

// UnmarshalJSON implements json.Unmarshaler
func (k *NodeKind) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	switch s {
	case "root":
		*k = NodeKindRoot
	case "chance":
		*k = NodeKindChance
	case "player":
		*k = NodeKindPlayer
	case "terminal":
		*k = NodeKindTerminal
	case "rollout":
		*k = NodeKindRollout
	default:
		return fmt.Errorf("invalid node kind: %s", s)
	}
	return nil
}

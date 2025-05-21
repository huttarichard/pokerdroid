package table

import (
	"bytes"
	"fmt"

	"github.com/pokerdroid/poker/encbin"
)

// MarshalBinary serializes a Game (its GameParams and the state chain)
// into binary form. It first writes the game parameters (using
// encbin.MarshalWithLen[uint16]), then writes the number of states as a uint16,
// and then for each state (ordered from the root to the final state) writes
// the state using encbin.MarshalWithLen[uint16].
func MarshalBinary(p GameParams, s *State) ([]byte, error) {
	var buf bytes.Buffer

	// First encode GameParams using the helper.
	err := encbin.MarshalWithLen[uint16](&buf, &p)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal game params: %w", err)
	}

	// Collect the chain of states from the final state back to the root.
	var states []*State
	for st := s; st != nil; st = st.Previous {
		states = append(states, st)
	}

	// Reverse the slice so that the earliest (root) state is first.
	for i, j := 0, len(states)-1; i < j; i, j = i+1, j-1 {
		states[i], states[j] = states[j], states[i]
	}

	// Write the number of states as a header (using uint16).
	count := uint16(len(states))

	err = encbin.MarshalValues(&buf, count)
	if err != nil {
		return nil, fmt.Errorf("failed to write state count: %w", err)
	}

	// For each state in the chain, marshal it with a length prefix.
	for i, st := range states {
		err := encbin.MarshalWithLen[uint16](&buf, st)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal state %d: %w", i, err)
		}
	}

	return buf.Bytes(), nil
}

// UnmarshalBinary deserializes the binary data into a Game instance.
// It first reads the game parameters (using encbin.UnmarshalWithLen[uint16]),
// then reads a uint16 that holds the number of states that follow,
// unmarshals each state (using encbin.UnmarshalWithLen[uint16]) and re‑links the
// chain via the Previous field. The final state becomes the Latest state in the Game.
func UnmarshalBinary(data []byte) (params GameParams, state *State, err error) {
	buf := bytes.NewReader(data)

	// Unmarshal game parameters.
	err = encbin.UnmarshalWithLen[uint16](buf, &params)
	if err != nil {
		err = fmt.Errorf("failed to unmarshal game params: %w", err)
		return
	}

	// Read the count of states.
	var count uint16
	err = encbin.UnmarshalValues(buf, &count)
	if err != nil {
		err = fmt.Errorf("failed to read state count: %w", err)
		return
	}

	if count == 0 {
		err = fmt.Errorf("no states in chain")
		return
	}

	// Read and unmarshal each state.
	states := make([]*State, count)
	for i := 0; i < int(count); i++ {
		st := &State{}
		err = encbin.UnmarshalWithLen[uint16](buf, st)
		if err != nil {
			err = fmt.Errorf("failed to unmarshal state %d: %w", i, err)
			return
		}
		states[i] = st
	}

	// Re-link the state chain: for each state (except the first),
	// set its Previous pointer to the prior state.
	for i := 1; i < int(count); i++ {
		states[i].Previous = states[i-1]
	}

	return params, states[count-1], nil
}

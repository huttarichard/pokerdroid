package bot

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/pokerdroid/poker"
	"github.com/pokerdroid/poker/card"
	"github.com/pokerdroid/poker/encbin"
	"github.com/pokerdroid/poker/table"
)

type Advisor interface {
	Advise(ctx context.Context, loggr poker.Logger, r State) (table.DiscreteAction, error)
}

type State struct {
	Params    table.GameParams `json:"params"`
	State     *table.State     `json:"state"`
	Hole      card.Cards       `json:"hole"`
	Community card.Cards       `json:"community"`
}

func (s State) Validate() error {
	if len(s.Hole) != 2 {
		return fmt.Errorf("invalid hole cards: %v", s.Hole)
	}

	if s.Hole[0] == card.Card00 || s.Hole[1] == card.Card00 {
		return fmt.Errorf("invalid hole cards: %v", s.Hole)
	}

	if len(s.Community) > 5 {
		return fmt.Errorf("invalid community cards: %v", s.Community)
	}

	return nil
}

func (s State) String() string {
	p := s.State.Path(s.Params.SbAmount)
	return fmt.Sprintf("Hole: %s, Community: %s, Path: %s", s.Hole, s.Community, p)
}

func (s State) MarshalBinary() ([]byte, error) {
	if err := s.Validate(); err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer(nil)

	data, err := table.MarshalBinary(s.Params, s.State)
	if err != nil {
		return nil, err
	}

	err = encbin.MarshalValues(buf, uint32(len(data)))
	if err != nil {
		return nil, err
	}

	_, err = buf.Write(data)
	if err != nil {
		return nil, err
	}

	_, err = buf.Write(s.Hole.Bytes())
	if err != nil {
		return nil, err
	}

	_, err = buf.Write(s.Community.Bytes())
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (s *State) UnmarshalBinary(data []byte) error {
	buf := bytes.NewBuffer(data)

	var size uint32
	err := encbin.UnmarshalValues(buf, &size)
	if err != nil {
		return err
	}

	res := make([]byte, size)
	_, err = io.ReadFull(buf, res)
	if err != nil {
		return err
	}

	s.Params, s.State, err = table.UnmarshalBinary(res)
	if err != nil {
		return err
	}

	hole := make([]byte, 2)
	_, err = io.ReadFull(buf, hole)
	if err != nil {
		return err
	}
	s.Hole = card.NewCardsFromBytes(hole)

	com := make([]byte, 5)
	n, err := io.ReadFull(buf, com)
	switch err {
	case nil:
	case io.EOF:
	case io.ErrUnexpectedEOF:
	default:
		return err
	}

	s.Community = card.NewCardsFromBytes(com[:n])
	return nil
}

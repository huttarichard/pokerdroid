package table

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/pokerdroid/poker/chips"
	"github.com/pokerdroid/poker/encbin"
)

// State represents the state of the game.
// Only params that are absolutely necessary to make
// next move are included in this.
type State struct {
	Players      Players `json:"players"`
	Street       Street  `json:"street"`
	TurnPos      uint8   `json:"turn_pos"`
	BtnPos       uint8   `json:"btn_pos"`
	BetAction    uint8   `json:"bet_action"`
	StreetAction uint8   `json:"street_action"`
	// todo make sure its max with player stack
	CallAmount chips.Chips `json:"call_amount"`
	// Biggest Street Commitment chips.Chips +
	// if it was raise or bet
	BSC struct {
		// If Action is Raise this is addition
		// If Action is Bet this is amount
		Amount   chips.Chips `json:"amount"`
		Addition chips.Chips `json:"addition"`
		Action   ActionKind  `json:"action"`
	} `json:"bsc"`

	// These are per street stats
	// they will reset on new street

	// Per Player Street Commitment
	PSC chips.List `json:"psc"`
	// Per Player Street Action Counter
	PSAC []uint8 `json:"psac"`
	// Per Player street last action
	PSLA []ActionKind `json:"psla"`

	// This is for debug purposes
	// should not be serialized or accounted for
	Previous *State `json:"-"`
}

func NewState(params GameParams) *State {
	players := make(Players, params.NumPlayers)
	for i := range players {
		players[i] = NewPlayer()
	}

	if params.MaxActionsPerRound == 0 {
		params.MaxActionsPerRound = uint8(params.NumPlayers * 3)
	}

	if params.TerminalStreet == NoStreet {
		params.TerminalStreet = River
	}

	return &State{
		Players:      players,
		Street:       Preflop,
		TurnPos:      params.BtnPos,
		BtnPos:       params.BtnPos,
		CallAmount:   chips.Zero,
		BetAction:    0,
		StreetAction: 0,
		PSC:          chips.NewListAlloc(params.NumPlayers),
		PSAC:         make([]uint8, params.NumPlayers),
		PSLA:         make([]ActionKind, params.NumPlayers),
		Previous:     nil,
	}
}

func (r *State) Next() *State {
	next := &State{
		Players:      r.Players.Clone(),
		Street:       r.Street,
		TurnPos:      r.TurnPos,
		BtnPos:       r.BtnPos,
		CallAmount:   r.CallAmount,
		BetAction:    r.BetAction,
		BSC:          r.BSC,
		PSC:          r.PSC.Copy(),
		StreetAction: r.StreetAction,
		PSAC:         make([]uint8, len(r.PSAC)),
		PSLA:         make([]ActionKind, len(r.PSLA)),
		Previous:     r,
	}
	copy(next.PSAC, r.PSAC)
	copy(next.PSLA, r.PSLA)
	return next
}

func (r *State) Finished() bool {
	return r.Street == Finished
}

func (r *State) String() string {
	var str strings.Builder
	w := str.WriteString
	s := fmt.Sprintf

	w(s("Street: %s\n", r.Street))
	if r.BSC.Action != NoAction {
		w(s("Biggest Street Commitment: %s %s\n", r.BSC.Amount, r.BSC.Action))
	}

	for i, p := range r.Players {
		w(s("Player %d", i))
		if i == int(r.TurnPos) {
			w(s(" - on turn"))
		}
		if i == int(r.BtnPos) {
			w(s(" - on btn"))
		}
		w(":\n")
		if i == int(r.TurnPos) {
			w(s("\tHas to call: %s\n", r.CallAmount.String()))
		}
		w(s("\tPaid: %s\n", p.Paid.String()))
		w(s("\tStatus: %s\n", p.Status.String()))
		w(s("\tCommited on street: %s\n", r.PSC[i].String()))
		w(s("\tAction count on street: %d\n", r.PSAC[i]))
		w(s("\tLast action on street: %s\n", r.PSLA[i].String()))
	}

	return str.String()
}

// Equal returns true if two states are equal, comparing all fields deeply
func (s *State) Equal(other *State) bool {
	if s == nil && other == nil {
		return true
	}
	if s == nil || other == nil {
		return false
	}

	// Compare basic fields
	if s.Street != other.Street ||
		s.TurnPos != other.TurnPos ||
		s.BtnPos != other.BtnPos ||
		s.BetAction != other.BetAction ||
		s.StreetAction != other.StreetAction {
		return false
	}

	// Compare chips
	if !s.CallAmount.Equal(other.CallAmount) {
		return false
	}

	// Compare BSC
	if !s.BSC.Amount.Equal(other.BSC.Amount) ||
		!s.BSC.Addition.Equal(other.BSC.Addition) ||
		s.BSC.Action != other.BSC.Action {
		return false
	}

	// Compare Players
	if len(s.Players) != len(other.Players) {
		return false
	}
	for i := range s.Players {
		if !s.Players[i].Paid.Equal(other.Players[i].Paid) ||
			s.Players[i].Status != other.Players[i].Status {
			return false
		}
	}

	// Compare PSC (Player Street Commitment)
	if len(s.PSC) != len(other.PSC) {
		return false
	}
	for i := range s.PSC {
		if !s.PSC[i].Equal(other.PSC[i]) {
			return false
		}
	}

	// Compare PSAC (Player Street Action Count)
	if len(s.PSAC) != len(other.PSAC) {
		return false
	}
	for i := range s.PSAC {
		if s.PSAC[i] != other.PSAC[i] {
			return false
		}
	}

	// Compare PSLA (Player Street Last Action)
	if len(s.PSLA) != len(other.PSLA) {
		return false
	}
	for i := range s.PSLA {
		if s.PSLA[i] != other.PSLA[i] {
			return false
		}
	}

	return true
}

func (s State) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)

	// Marshal basic values
	err := encbin.MarshalValues(
		buf,
		s.Street,
		s.TurnPos,
		s.BtnPos,
		s.BetAction,
		s.StreetAction,
		s.CallAmount,
		s.BSC.Amount,
		s.BSC.Addition,
		s.BSC.Action,
	)
	if err != nil {
		return nil, err
	}

	// Marshal Players slice
	err = encbin.MarshalSliceLen[Player, uint8](buf, s.Players)
	if err != nil {
		return nil, err
	}

	// Marshal PSC (chips.List)
	err = encbin.MarshalSliceLen[chips.Chips, uint8](buf, s.PSC)
	if err != nil {
		return nil, err
	}

	// Marshal PSAC ([]uint8)
	err = encbin.MarshalSliceLen[uint8, uint8](buf, s.PSAC)
	if err != nil {
		return nil, err
	}

	// Marshal PSLA ([]ActionKind)
	err = encbin.MarshalSliceLen[ActionKind, uint8](buf, s.PSLA)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (s *State) UnmarshalBinary(data []byte) error {
	buf := bytes.NewReader(data)

	// Unmarshal basic values
	err := encbin.UnmarshalValues(
		buf,
		&s.Street,
		&s.TurnPos,
		&s.BtnPos,
		&s.BetAction,
		&s.StreetAction,
		&s.CallAmount,
		&s.BSC.Amount,
		&s.BSC.Addition,
		&s.BSC.Action,
	)
	if err != nil {
		return err
	}

	// Unmarshal Players slice
	s.Players, err = encbin.UnmarhsalSliceLen[Player, uint8](buf)
	if err != nil {
		return err
	}

	// Unmarshal PSC (chips.List)
	s.PSC, err = encbin.UnmarhsalSliceLen[chips.Chips, uint8](buf)
	if err != nil {
		return err
	}

	// Unmarshal PSAC ([]uint8)
	s.PSAC, err = encbin.UnmarhsalSliceLen[uint8, uint8](buf)
	if err != nil {
		return err
	}

	// Unmarshal PSLA ([]ActionKind)
	s.PSLA, err = encbin.UnmarhsalSliceLen[ActionKind, uint8](buf)
	if err != nil {
		return err
	}

	return nil
}

// Size returns the total size in bytes needed to marshal the State
func (s *State) Size() uint64 {
	if s == nil {
		return 0
	}

	size := uint64(0)

	// Fixed size components
	size += 1 // Street (uint8)
	size += 1 // TurnPos (uint8)
	size += 1 // BtnPos (uint8)
	size += 1 // BetAction (uint8)
	size += 1 // StreetAction (uint8)
	size += 4 // CallAmount (chips.Chips - float32)

	// BSC struct
	size += 4 // Amount (chips.Chips)
	size += 4 // Addition (chips.Chips)
	size += 1 // Action (ActionKind - uint8)
	// Variable size components
	size += uint64(len(s.Players)) * 5 // Each Player: 4 bytes Paid + 1 byte Status
	size += uint64(len(s.PSC)) * 4     // PSC (chips.List - []float32)
	size += uint64(len(s.PSAC))        // PSAC ([]uint8)
	size += uint64(len(s.PSLA))        // PSLA ([]ActionKind - []uint8)

	// Length prefixes for slices
	size += 1 // Players length
	size += 1 // PSC length
	size += 1 // PSAC length
	size += 1 // PSLA length

	return size
}

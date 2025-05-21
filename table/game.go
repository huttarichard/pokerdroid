package table

import (
	"bytes"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/pokerdroid/poker/chips"
	"github.com/pokerdroid/poker/encbin"
)

type GameParams struct {
	NumPlayers         uint8       `json:"num_players"`
	MaxActionsPerRound uint8       `json:"max_actions_per_round"`
	BtnPos             uint8       `json:"btn_pos"`
	SbAmount           chips.Chips `json:"sb_amount"`
	BetSizes           [][]float32 `json:"bet_sizes"`
	InitialStacks      chips.List  `json:"initial_stacks"`
	TerminalStreet     Street      `json:"terminal_street"`
	MinBet             bool        `json:"min_bet"`
	Limp               bool        `json:"limp"`
	// Disable validation to improve performance
	DisableV bool `json:"disable_v"`
}

// Size returns the number of bytes needed to store GameParams
func (g GameParams) Size() uint64 {
	size := uint64(0)

	// Fixed size fields
	size += 1 // NumPlayers (uint8)
	size += 1 // MaxActionsPerRound (uint8)
	size += 1 // BtnPos (uint8)
	size += 4 // SbAmount (chips.Chips - float32)
	size += 1 // TerminalStreet (Street - uint8)
	size += 1 // DisableV (bool - uint8)
	size += 1 // MinBet (bool - uint8)

	size += 1 // BetSizes length
	for _, betSizes := range g.BetSizes {
		size += 1 + uint64(len(betSizes))*4 // Length prefix + float32 per bet size
	}

	size += 1 + uint64(len(g.InitialStacks))*4 // Length prefix + float32 per stack

	return size
}

// EffectiveStack returns the effective stack for a player.
// For example giving the following initial stacks:
// [100, 200, 300, 400]
// Effective stacks would be:
// [100, 200, 300, 300]
func (g GameParams) EffectiveStack(p uint8) chips.Chips {
	switch len(g.InitialStacks) {
	case 0:
		return 0
	case 1:
		return g.InitialStacks[0]
	default:
		sorted := make(chips.List, len(g.InitialStacks))
		copy(sorted, g.InitialStacks)
		sort.Slice(sorted, func(i, j int) bool {
			return sorted[i] < sorted[j]
		})
		return chips.Min(sorted[len(sorted)-2], g.InitialStacks[p])
	}
}

// String returns a string representation of GameParams.
func (g GameParams) String() string {
	var sb strings.Builder

	sb.WriteString("GameParams{\n")
	sb.WriteString(" Players:")
	sb.WriteString(strconv.Itoa(int(g.NumPlayers)))
	sb.WriteString("\n")
	sb.WriteString(" MaxActions:")
	sb.WriteString(strconv.Itoa(int(g.MaxActionsPerRound)))
	sb.WriteString("\n")
	sb.WriteString(" SB:")
	sb.WriteString(g.SbAmount.String())
	sb.WriteString("\n")
	sb.WriteString(" BetSizes:")
	sb.WriteString(fmt.Sprint(g.BetSizes))
	sb.WriteString("\n")
	sb.WriteString(" MinBet:")
	sb.WriteString(strconv.FormatBool(g.MinBet))
	sb.WriteString("\n")
	sb.WriteString(" Stacks:")
	for _, stack := range g.InitialStacks {
		sb.WriteString(stack.String())
		sb.WriteString(", ")
	}
	sb.WriteString("\n")
	sb.WriteString("}")

	return sb.String()
}

// NewGameParams will create a new GameParams with default values.
// Set np as number of players and bb is how many big blinds to start with.
func NewGameParams(np uint8, stack chips.Chips) GameParams {
	is := chips.NewListAlloc(np)
	for i := range is {
		is[i] = stack
	}

	g := GameParams{
		NumPlayers:         np,
		MaxActionsPerRound: uint8(np * 3),
		BtnPos:             0,
		SbAmount:           chips.New(1),
		InitialStacks:      is,
		TerminalStreet:     River,
		DisableV:           false,
		MinBet:             false,
	}

	g.SetBetSizes()

	return g
}

func (g GameParams) Clone() GameParams {
	prsm := GameParams{
		NumPlayers:         g.NumPlayers,
		MaxActionsPerRound: g.MaxActionsPerRound,
		BtnPos:             g.BtnPos,
		SbAmount:           g.SbAmount,
		TerminalStreet:     g.TerminalStreet,
		DisableV:           g.DisableV,
		MinBet:             g.MinBet,
	}
	prsm.BetSizes = make([][]float32, len(g.BetSizes))
	for i, betSizes := range g.BetSizes {
		prsm.BetSizes[i] = make([]float32, len(betSizes))
		copy(prsm.BetSizes[i], betSizes)
	}

	prsm.InitialStacks = make(chips.List, len(g.InitialStacks))
	copy(prsm.InitialStacks, g.InitialStacks)
	return prsm
}

var BetSizesDeep = [][]float32{
	// 200bb stacks
	{0.5, 1, 1.5, 3, 9, 15, 25, 50},
	{0.5, 1, 2, 4, 8, 16, 32},
	{0.5, 1, 3, 8},
	{1, 3},
	{1},
}

var BetSizesMedium = [][]float32{
	// 100bb stacks
	{0.5, 0.75, 1, 1.5, 2, 3, 5, 8, 16, 32}, // Initial bets
	{0.5, 0.75, 1, 1.5, 2, 3, 5, 8},         // Raises
	{0.5, 1, 2, 4},                          // 3-bets
	{1, 2},                                  // 4-bets
	{1},                                     // 5-bets
}

var BetSizesShallow = [][]float32{
	// 25bb stacks
	{0.25, 0.5, 0.75, 1, 2, 3, 4, 6, 8, 12, 16}, // Initial bets
	{0.25, 0.5, 1, 2, 4, 8, 16},                 // Raises
	{0.25, 0.5, 1, 2, 4, 8},                     // 3-bets
	{0.5, 1, 2, 4, 8},                           // 4-bets
}

func (g *GameParams) SetBetSizes() {
	depth := chips.Min(g.InitialStacks...).Div(g.SbAmount.Mul(2))

	if depth.GreaterThanOrEqual(200) {
		g.BetSizes = BetSizesDeep
		return
	}

	if depth.GreaterThanOrEqual(100) {
		g.BetSizes = BetSizesMedium
		return
	}

	g.BetSizes = BetSizesShallow
}

func (g GameParams) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)

	// Marshal basic values
	err := encbin.MarshalValues(
		buf,
		g.NumPlayers,
		g.MaxActionsPerRound,
		g.BtnPos,
		g.TerminalStreet,
		g.DisableV,
		g.SbAmount,
		g.MinBet,
	)
	if err != nil {
		return nil, err
	}

	// Write BetSizes length
	err = encbin.MarshalValues(buf, uint8(len(g.BetSizes)))
	if err != nil {
		return nil, err
	}

	// Marshal BetSizes slice
	for _, betSizes := range g.BetSizes {
		err = encbin.MarshalSliceLen[float32, uint8](buf, betSizes)
		if err != nil {
			return nil, err
		}
	}

	// Marshal InitialStacks
	err = encbin.MarshalSliceLen[chips.Chips, uint8](buf, g.InitialStacks)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (g *GameParams) UnmarshalBinary(data []byte) error {
	buf := bytes.NewReader(data)

	// Unmarshal basic values
	err := encbin.UnmarshalValues(
		buf,
		&g.NumPlayers,
		&g.MaxActionsPerRound,
		&g.BtnPos,
		&g.TerminalStreet,
		&g.DisableV,
		&g.SbAmount,
		&g.MinBet,
	)
	if err != nil {
		return err
	}

	// Unmarshal BetSizes length
	var betSizesLength uint8
	err = encbin.UnmarshalValues(buf, &betSizesLength)
	if err != nil {
		return err
	}

	g.BetSizes = make([][]float32, betSizesLength)
	for i := range g.BetSizes {
		g.BetSizes[i], err = encbin.UnmarhsalSliceLen[float32, uint8](buf)
		if err != nil {
			return err
		}
	}

	// Unmarshal InitialStacks
	g.InitialStacks, err = encbin.UnmarhsalSliceLen[chips.Chips, uint8](buf)
	if err != nil {
		return err
	}

	return nil
}

type Game struct {
	GameParams `json:"game_params"`

	Latest *State `json:"latest"`
}

func NewGame(params GameParams) (game *Game, err error) {
	state := NewState(params)

	state, err = MakeInitialBets(params, state)
	if err != nil {
		return nil, err
	}

	return &Game{
		GameParams: params,
		Latest:     state,
	}, nil
}

func NewGameFromState(params GameParams, state *State) *Game {
	return &Game{
		GameParams: params,
		Latest:     state,
	}
}

func (g *Game) Action(action Actioner) error {
	state, err := MakeAction(g.GameParams, g.Latest, action)
	if err != nil {
		return err
	}
	state, err = Move(g.GameParams, state)
	if err != nil {
		return err
	}
	g.Latest = state
	return nil
}

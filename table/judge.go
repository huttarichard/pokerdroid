package table

import (
	"bytes"
	"sort"

	"github.com/pokerdroid/poker/card"
	"github.com/pokerdroid/poker/chips"
	"github.com/pokerdroid/poker/encbin"
	"github.com/pokerdroid/poker/eval"
)

// POTS =================================

type Pot struct {
	Amount  chips.Chips
	Players []uint8
}

func (p Pot) Size() uint64 {
	size := uint64(0)
	size += 4
	size += 1
	size += uint64(len(p.Players))
	return size
}

// MarshalBinary implements the encoding.BinaryMarshaler interface
func (p Pot) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)

	// Marshal Amount
	err := encbin.MarshalValues(buf, p.Amount)
	if err != nil {
		return nil, err
	}

	err = encbin.MarshalSliceLen[uint8, uint8](buf, p.Players)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface
func (p *Pot) UnmarshalBinary(data []byte) error {
	buf := bytes.NewReader(data)

	// Unmarshal Amount
	err := encbin.UnmarshalValues(buf, &p.Amount)
	if err != nil {
		return err
	}

	// Initialize Players slice
	p.Players, err = encbin.UnmarhsalSliceLen[uint8, uint8](buf)
	if err != nil {
		return err
	}

	return nil
}

type Pots []Pot

func (pp Pots) Size() uint64 {
	size := uint64(0)
	size += 1
	for _, pot := range pp {
		size += 1
		size += pot.Size()
	}
	return size
}

// MarshalBinary implements the encoding.BinaryMarshaler interface
func (pp Pots) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)

	// Marshal length
	err := encbin.MarshalValues(buf, uint8(len(pp)))
	if err != nil {
		return nil, err
	}

	// Marshal each pot
	for i := range pp {
		err = encbin.MarshalWithLen[uint8](buf, pp[i])
		if err != nil {
			return nil, err
		}
	}

	return buf.Bytes(), nil
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface
func (pp *Pots) UnmarshalBinary(data []byte) error {
	buf := bytes.NewReader(data)

	// Read length
	var length uint8
	err := encbin.UnmarshalValues(buf, &length)
	if err != nil {
		return err
	}

	// Initialize slice
	*pp = make(Pots, length)

	// Unmarshal each pot
	for i := range *pp {
		p := &Pot{}
		err = encbin.UnmarshalWithLen[uint8](buf, p)
		if err != nil {
			return err
		}
		(*pp)[i] = *p
	}

	return nil
}

func (pp Pots) Sum() chips.Chips {
	sum := chips.Zero
	for _, pot := range pp {
		sum = sum.Add(pot.Amount)
	}
	return sum
}

func GetPots(p Players) (pots Pots) {
	allin := p.Indicies(IsAllIn)

	// Sort all-in players by their paid amount in ascending order
	sort.Slice(allin, func(i, j int) bool {
		return p[allin[i]].Paid.LessThan(p[allin[j]].Paid)
	})

	side := Pots{}
	for _, pi := range allin {
		sidePotSum := side.Sum()
		pot := p[pi].Paid.Sub(sidePotSum)
		eligible := []uint8{}
		for i, ep := range p {
			if ep.Paid.GreaterThanOrEqual(p[pi].Paid) && ep.Status != StatusFolded {
				eligible = append(eligible, uint8(i))
			}
		}
		side = append(side, Pot{
			Amount:  pot,
			Players: eligible,
		})
	}

	maxpay := p.PaidMax()
	eligible := []uint8{}

	for id, player := range p {
		if player.Paid.Equal(maxpay) && player.Status != StatusFolded {
			eligible = append(eligible, uint8(id))
		}
	}

	mp := Pot{
		Amount:  p.PaidSum().Sub(side.Sum()),
		Players: eligible,
	}

	pots = append(pots, mp)
	pots = append(pots, side...)

	return pots
}

// JUDGE =================================

type Judger interface {
	// Compare compares hands returns 2 if tie, 0 if win, 1 if loss for a
	// ab are indexes
	// Judge(a, b int) int

	// Compare hands of mupltiple players, pp are indexes of players
	// returns indexes of winners
	Judge(pp []uint8) []uint8
}

func GetWinners(ss *State, judge Judger) []uint8 {
	active := ss.Players.Indicies(IsActivePlayer)
	switch len(active) {
	case 1:
		return []uint8{active[0]}
	case 0:
		return []uint8{}
	}

	winners := judge.Judge(active)
	return winners
}

func GetWinnings(pp Players, judge Judger) chips.List {
	winnings := chips.NewListAlloc(len(pp))

	// Get pots, main and side pots
	for _, pot := range GetPots(pp) {
		// If someone folds
		if len(pot.Players) == 1 {
			winnings[pot.Players[0]] = winnings[pot.Players[0]].Add(pot.Amount)
			continue
		}

		winners := judge.Judge(pot.Players)

		// Most of the time there is only one winner
		// just optimalization
		if len(winners) == 1 {
			winnings[winners[0]] = winnings[winners[0]].Add(pot.Amount)
			continue
		}

		// If there is a tie, split the pot
		prize := pot.Amount.Div(chips.New(len(winners)))
		for _, winner := range winners {
			winnings[winner] = winnings[winner].Add(prize)
		}
	}
	return winnings
}

type Cards struct {
	Community card.Cards
	Players   []card.Cards
}

func (d *Cards) Judge(pp []uint8) []uint8 {
	var hands []card.Cards

	for _, p := range pp {
		hands = append(hands, d.Players[p].Clone())
	}

	return eval.MustJudgeBoard(hands, d.Community)
}

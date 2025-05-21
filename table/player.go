package table

import (
	"github.com/pokerdroid/poker/chips"
)

type Status uint8

const (
	StatusUnknown Status = iota
	StatusFolded
	StatusActive
	StatusAllIn
)

func (s Status) String() string {
	switch s {
	case StatusFolded:
		return "folded"
	case StatusActive:
		return "active"
	case StatusAllIn:
		return "all-in"
	}
	return "unknown"
}

type Player struct {
	Paid   chips.Chips `json:"paid"`
	Status Status      `json:"status"`
}

func (p Player) Size() uint64 {
	size := uint64(0)
	size += 4
	size += 1
	return size
}

func NewPlayer() Player {
	return Player{
		Paid:   chips.Zero,
		Status: StatusActive,
	}
}

type Players []Player

func (p Players) LastAlive(idx uint8) bool {
	if p[idx].Status == StatusFolded {
		return false
	}
	for i, x := range p {
		if uint8(i) == idx {
			continue
		}
		if x.Status != StatusFolded {
			return false
		}
	}
	return true
}

func (p Players) Clone() Players {
	pp := make(Players, len(p))
	copy(pp, p)
	return pp
}

func (p Players) PaidSum() chips.Chips {
	paid := chips.Zero
	for _, player := range p {
		paid = paid.Add(player.Paid)
	}
	return paid
}

func (p Players) PaidMax() chips.Chips {
	if len(p) == 0 {
		return chips.Zero
	}

	maxPaid := p[0].Paid
	for _, player := range p[1:] {
		if player.Paid.GreaterThan(maxPaid) {
			maxPaid = player.Paid
		}
	}
	return maxPaid
}

// FILTERS ================================

type Filter func(p Player) bool

func IsActivePlayer(p Player) bool {
	return p.Status != StatusFolded
}

func IsAllIn(p Player) bool {
	return p.Status == StatusAllIn
}

func IsWaitingAskPlayer(p Player) bool {
	return p.Status == StatusActive
}

func (p Players) Indicies(a Filter) (pp []uint8) {
	for ix, player := range p {
		if !a(player) {
			continue
		}
		pp = append(pp, uint8(ix))
	}
	return pp
}

func (p Players) Len(a Filter) (pp int) {
	for _, player := range p {
		if a(player) {
			pp++
		}
	}
	return pp
}

func (p Players) FromPos(startPos int) []int {
	var pp []int
	for ix := range p {
		pp = append(pp, ix)
	}
	px := append(pp, pp...)
	px = px[startPos : startPos+len(pp)]
	return px
}

func (p Players) FindPos(start int, filter Filter) int {
	for _, px := range p.FromPos(start) {
		if filter(p[px]) {
			return px
		}
	}
	return -1
}

func (p Players) FindActivePos(start int) int {
	return p.FindPos(start, IsActivePlayer)
}

func (p Players) FindWaitingPos(start int) int {
	return p.FindPos(start, IsWaitingAskPlayer)
}

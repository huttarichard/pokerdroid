package absp

import (
	"bytes"
	"os"
	"testing"

	"github.com/edsrzf/mmap-go"
	"github.com/google/uuid"
	"github.com/pokerdroid/poker/abs"
	"github.com/pokerdroid/poker/abs/flop"
	"github.com/pokerdroid/poker/abs/river"
	"github.com/pokerdroid/poker/abs/turn"
	"github.com/pokerdroid/poker/card"
	"github.com/pokerdroid/poker/encbin"
	"github.com/pokerdroid/poker/equity"
	"github.com/pokerdroid/poker/iso"
	"github.com/pokerdroid/poker/table"
)

type Iso struct{}

func NewIso() Iso {
	return Iso{}
}

func (Iso) Map(cds card.Cards) abs.Cluster {

	// pretty.Println(cds)
	switch len(cds) {
	case 2:
		return abs.Cluster(iso.Preflop.Index(cds))
	case 5:
		return abs.Cluster(iso.Flop.Index(cds))
	case 6:
		return abs.Cluster(iso.Turn.Index(cds))
	case 7:
		return abs.Cluster(iso.River.Index(cds))
	default:
		panic("invalid number of cards")
	}
}

type Abs struct {
	UID   uuid.UUID
	Flop  *flop.Abs
	Turn  *turn.Abs
	River *river.Abs
}

func NewFromFile(filename string) (*Abs, error) {
	f, err := os.OpenFile(filename, os.O_RDONLY, 0)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	data, err := mmap.Map(f, mmap.RDONLY, 0)
	if err != nil {
		return nil, err
	}
	defer data.Unmap()

	x := &Abs{
		Flop:  &flop.Abs{},
		Turn:  &turn.Abs{},
		River: &river.Abs{},
	}
	err = x.UnmarshalBinary(data)
	return x, err
}

func TstNewAbs(t testing.TB) *Abs {
	absp := os.Getenv("ABS_PATH")
	if absp == "" {
		t.Skip("ABS_PATH not set")
		return nil
	}
	abs, err := NewFromFile(absp)
	if err != nil {
		t.Fatal(err)
	}
	return abs
}

func (a *Abs) Map(cds card.Cards) abs.Cluster {
	switch len(cds) {
	case 2:
		return abs.Cluster(iso.Preflop.Index(cds))
	case 5:
		return a.Flop.Map(iso.Flop.Index(cds))
	case 6:
		return a.Turn.Map(iso.Turn.Index(cds))
	case 7:
		return a.River.Map(iso.River.Index(cds))
	default:
		panic("invalid number of cards")
	}
}

func (a *Abs) Equity(street table.Street, c abs.Cluster) equity.Equity {
	var eq float32
	switch street {
	case table.Preflop:
		return preflop[c]
	case table.Flop:
		eq = a.Flop.Equity[c]
	case table.Turn:
		eq = a.Turn.Equity[c]
	case table.River:
		return a.River.Equities[c]
	default:
		panic("invalid street")
	}
	return equity.NewEquity(eq, 0)
}

func (a *Abs) MarshalBinary() ([]byte, error) {
	buf := bytes.NewBuffer(nil)

	err := encbin.MarshalWithLen[uint16](buf, a.UID)
	if err != nil {
		return nil, err
	}

	err = encbin.MarshalWithLen[uint64](buf, a.Flop)
	if err != nil {
		return nil, err
	}

	err = encbin.MarshalWithLen[uint64](buf, a.Turn)
	if err != nil {
		return nil, err
	}

	err = encbin.MarshalWithLen[uint64](buf, a.River)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (a *Abs) UnmarshalBinary(data []byte) error {
	buf := bytes.NewBuffer(data)

	err := encbin.UnmarshalWithLen[uint16](buf, &a.UID)
	if err != nil {
		return err
	}

	err = encbin.UnmarshalWithLen[uint64](buf, a.Flop)
	if err != nil {
		return err
	}

	err = encbin.UnmarshalWithLen[uint64](buf, a.Turn)
	if err != nil {
		return err
	}

	err = encbin.UnmarshalWithLen[uint64](buf, a.River)
	if err != nil {
		return err
	}

	return nil
}

type AbsFn func(cds card.Cards) abs.Cluster

func (a AbsFn) Map(cds card.Cards) abs.Cluster {
	return a(cds)
}

package card

import (
	"bytes"
	"encoding/binary"
)

// This array is for Stringer index is the HandRank
var handRanksNames = []string{
	"Invalid",
	"High card",
	"One pair",
	"Two pairs",
	"Three of a kind",
	"Straight",
	"Flush",
	"Full house",
	"Four of a kind",
	"Straight flush",
}

func HandRankNames() []string {
	return append([]string{}, handRanksNames...)
}

// HandRankKind is a type of the hand
type HandRankKind int

const (
	HandRankNone          HandRankKind = 0
	HandRankHighCard      HandRankKind = 1
	HandRankOnePair       HandRankKind = 2
	HandRankTwoPairs      HandRankKind = 3
	HandRankThreeOfaKind  HandRankKind = 4
	HandRankStraight      HandRankKind = 5
	HandRankFlush         HandRankKind = 6
	HandRankFullHouse     HandRankKind = 7
	HandRankFourOfaKind   HandRankKind = 8
	HandRankStraightFlush HandRankKind = 9
)

// String returns string representation of hand type
func (ht HandRankKind) String() string {
	return handRanksNames[ht]
}

// HandRank contains information about hand type and rank
type HandRank struct {
	Kind HandRankKind `json:"kind"`
	Rank uint32       `json:"rank"`
}

func NewHandRank(ht HandRankKind, rank uint32) HandRank {
	return HandRank{Kind: ht, Rank: rank}
}

func (h HandRank) String() string {
	if h.Kind == HandRankNone {
		return "unknown"
	}
	return h.Kind.String()
}

func (h HandRank) Empty() bool {
	return h.Kind == HandRankNone
}

// Compare compares hands returns 2 if tie, 0 if win, 1 if loss
func (h HandRank) Compare(other HandRank) int {
	if h.Rank == other.Rank && h.Kind == other.Kind {
		return 2
	}

	if h.Kind > other.Kind {
		return 0
	}

	if h.Kind < other.Kind {
		return 1
	}

	if h.Rank > other.Rank {
		return 0
	}

	return 1
}

type handRankBinary struct {
	Kind uint8
	Rank uint32
}

// MarshalBinary implements the encoding.BinaryMarshaler interface.
func (h HandRank) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, handRankBinary{
		Kind: uint8(h.Kind),
		Rank: h.Rank,
	})
	return buf.Bytes(), err
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface.
func (h *HandRank) UnmarshalBinary(buf []byte) error {
	dd := handRankBinary{}
	err := binary.Read(bytes.NewBuffer(buf), binary.LittleEndian, &dd)
	h.Kind = HandRankKind(dd.Kind)
	h.Rank = dd.Rank
	return err
}

package eval

import (
	"bytes"
	"compress/gzip"
	"encoding/binary"
	"errors"
	"io"

	_ "embed"

	"github.com/pokerdroid/poker/card"
)

//go:embed ranks.bin
var ranksBin []byte
var ranks []byte

func init() {
	reader, _ := gzip.NewReader(bytes.NewReader(ranksBin))
	defer reader.Close()
	ranks, _ = io.ReadAll(reader)
}

// Eval evaluates poker hand, it can evaluate 5,6 and 7 hands.
// Returned rank is the weight of the hand. Bigger rank is better hand.
func Eval(cards ...card.Card) (card.HandRank, error) {
	size := len(cards) // len is just for shorten code
	if size != 7 && size != 6 && size != 5 && size != 2 {
		return card.HandRank{}, errors.New("cards can be 7,6,5 or 2 length")
	}

	if size == 2 {
		var r card.HandRankKind
		if cards[0].Rank() == cards[1].Rank() {
			r = card.HandRankOnePair
		} else {
			r = card.HandRankHighCard
		}
		return card.HandRank{Kind: r}, nil
	}

	var p uint32 = 53
	for i := 0; i < size; i++ {
		p = evalCard(p+uint32(cards[i]), ranks)
	}

	if size == 5 || size == 6 {
		p = evalCard(p, ranks)
	}

	tp := card.HandRankKind(p >> 12)
	if tp == 0 {
		return card.HandRank{}, errors.New("wrong cards")
	}

	rank := p & 0x00000fff
	return card.NewHandRank(tp, rank), nil
}

func Judge(ccc []card.Cards) ([]uint8, error) {
	ranks := make([]card.HandRank, 0, len(ccc))

	for _, cc := range ccc {
		rank, err := Eval(cc...)
		if err != nil {
			return nil, err
		}
		ranks = append(ranks, rank)
	}

	winners := []uint8{0}

	for i, r := range ranks[1:] {
		i += 1

		winner := r.Compare(ranks[winners[0]])
		if winner == 0 {
			winners = []uint8{uint8(i)}
			continue
		}
		if winner == 2 {
			winners = append(winners, uint8(i))
		}
	}

	return winners, nil
}

func MustJudge(ccc []card.Cards) []uint8 {
	winners, err := Judge(ccc)
	if err != nil {
		panic(err)
	}
	return winners
}

func JudgeBoard(hole []card.Cards, board card.Cards) ([]uint8, error) {
	var ccc []card.Cards
	for _, h := range hole {
		ccc = append(ccc, append(h, board...))
	}
	return Judge(ccc)
}

func MustJudgeBoard(hole []card.Cards, board card.Cards) []uint8 {
	winners, err := JudgeBoard(hole, board)
	if err != nil {
		panic(err)
	}
	return winners
}

func evalCard(card uint32, ranks []byte) uint32 {
	offset := uint32(card) * 4
	x := ranks[offset : offset+4]
	r := binary.LittleEndian.Uint32(x)
	return r
}

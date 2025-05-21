package mc

import (
	"fmt"

	"github.com/pokerdroid/poker/card"
)

type OpenRange struct {
	Range [13][13]float64
}

func NewOpenRange() OpenRange {
	openRange := OpenRange{}
	for i := range openRange.Range {
		for j := range openRange.Range[i] {
			openRange.Range[i][j] = 1
		}
	}
	openRange.Range[10][4] = 0.37
	openRange.Range[10][5] = 0
	openRange.Range[10][6] = 0
	openRange.Range[11][3] = 0
	for i := 4; i < 10; i++ {
		openRange.Range[11][i] = 0
	}
	openRange.Range[11][10] = 0.28
	openRange.Range[12][3] = 0
	for i := 4; i < 12; i++ {
		openRange.Range[12][i] = 0
	}
	return openRange
}

//  Ahlpha holdem opening range
//  13:13, matrix, each node represent probability of not folding in HU
//
//    [A    K    Q    J    T    9    8    7    6    5    4    3    2]
//  A [1    1    1    1    1    1    1    1    1    1    1    1    1]
//  K [1    1    1    1    1    1    1    1    1    1    1    1    1]
//  Q [1    1    1    1    1    1    1    1    1    1    1    1    1]
//  J [1    1    1    1    1    1    1    1    1    1    1    1    1]
//  T [1    1    1    1    1    1    1    1    1    1    1    1    1]
//  9 [1    1    1    1    1    1    1    1    1    1    1    1    1]
//  8 [1    1    1    1    1    1    1    1    1    1    1    1    1]
//  7 [1    1    1    1    1    1    1    1    1    1    1    1    1]
//  6 [1    1    1    1    1    1    1    1    1    1    1    1    1]
//  5 [1    1    1    1    1    1    1    1    1    1    1    1    1]
//  4 [1    1    1    1    0.37 0.   0.   1    1    1    1    1    1]
//  3 [1    1    1    0.   0.   0.   0.   0.   0.   0.   0.28 1    1]
//  2 [1    1    1    0.   0.   0.   0.   0.   0.   0.   0.   0.   1]

func (o *OpenRange) WeakRange(hole card.Cards) bool {
	if len(hole) != 2 {
		fmt.Println("invalid hole cards")
		return false
	}

	rx, ry := MatrixPos(hole)
	return o.Range[rx][ry] < 1
}

func MatrixPos(hole card.Cards) (int, int) {
	suited := hole[0].Suite() == hole[1].Suite()
	var sorted card.Cards
	if hole[0].Rank() < hole[1].Rank() {
		sorted = card.Cards{hole[1], hole[0]}
	} else {
		sorted = card.Cards{hole[0], hole[1]}
	}

	r0 := 13 - int(sorted[0].Rank())
	r1 := 13 - int(sorted[1].Rank())

	if suited {
		return r0, r1
	} else {
		return r1, r0
	}
}

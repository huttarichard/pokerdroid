package omp

import (
	"testing"
	"time"

	"github.com/pokerdroid/poker/card"
)

func TestOMP(t *testing.T) {
	tx := time.Now()
	// "8dAhKh"
	board := card.Cards{card.Card8D, card.CardAH, card.CardKH}
	x := Equity(card.Card9S, card.CardKS, board, 6)
	w := time.Since(tx)
	print(x.String(), w.String())
}

// Package eval provides poker hand evaluation.
//
// Inspired by [phe].
//
// Can evaluate 5,6 and 7 card combinations.
// It uses TwoPlusTwo algorithm with pre-build ranks table.
// It can evaluate 70 millions 7-cards hands per second (benchmarked on Apple M1).
// That is approximately 14ns per eval.
// Key ingredient is pre-build ranks table, which increases binary size.
//
//	r, err := Eval(card.Card2C, card.Card2D, card.Card2H, card.Card2S, card.Card3C)
//	r.String() // Four of a kind
//
// [phe]: https://github.com/spiritofsim/phe
package eval

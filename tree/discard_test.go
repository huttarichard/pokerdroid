package tree

import (
	"testing"
)

func TestDiscardBelowEpsilonSimple(t *testing.T) {
	// Create a simple tree:
	//       Root
	//        |
	//     Chance
	//        |
	//     Player (Check:0.7, Call:0.3)
	//     /     \
	//  Term    Player (Call:0.6, Fold:0.4)
	//           /    \
	//        Term   Term

	// TODO
}

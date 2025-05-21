package baselinenn

import (
	"testing"

	"github.com/pokerdroid/poker/frand"
)

func TestForModelParams(t *testing.T) {
	m := NewModel[float64](ModelParams{
		NumPlayers:         4,
		MaxActionsPerRound: 1,
		Layers:             1,
		HiddenSize:         5,
		BetSizing: []float32{
			0.25, 0.5, 0.75, 1, 1.5, 2, 2.5, 3, 4,
		},
	})

	m.InitRandom(frand.NewUnsafeInt(42))
}

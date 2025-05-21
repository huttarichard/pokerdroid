package slumbot

import (
	"context"
	"testing"

	"github.com/pokerdroid/poker/bot/mc"
	"github.com/stretchr/testify/require"
)

func TestRound(t *testing.T) {
	tk, err := Login("brownass", "brownass")
	require.NoError(t, err)

	for i := 0; i < 10; i++ {
		advisor := mc.NewAdvisor()
		ctx := context.Background()

		_, err = Run(ctx, tk, advisor)
		require.NoError(t, err)
	}
}

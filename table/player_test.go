package table

import (
	"testing"

	"github.com/pokerdroid/poker/chips"
	"github.com/stretchr/testify/require"
)

func TestPlayers_PaidSumAndMax(t *testing.T) {
	pl := Players{
		{Paid: chips.New(10), Status: StatusActive},
		{Paid: chips.New(20), Status: StatusFolded},
		{Paid: chips.New(15), Status: StatusActive},
	}
	require.Equal(t, chips.NewFromInt(45), pl.PaidSum())
	require.Equal(t, chips.NewFromInt(20), pl.PaidMax())
}

func TestPlayers_FindPos(t *testing.T) {
	pl := Players{
		{Paid: chips.New(5), Status: StatusFolded},
		{Paid: chips.New(10), Status: StatusActive},
		{Paid: chips.New(0), Status: StatusActive},
	}
	pos := pl.FindPos(1, IsActivePlayer)
	require.Equal(t, 1, pos)
	pos = pl.FindPos(2, IsActivePlayer)
	require.Equal(t, 2, pos)
	pos = pl.FindPos(3, IsActivePlayer)
	require.Equal(t, 1, pos) // wraps around
}

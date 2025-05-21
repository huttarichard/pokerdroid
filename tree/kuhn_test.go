package tree

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestKuhnMarshalUnmarshal(t *testing.T) {
	kuhn := NewKuhn()

	data, err := kuhn.MarshalBinary()
	require.NoError(t, err)

	kuhn2 := new(Root)
	err = kuhn2.UnmarshalBinary(data)
	require.NoError(t, err)

	MustVisit(kuhn2, -1, func(n Node, children []Node, depth int) bool {
		d, err := n.MarshalBinary()
		require.NoError(t, err)

		require.Equal(t, len(d), int(n.Size()))
		return true
	})

	require.Equal(t, len(data), int(kuhn.Size()))
}

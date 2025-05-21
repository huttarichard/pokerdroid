package table

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStreet_StringAndParse(t *testing.T) {
	cases := []struct {
		str      string
		expected Street
	}{
		{"preflop", Preflop},
		{"flop", Flop},
		{"turn", Turn},
		{"river", River},
		{"finished", Finished},
		{"", NoStreet},
	}

	for _, c := range cases {
		s, err := NewStreetFromString(c.str)
		require.NoError(t, err)
		require.Equal(t, c.expected, s)
		require.Equal(t, c.str == "" && s == NoStreet, s.String() == "no street" || s.String() == "unknown")
	}
}

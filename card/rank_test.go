package card

import "testing"

func TestRankToString(t *testing.T) {
	tests := []struct {
		rank Rank
		want string
	}{
		{Two, "2"},
		{Three, "3"},
		{Four, "4"},
		{Five, "5"},
		{Six, "6"},
		{Seven, "7"},
		{Eight, "8"},
		{Nine, "9"},
		{Ten, "T"},
		{Jack, "J"},
		{Queen, "Q"},
		{King, "K"},
		{Ace, "A"},
	}
	for _, test := range tests {
		if got := test.rank.String(); got != test.want {
			t.Errorf("Rank(%d).String() = %q, want %q", test.rank, got, test.want)
		}
	}
}

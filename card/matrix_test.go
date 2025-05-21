package card

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TestNewMatrixNoBlockers confirms that using no blockers returns a normalized matrix,
// which should be equivalent to the precomputed full-deck distribution.
func TestNewMatrixNoBlockers(t *testing.T) {
	m := NewMatrixFromBlockers(nil)

	total := m.Sum()

	require.InDelta(t, 1.0, total, 1e-5)

	// Also check that some cells are nonzero.
	nonZeroFound := false
	for i := 0; i < 13 && !nonZeroFound; i++ {
		for j := 0; j < 13 && !nonZeroFound; j++ {
			if m[i][j] > 0 {
				nonZeroFound = true
			}
		}
	}
	require.True(t, nonZeroFound)
}

// TestNewMatrixWithBlocker verifies that when a blocker is provided, the associated row and column are zeroed.
// For instance, for Ace of spades ("As"):
//   - Ace has a rank value of 13, hence its index is 13 - 13 = 0.
func TestNewMatrixWithBlocker(t *testing.T) {
	// Ace of spades
	blocker := Parse("as")
	m := NewMatrixFromBlockers(Cards{blocker})

	// For Ace: row and column at index 0 should be zero.
	for j := 0; j < 13; j++ {
		require.InDelta(t, 0, m[0][j], 1e-5)
	}
	for i := 0; i < 13; i++ {
		require.InDelta(t, 0, m[i][0], 1e-5)
	}
}

// TestNewMatrixWithMultipleBlockers uses two blockers: Ace of spades ("As") and Two of clubs ("2c").
// For Ace: index = 13 - 13 = 0.
// For Two: rank value is 1, so index = 13 - 1 = 12.
func TestNewMatrixWithMultipleBlockers(t *testing.T) {
	blockers := Cards{
		Parse("As"), // Ace of spades --> index 0
		Parse("2c"), // Two of clubs   --> index 12
	}
	m := NewMatrixFromBlockers(blockers)

	// Check Ace row and column, index 0.
	for j := 0; j < 13; j++ {
		require.InDelta(t, 0, m[0][j], 1e-5)
	}
	for i := 0; i < 13; i++ {
		require.InDelta(t, 0, m[i][0], 1e-5)
	}
	// Check Two row and column (index 12).
	for j := 0; j < 13; j++ {
		require.InDelta(t, 0, m[12][j], 1e-5)
	}
	for i := 0; i < 13; i++ {
		require.InDelta(t, 0, m[i][12], 1e-5)
	}
}

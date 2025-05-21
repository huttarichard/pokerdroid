package equity

import (
	"fmt"
	"math"
)

const (
	// Precision multiplier for converting between float and uint16
	// 10000 gives us 4 decimal places of precision (0.0001)
	Precision = 10000
)

// Equity represents win/tie probabilities as uint16 values
// Loss probability is derived as 1 - (win + tie)
type Equity [2]uint16

// NewCompactEquity creates a new CompactEquity from win, lose, tie values
func NewEquity(win, tie float32) Equity {
	return Equity{
		uint16(math.Round(float64(win)*Precision*100) / 100),
		uint16(math.Round(float64(tie)*Precision*100) / 100),
	}
}

// ToFloat32 returns the equity as [win, lose, tie] float32 values
func (c Equity) Float32() [3]float32 {
	win := float32(c[0]) / Precision
	tie := float32(c[1]) / Precision
	lose := 1.0 - win - tie
	return [3]float32{win, lose, tie}
}

// Win returns the win probability as float32
func (c Equity) Win() float32 {
	return float32(c[0]) / Precision
}

// Lose returns the lose probability as float32
func (c Equity) Lose() float32 {
	return 1.0 - c.Win() - c.Tie()
}

// Tie returns the tie probability as float32
func (c Equity) Tie() float32 {
	return float32(c[1]) / Precision
}

// String returns a human-readable representation of the equity
func (c Equity) String() string {
	values := c.Float32()
	return fmt.Sprintf("w:%.4f/l:%.4f/t:%.4f", values[0], values[1], values[2])
}

// ShortString returns a compact representation of the equity
func (c Equity) ShortString() string {
	values := c.Float32()
	return fmt.Sprintf("%.4f/%.4f/%.4f", values[0], values[1], values[2])
}

// Equals checks if two CompactEquity values are equal
func (c Equity) Equals(other Equity) bool {
	return c[0] == other[0] && c[1] == other[1]
}

// Distance calculates Euclidean distance between two CompactEquity values
func (c Equity) Distance(other Equity) float64 {
	v1 := c.Float32()
	v2 := other.Float32()

	dx := float64(v1[0] - v2[0])
	dy := float64(v1[1] - v2[1])
	dz := float64(v1[2] - v2[2])

	return math.Sqrt(dx*dx + dy*dy + dz*dz)
}

// WinHalfDraw returns the win probability plus half of the tie probability
func (c Equity) WinDraw() float32 {
	return c.Win() + c.Tie()/2
}

// Won calculates the expected value based on pot size
func (c Equity) Won(pot float32) float32 {
	return c.Win()*pot + c.Tie()*pot/2
}

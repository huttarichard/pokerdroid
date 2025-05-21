package card

import (
	"bytes"
	"fmt"
	"math"
	"strings"

	"github.com/pokerdroid/poker/encbin"
	"golang.org/x/exp/constraints"
)

type Matrix [13][13]float64

func NewMatrix[T constraints.Float | constraints.Integer](flat []T) Matrix {
	m := Matrix{}
	for i := 0; i < 13; i++ {
		for j := 0; j < 13; j++ {
			m[i][j] = float64(flat[i*13+j])
		}
	}
	return m
}

func NewOpenMatrix() Matrix {
	m := Matrix{}
	for _, c1 := range Combinations(2) {
		x, y := Coordinates(c1)
		m[x][y] = 1.0
	}
	return m.Normalize()
}

func NewMatrixFromBlockers(blockers Cards) Matrix {
	// Start with a copy of the precomputed full deck matrix.
	m := NewOpenMatrix()

	for _, b := range blockers {
		r := int(b.Rank())
		// Skip if the card has an invalid rank.
		if r == 0 {
			continue
		}
		// Convert the card's rank to the corresponding matrix index.
		idx := 13 - r
		// Zero out the entire row and column for this rank.
		for i := 0; i < 13; i++ {
			m[idx][i] = 0.0
			m[i][idx] = 0.0
		}
	}
	return m.Normalize()
}

func (m Matrix) Sub(n Matrix) Matrix {
	for i := 0; i < 13; i++ {
		for j := 0; j < 13; j++ {
			m[i][j] -= n[i][j]
		}
	}
	return m
}

func (m Matrix) Flat() []float64 {
	flat := make([]float64, 13*13)
	for i := 0; i < 13; i++ {
		for j := 0; j < 13; j++ {
			flat[i*13+j] = m[i][j]
		}
	}
	return flat
}

func (m Matrix) Empty() bool {
	for _, row := range m {
		for _, val := range row {
			if val > 0 {
				return false
			}
		}
	}
	return true
}

func (m Matrix) Sum() float64 {
	var sum float64
	for _, row := range m {
		for _, val := range row {
			sum += val
		}
	}
	return sum
}

func (m Matrix) AbsSum() float64 {
	var sum float64
	for _, row := range m {
		for _, val := range row {
			sum += math.Abs(val)
		}
	}
	return sum
}

func (m Matrix) Normalize() Matrix {
	var m2 Matrix
	total := m.Sum()
	for i, row := range m {
		for j, val := range row {
			m2[i][j] = val / total
		}
	}
	return m2
}

func (m Matrix) String() string {
	var sb strings.Builder

	// Matrix rows
	for i := 0; i < 13; i++ {
		sb.WriteString("│")
		for j := 0; j < 13; j++ {
			sb.WriteString(fmt.Sprintf(" %.4f", m[i][j]))
		}
		sb.WriteString(" │\n")
	}

	return sb.String()
}

// MarshalBinary returns the binary serialization of the matrix.
// It mimics the style in encoding.go by writing each element sequentially.
func (m Matrix) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)
	// Write elements in row-major order.
	for i := 0; i < 13; i++ {
		for j := 0; j < 13; j++ {
			// encbin.MarshalValues writes the value into buf.
			if err := encbin.MarshalValues(buf, m[i][j]); err != nil {
				return nil, err
			}
		}
	}
	return buf.Bytes(), nil
}

// UnmarshalBinary parses the binary representation into the matrix.
// It reads 13×13 float32 values, restoring the state of the matrix.
func (m *Matrix) UnmarshalBinary(data []byte) error {
	r := bytes.NewReader(data)
	for i := 0; i < 13; i++ {
		for j := 0; j < 13; j++ {
			if err := encbin.UnmarshalValues(r, &((*m)[i][j])); err != nil {
				return err
			}
		}
	}
	return nil
}

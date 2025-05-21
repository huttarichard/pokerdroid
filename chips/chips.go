package chips

import (
	"fmt"
	"math"
)

const epsilon = 1e-6

type Chips float32

var Zero = Chips(0)

func New[T int | uint8 | int64 | float32 | float64](val T) Chips {
	return Chips(val)
}

// NewFromInt creates a new Chips value from an int.
func NewFromInt(val int64) Chips {
	return Chips(float32(val))
}

// NewFromFloat64 creates a new Chips value from a float64.
func NewFromFloat64(val float64) Chips {
	return Chips(float32(val))
}

// NewFromFloat creates a new Chips value from a float64.
func NewFromFloat(val float64) Chips {
	return Chips(float32(val))
}

// NewFromFloat32 creates a new Chips value from a float32.
func NewFromFloat32(val float32) Chips {
	return Chips(val)
}

// NewFromString creates a new Chips value from a string.
// Assumes the string is a valid float number.
func NewFromString(val string) Chips {
	var f float32
	fmt.Sscanf(val, "%f", &f)
	return Chips(f)
}

// String returns a string representation of c.
func (c Chips) String() string {
	return fmt.Sprintf("%f", c)
}

// Add adds a and b and returns the result.
func (c Chips) Add(b Chips) Chips {
	return c + b
}

// Sub subtracts b from a and returns the result.
func (c Chips) Sub(b Chips) Chips {
	return c - b
}

// Mul multiplies a and b and returns the result.
func (c Chips) Mul(b Chips) Chips {
	return c * b
}

// Div divides a by b and returns the result.
func (c Chips) Div(b Chips) Chips {
	if b == 0 {
		return Chips(math.Inf(1)) // return +Inf on division by zero
	}
	return c / b
}

// Equal checks if c and b are equal within a tolerance level.
func (c Chips) Equal(b Chips) bool {
	return math.Abs(float64(c-b)) < epsilon
}

// GreaterThan checks if c is greater than b.
func (c Chips) GreaterThan(b Chips) bool {
	return c > b
}

// GreaterThanOrEqual checks if c is greater than or equal b.
func (c Chips) GreaterThanOrEqual(b Chips) bool {
	return c >= b
}

// LessThan checks if c is less than b.
func (c Chips) LessThan(b Chips) bool {
	return c < b
}

// LessThanOrEqual checks if c is less than b.
func (c Chips) LessThanOrEqual(b Chips) bool {
	return c <= b
}

func (c Chips) StringFixed(places int) string {
	format := fmt.Sprintf("%%.%df", places)
	return fmt.Sprintf(format, c)
}

// Abs returns the absolute value of chips.
func (c Chips) Abs() Chips {
	return Chips(math.Abs(float64(c)))
}

// Pow returns the power of i.
func (c Chips) Pow(i Chips) Chips {
	return Chips(math.Pow(float64(c), float64(i)))
}

// Float32 gives float32 representation
func (c Chips) Float32() float32 {
	return float32(c)
}

// Float64 gives float64 representation
func (c Chips) Float64() float64 {
	return float64(c)
}

func (c Chips) Round(places int) Chips {
	return Chips(math.Round(float64(c)*math.Pow(10, float64(places))) / math.Pow(10, float64(places)))
}

func (c Chips) RoundUp(places int) Chips {
	return Chips(math.Ceil(float64(c)*math.Pow(10, float64(places))) / math.Pow(10, float64(places)))
}

func Max(a ...Chips) Chips {
	if len(a) == 0 {
		return Zero
	}
	max := a[0]
	for _, v := range a {
		if v.GreaterThan(max) {
			max = v
		}
	}
	return max
}

func Min(a ...Chips) Chips {
	if len(a) == 0 {
		return Zero
	}
	min := a[0]
	for _, v := range a {
		if v.LessThan(min) {
			min = v
		}
	}
	return min
}

type List []Chips

func NewList(a ...Chips) List {
	return List(a)
}

func NewListAlloc[T int | uint8](a T) List {
	return make(List, a)
}

func (l List) Sum() Chips {
	sum := Zero
	for _, v := range l {
		sum = sum.Add(v)
	}
	return sum
}

func (l List) Max() Chips {
	return Max(l...)
}

func (l List) Min() Chips {
	return Min(l...)
}

func (l List) Copy() List {
	ll := make(List, len(l))
	copy(ll, l)
	return ll
}

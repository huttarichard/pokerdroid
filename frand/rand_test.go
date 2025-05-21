package frand

import (
	"math"
	"testing"
)

func TestRandomImplementations(t *testing.T) {
	implementations := map[string]func() Rand{
		"unsafe": NewUnsafe,
		"hash":   NewHash,
	}

	for name, newRng := range implementations {
		t.Run(name, func(t *testing.T) {
			testRngImplementation(t, newRng())
		})
	}
}

func testRngImplementation(t *testing.T, rng Rand) {
	t.Helper()

	// Test Float64 range
	t.Run("Float64_range", func(t *testing.T) {
		for i := 0; i < 1000; i++ {
			v := rng.Float64()
			if v < 0 || v >= 1 {
				t.Errorf("Float64() = %v, want [0,1)", v)
			}
		}
	})

	// Test Float32 range
	t.Run("Float32_range", func(t *testing.T) {
		for i := 0; i < 1000; i++ {
			v := rng.Float32()
			if v < 0 || v >= 1 {
				t.Errorf("Float32() = %v, want [0,1)", v)
			}
		}
	})

	// Test Int63 range
	// Test Int63 range
	t.Run("Int63_range", func(t *testing.T) {
		for i := 0; i < 1000; i++ {
			v := rng.Int63()
			if v < 0 || v >= (1<<63-1) { // Changed to max int63 value
				t.Errorf("Int63() = %v, want [0,%v)", v, int64(1<<63-1))
			}
		}
	})

	// Test Int63n
	t.Run("Int63n", func(t *testing.T) {
		n := int64(100)
		for i := 0; i < 1000; i++ {
			v := rng.Int63n(n)
			if v < 0 || v >= n {
				t.Errorf("Int63n(%v) = %v, want [0,%v)", n, v, n)
			}
		}
	})

	// Test Int31n
	t.Run("Int31n", func(t *testing.T) {
		n := int32(100)
		for i := 0; i < 1000; i++ {
			v := rng.Int31n(n)
			if v < 0 || v >= n {
				t.Errorf("Int31n(%v) = %v, want [0,%v)", n, v, n)
			}
		}
	})

	// Test Intn
	t.Run("Intn", func(t *testing.T) {
		n := 100
		for i := 0; i < 1000; i++ {
			v := rng.Intn(n)
			if v < 0 || v >= n {
				t.Errorf("Intn(%v) = %v, want [0,%v)", n, v, n)
			}
		}
	})

	// Test Perm
	t.Run("Perm", func(t *testing.T) {
		n := 10
		p := rng.Perm(n)
		if len(p) != n {
			t.Errorf("Perm(%v) length = %v, want %v", n, len(p), n)
		}
		seen := make([]bool, n)
		for _, v := range p {
			if v < 0 || v >= n {
				t.Errorf("Perm value out of range: %v", v)
			}
			if seen[v] {
				t.Errorf("Duplicate value in Perm: %v", v)
			}
			seen[v] = true
		}
	})

	// Test Shuffle
	t.Run("Shuffle", func(t *testing.T) {
		n := 10
		arr := make([]int, n)
		for i := range arr {
			arr[i] = i
		}
		rng.Shuffle(len(arr), func(i, j int) {
			arr[i], arr[j] = arr[j], arr[i]
		})
		seen := make([]bool, n)
		for _, v := range arr {
			if v < 0 || v >= n {
				t.Errorf("Shuffled value out of range: %v", v)
			}
			seen[v] = true
		}
		for i, v := range seen {
			if !v {
				t.Errorf("Missing value after Shuffle: %v", i)
			}
		}
	})
}

func TestSeed(t *testing.T) {
	t.Run("deterministic", func(t *testing.T) {
		seed := NewSeed(42)
		rng1 := NewUnsafeSeed(seed)
		rng2 := NewUnsafeSeed(seed)

		// Should generate same sequence
		for i := 0; i < 100; i++ {
			if rng1.Float64() != rng2.Float64() {
				t.Error("Same seed should generate same sequence")
			}
		}
	})

	t.Run("next_seed", func(t *testing.T) {
		seed := NewSeed(42)
		rng1 := NewUnsafeSeed(seed)
		rng2 := Clone(rng1)

		// Should generate different sequences
		different := false
		for i := 0; i < 100; i++ {
			if rng1.Float64() != rng2.Float64() {
				different = true
				break
			}
		}
		if !different {
			t.Error("Different seeds should generate different sequences")
		}
	})
}

func TestClone(t *testing.T) {
	implementations := map[string]func() Rand{
		"unsafe": NewUnsafe,
		"hash":   NewHash,
	}

	for name, newRng := range implementations {
		t.Run(name, func(t *testing.T) {
			original := newRng()
			clone := Clone(original)

			// Clones should generate different sequences
			different := false
			for i := 0; i < 100; i++ {
				if original.Float64() != clone.Float64() {
					different = true
					break
				}
			}
			if !different {
				t.Error("Clone should generate different sequence")
			}
		})
	}
}

// Distribution tests
func TestDistribution(t *testing.T) {
	implementations := map[string]func() Rand{
		"unsafe": NewUnsafe,
		"hash":   NewHash,
	}

	for name, newRng := range implementations {
		t.Run(name, func(t *testing.T) {
			rng := newRng()
			samples := 10000
			buckets := make([]int, 10)

			// Test Float64 distribution
			for i := 0; i < samples; i++ {
				v := rng.Float64()
				bucket := int(v * 10)
				if bucket == 10 {
					bucket = 9
				}
				buckets[bucket]++
			}

			// Check if roughly uniform
			expected := samples / 10
			tolerance := float64(expected) * 0.2 // 20% tolerance
			for i, count := range buckets {
				if math.Abs(float64(count-expected)) > tolerance {
					t.Errorf("Non-uniform distribution in bucket %d: got %d, expected %d±%d",
						i, count, expected, int(tolerance))
				}
			}
		})
	}
}

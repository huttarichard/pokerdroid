package river

import (
	"math"
	"os"
	"runtime"
	"sync"
	"testing"

	"github.com/pokerdroid/poker"
	"github.com/pokerdroid/poker/abs"
	"github.com/pokerdroid/poker/card"
	"github.com/pokerdroid/poker/equity"
	"github.com/pokerdroid/poker/equity/omp"
	"github.com/pokerdroid/poker/frand"
	"github.com/pokerdroid/poker/iso"
	"github.com/stretchr/testify/require"
)

func TestBuckets(t *testing.T) {
	cls := map[abs.Cluster]equity.Equity{
		abs.Cluster(0): {0, 1},
		abs.Cluster(1): {1, 2},
		abs.Cluster(2): {3, 4},
		abs.Cluster(3): {5, 6},
		abs.Cluster(4): {7, 8},
		abs.Cluster(5): {9, 10},
		abs.Cluster(6): {11, 12},
		abs.Cluster(7): {13, 14},
	}

	b := &Buckets{}
	for k, v := range cls {
		b.Add(k, v)
	}

	data, err := b.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}

	tmpfile, err := os.CreateTemp("", "test.db")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	_, err = tmpfile.Write(data)
	if err != nil {
		t.Fatal(err)
	}

	err = tmpfile.Close()
	if err != nil {
		t.Fatal(err)
	}

	lg := &poker.TestingLogger{T: t}

	fb, err := NewBucketsFromFile(tmpfile.Name(), lg)
	if err != nil {
		t.Fatal(err)
	}

	for k, v := range cls {
		if fb.Get(k) != v {
			t.Fatalf("expected %v, got %v", v, fb.Get(k))
		}
	}
}

func TestBucketsOMP(t *testing.T) {
	for i := 0; i < 1000; i++ {
		i := i * 1000
		cb := iso.River.Unindex(uint64(i))
		eq := ComputeEquity(cb)
		eq2 := omp.Equity(cb[0], cb[1], cb[2:], 2)

		diff := math.Max(float64(eq.Float32()[2]), 0.02)
		require.InDelta(t, eq.WinDraw(), eq2.WinDraw(), diff)
	}
}

func TestBucketsPrecision(t *testing.T) {
	// b := TstNewBuckets(t)

	var gg sync.WaitGroup
	w := runtime.NumCPU()
	gg.Add(w)
	r := frand.NewHash()

	var all float64

	for i := 0; i < w; i++ {
		go func() {
			var i int
			var dist float64
			for {
				c := card.RandomCards(r, 7)
				if len(c) > 7 {
					panic(c)
				}
				eq1 := omp.Equity(c[0], c[1], c[2:], 2)
				// eq2 := b.Get(abs.Cluster(iso.River.Index(c)))
				eq2 := ComputeEquity(c)

				dist += math.Abs(float64(eq1.Tie() - eq2.Tie()))
				i++
				if i > 1_000 {
					break
				}
			}
			all += dist / float64(i)
			gg.Done()
		}()
	}

	gg.Wait()
	t.Logf("avg dist: %f", all/float64(w))
}

func TestComputeEquity(t *testing.T) {
	cb := card.Cards{card.Card2H, card.Card3H, card.Card4H, card.Card5H, card.Card6H, card.Card7H, card.Card8H}
	eq := ComputeEquity(cb)
	eq2 := omp.Equity(cb[0], cb[1], cb[2:], 2)

	require.InDelta(t, eq.Tie(), eq2.Tie(), 0.01)
}

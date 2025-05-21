package river

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"runtime"
	"sync/atomic"
	"testing"
	"time"

	"github.com/edsrzf/mmap-go"
	"github.com/pokerdroid/poker"
	"github.com/pokerdroid/poker/abs"
	"github.com/pokerdroid/poker/card"
	"github.com/pokerdroid/poker/equity"
	"github.com/pokerdroid/poker/eval"
	"github.com/pokerdroid/poker/iso"
	"golang.org/x/sync/errgroup"
)

const Size = 123_156_254

// RankSize is the size in bytes of one key-value record: 8 bytes for the key
// plus 4 bytes for the value = 12 bytes total.
const RankSize = 4

type Buckets struct {
	db [Size]equity.Equity
}

// NewBucketsFromFile memory-maps the specified file and loads each (key:8 bytes, val:4 bytes)
// into the DB in parallel. This is purely for initialization-phase usage.
func NewBucketsFromFile(filename string, logger poker.Logger) (*Buckets, error) {
	start := time.Now()
	f, err := os.OpenFile(filename, os.O_RDONLY, 0)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	info, err := f.Stat()
	if err != nil {
		return nil, err
	}

	if info.Size()%RankSize != 0 {
		return nil, fmt.Errorf("file size %d is not multiple of %d", info.Size(), RankSize)
	}

	ww := runtime.NumCPU()
	logger.Printf("river load: workers=%d", ww)

	data, err := mmap.Map(f, mmap.RDONLY, 0)
	if err != nil {
		return nil, err
	}
	defer data.Unmap()

	buckets := new(Buckets)

	err = abs.IndexWorkers[int](Size, func(i int, done uint64) error {
		offset := i * RankSize

		v1 := binary.LittleEndian.Uint16(data[offset : offset+2])
		v2 := binary.LittleEndian.Uint16(data[offset+2 : offset+4])

		val := equity.Equity{v1, v2}
		buckets.db[i] = val

		return nil
	})
	if err != nil {
		return nil, err
	}

	duration := time.Since(start)
	rate := float64(Size) / duration.Seconds()
	logger.Printf("river load complete: entries=%d duration=%v rate=%.0f keys/sec",
		Size, duration.Round(time.Millisecond), rate)

	return buckets, nil
}

func TstNewBuckets(t *testing.T) *Buckets {
	buckets, err := NewBucketsFromFile(os.Getenv("BUCKETS_PATH"), &poker.TestingLogger{T: t})
	if err != nil {
		t.Fatal(err)
	}
	return buckets
}

func ComputeBuckets(logger poker.Logger) (*Buckets, error) {
	ww := runtime.NumCPU()
	logger.Printf("river: starting compute equities: workers=%d", ww)

	var wg errgroup.Group
	buckets := new(Buckets)
	epw := Size / ww

	var counter uint64

	// Worker goroutines
	for w := 0; w < ww; w++ {
		id := w
		wg.Go(func() error {
			start := id * epw
			end := start + epw

			if id == ww-1 {
				end = Size
			}

			for i := start; i < end; i++ {
				cb := iso.River.Unindex(uint64(i))
				buckets.db[abs.Cluster(i)] = ComputeEquity(cb)
				atomic.AddUint64(&counter, 1)

				if counter%10_000 == 0 {
					logger.Printf("river compute: %d/%d", counter, Size)
				}
			}
			return nil
		})
	}

	// Wait for completion
	err := wg.Wait()
	if err != nil {
		return nil, err
	}

	return buckets, nil
}

func ComputeEquity(cb card.Cards) equity.Equity {
	var chances [3]float32

	rank, err := eval.Eval(cb...)
	if err != nil {
		panic(err)
	}

	pbs := cb[2:]
	combos := card.CombinationsFrom(card.All(cb...), 2)
	cbsl := 1 / float32(len(combos))
	for _, c := range combos {
		if card.IsAnyMatch(c, pbs) {
			continue
		}
		otherRank, err := eval.Eval(append(c, pbs...)...)
		if err != nil {
			panic(err)
		}
		chances[rank.Compare(otherRank)] += cbsl
	}

	return equity.NewEquity(chances[0], chances[2])
}

func (c *Buckets) Get(k abs.Cluster) equity.Equity {
	return c.db[k]
}

func (c *Buckets) Add(k abs.Cluster, v equity.Equity) {
	c.db[k] = v
}

func (c *Buckets) Len() int {
	return len(c.db)
}

func (c *Buckets) Distance(c1, c2 abs.Cluster) float64 {
	eq1 := c.Get(c1)
	eq2 := c.Get(c2)
	return eq1.Distance(eq2)
}

func (c *Buckets) MarshalBinary() ([]byte, error) {
	bb := make([]byte, 0, 8*Size)

	for _, cv := range c.db {
		bb = binary.LittleEndian.AppendUint16(bb, uint16(cv[0]))
		bb = binary.LittleEndian.AppendUint16(bb, uint16(cv[1]))
	}
	return bb, nil
}

func (c *Buckets) UnmarshalBinary(data []byte) error {
	f := bytes.NewReader(data)
	length := len(data) / RankSize

	for i := 0; i < int(length); i++ {
		var v1 uint16
		var v2 uint16

		if err := binary.Read(f, binary.LittleEndian, &v1); err != nil {
			return err
		}
		if err := binary.Read(f, binary.LittleEndian, &v2); err != nil {
			return err
		}

		c.db[abs.Cluster(i)] = equity.Equity{v1, v2}
	}

	return nil
}

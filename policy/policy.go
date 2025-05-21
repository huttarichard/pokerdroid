package policy

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"sync"

	"github.com/pokerdroid/poker/encbin"
	"github.com/pokerdroid/poker/float/f64"
)

type Policy struct {
	sync.Mutex

	Iteration      uint64    `json:"iteration"`
	Strategy       []float64 `json:"strategy"`
	StrategyWeight float64   `json:"-"`
	RegretSum      []float64 `json:"regret_sum"`
	StrategySum    []float64 `json:"strategy_sum"`
	Baseline       []float64 `json:"baseline"`
}

func New(actions int) *Policy {
	return &Policy{
		Iteration:      0,
		Strategy:       f64.Uniform(actions),
		StrategyWeight: 0.0,
		RegretSum:      make([]float64, actions),
		StrategySum:    make([]float64, actions),
		Baseline:       make([]float64, actions),
	}
}

func (p *Policy) String() string {
	var buf bytes.Buffer
	buf.WriteString("Policy:\n")
	buf.WriteString(fmt.Sprintf("  Iteration: %d\n", p.Iteration))
	buf.WriteString(fmt.Sprintf("  StrategyWeight: %.6f\n", p.StrategyWeight))

	// Helper to format float slices
	formatSlice := func(name string, slice []float64) {
		buf.WriteString(fmt.Sprintf("  %s: [", name))
		for i, v := range slice {
			if i > 0 {
				buf.WriteString(", ")
			}
			buf.WriteString(fmt.Sprintf("%.4f", v))
		}
		buf.WriteString("]\n")
	}

	formatSlice("Strategy", p.Strategy)
	formatSlice("RegretSum", p.RegretSum)
	formatSlice("StrategySum", p.StrategySum)
	formatSlice("Baseline", p.Baseline)

	return buf.String()
}

func (p *Policy) AddRegrets(w float64, regrets []float64) {
	f64.AxpyUnitary(w, regrets, p.RegretSum)
}

func (p *Policy) AddStrategyWeight(w float64) {
	p.StrategyWeight += w
}

func (p *Policy) GetAverageStrategy() []float64 {
	avgStrat := make([]float64, len(p.StrategySum))

	total := f64.Sum(p.StrategySum)

	if total > 0 {
		f64.ScalUnitaryTo(avgStrat, 1.0/total, p.StrategySum)
		return avgStrat
	}

	for i := range avgStrat {
		avgStrat[i] = 1.0 / float64(len(avgStrat))
	}

	return avgStrat
}

func (p *Policy) BuildStrategy() {
	copy(p.Strategy, p.RegretSum)
	f64.MakePositive(p.Strategy)

	total := f64.Sum(p.Strategy)

	if total > 0 {
		f64.ScalUnitary(1.0/total, p.Strategy)
		return
	}
	p.Strategy = f64.Uniform(len(p.Strategy))
}

func (p *Policy) Calculate(gi uint64, dis Discounter) {
	d := dis(gi)
	// d := dis(p.Iteration)
	p.Iteration++

	if d.StrategySum != 1.0 {
		f64.ScalUnitary(d.StrategySum, p.StrategySum)
	}
	// Add strategy weight to strategy sum
	f64.AxpyUnitary(p.StrategyWeight, p.Strategy, p.StrategySum)
	// Apply regret matching
	f64.ScalUnitaryToUP(p.RegretSum, d.PositiveRegret, d.NegativeRegret, p.RegretSum)
	// Rebuild strategy
	p.BuildStrategy()
	// Reset strategy weight
	p.StrategyWeight = 0.0
}

func (p *Policy) Clone() *Policy {
	// Calculate total size needed for all slices
	totalLen := len(p.Strategy) * 4 // 4 slices of same length

	// Single allocation for all slices
	buffer := make([]float64, totalLen)

	// Partition buffer into slices
	strategy := buffer[:len(p.Strategy)]
	regretSum := buffer[len(p.Strategy) : len(p.Strategy)*2]
	strategySum := buffer[len(p.Strategy)*2 : len(p.Strategy)*3]
	baseline := buffer[len(p.Strategy)*3 : len(p.Strategy)*4]

	// Copy data
	copy(strategy, p.Strategy)
	copy(regretSum, p.RegretSum)
	copy(strategySum, p.StrategySum)
	copy(baseline, p.Baseline)

	return &Policy{
		Iteration:      p.Iteration,
		StrategyWeight: p.StrategyWeight,
		Strategy:       strategy,
		RegretSum:      regretSum,
		StrategySum:    strategySum,
		Baseline:       baseline,
	}
}

// Add after Clone() method
func (p *Policy) Size() uint64 {
	size := uint64(0)

	// Fixed size fields
	size += 8 // Iteration (uint64)
	size += 4 // length (uint32)

	// Each float32 slice (Strategy, RegretSum, StrategySum, Baseline)
	sliceSize := uint64(len(p.RegretSum)) * 8 // float32 = 4 bytes
	size += sliceSize * 3                     // 4 slices of same length

	return size
}

func (p *Policy) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)

	length := uint32(len(p.RegretSum))

	err := encbin.MarshalValues(buf, p.Iteration, length)
	if err != nil {
		return nil, err
	}

	for _, slice := range [][]float64{
		p.RegretSum,
		p.StrategySum,
		p.Baseline,
	} {
		if f64.IsNanInf(slice) {
			return nil, fmt.Errorf("slice contains NaN or Inf")
		}
		if err := binary.Write(buf, binary.LittleEndian, slice); err != nil {
			return nil, err
		}
	}

	return buf.Bytes(), nil
}

func (p *Policy) UnmarshalBinary(data []byte) error {
	r := bytes.NewReader(data)

	var length uint32

	err := encbin.UnmarshalValues(r, &p.Iteration, &length)
	if err != nil {
		return err
	}

	p.RegretSum = make([]float64, length)
	p.StrategySum = make([]float64, length)
	p.Baseline = make([]float64, length)
	p.Strategy = make([]float64, length)

	for _, slice := range [][]float64{
		p.RegretSum,
		p.StrategySum,
		p.Baseline,
	} {
		if err := binary.Read(r, binary.LittleEndian, slice); err != nil {
			return err
		}
	}

	p.BuildStrategy()
	return nil
}

package policy

import (
	"math"
	"math/bits"
)

type Discount struct {
	PositiveRegret float64
	NegativeRegret float64
	StrategySum    float64
}

type Discounter func(iter uint64) Discount

func CFRD(alpha, beta, gamma float64) Discounter {
	return func(iter uint64) (d Discount) {
		fi := float64(iter)

		// α = 2/3 - good starting point?
		// t^alpha / (t^alpha + 1)
		tA := math.Pow(fi, alpha)
		d.PositiveRegret = tA / (tA + 1.0)

		// t^beta / (t^beta + 1)
		tB := math.Pow(fi, beta)
		d.NegativeRegret = tB / (tB + 1.0)

		// t^gamma / (t^gamma + 1)
		tG := math.Pow(fi, gamma)
		d.StrategySum = tG / (tG + 1.0)

		msb := iter - MSBEven(iter)
		if msb == 0 {
			d.StrategySum = 0
		}
		return
	}
}

func CFRP(iter uint64) (d Discount) {
	sum := 1.0
	msb := iter - MSBEven(iter)
	if msb == 0 {
		sum = 0
	}

	d.PositiveRegret = 1.0
	d.NegativeRegret = 0
	d.StrategySum = sum

	return
}

func CFRL(iter uint64) (d Discount) {
	d.PositiveRegret = 1.0
	d.NegativeRegret = 0
	d.StrategySum = float64(iter) / float64(iter+1)
	return
}

// // γ = 3 - good starting point?
// // (t / (t+1)) ^ gamma
// msb := iter - MSBEven(iter)
// tG := float64(msb) / float64(msb+1)
// d.StrategySum = math.Pow(tG, gamma)

func MSBEven(iter uint64) uint64 {
	if iter == 0 {
		return 0
	}
	msbPos := 63 - bits.LeadingZeros64(iter)
	if msbPos%2 != 0 {
		msbPos--
	}
	return 1 << msbPos
}

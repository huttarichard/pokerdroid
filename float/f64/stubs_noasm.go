// Copyright ©2016 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !amd64 || noasm || appengine || safe
// +build !amd64 noasm appengine safe

package f64

// AxpyUnitary is add alpha to float64
func AxpyUnitary(alpha float64, x, y []float64) {
	for i, v := range x {
		y[i] += alpha * v
	}
}

// DotUnitary is dot product of float64
func DotUnitary(x, y []float64) (sum float64) {
	for i, v := range x {
		sum += y[i] * v
	}
	return sum
}

// ScalUnitary is multiply alpha to float64
func ScalUnitary(alpha float64, x []float64) {
	for i := range x {
		x[i] *= alpha
	}
}

// ScalUnitaryTo is multiply alpha to float64
func ScalUnitaryTo(dst []float64, alpha float64, x []float64) {
	for i, v := range x {
		dst[i] = alpha * v
	}
}

// ScalUnitaryToUP is multiply alphaUp or alphaDown
// based on the sign of float64
func ScalUnitaryToUP(dst []float64, alphaUp, alphaDown float64, x []float64) {
	for i, v := range x {
		if v > 0 {
			dst[i] = alphaUp * v
		} else if v < 0 {
			dst[i] = alphaDown * v
		} else {
			dst[i] = 0
		}
	}
}

// AddConst is add constant to float64
func AddConst(alpha float64, x []float64) {
	for i := range x {
		x[i] += alpha
	}
}

// Sum returns sum of float64
func Sum(x []float64) float64 {
	var sum float64
	for _, v := range x {
		sum += v
	}
	return sum
}

// MakePositive clears negative values
func MakePositive(v []float64) {
	for i := range v {
		if v[i] < 0 {
			v[i] = 0.0
		}
	}
}

// Preventzero will prevent value in float to beme zero
func PreventZero(v []float64, tol float64) {
	for i := range v {
		if v[i] == 0 {
			v[i] = tol
		}
	}
}

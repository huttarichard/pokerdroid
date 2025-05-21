// Copyright ©2016 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !amd64 || noasm || appengine || safe
// +build !amd64 noasm appengine safe

package f32

// AxpyUnitary is add alpha to float32
func AxpyUnitary(alpha float32, x, y []float32) {
	for i, v := range x {
		y[i] += alpha * v
	}
}

// DotUnitary is dot product of float32
func DotUnitary(x, y []float32) (sum float32) {
	for i, v := range x {
		sum += y[i] * v
	}
	return sum
}

// ScalUnitary is multiply alpha to float32
func ScalUnitary(alpha float32, x []float32) {
	for i := range x {
		x[i] *= alpha
	}
}

// ScalUnitaryToUP is multiply alphaUp or alphaDown
// based on the sign of float32
func ScalUnitaryToUP(dst []float32, alphaUp, alphaDown float32, x []float32) {
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

// AddConst is add constant to float32
func AddConst(alpha float32, x []float32) {
	for i := range x {
		x[i] += alpha
	}
}

// Sum returns sum of float32
func Sum(x []float32) float32 {
	var sum float32
	for _, v := range x {
		sum += v
	}
	return sum
}

// MakePositive clears negative values
func MakePositive(v []float32) {
	for i := range v {
		if v[i] < 0 {
			v[i] = 0.0
		}
	}
}

// Preventzero will prevent value in float to beme zero
func PreventZero(v []float32, tol float32) {
	for i := range v {
		if v[i] == 0 {
			v[i] = tol
		}
	}
}

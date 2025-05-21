// Copyright ©2016 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !noasm && !appengine && !safe
// +build !noasm,!appengine,!safe

package f64

// AxpyUnitary is
//
//	for i, v := range x {
//		y[i] += alpha * v
//	}
func AxpyUnitary(alpha float64, x, y []float64)

// DotUnitary is
//
//	for i, v := range x {
//		sum += y[i] * v
//	}
//	return sum
func DotUnitary(x, y []float64) (sum float64)

// ScalUnitary is
//
//	for i := range x {
//		x[i] *= alpha
//	}
func ScalUnitary(alpha float64, x []float64)

// ScalUnitaryTo is
//
//	for i, v := range x {
//		dst[i] = alpha * v
//	}
func ScalUnitaryTo(dst []float64, alpha float64, x []float64)

// ScalUnitaryToUP is
//
//	 for i, v := range x {
//			if v > 0 {
//				dst[i] = alphaUp * v
//			} else if v < 0 {
//				dst[i] = alphaDown * v
//			} else {
//				dst[i] = 0
//			}
//		}
func ScalUnitaryToUP(dst []float64, alphaUp, alphaDown float64, x []float64)

// AddConst is
//
//	for i := range x {
//		x[i] += alpha
//	}
func AddConst(alpha float64, x []float64)

// Sum is
//
//	var sum float64
//	for i := range x {
//	    sum += x[i]
//	}
func Sum(x []float64) float64

// MakePositive is
//
//	for i := range v {
//		if v[i] < 0 {
//			v[i] = 0.0
//		}
//	}
func MakePositive(v []float64)

// Preventzero is
//
//	for i := range v {
//		if v[i] == 0 {
//			v[i] = tol
//		}
//	}
func PreventZero(v []float64, tol float64)

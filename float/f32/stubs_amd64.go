// Copyright ©2016 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !noasm && !appengine && !safe
// +build !noasm,!appengine,!safe

package f32

// AxpyUnitary is
//
//	for i, v := range x {
//		y[i] += alpha * v
//	}
func AxpyUnitary(alpha float32, x, y []float32)

// DotUnitary is
//
//	for i, v := range x {
//		sum += y[i] * v
//	}
//	return sum
func DotUnitary(x, y []float32) (sum float32)

// ScalUnitary is
//
//	for i := range x {
//		x[i] *= alpha
//	}
func ScalUnitary(alpha float32, x []float32)

// ScalUnitaryTo is
//
//	for i, v := range x {
//		dst[i] = alpha * v
//	}
// func ScalUnitaryTo(dst []float32, alpha float32, x []float32)

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
func ScalUnitaryToUP(dst []float32, alphaUp, alphaDown float32, x []float32)

// AddConst is
//
//	for i := range x {
//		x[i] += alpha
//	}
func AddConst(alpha float32, x []float32)

// Sum is
//
//	var sum float32
//	for i := range x {
//	    sum += x[i]
//	}
func Sum(x []float32) float32

// MakePositive is
//
//	for i := range v {
//		if v[i] < 0 {
//			v[i] = 0.0
//		}
//	}
func MakePositive(v []float32)

// Preventzero is
//
//	for i := range v {
//		if v[i] == 0 {
//			v[i] = tol
//		}
//	}
func PreventZero(v []float32, tol float32)

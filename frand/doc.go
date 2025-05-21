// Package frand stand for fast random.
//
// It implements random number generator using maphash function.
// Original need arised from need to have fast concurrent safe
// number generation.
//
// While standart library's math/rand is fast, it is not safe
// for concurrent use.
//
// This package implements both safe random generator which
// can be seeded as well as hash generator which
// is safe and really fast but cannot be seeded.
package frand

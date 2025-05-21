package f64

import "sync"

type Slice struct {
	Slice []float64
}

func (s *Slice) HasCap(i int) bool {
	return cap(s.Slice) >= i
}

func (s *Slice) Clean(f float64) {
	for i := range s.Slice {
		s.Slice[i] = f
	}
}

type Pool struct {
	pool sync.Pool
	Min  int
}

func NewPool(min int) *Pool {
	px := &Pool{Min: min}
	px.pool.New = func() interface{} {
		return &Slice{Slice: make([]float64, min)}
	}
	return px
}

func (p *Pool) Alloc(n int) *Slice {
	slice := p.pool.Get().(*Slice)
	if !slice.HasCap(n) {
		nx := append(slice.Slice, make([]float64, n)...)
		return &Slice{Slice: nx[:n]}
	}
	slice.Slice = slice.Slice[:n]
	return slice
}

func (p *Pool) Free(s *Slice) {
	if s.HasCap(1) {
		s.Slice = s.Slice[:0]
		p.pool.Put(s)
	}
}

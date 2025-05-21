package policy

import (
	"sync"
)

type Update struct {
	uu []*Policy
}

func (r *Update) AddUpdate(n *Policy) {
	r.uu = append(r.uu, n)
}

func (r *Update) Len() int {
	return len(r.uu)
}

func (r *Update) Process(gi uint64, d Discounter) {
	for _, p := range r.uu {
		if p == nil {
			continue
		}
		p.Calculate(gi, d)
		p.Unlock()
	}
}

type UpdatePool struct {
	pool sync.Pool
	Min  int
}

func NewUpdatePool(min int) *UpdatePool {
	px := &UpdatePool{Min: min}
	px.pool.New = func() interface{} {
		return &Update{uu: make([]*Policy, 0, min)}
	}
	return px
}

func (p *UpdatePool) Alloc() *Update {
	s := p.pool.Get().(*Update)
	return s
}

func (p *UpdatePool) Free(s *Update) {
	s.uu = s.uu[:0]
	p.pool.Put(s)
}

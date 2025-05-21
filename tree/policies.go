package tree

import (
	"reflect"
	"sync"

	"github.com/pokerdroid/poker/abs"
	"github.com/pokerdroid/poker/policy"
)

type Policies struct {
	Map map[abs.Cluster]*policy.Policy
	mux sync.RWMutex
}

func NewPolicies() *Policies {
	return &Policies{
		Map: make(map[abs.Cluster]*policy.Policy, 0),
	}
}

func (p *Policies) Acquire(c abs.Cluster, size int) (*policy.Policy, bool) {
	px, ok := p.Get(c)
	if ok {
		px.Lock()
		return px, true
	}

	px = policy.New(size)
	px.Lock()
	p.Store(c, px)
	return px, false
}

func (p *Policies) Clone() *Policies {
	a := make(map[abs.Cluster]*policy.Policy, len(p.Map))
	for k, v := range p.Map {
		a[k] = v.Clone()
	}
	return &Policies{Map: a}
}

func (p *Policies) Equal(o *Policies) bool {
	return reflect.DeepEqual(p.Map, o.Map)
}

func (p *Policies) Get(cl abs.Cluster) (*policy.Policy, bool) {
	p.mux.RLock()
	px, ok := p.Map[cl]
	p.mux.RUnlock()
	return px, ok
}

func (p *Policies) Store(cl abs.Cluster, v *policy.Policy) {
	p.mux.Lock()
	p.Map[cl] = v
	p.mux.Unlock()
}

func (p *Policies) Delete(cl abs.Cluster) {
	p.mux.Lock()
	delete(p.Map, cl)
	p.mux.Unlock()
}

func (p *Policies) Len() uint32 {
	return uint32(len(p.Map))
}

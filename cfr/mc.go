package cfr

import (
	"fmt"
	"sync/atomic"

	"github.com/pokerdroid/poker/abs"
	"github.com/pokerdroid/poker/dealer"
	"github.com/pokerdroid/poker/float/f64"
	"github.com/pokerdroid/poker/policy"
	"github.com/pokerdroid/poker/policy/sampler"
	"github.com/pokerdroid/poker/tree"
)

type MCParams struct {
	Tree     *tree.Root
	Abs      abs.Mapper
	PS       sampler.Sampler
	TS       sampler.Sampler
	Discount policy.Discounter
	Sampler  dealer.Dealer
	BU       policy.BaselineUpdater
	Prune    float64
}

type MC struct {
	MCParams
	pool *f64.Pool
}

func NewMC(p MCParams) *MC {
	return &MC{
		MCParams: p,
		pool:     f64.NewPool(3),
	}
}

func (c *MC) Run(p Params) (ev float64, up uint64) {
	np := uint64(c.Tree.Params.NumPlayers)
	pl := policy.NewUpdatePool(128)

	for i := uint64(0); i < p.Iterations; i++ {
		sample, err := p.Sampler.Sample(p.Rng)
		if err != nil {
			panic(err)
		}

		update := pl.Alloc()
		tid := uint8(i % np)

		t := &Task{
			TraversingID: tid,
			Sample:       sample,
			Update:       update,
			Rng:          p.Rng,
		}

		ev += c.runHelper(c.Tree, t, 0, 1, 1)
		up += uint64(update.Len())

		// Perform discounting and free update.
		n := atomic.AddUint64(&c.Tree.Iteration, 1)
		update.Process(n, c.Discount)

		// Free nodes and sample.
		c.Sampler.Put(sample)
		pl.Free(update)
	}

	return ev / float64(p.Iterations), up
}

func (c *MC) runHelper(node tree.Node, t *Task, depth uint8, sample, reach float64) (ev float64) {
	tree.MustExpand(c.Tree, node)

	switch x := node.(type) {
	case *tree.Root:
		ev = c.runHelper(c.Tree.Next, t, depth, sample, reach)

	case *tree.Chance:
		t.Sample.Sample(x.State.Street)
		ev = c.runHelper(x.Next, t, depth, sample, reach)

	case *tree.Terminal:
		ev = t.Sample.Utility(x, t.TraversingID)

	case tree.DecisionPoint:
		turn := x.GetTurnPos()

		if turn == t.TraversingID {
			ev = c.traverse(x, t, depth+1, sample, reach)
		} else {
			ev = c.sampling(x, t, depth+1, sample, reach)
		}

	default:
		panic(fmt.Sprintf("unknown node: %T", x))
	}

	return ev
}

func (c *MC) traverse(node tree.DecisionPoint, t *Task, depth uint8, sample, reach float64) float64 {
	acts := node.Len()
	cluster := t.Sample.Cluster(node, c.Abs)

	px := node.Acquire(c.Tree, cluster)

	regrets := c.pool.Alloc(acts)
	qs := c.pool.Alloc(acts)

	c.PS.Sample(t.Rng, acts, px, c.Tree.Iteration, depth, qs.Slice)

	// Prune if negative and every 10% iterations allow to pass.
	prune := c.Prune < 0 && px.Iteration%10 != 0

	for i, q := range qs.Slice {
		var util float64
		uHat := px.Baseline[i]

		if q <= 0 {
			regrets.Slice[i] = uHat
			continue
		}

		if prune && px.RegretSum[i] < c.Prune {
			regrets.Slice[i] = uHat
			continue
		}

		util = c.runHelper(
			node.GetNode(i),
			t,
			depth,
			sample*float64(q),
			reach,
		)
		uHat += (util - uHat) / q

		regrets.Slice[i] = uHat
		c.BU(px, 1./q, i, util)
	}

	cfv := f64.DotUnitary(px.Strategy, regrets.Slice)
	f64.AddConst(-cfv, regrets.Slice)

	px.AddRegrets(float64(reach/sample), regrets.Slice)
	t.Update.AddUpdate(px)

	c.pool.Free(regrets)
	c.pool.Free(qs)
	return cfv
}

// Sample player action according to strategy, do not update policy.
// Save selected action so that they are reused if this infoset is hit again.
func (c *MC) sampling(node tree.DecisionPoint, t *Task, depth uint8, sample, reach float64) float64 {
	actionsLen := node.Len()

	cluster := t.Sample.Cluster(node, c.Abs)

	px := node.Acquire(c.Tree, cluster)

	// Update average strategy for this node.
	// We perform "stochastic" updates as described in the MC-CFR paper.
	if sample > 0 {
		px.AddStrategyWeight(1. / sample)
	}

	qs := c.pool.Alloc(actionsLen)
	regrets := c.pool.Alloc(actionsLen)

	c.TS.Sample(t.Rng, actionsLen, px, c.Tree.Iteration, depth, qs.Slice)

	for i, q := range qs.Slice {
		var util float64
		uHat := px.Baseline[i]

		if q <= 0 {
			regrets.Slice[i] = uHat
			continue
		}

		rch := reach * float64(px.Strategy[i])
		util = c.runHelper(
			node.GetNode(i),
			t,
			depth,
			sample*float64(q),
			rch,
		)
		uHat += (util - uHat) / q

		regrets.Slice[i] = uHat
		c.BU(px, 1./q, i, util)
	}

	cfv := f64.DotUnitary(px.Strategy, regrets.Slice)

	t.Update.AddUpdate(px)

	c.pool.Free(regrets)
	c.pool.Free(qs)

	return cfv
}

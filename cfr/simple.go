package cfr

import (
	"fmt"
	"sync/atomic"

	"github.com/pokerdroid/poker/abs"
	"github.com/pokerdroid/poker/float/f64"
	"github.com/pokerdroid/poker/frand"
	"github.com/pokerdroid/poker/policy"
	"github.com/pokerdroid/poker/tree"
)

type SimpleMCParams struct {
	Tree     *tree.Root
	Discount policy.Discounter
	Abs      abs.Mapper
	BU       policy.BaselineUpdater
}

// This is also variant of MC-SimpleMC but without sampling.
type SimpleMC struct {
	SimpleMCParams
	pool *f64.Pool
}

func NewSimpleMC(p SimpleMCParams) *SimpleMC {
	c := &SimpleMC{
		SimpleMCParams: p,
		pool:           f64.NewPool(3),
	}
	return c
}

func (c *SimpleMC) Run(p Params) (ev float64, up uint64) {
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

		ev += c.runHelper(c.Tree, tid, t)
		up += uint64(update.Len())

		// Perform discounting and free update.
		n := atomic.AddUint64(&c.Tree.Iteration, 1)
		update.Process(n, c.Discount)

		// Free nodes and sample.
		p.Sampler.Put(sample)
		pl.Free(update)
	}

	return ev / float64(p.Iterations), up
}

func (c *SimpleMC) runHelper(node tree.Node, lp uint8, t *Task) (ev float64) {
	tree.MustExpand(c.Tree, node)

	switch x := node.(type) {
	case *tree.Root:
		ev = c.runHelper(c.Tree.Next, lp, t)

	case *tree.Chance:
		t.Sample.Sample(x.State.Street)
		ev = c.runHelper(x.Next, lp, t)

	case *tree.Terminal:
		ev = t.Sample.Utility(x, lp)

	case tree.DecisionPoint:
		if x.GetTurnPos() == t.TraversingID {
			ev = c.traverse(x, t)
		} else {
			ev = c.sampling(x, t)
		}

	default:
		panic(fmt.Sprintf("unknown node: %T", x))
	}

	return ev
}

func (c *SimpleMC) traverse(node tree.DecisionPoint, t *Task) float64 {
	cluster := t.Sample.Cluster(node, c.Abs)
	turnp := node.GetTurnPos()

	px := node.Acquire(c.Tree, cluster)

	aln := node.Len()
	regrets := c.pool.Alloc(aln)

	for i := 0; i < aln; i++ {
		node := node.GetNode(i)
		uHat := px.Baseline[i]

		util := c.runHelper(node, turnp, t)
		uHat += (util - uHat)

		regrets.Slice[i] = uHat
		c.BU(px, 1.0, i, util)
	}

	cfv := f64.DotUnitary(px.Strategy, regrets.Slice)
	f64.AddConst(-cfv, regrets.Slice)

	px.AddRegrets(1., regrets.Slice)
	t.Update.AddUpdate(px)
	c.pool.Free(regrets)

	return cfv
}

// Sample player action according to strategy, do not update policy.
// Save selected action so that they are reused if this infoset is hit again.
func (c *SimpleMC) sampling(node tree.DecisionPoint, t *Task) float64 {
	cluster := t.Sample.Cluster(node, c.Abs)
	px := node.Acquire(c.Tree, cluster)

	// Update average strategy for this node.
	// We perform "stochastic" updates as described in the MC-CFR paper.
	px.AddStrategyWeight(1.)

	t.Update.AddUpdate(px)

	idx := frand.SampleIndex(t.Rng, px.Strategy, 0.0001)

	return c.runHelper(node.GetNode(idx), node.GetTurnPos(), t)
}

package tree

import (
	"runtime"
	"sync"
	"sync/atomic"
)

// DiscardBelowEpsilon traverses the game tree from root and discards (sets to nil)
// any subtrees whose “reach probability” is below eps. Now done via a parallel
// traversal with a small worker pool. We avoid closing the channel prematurely
// to fix the “send on closed channel” panic.
func DiscardBelowEpsilon(root *Root, eps float64, concurrency int) uint32 {
	if root == nil {
		return 0
	}

	// We do nothing if eps<=0 or if we already marked the root fully expanded
	if eps <= 0 || root.Full {
		return 0
	}

	start := root.States

	// Channel of nodes to process; each carries a “node” plus the
	// probability factor that got us to that node.
	if concurrency < 1 {
		concurrency = 1
	}

	// tasks is the queue of work to be done,
	// tasksWG counts how many nodeProbs we've enqueued but not yet finished.
	tasks := make(chan *Player, concurrency)

	// Worker pool: each worker loops over the tasks channel.
	// This wait group will ensure all workers finish before we exit.
	var workersWG sync.WaitGroup
	workersWG.Add(concurrency)

	var pruned atomic.Uint32

	for w := 0; w < concurrency; w++ {
		go func() {
			defer workersWG.Done()

			// Need to prone the policies
		}()
	}

	MustVisit(root, -1, func(n Node, children []Node, depth int) bool {
		p, ok := n.(*Player)
		if !ok {
			return true
		}

		if p.Actions == nil {
			return true
		}

		tasks <- p
		return true
	})

	close(tasks)

	// Wait for all workers to drain the channel and exit.
	workersWG.Wait()

	// Count states again. This is necessary because the workers
	// might have
	root.States -= pruned.Load()

	// Return how many nodes have been discarded.
	return start - root.States
}

// DiscardReferenceAtDepth traverses the game tree from root and
// discards all reference nodes below given level.
func DiscardReferenceAtDepth(root *Root, at int) {
	walk := func(n Node, children []Node, depth int) bool {
		rf, ok := n.(*Reference)
		if !ok {
			return true
		}
		if depth > at && at != -1 {
			rf.Node = nil
			return false
		}
		return true
	}
	MustVisit(root, -1, walk)
	runtime.GC()
}

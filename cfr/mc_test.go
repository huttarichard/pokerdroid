package cfr

import (
	"bytes"
	"context"
	"testing"

	"github.com/pokerdroid/poker"
	kuhndealer "github.com/pokerdroid/poker/dealer/kuhn"
	"github.com/pokerdroid/poker/frand"
	"github.com/pokerdroid/poker/policy"
	"github.com/pokerdroid/poker/policy/sampler"
	"github.com/pokerdroid/poker/tree"
	"github.com/pokerdroid/poker/tree/profiling"
)

func TestPureMCKuhn(t *testing.T) {
	root := tree.NewKuhn()

	leafs := tree.FindLeafRunes(root)
	t.Logf("leafs: %d", len(leafs))

	r := frand.NewUnsafeInt(0)
	dealer := kuhndealer.NewGameSampler(r)

	cfrmc := NewMC(MCParams{
		PS: sampler.NewExternal(),
		TS: sampler.NewExternal(),

		Tree:     root,
		Discount: policy.CFRP,
		Abs:      kuhndealer.Clusters,
		Sampler:  dealer,

		BU:    policy.BaselineEMA(0.01),
		Prune: 0,
	})

	rprms := NewRunParams(root, dealer, kuhndealer.Clusters)
	rprms.SetBatch(1000, 1)
	rprms.SetEpochs(100)
	rprms.Rng = r
	rprms.Logger = &poker.TestingLogger{T: t}

	Run(context.Background(), cfrmc, rprms)

	exploit := Exploit(context.Background(), ExploitParams{
		Root:       root,
		Params:     root.Params,
		Sampler:    dealer,
		Rng:        r,
		Iterations: 10_000,
		Abs:        kuhndealer.Clusters,
	})

	t.Logf("exploitability: %f", exploit)

	buf := new(bytes.Buffer)

	st := profiling.GetKuhnProfile(root)

	buf.WriteString("\nAverage policy: \n")
	st.WriteTable(profiling.FormatAveragePolicy, buf)

	buf.WriteString("\nBaseline\n")
	st.WriteTable(profiling.FormatBaseline, buf)

	buf.WriteString("\nCurrent policy: \n")
	st.WriteTable(profiling.FormatCurrentPolicy, buf)

	buf.WriteString("\nRegrets summary\n")
	st.WriteTable(profiling.FormatRegretSummary, buf)

	buf.WriteString("\nStrategy summary\n")
	st.WriteTable(profiling.FormatStrategySummary, buf)

	buf.WriteString("\n\n")

	t.Log("\n" + buf.String())
}

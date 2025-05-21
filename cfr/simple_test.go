package cfr

import (
	"bytes"
	"context"
	"testing"

	"github.com/pokerdroid/poker"
	"github.com/pokerdroid/poker/card"
	"github.com/pokerdroid/poker/frand"
	"github.com/pokerdroid/poker/policy"
	"github.com/pokerdroid/poker/tree"
	"github.com/pokerdroid/poker/tree/profiling"

	kuhndealer "github.com/pokerdroid/poker/dealer/kuhn"
)

func TestSimpleKuhn(t *testing.T) {
	root := tree.NewKuhn()

	leafs := tree.FindLeafRunes(root)
	t.Logf("leafs: %d", len(leafs))

	r := frand.NewUnsafeInt(0)
	dealer := kuhndealer.NewGameSampler(r)

	cfrmc := NewSimpleMC(SimpleMCParams{
		Tree:     root,
		Discount: policy.CFRP,
		Abs:      kuhndealer.Clusters,
		BU:       policy.BaselineEMA(0.5),
	})

	rprms := NewRunParams(root, dealer, kuhndealer.Clusters)
	rprms.SetBatch(10000, 8)
	rprms.SetEpochs(10)
	rprms.Logger = &poker.TestingLogger{T: t}
	rprms.Rng = r

	Run(context.Background(), cfrmc, rprms)

	exploit := Exploit(context.Background(), ExploitParams{
		Sampler:    dealer,
		Params:     root.Params,
		Rng:        r,
		Root:       root,
		Iterations: 1000,
		Abs:        kuhndealer.Clusters,
		Workers:    10,
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

// Average policy:
// PATH      | CARDS |  FOLD   | CHECK  |  CALL   | ALLIN | BET POT 1.00
// ----------------+-------+---------+--------+---------+-------+---------------
// r:n:p         | [jd]  | -       | 79.54% | -       | -     | 20.46%
// r:n:p         | [qd]  | -       | 99.74% | -       | -     | 0.26%
// r:n:p         | [kd]  | -       | 38.07% | -       | -     | 61.93%
// r:n:k:p       | [jd]  | -       | 65.86% | -       | -     | 34.14%
// r:n:k:p       | [qd]  | -       | 99.99% | -       | -     | 0.01%
// r:n:k:p       | [kd]  | -       | 0.01%  | -       | -     | 100.00%
// r:n:k:b2.00:p | [jd]  | 100.00% | -      | 0.01%   | -     | -
// r:n:k:b2.00:p | [qd]  | 46.04%  | -      | 53.96%  | -     | -
// r:n:k:b2.00:p | [kd]  | 0.01%   | -      | 100.00% | -     | -
// r:n:b2.00:p   | [jd]  | 100.00% | -      | 0.01%   | -     | -
// r:n:b2.00:p   | [qd]  | 65.78%  | -      | 34.22%  | -     | -
// r:n:b2.00:p   | [kd]  | 0.01%   | -      | 99.99%  | -     | -

func TestBestResponse(t *testing.T) {
	root := tree.NewKuhn()

	leafs := tree.FindLeafRunes(root)
	t.Logf("leafs: %d", len(leafs))

	r := frand.NewUnsafeInt(0)
	dealer := kuhndealer.NewGameSampler(r)

	cfrmc := NewSimpleMC(SimpleMCParams{
		Tree:     root,
		Discount: policy.CFRP,
		Abs:      kuhndealer.Clusters,
		BU:       policy.BaselineEMA(0.1),
	})

	rprms := NewRunParams(root, dealer, kuhndealer.Clusters)
	rprms.SetBatch(100, 1)
	rprms.SetEpochs(1000)
	rprms.Rng = r

	Run(context.Background(), cfrmc, rprms)

	st := profiling.GetKuhnProfile(root)

	buf := new(bytes.Buffer)
	buf.WriteString("\nAverage policy: \n")
	st.WriteTable(profiling.FormatAveragePolicy, buf)

	t.Log("\n" + buf.String())

	brr := &BR{
		game: root,
		abs:  kuhndealer.Clusters,
	}

	evr := &EV{
		game: root,
		abs:  kuhndealer.Clusters,
	}

	var samples int
	total := float64(0.0)

	for _, smpl := range []*kuhndealer.Sample{
		{Cards: card.Cards{card.CardKD, card.CardJD}},
		{Cards: card.Cards{card.CardJD, card.CardKD}},

		{Cards: card.Cards{card.CardKD, card.CardQD}},
		{Cards: card.Cards{card.CardQD, card.CardKD}},

		{Cards: card.Cards{card.CardQD, card.CardJD}},
		{Cards: card.Cards{card.CardJD, card.CardQD}},
	} {

		t.Logf("sample: %v", smpl.Cards.String())

		br0 := brr.Get(frand.NewUnsafeInt(0), smpl, 0)
		t.Logf("br0: %f", br0)

		ev0 := evr.Get(frand.NewUnsafeInt(0), smpl, 0)
		t.Logf("ev0: %f", ev0)

		br1 := brr.Get(frand.NewUnsafeInt(0), smpl, 1)
		t.Logf("br1: %f", br1)

		ev1 := evr.Get(frand.NewUnsafeInt(0), smpl, 1)
		t.Logf("ev1: %f", ev1)

		t.Logf("br0 - ev0: %f", br0-ev0)

		t.Logf("br1 - ev1: %f", br1-ev1)

		t.Logf("((br0 - ev0) + (br1 - ev1)) / 2: %f", ((br0-ev0)+(br1-ev1))/2)

		total += ((br0 - ev0) - (br1 - ev1)) / 2
		samples++
		t.Logf("exploitability: %f", total/float64(samples))
	}

	// 6 possible samples, 2 players
	t.Logf("total exploitability: %f", total/float64(samples)) // 0.055 <- correct
}

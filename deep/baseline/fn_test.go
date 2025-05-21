package baselinenn

import (
	_ "github.com/nlpodyssey/spago/optimizers/adam"
)

// func TestMCNLH(t *testing.T) {
// 	var err error

// 	ehs := p2.TstNew(t)
// 	abs := cluster.TstNewAbs(t)

// 	// r := frand.NewUnsafeInt(0)

// 	prms := table.NewGameParams(2, chips.NewFromInt(100))
// 	prms.BetSizes = []float32{0.5, 1, 1.5}
// 	// prms.MinBet = true
// 	// prms.BetSizes = []float32{1}
// 	prms.TerminalStreet = table.River
// 	prms.MaxActionsPerRound = 5
// 	prms.DisableV = true

// 	s := table.NewState(prms)

// 	s, err = table.MakeInitialBets(prms, s)
// 	require.NoError(t, err)

// 	root := &tree.Root{
// 		Params:    prms,
// 		Iteration: 0,
// 		State:     s,
// 		Next:      nil,
// 	}

// 	// px := holdemsamples.SamplerParams{
// 	// 	NumPlayers: 2,
// 	// 	Clusters:   absg,
// 	// }

// 	sampler := holdemsamples.NewSingleHanded(holdemsamples.SingleHandedParams{
// 		NumPlayers: 2,
// 		Clusters:   absg,
// 		Hand:       card.NewCardsFromString("ac ah"),
// 		Player:     0,
// 	})

// 	m := ModelConfig[float32]{
// 		NumPlayers:         2,
// 		MaxActionsPerRound: 5,
// 		HiddenSize:         512,
// 		HiddenLayers:       6,
// 		Dropout:            0,
// 	}

// 	model := m.Model()

// 	ff, err := os.Open("./model.gob")
// 	require.NoError(t, err)

// 	err = model.Load(ff)
// 	require.NoError(t, err)

// 	cfrmc := cfr.NewMC(cfr.MCParams{
// 		Tree: root,
//      PS: sampler.NewExternal(),
//      TS: sampler.NewExternal(),

// 		DiscountParams: policy.DiscountParams{
// 			DiscountAlpha: 1.5,
// 			DiscountBeta:  0,
// 			DiscountGamma: 3,
// 		},
// 		PruneT:  0,
// 		Limiter: Baseline[float32](5, abs.Groups, prms, model),
// 	})

// 	rprms := cfr.NewRunParams(root, sampler)
// 	rprms.SetBatch(2000, 1)
// 	rprms.SetEpochs(600)
// 	rprms.Logger = &poker.TestingLogger{T: t}

// 	cfr.Run(context.Background(), cfrmc, rprms)

// 	buf := new(bytes.Buffer)

// 	// exploit := cfr.Exploit(context.Background(), cfr.ExploitParams{
// 	// 	Root:       root,
// 	// 	Sampler:    sampler,
// 	// 	Rng:        r,
// 	// 	Iterations: 10_000,
// 	// })

// 	// t.Logf("exploitability: %f", exploit)

// 	st, err := profiling.New(profiling.Params{
// 		PlayerID: 0,
// 		Tree:     root,
// 		Depth:    3,
// 		Board:    card.Cards{},
// 		Clusters: absg,
// 		BetSizes: root.Params.BetSizes,
// 	})
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	profiling.SaveVisualization(st, "./tree.png")

// 	buf.WriteString("\nAverage policy: \n")
// 	st.WriteTable(profiling.FormatAveragePolicy, buf)

// 	buf.WriteString("\nCurrent policy: \n")
// 	st.WriteTable(profiling.FormatCurrentPolicy, buf)

// 	buf.WriteString("\nRegrets summary\n")
// 	st.WriteTable(profiling.FormatRegretSummary, buf)

// 	buf.WriteString("\nStrategy summary\n")
// 	st.WriteTable(profiling.FormatStrategySummary, buf)

// 	buf.WriteString("\nBaseline\n")
// 	st.WriteTable(profiling.FormatBaseline, buf)

// 	buf.WriteString("\tReach\n")
// 	st.WriteTable(profiling.FormatAvgReach, buf)

// 	t.Log("\n" + buf.String())
// }

// func TestBaselines(t *testing.T) {
// 	// ehs := h2.TstNew(t)
// 	abs := cluster.TstNewAbs(t)
// 	// absg := cluster.NewAbsGetter(ehs, abs.Map)

// 	f, err := os.Open("/Volumes/T7/experiments/50bb_p2/tree.bin")
// 	require.NoError(t, err)

// 	root, err := tree.NewRootFromReadSeeker(f)
// 	require.NoError(t, err)

// 	m := ModelConfig[float32]{
// 		NumPlayers:         2,
// 		MaxActionsPerRound: 5,
// 		HiddenSize:         512,
// 		HiddenLayers:       6,
// 		// Dropout:            0.3,
// 		// Dropout: 0,
// 	}
// 	model := m.Model()

// 	ff, err := os.Open("./model.gob")
// 	require.NoError(t, err)

// 	err = model.Load(ff)
// 	require.NoError(t, err)

// 	diff, samples := Validate(ValidateParams{
// 		Root:    root,
// 		Buckets: abs.Groups,
// 		Model:   model,
// 		Depth:   8,
// 		Rng:     frand.NewUnsafeInt(40),
// 		Rounds:  50,
// 	})

// 	pretty.Println(diff, samples)
// }

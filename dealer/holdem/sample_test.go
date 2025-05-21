package holdemdealer

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/pokerdroid/poker/card"
	"github.com/pokerdroid/poker/chips"
	"github.com/pokerdroid/poker/eval"
	"github.com/pokerdroid/poker/frand"
	"github.com/pokerdroid/poker/table"
	"github.com/pokerdroid/poker/tree"
	"github.com/stretchr/testify/require"
)

func TestSamples2(t *testing.T) {
	params := SamplerParams{}

	hnd := New(params)
	tm := time.Now()
	for i := 0; i < 1_000_000; i++ {
		_, err := hnd.Sample(frand.NewUnsafeInt(0))
		require.NoError(t, err)
	}
	t.Log(time.Since(tm))
}

func TestUtility1(t *testing.T) {
	params := SamplerParams{NumPlayers: 2}

	hnd := New(params)
	sx, err := hnd.Sample(frand.NewUnsafeInt(42))
	require.NoError(t, err)

	sml := sx.(*Sample)

	tm := &tree.Terminal{
		Pots: table.Pots{
			{Amount: chips.New(200), Players: []uint8{0, 1}},
		},
		Players: table.Players{
			{Paid: chips.New(100), Status: table.StatusAllIn},
			{Paid: chips.New(100), Status: table.StatusAllIn},
		},
	}

	t.Log(sml.Cards(0).String())
	t.Log(sml.Cards(1).String())

	require.Equal(t, sx.Utility(tm, 0), float64(100))
	require.Equal(t, sx.Utility(tm, 1), float64(-100))
}

func TestUtility2(t *testing.T) {
	s := &Sample{
		rng: frand.NewUnsafeInt(42),
		cur: table.Preflop,
	}

	tm := &tree.Terminal{
		Pots: table.Pots{
			{Amount: chips.New(3), Players: []uint8{0}},
		},
		Players: table.Players{
			{Paid: chips.New(1), Status: table.StatusActive},
			{Paid: chips.New(2), Status: table.StatusFolded},
		},
	}

	require.Equal(t, s.Utility(tm, 0), float64(2))
	require.Equal(t, s.Utility(tm, 1), float64(-2))
}

func TestUtility3(t *testing.T) {
	params := SamplerParams{NumPlayers: 2}
	rng := frand.NewUnsafeInt(1)
	hnd := New(params)

	for i := 0; i < 100; i++ {
		amount := rng.Intn(1000)

		smpl, err := hnd.Sample(rng)
		require.NoError(t, err)

		sml := smpl.(*Sample)
		sml.Sample(table.River)

		jdg, err := eval.Judge([]card.Cards{sml.Cards(0), sml.Cards(1)})
		require.NoError(t, err)

		if len(jdg) == 2 {
			continue
		}

		tm := &tree.Terminal{
			Pots: table.Pots{
				{Amount: chips.New(amount * 2), Players: []uint8{0, 1}},
			},
			Players: table.Players{
				{Paid: chips.New(amount), Status: table.StatusAllIn},
				{Paid: chips.New(amount), Status: table.StatusAllIn},
			},
		}

		util := smpl.Utility(tm, jdg[0])
		require.Equal(t, util, float64(amount))

		var other = 0
		if jdg[0] == 0 {
			other = 1
		}
		util = smpl.Utility(tm, uint8(other))
		require.Equal(t, util, float64(amount*-1))
	}
}

func TestSampleUniqueness(t *testing.T) {
	params := SamplerParams{NumPlayers: 2}

	s := New(params)

	const numSamples = 1000
	samples := make(map[string]struct{})
	var mu sync.Mutex

	r := frand.NewUnsafeInt(1)

	var wg sync.WaitGroup
	for i := 0; i < numSamples; i++ {
		wg.Add(1)
		go func(r frand.Rand) {
			defer wg.Done()
			sample, err := s.Sample(r)
			require.NoError(t, err)
			gs := sample.(*Sample)
			gs.Sample(table.River)

			key := fmt.Sprintf("%v%v", gs.Cards(0), gs.Cards(1))
			mu.Lock()
			if _, exists := samples[key]; exists {
				t.Errorf("Duplicate sample detected: %s", key)
			}
			samples[key] = struct{}{}
			mu.Unlock()

			s.Put(sample)
		}(frand.Clone(r))
	}

	wg.Wait()
}

func TestSampleChanceNode(t *testing.T) {
	params := SamplerParams{NumPlayers: 2}

	hnd := New(params)
	rng := frand.NewUnsafeInt(42)
	sx, err := hnd.Sample(rng)
	require.NoError(t, err)

	sample := sx.(*Sample)

	// Initially, the board should be empty
	require.Empty(t, sample.Board(), "Board should be empty initially")

	// Sample the flop
	sample.Sample(table.Flop)

	// After sampling the flop, we should have 3 cards on the board
	require.Len(t, sample.Board(), 3, "Board should have 3 cards after flop")

	// Sample the turn
	sample.Sample(table.Turn)

	// After sampling the turn, we should have 4 cards on the board
	require.Len(t, sample.Board(), 4, "Board should have 4 cards after turn")

	// Create a chance node for the river
	sample.Sample(table.River)

	// After sampling the river, we should have 5 cards on the board
	require.Len(t, sample.Board(), 5, "Board should have 5 cards after river")

	// Verify all cards are unique
	allCards := make(card.Cards, 0, 9) // 2 hole cards per player + 5 board cards
	allCards = append(allCards, sample.hands[0][:2]...)
	allCards = append(allCards, sample.hands[1][:2]...)
	allCards = append(allCards, sample.Board()...)

	require.Equal(t, 9, len(allCards), "Should have 9 total cards")

	// Check for duplicates
	seen := make(map[card.Card]bool)
	for _, c := range allCards {
		if seen[c] {
			t.Errorf("Duplicate card found: %v", c)
		}
		seen[c] = true
	}
}

func TestSampleChanceNodeSequential(t *testing.T) {
	params := SamplerParams{NumPlayers: 2}

	hnd := New(params)
	rng := frand.NewUnsafeInt(42)
	sx, err := hnd.Sample(rng)
	require.NoError(t, err)

	sample := sx.(*Sample)

	// Test sampling each street in sequence
	streets := []table.Street{table.Flop, table.Turn, table.River}
	expectedBoardSizes := []int{3, 4, 5}

	for i, street := range streets {
		sample.Sample(street)

		require.Len(t, sample.Board(), expectedBoardSizes[i],
			"Board should have %d cards after %s", expectedBoardSizes[i], street)
	}
}

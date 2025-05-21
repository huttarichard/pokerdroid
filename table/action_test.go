package table

import (
	"testing"

	"github.com/kr/pretty"
	"github.com/pokerdroid/poker/chips"
	"github.com/stretchr/testify/require"
)

func TestActionKind_Basics(t *testing.T) {
	tests := []struct {
		input    string
		expected ActionKind
	}{
		{"fold", Fold},
		{"check", Check},
		{"call", Call},
		{"raise", Raise},
		{"bet", Bet},
		{"sb", SmallBlind},
		{"bb", BigBlind},
		{"allin", AllIn},
	}
	for _, te := range tests {
		ak, err := NewActionFromString(te.input)
		require.NoError(t, err)
		require.Equal(t, te.expected, ak)
		require.NotEmpty(t, ak.String())
		require.NotEmpty(t, ak.Short())
	}

	_, err := NewActionFromString("invalid")
	require.Error(t, err)
}

func TestDiscreteAction_GetAction(t *testing.T) {
	p := GameParams{
		NumPlayers:    2,
		InitialStacks: chips.NewList(100, 100),
	}
	s := &State{
		Players:    []Player{{Paid: chips.NewFromInt(0)}, {Paid: chips.NewFromInt(0)}},
		TurnPos:    0,
		CallAmount: chips.NewFromInt(5),
	}

	// Fold
	ak, amt := DiscreteAction(DFold).GetAction(p, s)
	require.Equal(t, Fold, ak)
	require.True(t, amt.Equal(chips.Zero))

	// Check
	s.CallAmount = chips.Zero
	ak, amt = DiscreteAction(DCheck).GetAction(p, s)
	require.Equal(t, Check, ak)
	require.True(t, amt.Equal(chips.Zero))

	// Call
	s.CallAmount = chips.NewFromInt(10)
	ak, amt = DiscreteAction(DCall).GetAction(p, s)
	require.Equal(t, Call, ak)
	require.Equal(t, chips.NewFromInt(10), amt)

	// All In
	s.Players[s.TurnPos].Paid = chips.New(90)
	ak, amt = DiscreteAction(DAllIn).GetAction(p, s)
	require.Equal(t, AllIn, ak)
	// Because stack is 100 total, paid=90 => all-in is 10 left
	require.Equal(t, chips.NewFromInt(10), amt)

	// Raise fraction of pot
	s.CallAmount = chips.NewFromInt(0)
	s.Players[0].Paid = chips.New(0)
	pt := s.Players.PaidSum()
	require.True(t, pt.Equal(chips.Zero))
	ak, amt = DiscreteAction(0.5).GetAction(p, s)
	require.Equal(t, Bet, ak)
	require.True(t, amt.Equal(chips.Zero)) // 0.5 * 0 pot => 0
}

func TestDiscreteLegalActionsCreation(t *testing.T) {
	p := GameParams{
		NumPlayers:         2,
		InitialStacks:      chips.NewList(100, 100),
		BetSizes:           BetSizesDeep,
		TerminalStreet:     River,
		MaxActionsPerRound: 4,
		SbAmount:           chips.NewFromInt(1),
	}
	s := &State{
		Players: []Player{
			{Paid: chips.New(1), Status: StatusActive}, // SB
			{Paid: chips.New(2), Status: StatusActive}, // BB
		},
		TurnPos:    0,
		Street:     Preflop,
		CallAmount: chips.New(1), // SB needs 1 chip to call
	}
	da := NewDiscreteLegalActions(p, s)
	require.NotEmpty(t, da)
}

func TestDifferentStackSizes(t *testing.T) {
	// Create game params with different stack sizes
	p := GameParams{
		NumPlayers:         2,
		InitialStacks:      chips.NewList(180, 50),
		SbAmount:           chips.NewFromInt(1),
		BetSizes:           BetSizesDeep,
		MaxActionsPerRound: 4,
		TerminalStreet:     River,
	}

	// Create new game
	game, err := NewGame(p)
	require.NoError(t, err)

	// After blinds: p0: 179, p1: 49
	// p0 calls BB
	err = game.Action(ActionAmount{
		Action: Call,
		Amount: chips.NewFromInt(1),
	})
	require.NoError(t, err)

	// Check legal actions for p1
	legalActions := NewLegalActions(p, game.Latest)
	require.Equal(t, LegalActions{
		Bet:   chips.NewFromFloat(2),
		Check: chips.NewFromInt(0),
		AllIn: chips.NewFromInt(48),
	}, legalActions)

	// p1 checks
	err = game.Action(ActionAmount{
		Action: Check,
		Amount: chips.NewFromInt(0),
	})
	require.NoError(t, err)

	// p0 checks
	err = game.Action(ActionAmount{
		Action: Check,
		Amount: chips.NewFromInt(0),
	})
	require.NoError(t, err)

	legalActions = NewLegalActions(p, game.Latest)
	require.Equal(t, LegalActions{
		Bet:   chips.NewFromFloat(2),
		Check: chips.NewFromInt(0),
		AllIn: chips.NewFromInt(178),
	}, legalActions)

	// p0 bets 48
	err = game.Action(ActionAmount{
		Action: Bet,
		Amount: chips.NewFromInt(48),
	})
	require.NoError(t, err)

	legalActions = NewLegalActions(p, game.Latest)
	require.Equal(t, LegalActions{
		Fold: chips.NewFromInt(0),
		Call: chips.NewFromInt(48),
	}, legalActions)
}

func TestDifferentStackSizes2(t *testing.T) {
	p := GameParams{
		NumPlayers:         2,
		InitialStacks:      chips.NewList(258, 75),
		SbAmount:           chips.NewFromInt(1),
		BetSizes:           BetSizesDeep,
		MaxActionsPerRound: 4,
		TerminalStreet:     River,
	}

	game, err := NewGame(p)
	require.NoError(t, err)

	// Call BB
	err = game.Action(ActionAmount{
		Action: Call,
		Amount: chips.NewFromInt(1),
	})
	require.NoError(t, err)

	// Bet 14
	err = game.Action(ActionAmount{
		Action: Bet,
		Amount: chips.NewFromInt(14),
	})
	require.NoError(t, err)

	// Raise to 151
	err = game.Action(ActionAmount{
		Action: Raise,
		Amount: chips.NewFromInt(151),
	})
	require.NoError(t, err)

	legalActions := NewLegalActions(p, game.Latest)
	require.Equal(t, LegalActions{
		Fold: chips.NewFromInt(0),
		Call: chips.NewFromInt(59),
	}, legalActions)

	// Call 59
	err = game.Action(ActionAmount{
		Action: Call,
		Amount: chips.NewFromInt(59),
	})
	require.NoError(t, err)

	require.True(t, game.Latest.Finished())
}

func TestBettingPatternWithDeepStacks(t *testing.T) {
	p := GameParams{
		NumPlayers:         2,
		InitialStacks:      chips.NewList(87, 279),
		SbAmount:           chips.NewFromInt(1),
		BetSizes:           BetSizesDeep,
		MaxActionsPerRound: 4,
		TerminalStreet:     River,
	}

	game, err := NewGame(p)
	require.NoError(t, err)

	actions := []struct {
		action ActionKind
		amount chips.Chips
	}{
		{Raise, chips.NewFromInt(51)}, // P1 raises to 51
		{Call, chips.NewFromInt(50)},  // P0 calls

		// Flop
		{Check, chips.Zero}, // P0 checks
		{Check, chips.Zero}, // P1 checks

		// Turn
		{Check, chips.Zero}, // P0 checks
		{Check, chips.Zero}, // P1 checks

		// River
		{Check, chips.Zero},           // P0 checks
		{AllIn, chips.NewFromInt(35)}, // P1 bets 52
	}

	// Execute actions
	for i, a := range actions {
		err = game.Action(ActionAmount{
			Action: a.action,
			Amount: a.amount,
		})
		require.NoError(t, err, "Action %d failed", i)
	}

	require.Equal(t, River, game.Latest.Street)

	legal := NewDiscreteLegalActions(p, game.Latest)
	require.Equal(t, DiscreteLegalActions{
		DFold: chips.NewFromInt(0),
		DCall: chips.NewFromInt(35),
	}, legal)
}

// func TestDiscreteLegalActions(t *testing.T) {
// 	p := GameParams{
// 		NumPlayers:         2,
// 		InitialStacks:      chips.NewList(100, 100),
// 		SbAmount:           chips.NewFromInt(1),
// 		BetSizes:           BetSizesDeep,
// 		MaxActionsPerRound: 4,
// 		TerminalStreet:     River,
// 		Limp:               true,
// 	}

// 	game, err := NewGame(p)
// 	require.NoError(t, err)

// 	legal := NewLegalActions(p, game.Latest)
// 	require.Equal(t, LegalActions{
// 		Fold:  chips.NewFromInt(0),
// 		Call:  chips.NewFromInt(1),
// 		Raise: chips.NewFromInt(3),
// 		AllIn: chips.NewFromInt(99),
// 	}, legal)

// 	discrete := NewDiscreteLegalActions(p, game.Latest)

// 	require.Equal(t, DiscreteLegalActions{
// 		DAllIn: chips.NewFromInt(99),
// 		DFold:  chips.NewFromInt(0),
// 		DCall:  chips.NewFromInt(1),
// 		1:      chips.NewFromInt(3),
// 		1.5:    chips.NewFromFloat(4.5),
// 		2:      chips.NewFromInt(6),
// 		3:      chips.NewFromInt(9),
// 		9:      chips.NewFromInt(27),
// 		15:     chips.NewFromInt(45),
// 		25:     chips.NewFromInt(75),
// 	}, discrete)
// }

func TestStreetActionTracking(t *testing.T) {
	p := GameParams{
		NumPlayers:         2,
		InitialStacks:      chips.NewList(100, 100),
		SbAmount:           chips.NewFromInt(1),
		BetSizes:           BetSizesDeep,
		MaxActionsPerRound: 4,
		TerminalStreet:     River,
	}

	// Create initial state
	s := NewState(p)
	require.Equal(t, uint8(0), s.StreetAction)

	// Make initial bets (SB + BB)
	s, err := MakeInitialBets(p, s)
	require.NoError(t, err)
	// After SB and BB, street action should be 2
	require.Equal(t, uint8(2), s.StreetAction)

	// SB calls BB
	s, err = MakeAction(p, s, DCall)
	require.NoError(t, err)
	require.Equal(t, uint8(3), s.StreetAction)

	s, err = Move(p, s)
	require.NoError(t, err)

	// BB checks
	s, err = MakeAction(p, s, DCheck)
	require.NoError(t, err)
	require.Equal(t, uint8(4), s.StreetAction)

	// Move to flop - should reset street action
	s, err = Move(p, s)
	require.NoError(t, err)
	require.Equal(t, Flop, s.Street)
	require.Equal(t, uint8(0), s.StreetAction)

	// First player checks
	s, err = MakeAction(p, s, ActionAmount{
		Action: Check,
		Amount: chips.NewFromInt(0),
	})
	require.NoError(t, err)
	require.Equal(t, uint8(1), s.StreetAction)

	s, err = Move(p, s)
	require.NoError(t, err)

	// Second player checks
	s, err = MakeAction(p, s, ActionAmount{
		Action: Check,
		Amount: chips.NewFromInt(0),
	})
	require.NoError(t, err)
	require.Equal(t, uint8(2), s.StreetAction)

	// Move to turn - should reset street action
	s, err = Move(p, s)
	require.NoError(t, err)
	require.Equal(t, Turn, s.Street)
	require.Equal(t, uint8(0), s.StreetAction)

	// Verify PSAC (Per Street Action Counter) is reset
	for _, count := range s.PSAC {
		require.Equal(t, uint8(0), count)
	}

	// Verify PSLA (Per Street Last Action) is reset
	for _, action := range s.PSLA {
		require.Equal(t, NoAction, action)
	}

	// Verify PSC (Per Street Commitment) is reset
	for _, commitment := range s.PSC {
		require.True(t, commitment.Equal(chips.Zero))
	}

	// Verify BSC (Biggest Street Commitment) is reset
	require.True(t, s.BSC.Amount.Equal(chips.Zero))
	require.True(t, s.BSC.Addition.Equal(chips.Zero))
	require.Equal(t, NoAction, s.BSC.Action)
}

func TestAllInSituation(t *testing.T) {
	p := GameParams{
		NumPlayers:         2,
		InitialStacks:      chips.NewList(100, 100),
		SbAmount:           chips.NewFromInt(1),
		BetSizes:           BetSizesDeep,
		MaxActionsPerRound: 8,
		TerminalStreet:     River,
		BtnPos:             0,
	}

	// Create new game
	game, err := NewGame(p)
	require.NoError(t, err)

	// Sequence of actions to recreate the situation
	actions := []struct {
		action ActionKind
		amount chips.Chips
	}{
		{Call, chips.NewFromInt(1)},   // P0 calls BB
		{Bet, chips.NewFromInt(2)},    // P1 checks
		{AllIn, chips.NewFromInt(98)}, // P0 goes all-in with remaining stack
		{Call, chips.NewFromInt(96)},  // P1 calls 96
	}

	// Execute each action
	for _, a := range actions {
		state, err := MakeAction(p, game.Latest, ActionAmount{
			Action: a.action,
			Amount: a.amount,
		})
		require.NoError(t, err)
		state, err = Move(p, state)
		require.NoError(t, err)
		game.Latest = state
	}
}

func TestMinBettingRules(t *testing.T) {
	p := GameParams{
		NumPlayers:         2,
		InitialStacks:      chips.NewList(100, 100),
		SbAmount:           chips.NewFromInt(1),
		BetSizes:           BetSizesDeep,
		MaxActionsPerRound: 4,
		TerminalStreet:     River,
	}

	t.Run("preflop min raise rules", func(t *testing.T) {
		game, err := NewGame(p)
		require.NoError(t, err)

		// After SB(1) and BB(2), SB to act
		legal := NewLegalActions(p, game.Latest)
		// 1 to call, 1 was previous bet so 1 + 2 = 3
		require.Equal(t, chips.NewFromInt(3), legal[Raise], "min raise should be 2BB (3)")

		// SB calls
		err = game.Action(ActionAmount{Action: Call, Amount: chips.NewFromInt(1)})
		require.NoError(t, err)

		// BB can raise
		legal = NewLegalActions(p, game.Latest)
		require.Equal(t, chips.NewFromInt(2), legal[Bet], "after call, min raise should be 4")

		// BB raises to 8
		err = game.Action(ActionAmount{Action: Bet, Amount: chips.NewFromInt(8)})
		require.NoError(t, err)

		// SB's min raise should be previous bet(8) + size of last raise(6) = 14
		legal = NewLegalActions(p, game.Latest)
		require.Equal(t, chips.NewFromInt(16), legal[Raise], "min raise after 8 should be 14")
	})

	t.Run("flop min bet and raise rules", func(t *testing.T) {
		game, err := NewGame(p)
		require.NoError(t, err)

		// Get to flop with calls
		err = game.Action(ActionAmount{Action: Call, Amount: chips.NewFromInt(1)}) // SB calls
		require.NoError(t, err)

		err = game.Action(ActionAmount{Action: Check}) // BB checks
		require.NoError(t, err)

		require.Equal(t, Flop, game.Latest.Street)

		// First to act - min bet should be BB size
		legal := NewLegalActions(p, game.Latest)
		require.Equal(t, chips.NewFromInt(2), legal[Bet], "min bet should be BB size")

		// Bet 10
		err = game.Action(ActionAmount{Action: Bet, Amount: chips.NewFromInt(10)})
		require.NoError(t, err)

		// Next player - min raise should be 20 (previous bet 10 + size of bet 10)
		legal = NewLegalActions(p, game.Latest)
		require.Equal(t, chips.NewFromInt(20), legal[Raise], "min raise should be 20")

		// Raise to 25
		err = game.Action(ActionAmount{Action: Raise, Amount: chips.NewFromInt(25)})
		require.NoError(t, err)

		// First player - min raise should be 40 (previous bet 25 + size of raise 15)
		legal = NewLegalActions(p, game.Latest)
		require.Equal(t, chips.NewFromInt(40), legal[Raise], "min raise should be 40")
	})

	t.Run("all-in less than min raise", func(t *testing.T) {
		p := GameParams{
			NumPlayers:         2,
			InitialStacks:      chips.NewList(100, 30), // P1 has only 30
			SbAmount:           chips.NewFromInt(1),
			BetSizes:           BetSizesDeep,
			MaxActionsPerRound: 8,
			TerminalStreet:     River,
		}

		game, err := NewGame(p)
		require.NoError(t, err)

		// SB calls
		err = game.Action(ActionAmount{Action: Call, Amount: chips.NewFromInt(1)})
		require.NoError(t, err)

		// BB checks
		err = game.Action(ActionAmount{Action: Check})
		require.NoError(t, err)

		// BB checks
		err = game.Action(ActionAmount{Action: Check})
		require.NoError(t, err)

		// SB bets
		err = game.Action(ActionAmount{Action: Bet, Amount: chips.NewFromInt(20)})
		require.NoError(t, err)

		// BB goes all in
		err = game.Action(ActionAmount{Action: AllIn, Amount: chips.NewFromInt(28)})
		require.NoError(t, err)

		// P0's min raise should still be 40 (based on full raise amount)
		legal := NewLegalActions(p, game.Latest)
		require.Equal(t, chips.NewFromInt(8), legal[Call], "call should be 8")
	})

	t.Run("continuous reraises heads-up", func(t *testing.T) {
		p := GameParams{
			NumPlayers:         2,
			InitialStacks:      chips.NewList(500, 500),
			SbAmount:           chips.NewFromInt(1),
			BetSizes:           BetSizesDeep,
			MaxActionsPerRound: 8,
			TerminalStreet:     River,
		}

		game, err := NewGame(p)
		require.NoError(t, err)

		// SB calls BB
		err = game.Action(ActionAmount{Action: Call, Amount: chips.NewFromInt(1)})
		require.NoError(t, err)

		// BB bets 8
		err = game.Action(ActionAmount{Action: Bet, Amount: chips.NewFromInt(8)})
		require.NoError(t, err)

		// SB raises to 24 (+16)
		err = game.Action(ActionAmount{Action: Raise, Amount: chips.NewFromInt(24)})
		require.NoError(t, err)
		legal := NewLegalActions(p, game.Latest)
		require.Equal(t, chips.NewFromInt(40), legal[Raise],
			"min raise should be 40 (current bet 24 + previous raise 16)")

		// BB raises to 56 (+32)
		err = game.Action(ActionAmount{Action: Raise, Amount: chips.NewFromInt(56)})
		require.NoError(t, err)
		legal = NewLegalActions(p, game.Latest)
		require.Equal(t, chips.NewFromInt(88), legal[Raise],
			"min raise should be 88 (current bet 56 + previous raise 32)")

		// SB raises to 120 (+64)
		err = game.Action(ActionAmount{Action: Raise, Amount: chips.NewFromInt(120)})
		require.NoError(t, err)
		legal = NewLegalActions(p, game.Latest)

		require.Equal(t, chips.NewFromInt(184), legal[Raise],
			"min raise should be 184 (current bet 120 + previous raise 64)")
	})

	t.Run("multiplayer min raise after fold", func(t *testing.T) {
		p := GameParams{
			NumPlayers:         3,
			InitialStacks:      chips.NewList(100, 100, 100),
			SbAmount:           chips.NewFromInt(1),
			BetSizes:           BetSizesDeep,
			MaxActionsPerRound: 8,
			TerminalStreet:     River,
		}

		game, err := NewGame(p)
		require.NoError(t, err)

		// UTG (P2) calls
		err = game.Action(ActionAmount{Action: Call, Amount: chips.NewFromInt(2)})
		require.NoError(t, err)

		// SB (P0) calls
		err = game.Action(ActionAmount{Action: Call, Amount: chips.NewFromInt(1)})
		require.NoError(t, err)

		// BB (P1) bets 8
		err = game.Action(ActionAmount{Action: Bet, Amount: chips.NewFromInt(8)})
		require.NoError(t, err)

		// UTG (P2) raises to 24
		err = game.Action(ActionAmount{Action: Raise, Amount: chips.NewFromInt(24)})
		require.NoError(t, err)

		// SB (P0) folds
		err = game.Action(ActionAmount{Action: Fold})
		require.NoError(t, err)

		// BB's (P1) min raise should be to 40 (current bet 24 + previous raise 16)
		legal := NewLegalActions(p, game.Latest)
		require.Equal(t, chips.NewFromInt(40), legal[Raise],
			"min raise should be 40 after previous player folded")
	})
}

func TestBestSizes(t *testing.T) {
	p := GameParams{
		NumPlayers:         2,
		InitialStacks:      chips.NewList(400, 400),
		SbAmount:           chips.NewFromInt(1),
		BetSizes:           BetSizesDeep,
		MaxActionsPerRound: 8,
		TerminalStreet:     River,
	}

	game, err := NewGame(p)
	require.NoError(t, err)

	err = game.Action(DCall)
	require.NoError(t, err)

	// err = game.Action(ActionAmount{Action: Bet, Amount: chips.NewFromInt(10)})
	// require.NoError(t, err)

	// legal := NewDiscreteLegalActions(p, game.Latest)
	// pretty.Println(legal.String())

	err = game.Action(DiscreteAction(0.5))
	require.NoError(t, err)

	err = game.Action(DCall)
	require.NoError(t, err)

	legal := NewDiscreteLegalActions(p, game.Latest)
	pretty.Println(legal.String())
}

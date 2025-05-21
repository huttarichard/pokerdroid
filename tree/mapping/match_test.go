package mapping

import (
	"testing"

	"github.com/pokerdroid/poker/chips"
	"github.com/pokerdroid/poker/table"
	"github.com/stretchr/testify/require"
)

func TestMatchAction(t *testing.T) {
	testCases := []struct {
		desc string
		act  table.ActionAmount
		dis  []table.DiscreteAction
		pot  chips.Chips
		exp  int
	}{
		{
			desc: "exact match fold",
			act: table.ActionAmount{
				Action: table.Fold,
				Amount: chips.Zero,
			},
			dis: []table.DiscreteAction{
				table.DCall,
				table.DFold,
				table.DiscreteAction(1.0),
			},
			pot: chips.NewFromInt(10),
			exp: 1,
		},
		{
			desc: "exact match check",
			act: table.ActionAmount{
				Action: table.Check,
				Amount: chips.Zero,
			},
			dis: []table.DiscreteAction{
				table.DCall,
				table.DCheck,
				table.DiscreteAction(1.0),
			},
			pot: chips.NewFromInt(10),
			exp: 1,
		},
		{
			desc: "exact match call",
			act: table.ActionAmount{
				Action: table.Call,
				Amount: chips.NewFromInt(10),
			},
			dis: []table.DiscreteAction{
				table.DCall,
				table.DCheck,
				table.DiscreteAction(1.0),
			},
			pot: chips.NewFromInt(10),
			exp: 0,
		},
		{
			desc: "exact match all-in",
			act: table.ActionAmount{
				Action: table.AllIn,
				Amount: chips.NewFromInt(100),
			},
			dis: []table.DiscreteAction{
				table.DCall,
				table.DAllIn,
				table.DiscreteAction(1.0),
			},
			pot: chips.NewFromInt(10),
			exp: 1,
		},

		// POT:     $10
		// BET:     $15 - 1.5x
		// Actions: 0.5x, 1x, 2x
		// 		    $5    $10 $20
		//
		// (10+20+2x10x20)/(10+20+2) = $13.43
		// $15 > $13.43 = bet 2x
		{
			desc: "bet matching closest pot multiplier",
			act: table.ActionAmount{
				Action: table.Bet,
				Amount: chips.NewFromInt(15), // 1.5x pot
			},
			dis: []table.DiscreteAction{
				table.DCall,
				table.DiscreteAction(1.0), // pot-sized bet
				table.DiscreteAction(2.0), // 2x pot
				table.DiscreteAction(0.5), // half-pot
			},
			pot: chips.NewFromInt(10),
			exp: 2, // should match 1.0 as closest to 1.5
		},
		// POT:     $10
		// BET:     $25 - 2.5x
		// Actions: 2x, 3x
		// 		    $20 $30
		//
		// (20+30+2x20x30)/(20+30+2) = $24.03
		// $25 > $24.03 = bet 3x
		{
			desc: "raise matching closest pot multiplier",
			act: table.ActionAmount{
				Action: table.Raise,
				Amount: chips.NewFromInt(25), // 2.5x pot
			},
			dis: []table.DiscreteAction{
				table.DCall,
				table.DiscreteAction(2.0), // 2x pot
				table.DiscreteAction(3.0), // 3x pot
			},
			pot: chips.NewFromInt(10),
			exp: 2, // should match 2.0 as closest to 2.5
		},
		// POT:     $10
		// BET:     $8 - 0.8x
		// Actions: 0.5x, 4x
		// 		    $5 $40
		//
		// (5+40+2x5x40)/(5+40+2) = $9.46
		// $8 > $9.46 = bet 4x
		{
			desc: "raise matching closest pot multiplier",
			act: table.ActionAmount{
				Action: table.Raise,
				Amount: chips.NewFromInt(8), // 0.8x pot
			},
			dis: []table.DiscreteAction{
				table.DCall,
				table.DiscreteAction(0.5), // 0.5x pot
				table.DiscreteAction(4.0), // 4x pot
			},
			pot: chips.NewFromInt(10),
			exp: 1, // should match 2.0 as closest to 2.5
		},
		{
			desc: "no matching action found",
			act: table.ActionAmount{
				Action: table.Bet,
				Amount: chips.NewFromInt(10),
			},
			dis: []table.DiscreteAction{
				table.DCall,
				table.DCheck,
				table.DFold,
			},
			pot: chips.NewFromInt(10),
			exp: -1,
		},
		{
			desc: "empty discrete actions",
			act: table.ActionAmount{
				Action: table.Bet,
				Amount: chips.NewFromInt(10),
			},
			dis: []table.DiscreteAction{},
			pot: chips.NewFromInt(10),
			exp: -1,
		},
		{
			desc: "no pot multipliers available for bet",
			act: table.ActionAmount{
				Action: table.Bet,
				Amount: chips.NewFromInt(10),
			},
			dis: []table.DiscreteAction{
				table.DCall,
				table.DCheck,
				table.DFold,
				table.DAllIn,
			},
			pot: chips.NewFromInt(10),
			exp: -1,
		},
		{
			desc: "very small pot size",
			act: table.ActionAmount{
				Action: table.Bet,
				Amount: chips.NewFromInt(1),
			},
			dis: []table.DiscreteAction{
				table.DiscreteAction(0.5),
				table.DiscreteAction(1.0),
				table.DiscreteAction(2.0),
			},
			pot: chips.NewFromFloat32(0.1),
			exp: 2, // 1/0.1 = 10x pot, closest to 2.0
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			result := MatchAction(tC.act, tC.dis, tC.pot)
			require.Equal(t, tC.exp, result, "Expected index %d but got %d for %s", tC.exp, result, tC.desc)
		})
	}
}

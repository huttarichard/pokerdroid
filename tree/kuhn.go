package tree

import (
	"github.com/pokerdroid/poker/chips"
	"github.com/pokerdroid/poker/table"
)

func init() {
	KuhnChance.Parent = KuhnRoot
	KuhnP1.Parent = KuhnChance

	KuhnP1K.Parent = KuhnP1
	KuhnP1KK.Parent = KuhnP1K
	KuhnP1KB.Parent = KuhnP1K
	KuhnP1KBC.Parent = KuhnP1KB
	KuhnP1KBF.Parent = KuhnP1KB

	KuhnP1B.Parent = KuhnP1
	KuhnP1BC.Parent = KuhnP1B
	KuhnP1BF.Parent = KuhnP1B
}

func NewKuhn() (t *Root) {
	return KuhnRoot
}

var KuhnRoot = &Root{
	States: 11,
	Params: table.GameParams{
		NumPlayers: 2,
		BetSizes:   [][]float32{{1}},
	},
	State: &table.State{
		TurnPos: 0,
		PSC:     chips.List{1, 1},
		Players: table.Players{
			{Paid: chips.NewFromInt(1)},
			{Paid: chips.NewFromInt(1)},
		},
	},
	Next: KuhnChance,
	Full: true,
}

var KuhnChance = &Chance{
	Next: KuhnP1,
	State: &table.State{
		Street: table.Preflop,
	},
}

var KuhnP1Policies = NewStoreBacking()

var KuhnP1 = &Player{
	TurnPos: 0,
	Actions: &PlayerActions{
		Actions: []table.DiscreteAction{
			table.DCheck,
			table.DiscreteAction(1),
		},
		Nodes: []Node{
			KuhnP1K,
			KuhnP1B,
		},
		Policies: NewStoreBacking(),
	},
	State: &table.State{
		TurnPos: 0,
		PSC:     chips.List{0, 0},
		Players: table.Players{
			{Paid: chips.NewFromInt(1)},
			{Paid: chips.NewFromInt(1)},
		},
	},
}

// ========================================
// k

var KuhnP1K = &Player{
	TurnPos: 1,
	Actions: &PlayerActions{
		Actions: []table.DiscreteAction{
			table.DCheck,
			table.DiscreteAction(1),
		},
		Nodes: []Node{
			// k:k:[T]
			KuhnP1KK,
			// k:b2
			KuhnP1KB,
		},
		Policies: NewStoreBacking(),
	},
	State: &table.State{
		TurnPos: 1,
		PSC:     chips.List{0, 0},
		Players: table.Players{
			{Paid: chips.NewFromInt(1)},
			{Paid: chips.NewFromInt(1)},
		},
	},
}

var KuhnP1KK = &Terminal{
	Players: table.Players{
		{Paid: chips.NewFromInt(1), Status: table.StatusActive},
		{Paid: chips.NewFromInt(1), Status: table.StatusActive},
	},
	Pots: table.Pots{{
		Amount:  chips.NewFromInt(2),
		Players: []uint8{0, 1},
	}},
}

var KuhnP1KB = &Player{
	TurnPos: 0,
	Actions: &PlayerActions{
		Actions: []table.DiscreteAction{
			table.DCall,
			table.DFold,
		},
		Nodes: []Node{
			// k:b2:c:[T]
			KuhnP1KBC,
			// k:b2:f:[T]
			KuhnP1KBF,
		},
		Policies: NewStoreBacking(),
	},
	State: &table.State{
		TurnPos: 0,
		PSC:     chips.List{0, 1},
		Players: table.Players{
			{Paid: chips.NewFromInt(1)},
			{Paid: chips.NewFromInt(2)},
		},
	},
}

var KuhnP1KBC = &Terminal{
	Players: table.Players{
		{Paid: chips.NewFromInt(2), Status: table.StatusAllIn},
		{Paid: chips.NewFromInt(2), Status: table.StatusAllIn},
	},
	Pots: table.Pots{{
		Amount:  chips.NewFromInt(4),
		Players: []uint8{0, 1},
	}},
}

var KuhnP1KBF = &Terminal{
	Players: table.Players{
		{Paid: chips.NewFromInt(1), Status: table.StatusFolded},
		{Paid: chips.NewFromInt(2), Status: table.StatusAllIn},
	},
	Pots: table.Pots{{
		Amount:  chips.NewFromInt(3),
		Players: []uint8{1},
	}},
}

// ========================================
// b

var KuhnP1B = &Player{
	TurnPos: 1,
	Actions: &PlayerActions{
		Actions: []table.DiscreteAction{
			table.DCall,
			table.DFold,
		},
		Nodes: []Node{
			// b2:c:[T]
			KuhnP1BC,
			// k:b2:f:[T]
			KuhnP1BF,
		},
		Policies: NewStoreBacking(),
	},
	State: &table.State{
		TurnPos: 1,
		PSC:     chips.List{1, 0},
		Players: table.Players{
			{Paid: chips.NewFromInt(2)},
			{Paid: chips.NewFromInt(1)},
		},
	},
}

var KuhnP1BC = &Terminal{
	Players: table.Players{
		{Paid: chips.NewFromInt(2), Status: table.StatusAllIn},
		{Paid: chips.NewFromInt(2), Status: table.StatusAllIn},
	},
	Pots: table.Pots{{
		Amount:  chips.NewFromInt(4),
		Players: []uint8{0, 1},
	}},
}

var KuhnP1BF = &Terminal{
	Players: table.Players{
		{Paid: chips.NewFromInt(2), Status: table.StatusAllIn},
		{Paid: chips.NewFromInt(1), Status: table.StatusFolded},
	},
	Pots: table.Pots{{
		Amount:  chips.NewFromInt(3),
		Players: []uint8{0},
	}},
}

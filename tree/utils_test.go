package tree

import (
	"testing"

	"github.com/pokerdroid/poker/table"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLeafRunes(t *testing.T) {
	tx := &Terminal{
		Pots: table.Pots{
			0: table.Pot{
				Amount: 100,
			},
		},
		Players: table.Players{
			0: table.Player{
				Status: table.StatusActive,
			},
		},
	}

	px := &Player{
		State: &table.State{},
	}

	cx := &Chance{
		State: &table.State{},
	}

	acts := PlayerActions{
		Actions: []table.DiscreteAction{table.DiscreteAction(1), table.DCall, table.DFold},
		Nodes:   []Node{tx, px, cx},
	}

	p1 := &Player{
		TurnPos: 1,
		State:   &table.State{},
		Actions: &acts,
	}

	ch := &Chance{
		Next:  p1,
		State: &table.State{},
	}

	r := &Root{
		Next:   ch,
		Params: table.NewGameParams(2, 100),
		State:  &table.State{},
	}

	tx.Parent = p1
	px.Parent = p1
	cx.Parent = p1
	acts.Parent = p1
	p1.Parent = ch
	ch.Parent = r

	nn := FindLeafRunes(r)
	require.Len(t, nn, 3)

	require.Equal(t, nn[0].String(), "r:n:b0.00:t")
	require.Equal(t, nn[1].String(), "r:n:c:p")
	require.Equal(t, nn[2].String(), "r:n:f:n")
}

func TestExtractActions(t *testing.T) {
	t.Run("simple action sequence", func(t *testing.T) {
		// Build a simple tree with known actions
		terminal := &Terminal{}

		player2 := &Player{
			TurnPos: 1,
			State:   &table.State{},
			Actions: &PlayerActions{
				Actions: []table.DiscreteAction{table.DFold},
				Nodes:   []Node{terminal},
			},
		}
		terminal.Parent = player2

		player1 := &Player{
			TurnPos: 0,
			State:   &table.State{},
			Actions: &PlayerActions{
				Actions: []table.DiscreteAction{table.DCall, table.DFold},
				Nodes:   []Node{player2, terminal},
			},
		}
		player2.Parent = player1

		chance := &Chance{
			Next: player1,
		}
		player1.Parent = chance

		root := &Root{
			Next:   chance,
			Params: table.NewGameParams(2, 100),
		}
		chance.Parent = root

		// Test extraction from different nodes
		tests := []struct {
			name     string
			node     Node
			expected []Action
		}{
			{
				name: "from terminal",
				node: terminal,
				expected: []Action{
					{Action: table.DCall}, // player1's action
					{Action: table.DFold}, // player2's action
				},
			},
			{
				name: "from player2",
				node: player2,
				expected: []Action{
					{Action: table.DCall}, // player1's action
				},
			},
			{
				name:     "from player1",
				node:     player1,
				expected: []Action{},
			},
			{
				name:     "from chance",
				node:     chance,
				expected: []Action{},
			},
			{
				name:     "from root",
				node:     root,
				expected: []Action{},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				actions := ExtractActions(tt.node)
				require.Equal(t, len(tt.expected), len(actions), "action count mismatch")

				for i, expected := range tt.expected {
					require.Equal(t, expected.Action, actions[i].Action,
						"action type mismatch at index %d", i)
				}
			})
		}
	})

	t.Run("edge cases", func(t *testing.T) {
		// Test nil node
		actions := ExtractActions(nil)
		assert.Empty(t, actions)

		// Test node without parent
		player := &Player{
			TurnPos: 0,
		}
		actions = ExtractActions(player)
		assert.Empty(t, actions)

		// Test node with invalid LastAction index
		player = &Player{
			TurnPos: 0,
		}
		actions = ExtractActions(player)
		assert.Empty(t, actions)
	})
}

func TestExtractActions2(t *testing.T) {
	// Build a simple tree with known actions
	terminal := &Terminal{}

	player4 := &Player{
		TurnPos: 3,
		State:   &table.State{},
		Actions: &PlayerActions{
			Actions: []table.DiscreteAction{table.DFold},
			Nodes:   []Node{terminal},
		},
	}
	terminal.Parent = player4

	player3 := &Player{
		TurnPos: 2,
		State:   &table.State{},
		Actions: &PlayerActions{
			Actions: []table.DiscreteAction{table.DCall, table.DFold},
			Nodes:   []Node{player4, terminal},
		},
	}
	player4.Parent = player3

	chanceMiddle := &Chance{
		Next: player3,
	}
	player3.Parent = chanceMiddle

	player2 := &Player{
		TurnPos: 1,
		State:   &table.State{},
		Actions: &PlayerActions{
			Actions: []table.DiscreteAction{table.DCall, table.DFold},
			Nodes:   []Node{chanceMiddle, terminal},
		},
	}
	chanceMiddle.Parent = player2

	player1 := &Player{
		TurnPos: 0,
		State:   &table.State{},
		Actions: &PlayerActions{
			Actions: []table.DiscreteAction{table.DCall, table.DFold},
			Nodes:   []Node{player2, terminal},
		},
	}
	player2.Parent = player1

	chance := &Chance{
		Next: player1,
	}
	player1.Parent = chance

	root := &Root{
		Next:   chance,
		Params: table.NewGameParams(4, 100),
	}
	chance.Parent = root

	actions1 := ExtractActions(player3)
	last := actions1[len(actions1)-1]

	require.Equal(t, last.Node, chanceMiddle)
	require.Equal(t, GetPath(last.Node).String(), GetPath(chanceMiddle).String())
}

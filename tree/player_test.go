package tree

import (
	"testing"

	"github.com/pokerdroid/poker/chips"
	"github.com/pokerdroid/poker/policy"
	"github.com/pokerdroid/poker/table"
	"github.com/stretchr/testify/require"
)

func TestPlayer_Basic(t *testing.T) {
	// Create test state
	prms := table.NewGameParams(2, 100)
	state := table.NewState(prms)

	// Create test nodes
	parent := &Root{State: state}
	actions := &PlayerActions{
		Actions:  []table.DiscreteAction{table.DCall, table.DFold},
		Nodes:    make([]Node, 2),
		Policies: NewPolicies(),
	}

	player := &Player{
		Parent:  parent,
		TurnPos: 1,
		State:   state,
		Actions: actions,
	}

	// Test basic properties
	require.Equal(t, NodeKindPlayer, player.Kind())
	require.Equal(t, parent, player.GetParent())
	require.Equal(t, uint8(1), player.TurnPos)
	require.NotNil(t, player.Actions)
}

func TestPlayer_GetActionIdx(t *testing.T) {
	node1 := &Terminal{}
	node2 := &Terminal{}
	player := &Player{
		Actions: &PlayerActions{
			Nodes: []Node{node1, node2},
		},
	}

	idx, ok := player.GetActionIdx(node1)
	require.True(t, ok)
	require.Equal(t, 0, idx)

	idx, ok = player.GetActionIdx(node2)
	require.True(t, ok)
	require.Equal(t, 1, idx)

	idx, ok = player.GetActionIdx(&Terminal{})
	require.False(t, ok)
	require.Equal(t, -1, idx)
}

func TestPlayer_MarshalBinary(t *testing.T) {
	tests := []struct {
		name    string
		player  *Player
		wantErr bool
	}{
		{
			name: "basic player",
			player: &Player{
				TurnPos: 1,
				State:   table.NewState(table.NewGameParams(2, 100)),
				Actions: &PlayerActions{
					Actions:  []table.DiscreteAction{table.DCall},
					Nodes:    []Node{&Terminal{}},
					Policies: NewPolicies(),
				},
			},
			wantErr: false,
		},
		{
			name: "nil state",
			player: &Player{
				TurnPos: 1,
				State:   nil,
				Actions: nil,
			},
			wantErr: false,
		},
		{
			name: "nil actions",
			player: &Player{
				TurnPos: 1,
				State:   table.NewState(table.NewGameParams(2, 100)),
				Actions: nil,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Marshal
			data, err := tt.player.MarshalBinary()
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, data)

			// Unmarshal into new player
			newPlayer := &Player{}
			err = newPlayer.UnmarshalBinary(data)
			require.NoError(t, err)

			// Verify fields
			require.Equal(t, tt.player.TurnPos, newPlayer.TurnPos)
			if tt.player.State == nil {
				require.Nil(t, newPlayer.State)
			} else {
				require.NotNil(t, newPlayer.State)
			}
			if tt.player.Actions == nil {
				require.Nil(t, newPlayer.Actions)
			} else {
				require.NotNil(t, newPlayer.Actions)
				require.Equal(t, tt.player.Actions.Actions, newPlayer.Actions.Actions)
			}
		})
	}
}

func TestPlayer_UnmarshalBinary_InvalidKind(t *testing.T) {
	player := &Player{}
	err := player.UnmarshalBinary([]byte{byte(NodeKindTerminal)}) // Wrong node kind
	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid node kind")
}

func TestPlayer_RoundTrip(t *testing.T) {
	// Create a complex player node
	prms := table.NewGameParams(2, 100)
	state := table.NewState(prms)

	original := &Player{
		TurnPos: 1,
		State:   state,
		Actions: &PlayerActions{
			Actions:  []table.DiscreteAction{table.DCall, table.DFold},
			Nodes:    make([]Node, 2),
			Policies: NewPolicies(),
		},
	}

	original.Actions.Nodes[0] = &Terminal{
		Players: table.Players{
			{Status: table.StatusActive},
			{Status: table.StatusFolded},
		},
	}
	original.Actions.Nodes[1] = &Chance{State: state}

	// Marshal
	data, err := original.MarshalBinary()
	require.NoError(t, err)

	// Unmarshal
	reconstructed := &Player{}
	err = reconstructed.UnmarshalBinary(data)
	require.NoError(t, err)

	// Verify key properties
	require.Equal(t, original.TurnPos, reconstructed.TurnPos)
	require.NotNil(t, reconstructed.State)
	require.NotNil(t, reconstructed.Actions)
	require.Equal(t, original.Actions.Actions, reconstructed.Actions.Actions)
	require.Equal(t, len(original.Actions.Nodes), len(reconstructed.Actions.Nodes))

	// Verify nodes
	for i := range original.Actions.Nodes {
		require.Equal(t, original.Actions.Nodes[i].Kind(), reconstructed.Actions.Nodes[i].Kind())
	}

}

func TestActions_Basic(t *testing.T) {
	// Create test actions
	actions := &PlayerActions{
		Actions:  []table.DiscreteAction{table.DCall, table.DFold},
		Nodes:    make([]Node, 2),
		Policies: NewPolicies(),
	}

	// Add test nodes
	actions.Nodes[0] = &Terminal{
		Pots: table.Pots{{Amount: chips.NewFromInt(100)}},
		Players: table.Players{
			{Paid: chips.NewFromInt(50), Status: table.StatusActive},
		},
	}
	actions.Nodes[1] = nil // Test nil node

	// Test basic properties
	require.Equal(t, 2, len(actions.Actions))
	require.Equal(t, 2, len(actions.Nodes))
	require.NotNil(t, actions.Policies)
}

func TestActions_MarshalBinary(t *testing.T) {
	tests := []struct {
		name    string
		actions *PlayerActions
		wantErr bool
	}{
		{
			name: "basic actions",
			actions: &PlayerActions{
				Actions:  []table.DiscreteAction{table.DCall, table.DFold},
				Nodes:    []Node{&Terminal{}, nil},
				Policies: NewPolicies(),
			},
			wantErr: false,
		},
		{
			name: "empty actions",
			actions: &PlayerActions{
				Actions:  []table.DiscreteAction{},
				Nodes:    []Node{},
				Policies: NewPolicies(),
			},
			wantErr: false,
		},
		{
			name: "nil node in actions",
			actions: &PlayerActions{
				Actions:  []table.DiscreteAction{table.DCall},
				Nodes:    []Node{nil},
				Policies: NewPolicies(),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Marshal
			data, err := tt.actions.MarshalBinary()
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, data)

			// Unmarshal into new actions
			newActions := &PlayerActions{}
			err = newActions.UnmarshalBinary(data)
			require.NoError(t, err)

			// Compare fields
			require.Equal(t, tt.actions.Actions, newActions.Actions)
			require.Equal(t, len(tt.actions.Nodes), len(newActions.Nodes))
			require.True(t, tt.actions.Policies.Equal(newActions.Policies))

			// Compare nodes
			for i := range tt.actions.Nodes {
				if tt.actions.Nodes[i] == nil {
					require.Nil(t, newActions.Nodes[i])
					continue
				}
				require.Equal(t, tt.actions.Nodes[i].Kind(), newActions.Nodes[i].Kind())
			}
		})
	}
}

func TestActions_RoundTrip(t *testing.T) {
	// Create a complex actions node
	original := &PlayerActions{
		Actions:  []table.DiscreteAction{table.DCall, table.DFold, table.DCheck},
		Nodes:    make([]Node, 3),
		Policies: NewPolicies(),
	}

	original.Nodes[0] = &Terminal{
		Players: table.Players{
			{Status: table.StatusActive},
			{Status: table.StatusFolded},
		},
	}
	original.Nodes[1] = nil
	original.Nodes[2] = &Chance{State: table.NewState(table.NewGameParams(2, 100))}

	// Marshal
	data, err := original.MarshalBinary()
	require.NoError(t, err)

	// Unmarshal
	reconstructed := &PlayerActions{}
	err = reconstructed.UnmarshalBinary(data)
	require.NoError(t, err)

	// Verify key properties
	require.Equal(t, original.Actions, reconstructed.Actions)
	require.Equal(t, len(original.Nodes), len(reconstructed.Nodes))
	require.True(t, original.Policies.Equal(reconstructed.Policies))

	// Verify nodes
	for i := range original.Nodes {
		if original.Nodes[i] == nil {
			require.Nil(t, reconstructed.Nodes[i])
			continue
		}
		require.Equal(t, original.Nodes[i].Kind(), reconstructed.Nodes[i].Kind())
	}
}

func TestActions_Size(t *testing.T) {
	tests := []struct {
		name    string
		setup   func() *PlayerActions
		actions int
	}{
		{
			name: "empty actions",
			setup: func() *PlayerActions {
				return &PlayerActions{
					Actions:  []table.DiscreteAction{},
					Nodes:    []Node{},
					Policies: NewPolicies(),
				}
			},
			actions: 0,
		},
		{
			name: "single action no policy",
			setup: func() *PlayerActions {
				return &PlayerActions{
					Actions:  []table.DiscreteAction{table.DCheck},
					Nodes:    []Node{nil},
					Policies: NewPolicies(),
				}
			},
			actions: 1,
		},
		{
			name: "multiple actions with policies",
			setup: func() *PlayerActions {
				a := &PlayerActions{
					Actions: []table.DiscreteAction{
						table.DCheck,
						table.DCall,
						table.DFold,
					},
					Nodes:    make([]Node, 3),
					Policies: NewPolicies(),
				}
				// Add some policies
				a.Policies.Store(1, policy.New(3))
				a.Policies.Store(2, policy.New(3))
				return a
			},
			actions: 3,
		},
		{
			name: "actions with terminal nodes",
			setup: func() *PlayerActions {
				a := &PlayerActions{
					Actions: []table.DiscreteAction{
						table.DCheck,
						table.DCall,
					},
					Nodes: []Node{
						&Terminal{},
						&Terminal{},
					},
					Policies: NewPolicies(),
				}
				return a
			},
			actions: 2,
		},
		{
			name: "actions with terminal nodes 1",
			setup: func() *PlayerActions {
				a := &PlayerActions{
					Actions: []table.DiscreteAction{
						table.DCheck,
						table.DCall,
						table.DAllIn,
					},
					Nodes: []Node{
						&Terminal{},
						&Terminal{},
						&Terminal{},
					},
					Policies: NewPolicies(),
				}
				return a
			},
			actions: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actions := tt.setup()
			require.NoError(t, actions.Validate())

			reportedSize := actions.Size()
			data, err := actions.MarshalBinary()
			require.NoError(t, err)

			require.Equal(t, reportedSize, uint64(len(data)),
				"marshaled size (%d) does not match reported size (%d)",
				len(data), reportedSize)

			// Test round trip
			newActions := &PlayerActions{}
			err = newActions.UnmarshalBinary(data)
			require.NoError(t, err)

			// Verify basic properties
			require.Equal(t, len(actions.Actions), len(newActions.Actions))
			require.Equal(t, len(actions.Nodes), len(newActions.Nodes))

			// Compare actions
			for i, act := range actions.Actions {
				require.Equal(t, act, newActions.Actions[i])
			}

			// Compare policies
			require.Equal(t, actions.Policies.Len(), newActions.Policies.Len())

			for cl, pol := range actions.Policies.Map {
				newPol, ok := newActions.Policies.Get(cl)
				require.True(t, ok)
				require.Equal(t, pol.Strategy, newPol.Strategy)
			}
		})
	}
}

package tree

import (
	"bytes"
	"encoding/binary"
	"testing"

	"github.com/pokerdroid/poker/chips"
	"github.com/pokerdroid/poker/table"
	"github.com/stretchr/testify/require"
)

func TestChance_Basic(t *testing.T) {
	// Create test state
	prms := table.NewGameParams(2, 100)
	state := table.NewState(prms)

	// Create test nodes
	parent := &Root{State: state}
	next := &Terminal{Parent: parent}

	chance := &Chance{
		Next:   next,
		Parent: parent,
		State:  state,
	}

	// Test basic properties
	require.Equal(t, NodeKindChance, chance.Kind())
	require.Equal(t, parent, chance.GetParent())
	require.Equal(t, next, chance.Next)
}

func TestChance_MarshalBinary(t *testing.T) {
	tests := []struct {
		name    string
		chance  *Chance
		wantErr bool
	}{
		{
			name: "basic chance node",
			chance: &Chance{
				State: table.NewState(table.NewGameParams(2, 100)),
				Next:  &Terminal{},
			},
			wantErr: false,
		},
		{
			name: "nil state",
			chance: &Chance{
				State: nil,
				Next:  &Terminal{},
			},
			wantErr: false,
		},
		{
			name: "nil next",
			chance: &Chance{
				State: table.NewState(table.NewGameParams(2, 100)),
				Next:  nil,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Marshal
			data, err := tt.chance.MarshalBinary()
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, data)

			// Unmarshal into new chance node
			newChance := &Chance{}
			err = newChance.UnmarshalBinary(data)
			require.NoError(t, err)

			// Verify state
			if tt.chance.State == nil {
				require.Nil(t, newChance.State)
			} else {
				require.NotNil(t, newChance.State)
				// Add deeper state comparison if needed
			}

			// Verify next node
			if tt.chance.Next == nil {
				require.Nil(t, newChance.Next)
			} else {
				require.NotNil(t, newChance.Next)
				require.Equal(t, tt.chance.Next.Kind(), newChance.Next.Kind())
			}
		})
	}
}

func TestChance_UnmarshalBinary_InvalidKind(t *testing.T) {
	chance := &Chance{}
	err := chance.UnmarshalBinary([]byte{byte(NodeKindTerminal)}) // Wrong node kind
	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid node kind")
}

func TestChance_RoundTrip(t *testing.T) {
	// Create a complex chance node
	prms := table.NewGameParams(2, 100)
	state := table.NewState(prms)

	original := &Chance{
		State: state,
		Next: &Terminal{
			Players: table.Players{
				{Status: table.StatusActive},
				{Status: table.StatusFolded},
			},
		},
	}

	// Marshal
	data, err := original.MarshalBinary()
	require.NoError(t, err)

	// Unmarshal
	reconstructed := &Chance{}
	err = reconstructed.UnmarshalBinary(data)
	require.NoError(t, err)

	// Verify key properties
	require.Equal(t, original.Kind(), reconstructed.Kind())
	require.NotNil(t, reconstructed.State)
	require.NotNil(t, reconstructed.Next)
	require.Equal(t, original.Next.Kind(), reconstructed.Next.Kind())
}

func TestChance_Size(t *testing.T) {
	tests := []struct {
		name  string
		setup func() *Chance
	}{
		{
			name: "empty chance",
			setup: func() *Chance {
				return &Chance{
					Next:  nil,
					State: nil,
				}
			},
		},
		{
			name: "chance with state",
			setup: func() *Chance {
				return &Chance{
					Next: nil,
					State: &table.State{
						Players: table.Players{
							{Status: table.StatusActive},
							{Status: table.StatusFolded},
						},
						Street:       table.Preflop,
						TurnPos:      0,
						BtnPos:       1,
						StreetAction: 0,
					},
				}
			},
		},
		{
			name: "chance with terminal next",
			setup: func() *Chance {
				return &Chance{
					Next: &Terminal{
						Players: table.Players{
							{Status: table.StatusActive},
						},
						Pots: table.Pots{
							{Amount: chips.NewFromInt(100)},
						},
					},
					State: nil,
				}
			},
		},
		{
			name: "chance with state and next",
			setup: func() *Chance {
				return &Chance{
					Next: &Terminal{
						Players: table.Players{
							{Status: table.StatusActive},
						},
					},
					State: &table.State{
						Players: table.Players{
							{Status: table.StatusActive},
						},
						Street:  table.Flop,
						TurnPos: 0,
						BtnPos:  0,
					},
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			chance := tt.setup()

			// Get size and marshal
			reportedSize := chance.Size()
			data, err := chance.MarshalBinary()
			require.NoError(t, err)

			// Analyze binary structure
			buf := bytes.NewBuffer(data)

			// Read State length
			var stateLen uint16
			binary.Read(bytes.NewReader(buf.Next(2)), binary.LittleEndian, &stateLen)

			// Skip state data
			if stateLen > 0 {
				buf.Next(int(stateLen))
			}

			// Read Next node length
			var nextLen uint64
			binary.Read(bytes.NewReader(buf.Next(8)), binary.LittleEndian, &nextLen)

			// Skip next node data
			if nextLen > 0 {
				buf.Next(int(nextLen))
			}

			require.Equal(t, reportedSize, uint64(len(data)),
				"marshaled size (%d) does not match reported size (%d)",
				len(data), reportedSize)
		})
	}
}

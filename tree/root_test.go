package tree

import (
	"bytes"
	"encoding/binary"
	"testing"

	"github.com/pokerdroid/poker/chips"
	"github.com/pokerdroid/poker/table"
	"github.com/stretchr/testify/require"
)

func TestRoot_Basic(t *testing.T) {
	// Create test root
	prms := table.NewGameParams(2, 100)
	state := table.NewState(prms)
	root := &Root{
		States:    5,
		Next:      &Chance{},
		Params:    prms,
		Iteration: 42,
		State:     state,
	}

	// Test basic properties
	require.Equal(t, NodeKindRoot, root.Kind())
	require.Nil(t, root.GetParent())
	require.Equal(t, uint32(5), root.States)
	require.Equal(t, uint64(42), root.Iteration)
	require.NotNil(t, root.Next)
	require.NotNil(t, root.State)
}

func TestRoot_MarshalBinary(t *testing.T) {
	tests := []struct {
		name    string
		root    *Root
		wantErr bool
	}{
		{
			name: "basic root",
			root: &Root{
				States:    10,
				Next:      &Chance{},
				Params:    table.NewGameParams(2, 100),
				Iteration: 1,
				State:     table.NewState(table.NewGameParams(2, 100)),
			},
			wantErr: false,
		},
		{
			name: "nil next",
			root: &Root{
				States:    5,
				Next:      nil,
				Params:    table.NewGameParams(2, 100),
				Iteration: 0,
				State:     table.NewState(table.NewGameParams(2, 100)),
			},
			wantErr: false,
		},
		{
			name: "zero nodes",
			root: &Root{
				States:    0,
				Next:      &Chance{},
				Params:    table.NewGameParams(2, 100),
				Iteration: 100,
				State:     table.NewState(table.NewGameParams(2, 100)),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Marshal
			data, err := tt.root.MarshalBinary()
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, data)

			// Unmarshal into new root
			newRoot := &Root{}
			err = newRoot.UnmarshalBinary(data)
			require.NoError(t, err)

			// Verify fields
			require.Equal(t, tt.root.States, newRoot.States)
			require.Equal(t, tt.root.Iteration, newRoot.Iteration)
			if tt.root.Next == nil {
				require.Nil(t, newRoot.Next)
			} else {
				require.NotNil(t, newRoot.Next)
				require.Equal(t, tt.root.Next.Kind(), newRoot.Next.Kind())
			}
		})
	}
}

func TestRoot_UnmarshalBinary_InvalidKind(t *testing.T) {
	root := &Root{}
	err := root.UnmarshalBinary([]byte{byte(NodeKindPlayer)}) // Wrong node kind
	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid node kind")
}

func TestRoot_RoundTrip(t *testing.T) {
	// Create a complex root node
	prms := table.NewGameParams(2, 100)
	state := table.NewState(prms)

	original := &Root{
		States:    42,
		Next:      &Chance{State: state},
		Params:    prms,
		Iteration: 100,
		State:     state,
	}

	// Marshal
	data, err := original.MarshalBinary()
	require.NoError(t, err)

	// Unmarshal
	reconstructed := &Root{}
	err = reconstructed.UnmarshalBinary(data)
	require.NoError(t, err)

	// Verify key properties
	require.Equal(t, original.States, reconstructed.States)
	require.Equal(t, original.Iteration, reconstructed.Iteration)
	require.NotNil(t, reconstructed.Next)
	require.Equal(t, original.Next.Kind(), reconstructed.Next.Kind())
	require.NotNil(t, reconstructed.State)
}

// TestFindClosestRoot verifies that among a set of roots the one returned by
// FindClosestRoot (or our fixed version) is the one whose effective stack (for a given perspective)
// is closest to that of the query game params.
func TestFindClosestRoot(t *testing.T) {
	rootA := &Root{
		Params: table.GameParams{
			SbAmount: 1,
			InitialStacks: chips.List{
				chips.Chips(100),
				chips.Chips(100),
				chips.Chips(100),
				chips.Chips(100),
			},
			NumPlayers: 4,
		},
	}

	rootB := &Root{
		Params: table.GameParams{
			SbAmount: 1,
			InitialStacks: chips.List{
				chips.Chips(200),
				chips.Chips(200),
				chips.Chips(200),
				chips.Chips(200),
			},
			NumPlayers: 4,
		},
	}

	rootC := &Root{
		Params: table.GameParams{
			SbAmount: 1,
			InitialStacks: chips.List{
				chips.Chips(300),
				chips.Chips(300),
				chips.Chips(300),
				chips.Chips(300),
			},
			NumPlayers: 4,
		},
	}

	rootD := &Root{
		Params: table.GameParams{
			SbAmount: 1,
			InitialStacks: chips.List{
				chips.Chips(400),
				chips.Chips(400),
				chips.Chips(400),
				chips.Chips(400),
			},
			NumPlayers: 4,
		},
	}

	roots := []*Root{rootA, rootB, rootC, rootD}

	// Use query game parameters identical to rootA.
	queryParams := table.GameParams{
		SbAmount: 1,
		InitialStacks: chips.List{
			chips.Chips(100),
			chips.Chips(200),
			chips.Chips(300),
			chips.Chips(400),
		},
		NumPlayers: 4,
	}

	closest := FindClosestRoot(roots, queryParams, 0)
	require.Equal(t, rootA, closest)

	closest = FindClosestRoot(roots, queryParams, 1)
	require.Equal(t, rootB, closest)

	closest = FindClosestRoot(roots, queryParams, 2)
	require.Equal(t, rootC, closest)

	closest = FindClosestRoot(roots, queryParams, 3)
	require.Equal(t, rootC, closest)
}

func TestRoot_Size(t *testing.T) {
	tests := []struct {
		name  string
		setup func() *Root
	}{
		{
			name: "root with state",
			setup: func() *Root {
				prms := table.NewGameParams(2, chips.NewFromInt(100))
				s := table.NewState(prms)
				s, _ = table.MakeInitialBets(prms, s)
				return &Root{
					States:    1,
					Nodes:     1,
					Next:      nil,
					Params:    prms,
					Iteration: 1,
					State:     s,
				}
			},
		},
		{
			name: "root with terminal next",
			setup: func() *Root {
				prms := table.NewGameParams(2, chips.NewFromInt(100))
				return &Root{
					States: 1,
					Nodes:  2,
					Next: &Terminal{
						Players: table.Players{
							{Status: table.StatusActive, Paid: chips.NewFromInt(100)},
							{Status: table.StatusFolded, Paid: chips.NewFromInt(50)},
						},
						Pots: table.Pots{
							{Amount: chips.NewFromInt(150)},
						},
					},
					Params:    prms,
					Iteration: 1,
					State:     nil,
				}
			},
		},
		{
			name: "root with state and chance next",
			setup: func() *Root {
				prms := table.NewGameParams(2, chips.NewFromInt(100))
				s := table.NewState(prms)
				s, _ = table.MakeInitialBets(prms, s)
				return &Root{
					States: 2,
					Nodes:  3,
					Next: &Chance{
						State: s,
						Next: &Terminal{
							Players: table.Players{
								{Status: table.StatusActive, Paid: chips.NewFromInt(100)},
							},
							Pots: table.Pots{
								{Amount: chips.NewFromInt(100)},
							},
						},
					},
					Params:    prms,
					Iteration: 2,
					State:     s,
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := tt.setup()

			// Get size and marshal
			reportedSize := root.Size()
			data, err := root.MarshalBinary()
			require.NoError(t, err)

			// Analyze binary structure
			buf := bytes.NewBuffer(data)

			// Read fixed fields
			var states uint32
			var nodes uint32
			var iteration uint64
			binary.Read(bytes.NewReader(buf.Next(4)), binary.LittleEndian, &states)
			binary.Read(bytes.NewReader(buf.Next(4)), binary.LittleEndian, &nodes)
			binary.Read(bytes.NewReader(buf.Next(8)), binary.LittleEndian, &iteration)

			// Read Params length
			var paramsLen uint64
			binary.Read(bytes.NewReader(buf.Next(8)), binary.LittleEndian, &paramsLen)

			// Skip params data
			if paramsLen > 0 {
				buf.Next(int(paramsLen))
			}

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

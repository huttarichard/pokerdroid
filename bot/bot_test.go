package bot_test

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/pokerdroid/poker/bot"
	"github.com/pokerdroid/poker/card"
	"github.com/pokerdroid/poker/chips"
	"github.com/pokerdroid/poker/table"
	"github.com/stretchr/testify/require"
)

func TestBotStateMarshalUnmarshal(t *testing.T) {
	// Create dummy game parameters.
	params := table.GameParams{
		NumPlayers:         2,
		MaxActionsPerRound: 5,
		BtnPos:             0,
		SbAmount:           1,                    // Assuming chips.Chips is compatible with int
		BetSizes:           [][]float32{{1.0}},   // Dummy bet sizes
		InitialStacks:      chips.List{100, 100}, // Assuming chips.List is []int
		TerminalStreet:     table.Preflop,        // Use Preflop as a dummy value
		DisableV:           false,
	}

	// Create a minimal table state.
	// Only a few essential fields are set so that the underlying table.MarshalBinary works.
	ts := &table.State{
		Street:       table.Preflop,
		TurnPos:      0,
		BtnPos:       0,
		StreetAction: 0,
		CallAmount:   0,
		Previous:     nil,
	}

	// Create hole and community cards.
	hole := card.NewCardsFromString("as ks")
	com := card.NewCardsFromString("qs js ts 9s 8s")

	// Construct the original bot state.
	orig := &bot.State{
		Params:    params,
		State:     ts,
		Hole:      hole,
		Community: com,
	}

	// Marshal the original state to binary.
	data, err := orig.MarshalBinary()
	if err != nil {
		t.Fatalf("MarshalBinary failed: %v", err)
	}

	// Unmarshal the binary data into a new State.
	var got bot.State
	if err := got.UnmarshalBinary(data); err != nil {
		t.Fatalf("UnmarshalBinary failed: %v", err)
	}

	// Compare game parameters.
	if !reflect.DeepEqual(orig.Params, got.Params) {
		t.Errorf("Params mismatch:\ngot  %#v\nwant %#v", got.Params, orig.Params)
	}

	// Compare hole cards using their serialized bytes.
	if !bytes.Equal(orig.Hole.Bytes(), got.Hole.Bytes()) {
		t.Errorf("Hole mismatch:\ngot  %v\nwant %v", got.Hole, orig.Hole)
	}

	// Compare community cards.
	if !bytes.Equal(orig.Community.Bytes(), got.Community.Bytes()) {
		t.Errorf("Community mismatch:\ngot  %v\nwant %v", got.Community, orig.Community)
	}

	// Compare a key field from the table state (for example, Street).
	if orig.State.Street != got.State.Street {
		t.Errorf("Table state Street mismatch:\ngot  %v\nwant %v", got.State.Street, orig.State.Street)
	}
}

var exampleState = bot.State{
	Params: table.GameParams{
		NumPlayers:         0x2,
		MaxActionsPerRound: 0x6,
		BtnPos:             0x1,
		SbAmount:           0.019999999552965164,
		BetSizes:           [][]float32{{0.5, 1, 1.5}},
		InitialStacks:      chips.List{28.559999465942383, 13.289999961853027},
		TerminalStreet:     0x4,
		DisableV:           false,
	},
	State: &table.State{
		Players: table.Players{
			{Paid: 0.03999999910593033, Status: 0x1},
			{Paid: 0.11999999731779099, Status: 0x2},
		},
		Street:       table.Street(0x5),
		TurnPos:      0x0,
		BtnPos:       0x1,
		StreetAction: 0x3,
		CallAmount:   chips.Chips(0.07999999821186066),
		BSC: struct {
			Amount   chips.Chips      "json:\"amount\""
			Addition chips.Chips      "json:\"addition\""
			Action   table.ActionKind "json:\"action\""
		}{Amount: 0.07999999821186066, Addition: 0.07999999821186066, Action: table.ActionKind(0x1)},
		PSC:  chips.List{0, 0.07999999821186066},
		PSAC: []uint8{0x2, 0x1},
		PSLA: []table.ActionKind{0x4, 0x1},
		Previous: &table.State{
			Players: table.Players{
				{Paid: 0.03999999910593033, Status: 0x1},
				{Paid: 0.11999999731779099, Status: 0x2},
			},
			Street:       table.Street(0x2),
			TurnPos:      0x0,
			BtnPos:       0x1,
			StreetAction: 0x3,
			CallAmount:   chips.Chips(0.07999999821186066),
			BSC: struct {
				Amount   chips.Chips      "json:\"amount\""
				Addition chips.Chips      "json:\"addition\""
				Action   table.ActionKind "json:\"action\""
			}{Amount: 0.07999999821186066, Addition: 0.07999999821186066, Action: table.ActionKind(0x1)},
			PSC:  chips.List{0, 0.07999999821186066},
			PSAC: []uint8{0x2, 0x1},
			PSLA: []table.ActionKind{0x4, 0x1},
			Previous: &table.State{
				Players: table.Players{
					{Paid: 0.03999999910593033, Status: table.Status(0x2)},
					{Paid: 0.11999999731779099, Status: table.Status(0x2)},
				},
				Street:       table.Street(0x2),
				TurnPos:      0x0,
				BtnPos:       0x1,
				StreetAction: 0x2,
				CallAmount:   chips.Chips(0.07999999821186066),
				BSC: struct {
					Amount   chips.Chips      "json:\"amount\""
					Addition chips.Chips      "json:\"addition\""
					Action   table.ActionKind "json:\"action\""
				}{Amount: 0.07999999821186066, Addition: 0.07999999821186066, Action: table.ActionKind(0x1)},
				PSC:  chips.List{0, 0.07999999821186066},
				PSAC: []uint8{0x1, 0x1},
				PSLA: []table.ActionKind{0x5, 0x1},
				Previous: &table.State{
					Players: table.Players{
						{Paid: 0.03999999910593033, Status: table.Status(0x2)},
						{Paid: 0.11999999731779099, Status: table.Status(0x2)},
					},
					Street:       table.Street(0x2),
					TurnPos:      0x1,
					BtnPos:       0x1,
					StreetAction: 0x2,
					CallAmount:   chips.Chips(0),
					BSC: struct {
						Amount   chips.Chips      "json:\"amount\""
						Addition chips.Chips      "json:\"addition\""
						Action   table.ActionKind "json:\"action\""
					}{Amount: 0.07999999821186066, Addition: 0.07999999821186066, Action: table.ActionKind(0x1)},
					PSC:  chips.List{0, 0.07999999821186066},
					PSAC: []uint8{0x1, 0x1},
					PSLA: []table.ActionKind{0x5, 0x1},
					Previous: &table.State{
						Players: table.Players{
							{Paid: 0.03999999910593033, Status: table.Status(0x2)},
							{Paid: 0.03999999910593033, Status: table.Status(0x2)},
						},
						Street:       table.Street(0x2),
						TurnPos:      0x1,
						BtnPos:       0x1,
						StreetAction: 0x1,
						CallAmount:   chips.Chips(0),
						BSC: struct {
							Amount   chips.Chips      "json:\"amount\""
							Addition chips.Chips      "json:\"addition\""
							Action   table.ActionKind "json:\"action\""
						}{},
						PSC:  chips.List{0, 0},
						PSAC: []uint8{0x1, 0x0},
						PSLA: []table.ActionKind{0x5, 0x0},
						Previous: &table.State{
							Players: table.Players{
								{Paid: 0.03999999910593033, Status: table.Status(0x2)},
								{Paid: 0.03999999910593033, Status: table.Status(0x2)},
							},
							Street:       table.Street(0x2),
							TurnPos:      0x0,
							BtnPos:       0x1,
							StreetAction: 0x1,
							CallAmount:   chips.Chips(0),
							BSC: struct {
								Amount   chips.Chips      "json:\"amount\""
								Addition chips.Chips      "json:\"addition\""
								Action   table.ActionKind "json:\"action\""
							}{},
							PSC:  chips.List{0, 0},
							PSAC: []uint8{0x1, 0x0},
							PSLA: []table.ActionKind{0x5, 0x0},
							Previous: &table.State{
								Players: table.Players{
									{Paid: 0.03999999910593033, Status: table.Status(0x2)},
									{Paid: 0.03999999910593033, Status: table.Status(0x2)},
								},
								Street:       table.Street(0x2),
								TurnPos:      0x0,
								BtnPos:       0x1,
								StreetAction: 0x0,
								CallAmount:   chips.Chips(0),
								BSC: struct {
									Amount   chips.Chips      "json:\"amount\""
									Addition chips.Chips      "json:\"addition\""
									Action   table.ActionKind "json:\"action\""
								}{},
								PSC:  chips.List{0, 0},
								PSAC: []uint8{0x0, 0x0},
								PSLA: []table.ActionKind{0x0, 0x0},
								Previous: &table.State{
									Players: table.Players{
										{Paid: 0.03999999910593033, Status: table.Status(0x2)},
										{Paid: 0.03999999910593033, Status: table.Status(0x2)},
									},
									Street:       table.Street(0x1),
									TurnPos:      0x0,
									BtnPos:       0x1,
									StreetAction: 0x4,
									CallAmount:   chips.Chips(0),
									BSC: struct {
										Amount   chips.Chips      "json:\"amount\""
										Addition chips.Chips      "json:\"addition\""
										Action   table.ActionKind "json:\"action\""
									}{Amount: 0.03999999910593033, Addition: 0.019999999552965164, Action: table.ActionKind(0x3)},
									PSC:  chips.List{0.03999999910593033, 0.03999999910593033},
									PSAC: []uint8{0x2, 0x2},
									PSLA: []table.ActionKind{0x5, 0x6},
									Previous: &table.State{
										Players: table.Players{
											{Paid: 0.03999999910593033, Status: table.Status(0x2)},
											{Paid: 0.03999999910593033, Status: table.Status(0x2)},
										},
										Street:       table.Street(0x1),
										TurnPos:      0x0,
										BtnPos:       0x1,
										StreetAction: 0x3,
										CallAmount:   chips.Chips(0),
										BSC: struct {
											Amount   chips.Chips      "json:\"amount\""
											Addition chips.Chips      "json:\"addition\""
											Action   table.ActionKind "json:\"action\""
										}{Amount: 0.03999999910593033, Addition: 0.019999999552965164, Action: table.ActionKind(0x3)},
										PSC:  chips.List{0.03999999910593033, 0.03999999910593033},
										PSAC: []uint8{0x1, 0x2},
										PSLA: []table.ActionKind{0x3, 0x6},
										Previous: &table.State{
											Players: table.Players{
												{Paid: 0.03999999910593033, Status: table.Status(0x2)},
												{Paid: 0.03999999910593033, Status: table.Status(0x2)},
											},
											Street:       table.Street(0x1),
											TurnPos:      0x1,
											BtnPos:       0x1,
											StreetAction: 0x3,
											CallAmount:   chips.Chips(0.019999999552965164),
											BSC: struct {
												Amount   chips.Chips      "json:\"amount\""
												Addition chips.Chips      "json:\"addition\""
												Action   table.ActionKind "json:\"action\""
											}{Amount: 0.03999999910593033, Addition: 0.019999999552965164, Action: table.ActionKind(0x3)},
											PSC:      chips.List{0.03999999910593033, 0.03999999910593033},
											PSAC:     []uint8{0x1, 0x2},
											PSLA:     []table.ActionKind{0x3, 0x6},
											Previous: nil,
										},
									},
								},
							},
						},
					},
				},
			},
		},
	},
	Hole:      card.Cards{0x33, 0x2f},
	Community: card.Cards{0x2f, 0x1e, 0x2e},
}

func TestMarshalUnmarshal(t *testing.T) {
	s, err := exampleState.MarshalBinary()
	if err != nil {
		t.Fatalf("MarshalBinary failed: %v", err)
	}

	var got bot.State
	if err := got.UnmarshalBinary(s); err != nil {
		t.Fatalf("UnmarshalBinary failed: %v", err)
	}

	require.Equal(t, exampleState, got)
}

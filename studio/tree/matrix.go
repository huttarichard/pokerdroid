package studiotree

import (
	"fmt"

	"github.com/pokerdroid/poker/abs"
	absp "github.com/pokerdroid/poker/abs/pack"
	"github.com/pokerdroid/poker/card"
	"github.com/pokerdroid/poker/policy"
	"github.com/pokerdroid/poker/table"
	"github.com/pokerdroid/poker/tree"
)

type Cluster struct {
	Cards      []card.Cards     `json:"cards"`
	Policies   []*policy.Policy `json:"policies"`
	Policy     *policy.Policy   `json:"policy"`
	AvgStrat   []float64        `json:"avg_strat"`
	ReachProbs []float64        `json:"reach_probs"`
	Reach      float64          `json:"reach"`
	Clusters   []abs.Cluster    `json:"clusters"`
}

type Matrix [13][13]Cluster

type MatrixBuilder struct {
	data   [13][13]Cluster
	player *tree.Player
	board  card.Cards
	abs    *absp.Abs
	moves  []PlayerAction
}

func NewMatrixBuilder(player *tree.Player, board card.Cards, abs *absp.Abs, moves []PlayerAction) *MatrixBuilder {
	m := &MatrixBuilder{
		player: player,
		board:  board,
		abs:    abs,
		moves:  moves,
	}

	return m
}

func (m *MatrixBuilder) Build() (Matrix, error) {
	for x := 0; x < 13; x++ {
		for y := 0; y < 13; y++ {
			if err := m.process(x, y); err != nil {
				return m.data, err
			}
		}
	}
	return m.data, m.average()
}

func (m *MatrixBuilder) process(x, y int) error {
	possibleCards := card.CardsInCoordsWithBlockersAt(x, y, m.board)
	if len(possibleCards) == 0 {
		return nil // No valid hands at these coordinates with given blockers
	}

	var t int
	var reach float64

	// Process each possible hand at these coordinates
	for _, cards := range possibleCards {
		handCards := append(cards, m.board...)
		cl := m.abs.Map(handCards)

		m.data[x][y].Cards = append(m.data[x][y].Cards, cards)
		m.data[x][y].Clusters = append(m.data[x][y].Clusters, cl)

		if pol, ok := m.player.Actions.Policies.Get(cl); ok {
			m.data[x][y].Policies = append(m.data[x][y].Policies, pol)
		} else {
			m.data[x][y].Policies = append(m.data[x][y].Policies, nil)
		}

		r := m.reach(cards)
		m.data[x][y].ReachProbs = append(m.data[x][y].ReachProbs, r)
		reach += r
		t++
	}

	m.data[x][y].Reach = reach / float64(t)

	return nil
}

func (m *MatrixBuilder) average() error {
	for x := 0; x < 13; x++ {
		for y := 0; y < 13; y++ {
			if len(m.data[x][y].Policies) == 0 {
				continue
			}

			p := policy.New(len(m.player.Actions.Actions))
			count := float64(len(m.data[x][y].Policies))

			for _, p2 := range m.data[x][y].Policies {
				if p2 == nil {
					continue
				}
				for i := range p.Strategy {
					p.Strategy[i] += p2.Strategy[i]
					p.RegretSum[i] += p2.RegretSum[i]
					p.StrategySum[i] += p2.StrategySum[i]
					p.Baseline[i] += p2.Baseline[i]
				}
				p.Iteration += p2.Iteration
			}

			for i := range p.Strategy {
				p.Strategy[i] /= count
				p.RegretSum[i] /= count
				p.StrategySum[i] /= count
				p.Baseline[i] /= count
			}

			m.data[x][y].Policy = p
			m.data[x][y].AvgStrat = p.GetAverageStrategy()
		}
	}
	return nil
}

func (m *MatrixBuilder) reach(cards card.Cards) float64 {
	prob := float64(1.0)

	for _, move := range m.moves {
		if move.Player.TurnPos != m.player.TurnPos {
			continue
		}

		var expected int
		switch move.Player.State.Street {
		case table.Preflop:
			expected = 2
		case table.Flop:
			expected = 5
		case table.Turn:
			expected = 6
		case table.River:
			expected = 7
		default:
			panic("unrecognized street")
		}

		// Calculate number of board cards needed.
		req := expected - len(cards)
		if len(m.board) < req {
			panic(fmt.Sprintf("insufficient board cards: expected %d, got %d", req, len(m.board)))
		}

		handCards := append(cards, m.board[:req]...)
		if len(handCards) != expected {
			panic(fmt.Sprintf("illegal card count: got %d, expected %d", len(handCards), expected))
		}

		cl := m.abs.Map(handCards)

		pol, ok := move.Player.Actions.Policies.Get(cl)
		if !ok {
			continue
		}

		strat := pol.GetAverageStrategy()
		if move.Action < 0 || move.Action >= len(strat) {
			panic("invalid action index in reach calculation")
		}
		prob *= strat[move.Action]
	}
	return prob
}

func (m *MatrixBuilder) Get() [13][13]Cluster {
	return m.data
}

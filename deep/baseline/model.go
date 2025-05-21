package baselinenn

import (
	"encoding/gob"
	"fmt"
	"io"

	"github.com/nlpodyssey/spago/mat"
	"github.com/nlpodyssey/spago/mat/float"
	"github.com/nlpodyssey/spago/nn"
	"github.com/nlpodyssey/spago/nn/activation"
	"github.com/nlpodyssey/spago/nn/linear"
	"github.com/nlpodyssey/spago/nn/normalization/batchnorm"
	"github.com/pokerdroid/poker/deep"
	"github.com/pokerdroid/poker/frand"
	"github.com/pokerdroid/poker/table"
)

func init() {
	gob.Register(&Model{})
}

type StandardModule = nn.ModuleList[nn.StandardModel]

// Model defines a simple neural network for sine function approximation
type Model struct {
	nn.Module
	Layers  StandardModule
	Actions []table.DiscreteAction
}

func NewFromFile[T float.DType](r io.Reader) (*Model, error) {
	dec := gob.NewDecoder(r)
	var m Model
	if err := dec.Decode(&m); err != nil {
		return nil, err
	}
	return &m, nil
}

type ModelParams struct {
	NumPlayers         int       `json:"num_players"`
	MaxActionsPerRound int       `json:"max_actions_per_round"`
	Layers             int       `json:"layers"`
	HiddenSize         int       `json:"hidden_size"`
	BetSizing          []float32 `json:"bets"`
}

// NewModel creates a new model for sine approximation
func NewModel[T float.DType](p ModelParams) *Model {
	total, _ := EncodeFeatures[T](
		uint8(p.NumPlayers),
		uint8(p.MaxActionsPerRound),
	)

	// Create layers with proper activation modules
	ll := StandardModule{
		linear.New[T](total, p.HiddenSize),
	}

	for i := 0; i < p.Layers; i++ {
		ll = append(ll,
			linear.New[T](p.HiddenSize, p.HiddenSize),
			// todo make batchnorm work
			batchnorm.New[T](p.HiddenSize),
			activation.New(activation.ReLU),
		)
	}

	bets := []table.DiscreteAction{
		table.DAllIn,
		table.DFold,
		table.DCall,
		table.DCheck,
	}

	for _, b := range p.BetSizing {
		bets = append(bets, table.DiscreteAction(b))
	}

	ll = append(ll,
		linear.New[T](p.HiddenSize, len(bets)),
		activation.New(activation.Tanh),
	)

	return &Model{Layers: ll, Actions: bets}
}

func NewLinear[T float.DType](inputSize, layers, hiddenSize, outputSize int) StandardModule {
	// Create layers with proper activation modules
	ll := StandardModule{
		linear.New[T](inputSize, hiddenSize),
	}

	for i := 0; i < layers; i++ {
		ll = append(ll,
			linear.New[T](hiddenSize, hiddenSize),
			// todo make batchnorm work
			batchnorm.New[T](hiddenSize),
			activation.New(activation.ReLU),
		)
	}

	ll = append(
		ll,
		linear.New[T](hiddenSize, outputSize),
		activation.New(activation.Tanh),
	)

	return ll
}

// Forward performs the forward pass
func (m *Model) Forward(xs ...mat.Tensor) []mat.Tensor {
	// ModuleList.Forward handles the sequential processing
	out := m.Layers.Forward(xs...)
	return out
}

// InitRandom initializes the model weights
func (m *Model) InitRandom(rng frand.Rand) *Model {
	deep.ApplyXavierUniform(m, rng)
	return m
}

// Save saves the model to a file
func (m *Model) Save(w io.Writer) error {
	// Create encoder and encode the model
	enc := gob.NewEncoder(w)

	if err := enc.Encode(m); err != nil {
		return fmt.Errorf("failed to encode model: %w", err)
	}

	return nil
}

func (m *Model) Load(r io.Reader) error {
	dec := gob.NewDecoder(r)
	return dec.Decode(m)
}

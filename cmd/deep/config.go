package cmddeep

import (
	"encoding/json"
	"os"

	absp "github.com/pokerdroid/poker/abs/pack"
	"github.com/pokerdroid/poker/deep"
	"github.com/pokerdroid/poker/tree"
)

type SolutionConfig struct {
	Abs string `json:"abs"`
	Dir string `json:"dir"`
}

func NewRootsFromConfig(cfg SolutionConfig) (deep.Roots, error) {
	abs, err := absp.NewFromFile(cfg.Abs)
	if err != nil {
		return nil, err
	}

	roots, err := tree.NewFileRootsFromDir(cfg.Dir)
	if err != nil {
		return nil, err
	}

	solutions := make(deep.Roots, len(roots))
	for i, root := range roots {
		solutions[i] = deep.Root{
			Abs:      abs,
			FileRoot: root,
		}
	}

	return solutions, nil
}

type SolutionConfigs []SolutionConfig

func NewRootsFromConfigs(cfgs SolutionConfigs) (ss deep.Roots, err error) {
	for _, cfg := range cfgs {
		xs, err := NewRootsFromConfig(cfg)
		if err != nil {
			return nil, err
		}
		ss = append(ss, xs...)
	}
	return ss, nil
}

type TrainConfig[T any] struct {
	// Game configuration
	Model T `json:"model"`
	// Solutions configuration
	Solutions SolutionConfigs `json:"solutions"`
}

func NewTrainConfig[T any](config string) (TrainConfig[T], error) {
	cfg := TrainConfig[T]{}

	f, err := os.Open(config)
	if err != nil {
		return TrainConfig[T]{}, err
	}
	defer f.Close()

	err = json.NewDecoder(f).Decode(&cfg)
	if err != nil {
		return TrainConfig[T]{}, err
	}
	return cfg, nil
}

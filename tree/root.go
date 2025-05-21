package tree

import (
	"io"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/pokerdroid/poker/chips"
	"github.com/pokerdroid/poker/table"
)

type Root struct {
	AbsID     uuid.UUID        `json:"abs_id"`
	States    uint32           `json:"states"`
	Nodes     uint32           `json:"nodes"`
	Next      Node             `json:"-"`
	Params    table.GameParams `json:"params"`
	Iteration uint64           `json:"iteration"`
	State     *table.State     `json:"-"`
	Full      bool             `json:"-"`
}

func NewRoot(prms table.GameParams) (r *Root, err error) {
	r = &Root{
		Params: prms,
		State:  table.NewState(prms),
	}

	r.State, err = table.MakeInitialBets(prms, r.State)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func NewRootFromReadSeeker(rs io.ReadSeeker) (*Root, error) {
	root := &Root{}
	return root, root.ReadBinary(rs)
}

// FindClosestRoot finds the root with the closest effective
// stack to the given turn position.
//
// This is handy if you have bunch of roots and you want to find
// the one that is closest to the current game state.
func FindClosestRoot(roots []*Root, p table.GameParams, turn uint8) *Root {
	rootsP := make([]*Root, 0, len(roots))
	for _, r := range roots {
		if r.Params.NumPlayers != p.NumPlayers {
			continue
		}
		rootsP = append(rootsP, r)
	}

	if len(rootsP) == 0 {
		return nil
	}

	ratio := chips.NewFromFloat(1).Div(p.SbAmount)
	effs := p.EffectiveStack(turn).Mul(ratio)
	closest := rootsP[0]
	diff := closest.Params.EffectiveStack(turn).Sub(effs).Abs()

	for _, r := range rootsP[1:] {
		std := r.Params.EffectiveStack(turn)
		dff := std.Sub(effs).Abs()

		if dff < diff {
			diff = dff
			closest = r
		}
	}
	return closest
}

func NewFromFile(path string) (*Root, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	root := &Root{}
	err = root.UnmarshalBinary(data)
	if err != nil {
		return root, err
	}
	return root, nil
}

func (ch *Root) Kind() NodeKind {
	return NodeKindRoot
}

func (ch *Root) GetParent() Node {
	return nil
}

func TstRootFromEnv(t *testing.T, env string) *Root {
	es := os.Getenv(env)
	if es == "" {
		t.Skipf("env %s is not set", env)
	}
	root, err := NewFromFile(es)
	if err != nil {
		t.Fatal(err)
	}
	return root
}

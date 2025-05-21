package deep

import (
	absp "github.com/pokerdroid/poker/abs/pack"
	"github.com/pokerdroid/poker/tree"
)

type Root struct {
	Abs      *absp.Abs
	FileRoot *tree.FileRoot
}

type Roots []Root

func (s Roots) Close() error {
	for _, s := range s {
		err := s.FileRoot.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

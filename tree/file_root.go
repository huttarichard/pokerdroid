package tree

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type FileRoot struct {
	Root *Root
	File *os.File
}

func (r *FileRoot) NewRoot() (fr *Root, err error) {
	file, err := os.Open(r.File.Name())
	if err != nil {
		return nil, err
	}

	return NewRootFromReadSeeker(file)
}

func (r *FileRoot) Close() error {
	return r.File.Close()
}

type FileRoots []*FileRoot

func NewFileRootsFromDir(dir string) (FileRoots, error) {
	var roots FileRoots

	// Find all tree*.bin files recursively
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		match, err := filepath.Match("tree*.bin", filepath.Base(path))
		if err != nil {
			return fmt.Errorf("matching %s: %w", path, err)
		}

		if info.IsDir() || !match {
			return nil
		}

		// Open file
		f, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("opening %s: %w", path, err)
		}

		// Load root from file
		root, err := NewRootFromReadSeeker(f)
		if err != nil {
			return fmt.Errorf("loading root from %s: %w", path, err)
		}
		roots = append(roots, &FileRoot{Root: root, File: f})
		return nil
	})

	return roots, err
}

func (f FileRoots) Roots() []*Root {
	roots := make([]*Root, len(f))
	for i, r := range f {
		roots[i] = r.Root
	}
	return roots
}

func (f FileRoots) String() string {
	var s strings.Builder
	ws := s.WriteString
	for _, r := range f {
		ws("Solution:\n")
		ws(fmt.Sprintf("\tInitialStacks: %v\n", r.Root.Params.InitialStacks))
		ws(fmt.Sprintf("\tBetSizes: %v\n", r.Root.Params.BetSizes))
		ws(fmt.Sprintf("\tMaxActionsPerRound: %d\n", r.Root.Params.MaxActionsPerRound))
		ws("\n")
	}
	return s.String()
}

func (f FileRoots) Close() error {
	for _, r := range f {
		err := r.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

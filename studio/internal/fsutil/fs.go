package fsutil

import (
	"errors"
	"io/fs"
	"path/filepath"
)

type DirMakerFS interface {
	MkdirAll(path string) error
}

func MkdirAll(fsx fs.FS, path string) error {
	fsw, ok := fsx.(DirMakerFS)
	if !ok {
		return errors.New("unsupported fs type")
	}
	return fsw.MkdirAll(path)
}

type RemoverFS interface {
	RemoveAll(path string) error
}

func RemoveAll(fsx fs.FS, path string) error {
	fsw, ok := fsx.(RemoverFS)
	if !ok {
		return errors.New("unsupported fs type")
	}
	return fsw.RemoveAll(path)
}

type WriteFileFS interface {
	WriteFile(path string, data []byte, perm fs.FileMode) (err error)
}

func WriteFileSimple(fsx fs.FS, path string, data []byte) (err error) {
	fsw, ok := fsx.(WriteFileFS)
	if !ok {
		return errors.New("unsupported fs type")
	}

	fsmk, ok := fsx.(DirMakerFS)
	if ok {
		dir := filepath.Dir(path)
		fsmk.MkdirAll(dir)
	}

	return fsw.WriteFile(path, data, fs.ModePerm)
}

type PathFS interface {
	fs.FS
	Path() string
}

type pathInFS struct {
	fs.FS
	path string
}

func PathInFS(fs fs.FS, path string) *pathInFS {
	return &pathInFS{fs, path}
}

func (p *pathInFS) Path() string {
	return p.path
}

func GetPathInFS(fs fs.FS, def ...string) string {
	fsw, ok := fs.(PathFS)
	if !ok {
		if len(def) > 0 {
			return def[0]
		}
		return "."
	}
	return fsw.Path()
}

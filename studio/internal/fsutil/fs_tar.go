package fsutil

import (
	"io"
	"io/fs"

	"github.com/nlepage/go-tarfs"
)

type TarFS struct {
	fs fs.FS
}

var (
	_ fs.FS         = (*TarFS)(nil)
	_ fs.ReadDirFS  = (*TarFS)(nil)
	_ fs.ReadFileFS = (*TarFS)(nil)
	_ fs.StatFS     = (*TarFS)(nil)
	_ fs.GlobFS     = (*TarFS)(nil)
	_ fs.SubFS      = (*TarFS)(nil)
)

func NewTarFS(r io.Reader) (*TarFS, error) {
	t := &TarFS{}
	tfs, err := tarfs.New(r)
	if err != nil {
		return nil, err
	}
	t.fs = tfs
	return t, nil
}

// Open implements fs.FS.
func (md *TarFS) Open(filename string) (fs.File, error) {
	return md.fs.Open(filename)
}

func (md *TarFS) ReadFile(name string) ([]byte, error) {
	return fs.ReadFile(md.fs, name)
}

func (md *TarFS) ReadDir(name string) ([]fs.DirEntry, error) {
	return fs.ReadDir(md.fs, name)
}

func (md *TarFS) Stat(name string) (fs.FileInfo, error) {
	return fs.Stat(md.fs, name)
}

func (md *TarFS) Glob(name string) ([]string, error) {
	return fs.Glob(md.fs, name)
}

func (md *TarFS) Sub(name string) (fs.FS, error) {
	return fs.Sub(md.fs, name)
}

package fsutil

import (
	"embed"
	"io/fs"
)

type EmbedFS struct {
	fs embed.FS
}

var (
	_ fs.FS         = (*EmbedFS)(nil)
	_ fs.ReadDirFS  = (*EmbedFS)(nil)
	_ fs.ReadFileFS = (*EmbedFS)(nil)
	_ fs.StatFS     = (*EmbedFS)(nil)
	_ fs.GlobFS     = (*EmbedFS)(nil)
	_ fs.SubFS      = (*EmbedFS)(nil)
)

func NewEmbedFS(r embed.FS) *EmbedFS {
	return &EmbedFS{fs: r}
}

// Open implements fs.FS.
func (md *EmbedFS) Open(filename string) (fs.File, error) {
	return md.fs.Open(filename)
}

func (md *EmbedFS) ReadFile(name string) ([]byte, error) {
	return fs.ReadFile(md.fs, name)
}

func (md *EmbedFS) ReadDir(name string) ([]fs.DirEntry, error) {
	return fs.ReadDir(md.fs, name)
}

func (md *EmbedFS) Stat(name string) (fs.FileInfo, error) {
	return fs.Stat(md.fs, name)
}

func (md *EmbedFS) Glob(name string) ([]string, error) {
	return fs.Glob(md.fs, name)
}

func (md *EmbedFS) Sub(name string) (fs.FS, error) {
	return fs.Sub(md.fs, name)
}

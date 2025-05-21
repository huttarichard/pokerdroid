package fsutil

import (
	"bytes"
	"io/fs"
	"os"
	"time"

	"github.com/liamg/memoryfs"
)

// MountedDir implements a read-only FS in a mounted directory.
// File is writable.
type MemFS struct {
	fs *memoryfs.FS
}

var (
	_ fs.FS         = (*MemFS)(nil)
	_ fs.ReadDirFS  = (*MemFS)(nil)
	_ fs.ReadFileFS = (*MemFS)(nil)
	_ fs.StatFS     = (*MemFS)(nil)
	_ fs.GlobFS     = (*MemFS)(nil)
	_ fs.SubFS      = (*MemFS)(nil)
	_ DirMakerFS    = (*MemFS)(nil)
	_ RemoverFS     = (*MemFS)(nil)
)

func NewMemFS() *MemFS {
	m := &MemFS{memoryfs.New()}
	return m
}

// Open implements fs.FS.
// Open opens the underlying file which must exist and not be a directory for
// reading and writing with ModePerm. If an error occurs it is returned.
func (md *MemFS) Open(filename string) (fs.File, error) {
	x, err := md.fs.Open(filename)
	switch err.(type) {
	case *fs.PathError:
		err = md.fs.WriteFile(filename, []byte{}, os.ModePerm)
		if err != nil {
			return x, err
		}
		x, err = md.fs.Open(filename)
	}
	return x, err
}

func (md *MemFS) ReadFile(name string) ([]byte, error) {
	return fs.ReadFile(md.fs, name)
}

func (md *MemFS) ReadDir(name string) ([]fs.DirEntry, error) {
	return fs.ReadDir(md.fs, name)
}

func (md *MemFS) Stat(name string) (fs.FileInfo, error) {
	return fs.Stat(md.fs, name)
}

func (md *MemFS) Glob(name string) ([]string, error) {
	return fs.Glob(md.fs, name)
}

func (md *MemFS) Sub(name string) (fs.FS, error) {
	return fs.Sub(md.fs, name)
}

func (md *MemFS) WriteFile(path string, data []byte, perm fs.FileMode) (err error) {
	return md.fs.WriteFile(path, data, perm)
}

func (md *MemFS) RemoveAll(path string) error {
	return md.fs.RemoveAll(path)
}

func (md *MemFS) MkdirAll(path string) error {
	return md.fs.MkdirAll(path, os.ModePerm)
}

type FileInfo struct {
	size    int64
	path    string
	modtime time.Time
	isDir   bool
	sys     any
}

var _ fs.FileInfo = (*FileInfo)(nil)

func (m *FileInfo) Name() string {
	return m.path
}

func (m *FileInfo) Size() int64 {
	return m.size
}

func (m *FileInfo) Mode() fs.FileMode {
	return fs.ModePerm
}

func (m *FileInfo) ModTime() time.Time {
	return m.modtime
}

func (m *FileInfo) IsDir() bool {
	return m.isDir
}

func (m *FileInfo) Sys() any {
	return m.sys
}

type MemFile struct {
	Path     string
	Contents *bytes.Buffer
	info     *FileInfo
}

func NewMemFile(path string, buf *bytes.Buffer) *MemFile {
	info := &FileInfo{
		size:    int64(buf.Len()),
		path:    path,
		modtime: time.Now(),
		isDir:   false,
		sys:     nil,
	}
	return &MemFile{
		Path:     path,
		Contents: buf,
		info:     info,
	}
}

var _ fs.File = (*MemFile)(nil)

func (m *MemFile) Stat() (fs.FileInfo, error) {
	return m.info, nil
}

func (m *MemFile) Read(p []byte) (n int, err error) {
	return m.Contents.Read(p)
}

func (m *MemFile) Close() error {
	m.Contents = nil
	return nil
}

package fsutil

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// MountedDir implements a writable FS in a mounted directory.
// File is writable.
type DirFS struct {
	// root is the path to the mounted directory.
	dir string
	fs  fs.FS
}

var (
	_ fs.FS         = (*DirFS)(nil)
	_ fs.ReadDirFS  = (*DirFS)(nil)
	_ fs.ReadFileFS = (*DirFS)(nil)
	_ fs.StatFS     = (*DirFS)(nil)
	_ fs.GlobFS     = (*DirFS)(nil)
	_ fs.SubFS      = (*DirFS)(nil)
	_ DirMakerFS    = (*DirFS)(nil)
	_ RemoverFS     = (*DirFS)(nil)
)

// NewMountedDir returns a new DirFS instance
func NewDirFS(root string) *DirFS {
	return &DirFS{dir: root, fs: os.DirFS(root)}
}

// Open implements fs.FS.
// Open opens the underlying file which must exist and not be a directory for
// reading and writing with ModePerm. If an error occurs it is returned.
func (md *DirFS) Open(path string) (fs.File, error) {
	pth, err := clean(md.dir, path)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(pth)
	switch {
	case errors.Is(err, os.ErrNotExist):
		return nil, &fs.PathError{Op: "open", Path: pth, Err: fs.ErrNotExist}
	case err != nil:
		return nil, &fs.PathError{Op: "open", Path: pth, Err: err}
	case fi.IsDir():
		return &Dir{f: fi, base: pth}, nil
	default:
		pth, _ = strings.CutPrefix(pth, md.dir)
		pth = strings.TrimPrefix(pth, "/")
		return md.fs.Open(pth)
	}
}

func (md *DirFS) ReadFile(name string) ([]byte, error) {
	return fs.ReadFile(md.fs, name)
}

func (md *DirFS) ReadDir(name string) ([]fs.DirEntry, error) {
	return fs.ReadDir(md.fs, name)
}

func (md *DirFS) Stat(name string) (fs.FileInfo, error) {
	return fs.Stat(md.fs, name)
}

func (md *DirFS) Glob(name string) ([]string, error) {
	return fs.Glob(md.fs, name)
}

func (md *DirFS) Sub(name string) (fs.FS, error) {
	sub, err := fs.Sub(md.fs, name)
	if err != nil {
		return nil, err
	}
	return &DirFS{dir: filepath.Join(md.dir, name), fs: sub}, nil
}

func (md *DirFS) WriteFile(path string, data []byte, perm fs.FileMode) (err error) {
	pth, err := clean(md.dir, path)
	if err != nil {
		return err
	}
	var f *os.File
	var d bool
	stat, err := os.Stat(pth)
	switch err.(type) {
	case *fs.PathError:
	case nil:
		d = stat.IsDir()
	default:
		return err
	}
	if d {
		return errors.New("path is a directory")
	}
	f, err = os.OpenFile(pth, os.O_RDWR|os.O_CREATE, perm)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(data)
	return
}

func (md *DirFS) RemoveAll(path string) error {
	pth, err := clean(md.dir, path)
	if err != nil {
		return err
	}
	return os.RemoveAll(pth)
}

func (md *DirFS) MkdirAll(path string) error {
	pth, err := clean(md.dir, path)
	if err != nil {
		return err
	}
	return os.MkdirAll(pth, os.ModePerm)
}

func clean(base, ptx string) (string, error) {
	base = filepath.Clean(base)

	var pth string
	pth = filepath.Join(base, ptx)
	pth = filepath.Clean(pth)

	if !strings.HasPrefix(pth, base) {
		return "", errors.New("path is outside of the mounted directory")
	}

	return pth, nil
}

// An Dir is a directory open for reading.
type Dir struct {
	f    fs.FileInfo // the directory file itself
	base string
}

func (d *Dir) Read([]byte) (int, error) {
	return 0, &fs.PathError{Op: "read", Path: d.f.Name(), Err: errors.New("is a directory")}
}
func (d *Dir) Close() error               { return nil }
func (d *Dir) Stat() (fs.FileInfo, error) { return d.f, nil }

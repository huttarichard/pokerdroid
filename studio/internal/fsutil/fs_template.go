package fsutil

import (
	"bytes"
	"html/template"
	"io/fs"
)

type TemplateFS struct {
	fs       fs.FS
	provider func(f string) (any, error)
	tmpl     *template.Template
}

var (
	_ fs.FS         = (*TemplateFS)(nil)
	_ fs.ReadDirFS  = (*TemplateFS)(nil)
	_ fs.ReadFileFS = (*TemplateFS)(nil)
	_ fs.StatFS     = (*TemplateFS)(nil)
	_ fs.GlobFS     = (*TemplateFS)(nil)
	_ fs.SubFS      = (*TemplateFS)(nil)
)

func NewTemplate(fsx fs.FS, pattern []string, provider func(f string) (any, error)) (fs.FS, error) {
	t, err := template.ParseFS(fsx, pattern...)
	if err != nil {
		return nil, err
	}

	return &TemplateFS{
		fs:       fsx,
		provider: provider,
		tmpl:     t,
	}, nil
}

func (a *TemplateFS) Open(name string) (fs.File, error) {
	data, err := a.fs.Open(name)
	if err != nil {
		return nil, err
	}

	var found bool
	for _, p := range a.tmpl.Templates() {
		if p.Name() == name {
			found = true
			break
		}
	}
	if !found {
		return data, nil
	}

	buf := bytes.NewBuffer(nil)

	payload, err := a.provider(name)
	if err != nil {
		return nil, err
	}

	err = a.tmpl.ExecuteTemplate(buf, name, payload)
	if err != nil {
		return nil, err
	}

	return NewMemFile(name, buf), nil
}

func (md *TemplateFS) ReadFile(name string) ([]byte, error) {
	return fs.ReadFile(md.fs, name)
}

func (md *TemplateFS) ReadDir(name string) ([]fs.DirEntry, error) {
	return fs.ReadDir(md.fs, name)
}

func (md *TemplateFS) Stat(name string) (fs.FileInfo, error) {
	return fs.Stat(md.fs, name)
}

func (md *TemplateFS) Glob(name string) ([]string, error) {
	return fs.Glob(md.fs, name)
}

func (md *TemplateFS) Sub(name string) (fs.FS, error) {
	return fs.Sub(md.fs, name)
}

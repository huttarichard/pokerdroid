package jsc

import (
	"encoding/json"
	"io"
	"io/fs"
)

type Import struct {
	Path     string `json:"path"`
	Kind     string `json:"kind"`
	Original string `json:"original"`
}

type MetaInput struct {
	Bytes   int64    `json:"bytes"`
	Imports []Import `json:"imports"`
	Format  string   `json:"format"`
}

type MetaOuputInput struct {
	Bytes int64 `json:"bytesInOutput"`
}

type MetaOuput struct {
	Bytes      int64                     `json:"bytes"`
	Entrypoint string                    `json:"entrypoint"`
	Imports    []Import                  `json:"imports"`
	Inputs     map[string]MetaOuputInput `json:"inputs"`
}

type MetaOutputMap map[string]MetaOuput

// FindByEntrypoint returns the key of the output file that has the given entrypoint.
//
// Example usage:
//   fs := fs.FS(someFS)
//   meta, err := jsc.ReadMetafile(fs, "dist/meta.json")
//   if err != nil {
// 	  return err
//   }
//   file, ok := meta.Outputs.FindByEntrypoint("src/main.tsx")
//   if !ok {
// 	  return errors.New("no entrypoint file found")
//   }

func (m MetaOutputMap) FindByEntrypoint(entrypoint string) (string, bool) {
	for k, r := range m {
		if r.Entrypoint == entrypoint {
			return k, true
		}
	}
	return "", false
}

type Metafile struct {
	Inputs  map[string]MetaInput `json:"inputs"`
	Outputs MetaOutputMap        `json:"outputs"`
}

func ReadMetafile(fs fs.FS, path string) (Metafile, error) {
	f, err := fs.Open(path)
	if err != nil {
		return Metafile{}, err
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		return Metafile{}, err
	}

	var ff Metafile
	err = json.Unmarshal(data, &ff)
	if err != nil {
		return ff, err
	}

	return ff, nil
}

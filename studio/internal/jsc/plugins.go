package jsc

import (
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/evanw/esbuild/pkg/api"
)

// EsbuildLogger will replace logger of esbuild.
// Please use LogLevel: api.LogLevelSilent to not duplicate logs
func EsbuildLogger() api.Plugin {
	logger := log.Default()

	printMsg := func(msg api.Message) {
		if msg.Location != nil {
			logger.Printf("error (%s): %s", msg.Location.File, msg.Text)
			return
		}
		logger.Printf("error: %s", msg.Text)
	}

	// We use this as its hard to override
	// logger of esbuild
	onCompilerEnd := func(result *api.BuildResult) (api.OnEndResult, error) {
		for _, err := range result.Errors {
			printMsg(err)
		}
		for _, err := range result.Warnings {
			printMsg(err)
		}
		if len(result.Errors) == 0 {
			logger.Println("build succeeded")
		}
		return api.OnEndResult{}, nil
	}

	return api.Plugin{
		Name: "logger",
		Setup: func(build api.PluginBuild) {
			build.OnEnd(onCompilerEnd)
		},
	}
}

// EsbuildMetaWriter will write metafile to path.
func EsbuildMetaWriter(path string) api.Plugin {
	writer := func(result *api.BuildResult) (api.OnEndResult, error) {
		if result.Metafile != "" {
			os.WriteFile(path, []byte(result.Metafile), 0644)
		}
		return api.OnEndResult{}, nil
	}

	return api.Plugin{
		Name: "metawriter",
		Setup: func(build api.PluginBuild) {
			build.OnEnd(writer)
		},
	}
}

func ResolveWithFs(prefix string, fs fs.FS) api.Plugin {
	prefix = strings.ToLower(prefix)
	prefix = strings.TrimSuffix(prefix, "/")

	onLoad := func(args api.OnLoadArgs) (api.OnLoadResult, error) {
		file, err := fs.Open(args.Path)
		if err != nil {
			return api.OnLoadResult{}, err
		}
		defer file.Close()

		data, err := io.ReadAll(file)
		if err != nil {
			return api.OnLoadResult{}, err
		}
		contents := string(data)

		return api.OnLoadResult{
			Contents: &contents,
			Loader:   api.LoaderJS,
		}, nil
	}

	onResolve := func(args api.OnResolveArgs) (api.OnResolveResult, error) {
		args.Path = strings.TrimPrefix(args.Path, prefix+"/")

		path := strings.Split(args.Path, "/")

		if len(path) == 0 {
			return api.OnResolveResult{}, fmt.Errorf("invalid path: %q", args.Path)
		}

		if len(path) == 1 && len(path[0]) == 0 {
			return api.OnResolveResult{}, fmt.Errorf("invalid path: %q", args.Path)
		}

		if len(path) == 1 && len(path[0]) > 0 {
			args.Path = filepath.Join(path[0], "index.js")
		}

		return api.OnResolveResult{
			Path:      strings.TrimPrefix(args.Path, prefix+"/"),
			Namespace: "std",
		}, nil
	}

	return api.Plugin{
		Name: "resolve_with_fs",
		Setup: func(build api.PluginBuild) {
			build.OnLoad(api.OnLoadOptions{Namespace: prefix, Filter: ".*"}, onLoad)
			build.OnResolve(api.OnResolveOptions{Filter: "^" + prefix + "/"}, onResolve)
		},
	}
}

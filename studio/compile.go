package studio

import (
	"path/filepath"

	"github.com/evanw/esbuild/pkg/api"
	fe "github.com/pokerdroid/poker/studio/frontend"
	"github.com/pokerdroid/poker/studio/internal/fsutil"
	"github.com/pokerdroid/poker/studio/internal/jsc"
)

type BuildContextParams struct {
	Production bool
}

func NewBuildContext(c BuildContextParams) (bctx api.BuildContext, err error) {
	wp, err := fsutil.ResolveProjectPath(fe.Dir)
	if err != nil {
		return nil, err
	}

	distdir := filepath.Join(wp, "dist")

	return jsc.NewBuildContext(api.BuildOptions{
		EntryPoints: []string{
			"src/main.tsx",
		},
		Outdir:        distdir,
		Bundle:        true,
		Write:         true,
		AbsWorkingDir: wp,
		JSX:           api.JSXAutomatic,
		Loader: map[string]api.Loader{
			".css": api.LoaderGlobalCSS,
			".ttf": api.LoaderFile,
			".svg": api.LoaderFile,
			".js":  api.LoaderJSX,
			".ts":  api.LoaderTSX,
		},
		MinifyWhitespace:  c.Production,
		MinifySyntax:      c.Production,
		MinifyIdentifiers: c.Production,
		Define:            map[string]string{},
		Platform:          api.PlatformBrowser,
		TreeShaking:       api.TreeShakingTrue,
		Format:            api.FormatCommonJS,
		Target:            api.ES2024,
		Splitting:         false,
		Sourcemap:         api.SourceMapExternal,
		Plugins: []api.Plugin{
			jsc.EsbuildMetaWriter(filepath.Join(distdir, "meta.json")),
		},
		Metafile: true,
		// PublicPath: "dist/",
		// EntryNames: "[name]",
	})
}

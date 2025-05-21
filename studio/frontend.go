package studio

import (
	"context"
	"fmt"
	"io/fs"
	"net"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	app "github.com/pokerdroid/poker"
	fe "github.com/pokerdroid/poker/studio/frontend"
	"github.com/pokerdroid/poker/studio/internal/fsutil"
	"github.com/pokerdroid/poker/studio/internal/jsc"
	webview "github.com/webview/webview_go"
	"golang.org/x/sync/errgroup"
)

type Frontend struct {
	Router chi.Router
	FS     fs.FS
	WebV   webview.WebView
	server *http.Server
	goes   []func(ctx context.Context) error
}

func New() *Frontend {
	fs := NewStaticAssetsFS()

	w := webview.New(false)
	w.SetTitle(app.Name)
	w.SetSize(1280, 720, webview.HintNone)
	w.SetSize(1280, 720, webview.HintMin)

	router := chi.NewRouter()
	server := &http.Server{Handler: router}

	return &Frontend{
		FS:     fs,
		Router: router,
		WebV:   w,
		server: server,
	}
}

func NewDevelopment() (*Frontend, error) {
	pth, err := fsutil.ResolveProjectPath(fe.Dir)
	if err != nil {
		return nil, err
	}

	fs := fsutil.NewDirFS(pth)

	wb := webview.New(true)
	wb.SetTitle(app.Name)
	wb.SetSize(1280, 720, webview.HintNone)
	wb.SetSize(1280, 720, webview.HintMin)

	router := chi.NewRouter()
	server := &http.Server{Handler: router}

	bc, err := NewBuildContext(BuildContextParams{
		Production: false,
	})
	if err != nil {
		return nil, err
	}

	w := fsutil.NewWatch(fs, time.Millisecond*200, "dist/*.js")

	goe := func(ctx context.Context) error {
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		_, err = jsc.WatchChanges(ctx, bc)
		if err != nil {
			return err
		}

		return w.WatchCb(ctx, func(changed []string) error {
			wb.Dispatch(func() {
				wb.Eval(`location.reload()`)
			})
			return nil
		})
	}

	return &Frontend{
		FS:     fs,
		Router: router,
		WebV:   wb,
		goes:   []func(ctx context.Context) error{goe},
		server: server,
	}, nil
}

func (f *Frontend) Run(ctx context.Context, l net.Listener) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	g, ctx := errgroup.WithContext(ctx)

	hnd, err := NewHttpHandler(f.FS)
	if err != nil {
		return err
	}

	f.Router.NotFound(hnd)

	g.Go(func() error {
		<-ctx.Done()

		err := f.server.Shutdown(context.Background())
		if err != nil {
			return err
		}
		f.WebV.Terminate()
		return nil
	})

	g.Go(func() error {
		err := f.server.Serve(l)
		switch err {
		case http.ErrServerClosed:
			return nil
		case nil:
			return nil
		default:
			return err
		}
	})

	for _, gf := range f.goes {
		g.Go(func() error {
			return gf(ctx)
		})
	}

	url := fmt.Sprintf("http://%s", l.Addr())
	f.WebV.Navigate(url)
	f.WebV.Run()

	return g.Wait()

}

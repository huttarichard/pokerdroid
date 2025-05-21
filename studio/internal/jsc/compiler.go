package jsc

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/evanw/esbuild/pkg/api"
)

// NewBuildContext creates a new esbuild build context.
func NewBuildContext(copts api.BuildOptions) (bctx api.BuildContext, err error) {
	if len(copts.Outdir) > 0 {
		err = CleanDistDirectory(copts.Outdir)
		if err != nil {
			return nil, err
		}
	}

	// Silent the logger as we are using our own logger
	copts.LogLevel = api.LogLevelSilent
	copts.Plugins = append(copts.Plugins, EsbuildLogger())

	bc, cerr := api.Context(copts)
	if cerr != nil {
		return nil, errors.New(cerr.Error())
	}

	return bc, nil
}

// WatchChanges will watch for changes and rebuild the project.
func WatchChanges(ctx context.Context, bctx api.BuildContext) (context.CancelFunc, error) {
	cancel := func() {
		bctx.Cancel()
		bctx.Dispose()
	}

	// Cancel context when parent is done
	go func() {
		<-ctx.Done()
		cancel()
	}()

	// bctx.Rebuild()

	err := bctx.Watch(api.WatchOptions{})
	if err != nil {
		return cancel, err
	}

	return cancel, nil
}

// CleanDistDirectory will clean and create dist directory.
// Will keep .gitkeep file in dist directory.
func CleanDistDirectory(dist string) error {
	err := os.RemoveAll(dist)
	if err != nil {
		return err
	}

	err = os.Mkdir(dist, 0755)
	if err != nil {
		return err
	}

	f, err := os.Create(filepath.Join(dist, ".gitkeep"))
	if err != nil {
		return err
	}
	f.Close()
	return nil
}

type BuildError struct {
	Messages []string
	Severity api.MessageKind
}

func NewBuildError(severity api.MessageKind, errors []api.Message) BuildError {
	errs := api.FormatMessages(errors, api.FormatMessagesOptions{
		TerminalWidth: 0,
		Kind:          severity,
		Color:         false,
	})

	return BuildError{Messages: errs, Severity: severity}
}

func (e BuildError) Error() string {
	return strings.Join(e.Messages, "\n")
}

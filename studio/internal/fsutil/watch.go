package fsutil

import (
	"context"
	"io/fs"
	"sort"
	"time"
)

type fileW struct {
	mod time.Time
	crc uint32
}

type Watch struct {
	fs            fs.FS
	pollTimeout   time.Duration
	watchPatterns []string
	mods          map[string]fileW
}

func NewWatch(filesystem fs.FS, pollTimeout time.Duration, patterns ...string) *Watch {
	return &Watch{
		fs:            filesystem,
		pollTimeout:   pollTimeout,
		watchPatterns: patterns,
	}
}

func (w *Watch) Compare() ([]string, error) {
	var changedFiles []string
	cf := make(map[string]fileW)

	for _, pattern := range w.watchPatterns {
		matches, err := fs.Glob(w.fs, pattern)
		if err != nil {
			return nil, err
		}

		for _, match := range matches {
			fi, err := fs.Stat(w.fs, match)
			if err != nil {
				return nil, err
			}

			// First compare last modified time
			cm := fi.ModTime()

			lm, exists := w.mods[match]
			if exists && cm.Equal(lm.mod) {
				cf[match] = lm
				continue
			}

			checksum, err := CRC32(w.fs, match)
			if err != nil {
				return nil, err
			}

			cf[match] = fileW{mod: cm, crc: checksum}

			m, exists := w.mods[match]
			if !exists || m.crc != checksum {
				changedFiles = append(changedFiles, match)
			}
		}
	}

	// Check for deleted files
	for file := range w.mods {
		_, exists := cf[file]
		if exists {
			continue
		}
		changedFiles = append(changedFiles, file)
	}

	w.mods = cf
	sort.Strings(changedFiles)
	return changedFiles, nil
}

func (w *Watch) Watch(ctx context.Context, ch chan []string) error {
	return w.WatchCb(ctx, func(files []string) error {
		ch <- files
		return nil
	})
}

func (w *Watch) WatchCb(ctx context.Context, fn func([]string) error) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(w.pollTimeout):
			changedFiles, err := w.Compare()
			if err != nil {
				return err
			}
			if len(changedFiles) > 0 {
				fn(changedFiles)
			}
		}
	}
}

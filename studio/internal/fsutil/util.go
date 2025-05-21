package fsutil

import (
	"fmt"
	"hash/crc32"
	"io"
	"io/fs"
)

func CopyFile(from, to fs.FS, path string) error {
	// Open source file
	srcFile, err := from.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open source file %s: %w", path, err)
	}
	defer srcFile.Close()

	data, err := io.ReadAll(srcFile)
	if err != nil {
		return fmt.Errorf("failed to read source file %s: %w", path, err)
	}

	err = WriteFileSimple(to, path, data)
	if err != nil {
		return err
	}

	return nil
}

func Copy(from, to fs.FS, px string) error {
	return fs.WalkDir(from, px, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		// Create directories
		if d.IsDir() {
			return MkdirAll(to, path)
		}
		// Copy files
		return CopyFile(from, to, path)
	})
}

func CRC32(fs fs.FS, path string) (uint32, error) {
	file, err := fs.Open(path)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	hash := crc32.NewIEEE()
	if _, err := io.Copy(hash, file); err != nil {
		return 0, err
	}

	return hash.Sum32(), nil
}

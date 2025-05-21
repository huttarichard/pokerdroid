// Package fsutil provides utilities for working with the filesystem.
// Majority of these should not be used for production code.
// Usecase for this is rather tooling or development.
package fsutil

import (
	"errors"
	"go/build"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
)

// RelativeToRoot says how much above this package is from the root.
const RelativeToRoot = "../../.."

// ResolveProjectPath resolves the path to be
// either absolute or relative to the working path.
func ResolveProjectPath(pth string) (string, error) {
	if filepath.IsAbs(pth) {
		_, err := os.Stat(pth)
		return pth, err
	}
	wp, err := GetProjectPath()
	if err != nil {
		return pth, err
	}
	pth = filepath.Join(wp, pth)
	_, err = os.Stat(pth)
	return pth, err
}

// MustResolveProjectPath resolves the path to be
func MustResolveProjectPath(pth string) string {
	res, err := ResolveProjectPath(pth)
	if err != nil {
		panic(err)
	}
	return res
}

// GetProjectPath returns the working path.
// That is either runtime path or build path.
func GetProjectPath() (string, error) {
	path, err := GetBuildPath()
	if err == nil {
		return path, nil
	}
	path, err = GetCallerPath()
	if err == nil {
		return path, nil
	}
	return "", errors.New("could not get working path")
}

// GetCallerPath returns the runtime path.
func GetCallerPath() (string, error) {
	_, f, _, ok := runtime.Caller(0)
	if !ok {
		return "", errors.New("could not get runtime path")
	}
	pth := filepath.Join(filepath.Dir(f), RelativeToRoot)
	pth = filepath.Clean(pth)
	return verifymod(pth)
}

// GetBuildPath returns the build path.
func GetBuildPath() (string, error) {
	pkg, err := build.Default.Import(Pkg(), "", build.FindOnly)
	if err != nil {
		return "", err
	}
	return verifymod(pkg.Dir)
}

// Pkg returns the root package path.
func Pkg() string {
	type tp struct{}
	pth := reflect.TypeOf(tp{}).PkgPath()
	pth = filepath.Clean(filepath.Join(pth, RelativeToRoot))
	return pth
}

// verifymod verifies if the go.mod file exists.
func verifymod(pth string) (string, error) {
	_, err := os.Stat(filepath.Join(pth, "go.mod"))
	if err != nil {
		return "", errors.New("could not find go.mod")
	}
	return pth, nil
}

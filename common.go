package load

import (
	"os"
	"path/filepath"
)

const (
	// StarExt starlark source code files default extension.
	StarExt = ".star"
	// ModuleIndexFile if the module is a directory it load the index file.
	ModuleIndexFile = "main" + StarExt
)

// SourceCodeFinder helper to find a starlark file in the a set directories.
type SourceCodeFinder struct {
	path []string
}

// NewSourceCodeFinder returns a new ModuleFinder.
func NewSourceCodeFinder(path ...string) *SourceCodeFinder {
	return &SourceCodeFinder{
		path: path,
	}
}

// Find find a given module in the path, it tries different variants.
func (l *SourceCodeFinder) Find(module string) (string, error) {
	possible := []string{
		module,
		module + StarExt,
		filepath.Join(module, ModuleIndexFile),
	}

	for _, path := range possible {
		if fullpath, ok := l.exists(path); ok {
			return fullpath, nil
		}
	}

	return "", ErrModuleNotFound
}
func (l *SourceCodeFinder) exists(filename string) (string, bool) {
	for _, path := range l.path {
		fullpath := filepath.Join(path, filename)
		if l.doExists(fullpath) {
			return fullpath, true
		}
	}

	return "", false
}

func (l *SourceCodeFinder) doExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}

	return !info.IsDir()
}

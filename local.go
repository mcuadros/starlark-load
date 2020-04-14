package load

import (
	"os"
	"path/filepath"
	"strings"
)

// LocalStrategy defines a load method based on local filesystem.
type LocalStrategy struct {
	path   []string
	finder *SourceCodeFinder
}

// NewLocalMethod returns a new LocalStrategy, this strategy search of modules
// on the given path. Path may contain several directories separated by a colon.
func NewLocalMethod(path string) *LocalStrategy {
	dir, _ := os.Getwd()
	paths := append([]string{dir}, strings.Split(path, ":")...)

	return &LocalStrategy{
		finder: NewSourceCodeFinder(paths...),
	}
}

// Resolve resolves a module name into a Module searching the given module in
// the path.
func (l *LocalStrategy) Resolve(module string) (Module, error) {
	fullpath, err := l.finder.Find(module)
	if err != nil {
		return nil, err
	}

	return NewStarlarkModule(l.cleanName(module), fullpath, nil), nil
}

func (l *LocalStrategy) cleanName(module string) string {
	if filepath.Base(module) == ModuleIndexFile {
		return filepath.Clean(module[:len(module)-len(ModuleIndexFile)])
	}

	ext := filepath.Ext(module)
	module = module[:len(module)-len(ext)]

	return module
}

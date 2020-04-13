package load

import (
	"errors"

	"go.starlark.net/starlark"
)

var (
	ErrModuleNotFound = errors.New("module not found")
	ErrCycleLoad      = errors.New("cycle in load graph")
)

const (
	ModInfoKey = "modinfo"
)

type Loader struct {
	predeclared starlark.StringDict
	modules     map[string]*moduleCache
	Methods     []LoadMethod
}

func NewLoader(predeclared starlark.StringDict) *Loader {
	return &Loader{
		predeclared: predeclared,
		modules:     make(map[string]*moduleCache),
	}
}

func (l *Loader) Load(t *starlark.Thread, module string) (starlark.StringDict, error) {
	m, err := l.resolve(module)
	if err != nil {
		return nil, err
	}

	cache, ok := l.modules[m.Name()]
	if cache == nil {
		if ok {
			// request for package whose loading is in progress
			return nil, ErrCycleLoad
		}

		// Add a placeholder to indicate "load in progress".
		l.modules[m.Name()] = nil

		globals, err := m.Load(t, l.predeclared)
		cache = &moduleCache{
			Module:  m,
			Globals: globals,
			LoadErr: err,
		}

		l.modules[m.Name()] = cache
	}

	return cache.Globals, cache.LoadErr
}

func (l *Loader) resolve(module string) (m Module, err error) {
	for _, method := range l.Methods {
		m, err = method.Resolve(module)
		if err == nil {
			break
		}

		if err != nil && err != ErrModuleNotFound {
			return nil, err
		}
	}

	if m == nil {
		err = ErrModuleNotFound
	}

	return
}

type moduleCache struct {
	Module  Module
	Globals starlark.StringDict
	LoadErr error
}

type Module interface {
	Name() string
	Path() string
	Load(t *starlark.Thread, predeclared starlark.StringDict) (starlark.StringDict, error)
}

type LoadMethod interface {
	Resolve(module string) (Module, error)
}

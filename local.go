package load

import (
	"os"
	"path/filepath"
	"strings"

	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
)

const (
	StarExt         = ".star"
	ModuleIndexFile = "main" + StarExt
)

type LocalLoadMethod struct {
	path []string
}

func NewLocalMethod(path string) *LocalLoadMethod {
	dir, _ := os.Getwd()

	return &LocalLoadMethod{
		path: append([]string{dir}, strings.Split(path, ":")...),
	}
}

func (l *LocalLoadMethod) Resolve(module string) (Module, error) {
	fullpath, err := l.findModule(module)
	if err != nil {
		return nil, err
	}

	return &LocalModule{
		name: l.cleanName(module),
		path: fullpath,
	}, nil
}

func (l *LocalLoadMethod) cleanName(module string) string {
	if filepath.Base(module) == ModuleIndexFile {
		return filepath.Clean(module[:len(module)-len(ModuleIndexFile)])
	}

	ext := filepath.Ext(module)
	module = module[:len(module)-len(ext)]

	return module
}

func (l *LocalLoadMethod) findModule(module string) (string, error) {
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
func (l *LocalLoadMethod) exists(filename string) (string, bool) {
	for _, path := range l.path {
		fullpath := filepath.Join(path, filename)
		if l.doExists(fullpath) {
			return fullpath, true
		}
	}

	return "", false
}

func (l *LocalLoadMethod) doExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}

	return !info.IsDir()
}

type LocalModule struct {
	name string
	path string
}

func (m *LocalModule) Name() string {
	return m.name
}

func (m *LocalModule) Path() string {
	return m.path
}

func (m *LocalModule) Load(t *starlark.Thread, predeclared starlark.StringDict) (starlark.StringDict, error) {
	thread := &starlark.Thread{Name: "exec " + m.name, Load: t.Load}

	globals, err := starlark.ExecFile(thread, m.path, nil, predeclared)
	globals[ModInfoKey] = m.modInfo()

	return globals, err
}

func (m *LocalModule) modInfo() starlark.Value {
	return starlarkstruct.FromStringDict(starlarkstruct.Default, starlark.StringDict{
		"name": starlark.String(m.name),
		"path": starlark.String(m.path),
	})
}

package load

import (
	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
)

const (
	// ModInfoKey key of the struct containing meta information about the loaded
	// module. ModInfo is injected on every loaded module.
	ModInfoKey = "modinfo"
)

// Module defines a starlark module.
type Module interface {
	Name() string
	Load(t *starlark.Thread, predeclared starlark.StringDict) (starlark.StringDict, error)
}

// StarlarkModule defines a Module, based on starlark code.
type StarlarkModule struct {
	name    string
	path    string
	modInfo map[string]string
}

// NewStarlarkModule returns a new StarlarkModule.
func NewStarlarkModule(name, path string, modInfo map[string]string) *StarlarkModule {
	return &StarlarkModule{
		name:    name,
		path:    path,
		modInfo: modInfo,
	}
}

// Name honors the Module interface.
func (m *StarlarkModule) Name() string {
	return m.name
}

// Load honors the Module interface.
func (m *StarlarkModule) Load(t *starlark.Thread, predeclared starlark.StringDict) (starlark.StringDict, error) {
	thread := &starlark.Thread{Name: "exec " + m.name, Load: t.Load}

	globals, err := starlark.ExecFile(thread, m.path, nil, predeclared)
	globals[ModInfoKey] = m.calculateModInfo()

	return globals, err
}

func (m *StarlarkModule) calculateModInfo() starlark.Value {
	dict := starlark.StringDict{
		"name": starlark.String(m.name),
		"path": starlark.String(m.path),
	}

	for k, value := range m.modInfo {
		dict[k] = starlark.String(value)
	}

	return starlarkstruct.FromStringDict(starlarkstruct.Default, dict)
}

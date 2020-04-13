package load

import (
	"errors"
	"math/rand"
	"testing"
	"time"

	"github.com/oklog/ulid"
	"github.com/stretchr/testify/require"
	"go.starlark.net/resolve"
	"go.starlark.net/starlark"
	"go.starlark.net/starlarktest"
)

func TestLoaderLoad_Single(t *testing.T) {
	_, err := doTestFile(t, "testdata/single.star")
	require.NoError(t, err)
}

func TestLoaderLoad_NotFound(t *testing.T) {
	_, err := doTestFile(t, "testdata/not_found.star")
	require.NotNil(t, err)
	require.Equal(t, errors.Is(err, ErrModuleNotFound), true)
}

func TestLoaderLoad_Cyclic(t *testing.T) {
	_, err := doTestFile(t, "testdata/cyclic.star")

	require.NotNil(t, err)
	require.Equal(t, errors.Is(err, ErrCycleLoad), true)
}

func doTestFile(t *testing.T, filename string) (starlark.StringDict, error) {
	resolve.AllowFloat = true
	resolve.AllowGlobalReassign = true
	resolve.AllowLambda = true

	assert, err := starlarktest.LoadAssertModule()
	require.NoError(t, err)

	predeclared := starlark.StringDict{
		"rand":   BuiltinRand(),
		"assert": assert["assert"],
	}

	loader := NewLoader(predeclared)
	loader.Methods = append(loader.Methods, NewLocalMethod("fixtures"))
	thread := &starlark.Thread{
		Load: loader.Load,
	}

	starlarktest.SetReporter(thread, t)
	return starlark.ExecFile(thread, filename, nil, predeclared)
}

func BuiltinRand() starlark.Value {
	return starlark.NewBuiltin("rand", func(t *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		entropy := ulid.Monotonic(rand.New(rand.NewSource(time.Now().UnixNano())), 0)
		ulid := ulid.MustNew(ulid.Timestamp(time.Now()), entropy)

		return starlark.String(ulid.String()), nil
	})
}

package load

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRemoteModuleLoad(t *testing.T) {
	_, err := doTestFile(t, "testdata/remote.star")
	require.NoError(t, err)
}

func TestNewRemoteModuleName(t *testing.T) {
	testCases := []struct {
		raw  string
		want RemoteModuleName
		err  error
	}{{
		"foo/qux.v2", RemoteModuleName{
			Server:       DefaultGitServer,
			Organization: "foo",
			Repository:   "qux",
			Version:      "2",
		}, nil,
	}, {
		"github.com/foo/qux.v2", RemoteModuleName{
			Server:       "github.com",
			Organization: "foo",
			Repository:   "qux",
			Version:      "2",
		}, nil,
	}, {
		"github.com/foo/qux.v2-aa221cde7f29b66b3cc3ddc2c8716a40d7378fb3", RemoteModuleName{
			Server:       "github.com",
			Organization: "foo",
			Repository:   "qux",
			Commit:       "aa221cde7f29b66b3cc3ddc2c8716a40d7378fb3",
			Version:      "2",
		}, nil,
	}, {
		"github.com/foo/qux.v2-aa221c", RemoteModuleName{
			Server:       "github.com",
			Organization: "foo",
			Repository:   "qux",
			Commit:       "aa221c",
			Version:      "2",
		}, nil,
	}, {
		"github.com/foo/qux.v2-aa221c/path/to/module", RemoteModuleName{
			Server:       "github.com",
			Organization: "foo",
			Repository:   "qux",
			Commit:       "aa221c",
			Version:      "2",
			Path:         "/path/to/module",
		}, nil,
	}, {
		"github.com/foo/qux.v2/path/to/module", RemoteModuleName{
			Server:       "github.com",
			Organization: "foo",
			Repository:   "qux",
			Version:      "2",
			Path:         "/path/to/module",
		}, nil,
	},
		{"github.com/foo/qux/path/to/module", RemoteModuleName{}, ErrMalformedModuleName},
		{"github.com/foo/qux.vA", RemoteModuleName{}, ErrMalformedModuleName},
		{"github.com/foo/qux.v4-foo", RemoteModuleName{}, ErrMalformedModuleName},
		{"github.com/foo/qux.v4/ff[", RemoteModuleName{}, ErrMalformedModuleName},
		{"github.com/foo/quxv4", RemoteModuleName{}, ErrMalformedModuleName},
	}

	for _, tc := range testCases {
		got, err := NewRemoteModuleName(tc.raw)
		require.Equal(t, err, tc.err)
		require.Equal(t, tc.want, got)
	}
}

package load

import (
	"testing"

	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/stretchr/testify/require"
)

func TestNewVersions(t *testing.T) {
	refs := make(memory.ReferenceStorage, 0)
	refs.SetReference(plumbing.NewHashReference("refs/heads/master", plumbing.NewHash("")))
	refs.SetReference(plumbing.NewHashReference("refs/tags/v1.0.0", plumbing.NewHash("")))
	refs.SetReference(plumbing.NewHashReference("refs/tags/1.1.2", plumbing.NewHash("")))
	refs.SetReference(plumbing.NewHashReference("refs/tags/1.1.3", plumbing.NewHash("")))
	refs.SetReference(plumbing.NewHashReference("refs/tags/v1.0.3", plumbing.NewHash("")))
	refs.SetReference(plumbing.NewHashReference("refs/tags/v2.0.3", plumbing.NewHash("")))
	refs.SetReference(plumbing.NewHashReference("refs/tags/v4.0.0-rc1", plumbing.NewHash("")))

	v := NewVersions(refs)
	require.Equal(t, v.Match("0").Name().String(), "refs/heads/master")
	require.Equal(t, v.Match("v0").Name().String(), "refs/heads/master")
	require.Equal(t, v.Match("v1.1").Name().String(), "refs/tags/1.1.3")
	require.Equal(t, v.Match("1.1").Name().String(), "refs/tags/1.1.3")
	require.Equal(t, v.Match("1.1.2").Name().String(), "refs/tags/1.1.2")
	require.Equal(t, v.Match("2").Name().String(), "refs/tags/v2.0.3")
	require.Equal(t, v.Match("4").Name().String(), "refs/tags/v4.0.0-rc1")
	require.Equal(t, v.Match("master").Name().String(), "refs/heads/master")
	require.Nil(t, v.Match("foo"))

	refs.SetReference(plumbing.NewHashReference("refs/tags/v0.0.0", plumbing.NewHash("")))

	v = NewVersions(refs)
	require.Equal(t, v.Match("0").Name().String(), "refs/tags/v0.0.0")
	require.Equal(t, v.Match("v0").Name().String(), "refs/tags/v0.0.0")
}

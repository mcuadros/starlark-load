package load

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLocalModuleLoad(t *testing.T) {
	_, err := doTestFile(t, "testdata/local.star")
	require.NoError(t, err)
}

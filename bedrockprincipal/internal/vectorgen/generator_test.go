package vectorgen

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCheckedInCoreVectorsRegenerateExactly(t *testing.T) {
	path := "../../testdata/v2/core-vectors.json"
	checked, err := os.ReadFile(path)
	require.NoError(t, err)
	regenerated, err := Regenerate(checked)
	require.NoError(t, err)
	require.Equal(t, string(checked), string(regenerated))
}

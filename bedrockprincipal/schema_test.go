package bedrockprincipal_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestV2SchemasAreClosedAndFreezeNonceJTI(t *testing.T) {
	payload := loadJSONObject(t, "schema/v2.schema.json")
	require.Equal(t, false, payload["additionalProperties"])
	properties := payload["properties"].(map[string]any)
	for _, name := range []string{"connect_session_nonce", "jti"} {
		property := properties[name].(map[string]any)
		require.Equal(t, float64(22), property["minLength"], name)
		require.Equal(t, float64(22), property["maxLength"], name)
		require.Equal(t, "^[A-Za-z0-9_-]{22}$", property["pattern"], name)
	}
	require.Equal(t, float64(1), properties["policy_revision"].(map[string]any)["minimum"])

	metadata := loadJSONObject(t, "schema/metadata-v2.schema.json")
	require.Equal(t, false, metadata["additionalProperties"])
	metadataProperties := metadata["properties"].(map[string]any)
	require.Equal(t, float64(16), metadataProperties["keys"].(map[string]any)["maxItems"])
	require.Equal(t, float64(300), metadataProperties["cache_max_age_seconds"].(map[string]any)["maximum"])
}

func TestCoreVectorsHaveOnlyLiteralVerificationInputsAndOutcomes(t *testing.T) {
	raw, err := os.ReadFile("testdata/v2/core-vectors.json")
	require.NoError(t, err)

	var vectors []map[string]json.RawMessage
	require.NoError(t, json.Unmarshal(raw, &vectors))
	require.NotEmpty(t, vectors)
	allowed := map[string]bool{
		"name": true, "compact_jws": true, "trusted_context": true,
		"verification_time_unix": true, "expected_error": true,
		"expected_principal": true,
	}
	for _, vector := range vectors {
		for key := range vector {
			require.True(t, allowed[key], "unexpected vector member %q", key)
		}
	}
}

func loadJSONObject(t *testing.T, path string) map[string]any {
	t.Helper()
	raw, err := os.ReadFile(path)
	require.NoError(t, err)
	var value map[string]any
	require.NoError(t, json.Unmarshal(raw, &value))
	return value
}

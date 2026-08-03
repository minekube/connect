package bedrockprincipal_test

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCandidateManifestPinsExactLocalArtifactsAndNoRelease(t *testing.T) {
	raw, err := os.ReadFile("testdata/v2/connect-sdk-candidate.json")
	require.NoError(t, err)
	var manifest struct {
		SchemaVersion        int    `json:"schema_version"`
		ApprovalState        string `json:"approval_state"`
		Commit               any    `json:"commit"`
		SourceTree           any    `json:"source_tree"`
		Tag                  any    `json:"tag"`
		CorpusTree           any    `json:"corpus_tree"`
		ArtifactSHA256       any    `json:"artifact_sha256"`
		GoModule             string `json:"go_module"`
		CorpusDigestSHA256   string `json:"corpus_digest_sha256"`
		SchemaSHA256         string `json:"schema_sha256"`
		MetadataSchemaSHA256 string `json:"metadata_schema_sha256"`
	}
	require.NoError(t, json.Unmarshal(raw, &manifest))
	require.Equal(t, 1, manifest.SchemaVersion)
	require.Equal(t, "candidate", manifest.ApprovalState)
	require.Nil(t, manifest.Commit)
	require.Nil(t, manifest.SourceTree)
	require.Nil(t, manifest.Tag)
	require.Nil(t, manifest.CorpusTree)
	require.Nil(t, manifest.ArtifactSHA256)
	require.Equal(t, "go.minekube.com/connect", manifest.GoModule)
	require.Equal(t, fileSHA256(t, "testdata/v2/core-vectors.json"), manifest.CorpusDigestSHA256)
	require.Equal(t, fileSHA256(t, "schema/v2.schema.json"), manifest.SchemaSHA256)
	require.Equal(t, fileSHA256(t, "schema/metadata-v2.schema.json"), manifest.MetadataSchemaSHA256)
}

func fileSHA256(t *testing.T, path string) string {
	t.Helper()
	raw, err := os.ReadFile(path)
	require.NoError(t, err)
	digest := sha256.Sum256(raw)
	return hex.EncodeToString(digest[:])
}

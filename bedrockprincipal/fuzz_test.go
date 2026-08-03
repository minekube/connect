package bedrockprincipal

import (
	"context"
	"encoding/json"
	"os"
	"testing"
)

func FuzzStrictEnvelope(f *testing.F) {
	raw, err := os.ReadFile("testdata/v2/core-vectors.json")
	if err != nil {
		f.Fatal(err)
	}
	var vectors []struct {
		CompactJWS string `json:"compact_jws"`
	}
	if err := json.Unmarshal(raw, &vectors); err != nil {
		f.Fatal(err)
	}
	f.Add(vectors[0].CompactJWS)
	f.Add(".")
	f.Add("a.b.c")
	f.Fuzz(func(t *testing.T, compact string) {
		envelope, err := NewSignedPrincipalEnvelope([]byte(compact))
		if err != nil {
			return
		}
		_, _ = verifierAt(testNowUnix).VerifyAndConsume(context.Background(), envelope, trustedContext())
	})
}

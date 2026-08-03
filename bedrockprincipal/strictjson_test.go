package bedrockprincipal

import (
	"context"
	"crypto/ed25519"
	"encoding/base64"
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStrictEnvelopeRejectsDuplicateUnknownTrailingAndArrayJSON(t *testing.T) {
	validHeader := `{"alg":"EdDSA","typ":"connect-bedrock-principal+jws;v=2","kid":"connect-v2-test"}`
	validBody := mustPayloadJSON(t, validPayload())
	missingVersion := withoutPayloadField(t, validBody, "version")
	missingIssuer := withoutPayloadField(t, validBody, "issuer")
	cases := []struct{ name, header, body string }{
		{"duplicate-header", `{"alg":"EdDSA","alg":"EdDSA","typ":"connect-bedrock-principal+jws;v=2","kid":"connect-v2-test"}`, validBody},
		{"unknown-header", `{"alg":"EdDSA","typ":"connect-bedrock-principal+jws;v=2","kid":"connect-v2-test","extra":true}`, validBody},
		{"duplicate-payload", validHeader, strings.Replace(validBody, `"version":2`, `"version":2,"version":2`, 1)},
		{"unknown-payload", validHeader, strings.TrimSuffix(validBody, "}") + `,"extra":true}`},
		{"trailing-payload", validHeader, validBody + ` {}`},
		{"header-array", `[]`, validBody},
		{"payload-array", validHeader, `[]`},
		{"nul-string", validHeader, strings.Replace(validBody, "BedrockOne", `Bedrock\u0000One`, 1)},
		{"unpaired-surrogate", validHeader, strings.Replace(validBody, "BedrockOne", `Bedrock\ud800One`, 1)},
		{"missing-version", validHeader, missingVersion},
		{"missing-issuer", validHeader, missingIssuer},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := verifierAt(testNowUnix).VerifyAndConsume(context.Background(), mustEnvelope(t, signRaw(tc.header, tc.body)), trustedContext())
			require.ErrorIs(t, err, Malformed)
		})
	}
}

func signRaw(header, payload string) string {
	protected := base64.RawURLEncoding.EncodeToString([]byte(header))
	body := base64.RawURLEncoding.EncodeToString([]byte(payload))
	input := protected + "." + body
	signature := ed25519.Sign(testPrivateKey, []byte(input))
	return input + "." + base64.RawURLEncoding.EncodeToString(signature)
}

func mustPayloadJSON(t *testing.T, payload map[string]any) string {
	t.Helper()
	compact := signPayload(t, payload)
	parts := strings.Split(compact, ".")
	decoded, err := base64.RawURLEncoding.DecodeString(parts[1])
	require.NoError(t, err)
	return string(decoded)
}

func withoutPayloadField(t *testing.T, body, field string) string {
	t.Helper()
	var payload map[string]any
	require.NoError(t, json.Unmarshal([]byte(body), &payload))
	delete(payload, field)
	encoded, err := json.Marshal(payload)
	require.NoError(t, err)
	return string(encoded)
}

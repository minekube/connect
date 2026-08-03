package vectorgen

import (
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
)

const fixedNow = int64(1785492000)

type vector struct {
	Name                 string          `json:"name"`
	CompactJWS           string          `json:"compact_jws"`
	TrustedContext       json.RawMessage `json:"trusted_context"`
	VerificationTimeUnix int64           `json:"verification_time_unix"`
	ExpectedError        string          `json:"expected_error"`
	ExpectedPrincipal    json.RawMessage `json:"expected_principal"`
}

// Regenerate signs only the predefined payload associated with each reviewed
// vector name. Trusted context and expected outcomes remain checked-in literals.
func Regenerate(checked []byte) ([]byte, error) {
	var vectors []vector
	if err := json.Unmarshal(checked, &vectors); err != nil {
		return nil, err
	}
	if len(vectors) != len(definitions()) {
		return nil, fmt.Errorf("expected %d vectors, got %d", len(definitions()), len(vectors))
	}
	seen := make(map[string]bool, len(vectors))
	for i := range vectors {
		if seen[vectors[i].Name] {
			return nil, fmt.Errorf("duplicate vector %q", vectors[i].Name)
		}
		seen[vectors[i].Name] = true
		definition, ok := definitions()[vectors[i].Name]
		if !ok {
			return nil, fmt.Errorf("unknown vector %q", vectors[i].Name)
		}
		vectors[i].CompactJWS = sign(definition.payload, definition.tamperSignature)
	}
	for name := range definitions() {
		if !seen[name] {
			return nil, fmt.Errorf("missing vector %q", name)
		}
	}
	regenerated, err := json.MarshalIndent(vectors, "", "  ")
	if err != nil {
		return nil, err
	}
	return append(regenerated, '\n'), nil
}

type definition struct {
	payload         map[string]any
	tamperSignature bool
}

func definitions() map[string]definition {
	validUnlinked := payload(1)
	validLinked := payload(2)
	validLinked["subject_kind"] = "bedrock_linked_java"
	validLinked["verification_method"] = "minecraft_full_jwks+client_jwt+ecdh_v1"
	validLinked["linked_java"] = map[string]any{
		"uuid": "123e4567-e89b-12d3-a456-426614174000",
		"name": "JavaOne",
		"provenance": map[string]any{
			"provider": "moxy_account_link_v1", "record_id": "record-test-1",
			"revision": 3, "verified_at": fixedNow - 10,
		},
	}
	malformedJTI := payload(3)
	malformedJTI["jti"] = "AAAAAAAAAAAAAAAAAAAAAB"
	policy := payload(4)
	lifetime := payload(5)
	lifetime["exp"] = fixedNow + 31
	tampered := payload(6)
	return map[string]definition{
		"valid-unlinked":           {payload: validUnlinked},
		"valid-linked":             {payload: validLinked},
		"malformed-jti-tail-bits":  {payload: malformedJTI},
		"policy-revision-mismatch": {payload: policy},
		"lifetime-thirty-one":      {payload: lifetime},
		"tampered-signature":       {payload: tampered, tamperSignature: true},
	}
}

func payload(jtiByte byte) map[string]any {
	return map[string]any{
		"version":                 2,
		"issuer":                  "minekube-connect-test",
		"trust_domain":            "urn:minekube:connect:test:corpus-v2",
		"audience":                "urn:minekube:connect:test:bedrock-principal:v2",
		"subject_kind":            "bedrock_xuid",
		"canonical_xuid":          "1",
		"canonical_unlinked_uuid": "00000000-0000-0000-0000-000000000001",
		"bedrock_display_name":    "BedrockOne",
		"endpoint_id":             "endpoint-test",
		"organization_id":         "organization-test",
		"connect_session_id":      "session-test",
		"connect_session_nonce":   "AAAAAAAAAAAAAAAAAAAAAA",
		"policy_revision":         7,
		"source_protocol":         "bedrock",
		"source_protocol_version": 776,
		"iat":                     fixedNow,
		"nbf":                     fixedNow,
		"exp":                     fixedNow + 30,
		"jti":                     base64.RawURLEncoding.EncodeToString(repeated(jtiByte, 16)),
		"verification_method":     "minecraft_legacy_chain+client_jwt+ecdh_v1",
	}
}

func sign(payload map[string]any, tamper bool) string {
	header := []byte(`{"alg":"EdDSA","typ":"connect-bedrock-principal+jws;v=2","kid":"connect-v2-test"}`)
	body, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}
	protected := base64.RawURLEncoding.EncodeToString(header)
	encodedPayload := base64.RawURLEncoding.EncodeToString(body)
	signingInput := protected + "." + encodedPayload
	seed := sha256.Sum256([]byte("connect-bedrock-principal-v2-test-only"))
	privateKey := ed25519.NewKeyFromSeed(seed[:])
	signature := ed25519.Sign(privateKey, []byte(signingInput))
	if tamper {
		signature[0] ^= 1
	}
	return signingInput + "." + base64.RawURLEncoding.EncodeToString(signature)
}

func repeated(value byte, count int) []byte {
	result := make([]byte, count)
	for i := range result {
		result[i] = value
	}
	return result
}

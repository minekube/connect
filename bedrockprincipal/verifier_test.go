package bedrockprincipal

import (
	"context"
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

const (
	testNowUnix = int64(1785492000)
	testNonce   = "AAAAAAAAAAAAAAAAAAAAAA"
	testJTI     = "AQEBAQEBAQEBAQEBAQEBAQ"
)

var (
	testSeed       = sha256.Sum256([]byte("connect-bedrock-principal-v2-test-only"))
	testPrivateKey = ed25519.NewKeyFromSeed(testSeed[:])
	testPublicKey  = testPrivateKey.Public().(ed25519.PublicKey)
)

func TestPrincipalCannotBeConstructedOutsideVerifier(t *testing.T) {
	cmd := exec.Command("go", "test", "./testdata/compilefail")
	out, err := cmd.CombinedOutput()
	require.Error(t, err)
	require.Contains(t, string(out), "verifiedPrincipal")
}

func TestVerifyReturnsTypedUnlinkedPrincipal(t *testing.T) {
	payload := validPayload()
	principal, err := verifierAt(testNowUnix).VerifyAndConsume(context.Background(), mustEnvelope(t, signPayload(t, payload)), trustedContext())
	require.NoError(t, err)
	require.Equal(t, BedrockXUID, principal.SubjectKind())
	require.Equal(t, CanonicalXUID("1"), principal.XUID())
	require.Equal(t, uuid.MustParse("00000000-0000-0000-0000-000000000001"), principal.CanonicalUnlinkedUUID())
	_, linked := principal.LinkedJava()
	require.False(t, linked)
	require.Equal(t, GameProfile{UUID: uuid.MustParse("00000000-0000-0000-0000-000000000001"), Name: "BedrockOne"}, principal.EffectiveGameProfile())
	require.Equal(t, int64(7), principal.Bindings().PolicyRevision)
}

func TestVerifierRequiresCanonicalUUIDSpelling(t *testing.T) {
	cases := []struct {
		name   string
		mutate func(map[string]any)
		want   PrincipalError
	}{
		{name: "unlinked-braced", mutate: func(payload map[string]any) {
			payload["canonical_unlinked_uuid"] = "{00000000-0000-0000-0000-000000000001}"
		}, want: IdentityInvalid},
		{name: "linked-uppercase", mutate: func(payload map[string]any) {
			payload = linkedPayload(payload)
			linked := payload["linked_java"].(map[string]any)
			linked["uuid"] = "123E4567-E89B-12D3-A456-426614174000"
		}, want: LinkInvalid},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			payload := validPayload()
			tc.mutate(payload)
			_, err := verifierAt(testNowUnix).VerifyAndConsume(context.Background(), mustEnvelope(t, signPayload(t, payload)), trustedContext())
			require.ErrorIs(t, err, tc.want)
		})
	}
}

func TestVerifierEnforcesLinkedJavaSchemaBounds(t *testing.T) {
	cases := []struct {
		name   string
		mutate func(map[string]any)
	}{
		{name: "name-character-set", mutate: func(payload map[string]any) {
			payload["linked_java"].(map[string]any)["name"] = "Java One"
		}},
		{name: "name-length", mutate: func(payload map[string]any) {
			payload["linked_java"].(map[string]any)["name"] = strings.Repeat("J", 17)
		}},
		{name: "record-id-length", mutate: func(payload map[string]any) {
			payload["linked_java"].(map[string]any)["provenance"].(map[string]any)["record_id"] = strings.Repeat("r", 129)
		}},
		{name: "verified-at-before-epoch", mutate: func(payload map[string]any) {
			payload["linked_java"].(map[string]any)["provenance"].(map[string]any)["verified_at"] = int64(-1)
		}},
		{name: "verified-at-after-maximum", mutate: func(payload map[string]any) {
			payload["linked_java"].(map[string]any)["provenance"].(map[string]any)["verified_at"] = int64(253402300800)
		}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			payload := linkedPayload(validPayload())
			tc.mutate(payload)
			_, err := verifierAt(testNowUnix).VerifyAndConsume(context.Background(), mustEnvelope(t, signPayload(t, payload)), trustedContext())
			require.ErrorIs(t, err, LinkInvalid)
		})
	}
}

func TestVerifierEnforcesFrozenBedrockBindingSchema(t *testing.T) {
	cases := []struct {
		name   string
		mutate func(map[string]any, *TrustedProposalContext)
	}{
		{name: "source-protocol", mutate: func(payload map[string]any, context *TrustedProposalContext) {
			payload["source_protocol"] = "java"
			context.SourceProtocol = "java"
		}},
		{name: "source-protocol-version", mutate: func(payload map[string]any, context *TrustedProposalContext) {
			payload["source_protocol_version"] = int32(0)
			context.SourceProtocolVersion = 0
		}},
		{name: "endpoint-id-length", mutate: func(payload map[string]any, context *TrustedProposalContext) {
			value := strings.Repeat("e", 129)
			payload["endpoint_id"] = value
			context.EndpointID = value
		}},
		{name: "organization-id-length", mutate: func(payload map[string]any, context *TrustedProposalContext) {
			value := strings.Repeat("o", 129)
			payload["organization_id"] = value
			context.OrganizationID = value
		}},
		{name: "connect-session-id-length", mutate: func(payload map[string]any, context *TrustedProposalContext) {
			value := strings.Repeat("s", 129)
			payload["connect_session_id"] = value
			context.ConnectSessionID = value
		}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			payload := validPayload()
			expectedContext := trustedContext()
			tc.mutate(payload, &expectedContext)
			_, err := verifierAt(testNowUnix).VerifyAndConsume(context.Background(), mustEnvelope(t, signPayload(t, payload)), expectedContext)
			require.ErrorIs(t, err, BindingMismatch)
		})
	}
}

func TestNonceAndJTIRequireCanonicalSixteenByteBase64URL(t *testing.T) {
	encode := base64.RawURLEncoding.EncodeToString
	cases := []struct {
		name  string
		field string
		value string
	}{
		{name: "nonce-15-decoded-bytes", field: "connect_session_nonce", value: encode(make([]byte, 15))},
		{name: "nonce-17-decoded-bytes", field: "connect_session_nonce", value: encode(make([]byte, 17))},
		{name: "jti-21-characters", field: "jti", value: strings.Repeat("A", 21)},
		{name: "jti-23-characters", field: "jti", value: strings.Repeat("A", 23)},
		{name: "nonce-padding", field: "connect_session_nonce", value: testNonce + "=="},
		{name: "jti-standard-alphabet", field: "jti", value: strings.Repeat("A", 21) + "+"},
		{name: "jti-nonzero-tail-bits", field: "jti", value: strings.Repeat("A", 21) + "B"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			payload := validPayload()
			payload[tc.field] = tc.value
			_, err := verifierAt(testNowUnix).VerifyAndConsume(context.Background(), mustEnvelope(t, signPayload(t, payload)), trustedContext())
			require.ErrorIs(t, err, Malformed)
		})
	}
}

func TestNonceClaimMustMatchTrustedContextBytes(t *testing.T) {
	ctx := trustedContext()
	ctx.ConnectSessionNonce[15] = 1
	_, err := verifierAt(testNowUnix).VerifyAndConsume(context.Background(), mustEnvelope(t, signPayload(t, validPayload())), ctx)
	require.ErrorIs(t, err, BindingMismatch)
}

func TestLifetimeAndSkewEqualityBoundaries(t *testing.T) {
	cases := []struct {
		name          string
		iat, nbf, exp int64
		want          PrincipalError
	}{
		{name: "now-plus-five-equals-nbf", iat: testNowUnix + 5, nbf: testNowUnix + 5, exp: testNowUnix + 30},
		{name: "now-minus-five-equals-exp", iat: testNowUnix - 35, nbf: testNowUnix - 35, exp: testNowUnix - 5},
		{name: "iat-equals-now-plus-five", iat: testNowUnix + 5, nbf: testNowUnix, exp: testNowUnix + 30},
		{name: "lifetime-equals-thirty", iat: testNowUnix, nbf: testNowUnix, exp: testNowUnix + 30},
		{name: "nbf-one-second-outside", iat: testNowUnix + 6, nbf: testNowUnix + 6, exp: testNowUnix + 30, want: TimeInvalid},
		{name: "exp-one-second-outside", iat: testNowUnix - 36, nbf: testNowUnix - 36, exp: testNowUnix - 6, want: TimeInvalid},
		{name: "iat-one-second-outside", iat: testNowUnix + 6, nbf: testNowUnix, exp: testNowUnix + 30, want: TimeInvalid},
		{name: "lifetime-one-second-outside", iat: testNowUnix, nbf: testNowUnix, exp: testNowUnix + 31, want: TimeInvalid},
	}
	for i, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			payload := validPayload()
			payload["iat"], payload["nbf"], payload["exp"] = tc.iat, tc.nbf, tc.exp
			payload["jti"] = base64.RawURLEncoding.EncodeToString([]byte(fmt.Sprintf("%016d", i)))
			_, err := verifierAt(testNowUnix).VerifyAndConsume(context.Background(), mustEnvelope(t, signPayload(t, payload)), trustedContext())
			if tc.want == "" {
				require.NoError(t, err)
			} else {
				require.ErrorIs(t, err, tc.want)
			}
		})
	}
}

func verifierAt(now int64) Verifier {
	return NewVerifier(VerifierConfiguration{
		Keys:   map[string]ed25519.PublicKey{"connect-v2-test": testPublicKey},
		Replay: NewMemoryReplayConsumer(65536, func() time.Time { return time.Unix(now, 0) }),
		Now:    func() time.Time { return time.Unix(now, 0) },
	})
}

func trustedContext() TrustedProposalContext {
	return TrustedProposalContext{
		Issuer: "minekube-connect-test", TrustDomain: "urn:minekube:connect:test:corpus-v2",
		Audience: "urn:minekube:connect:test:bedrock-principal:v2", EndpointID: "endpoint-test",
		OrganizationID: "organization-test", ConnectSessionID: "session-test",
		ConnectSessionNonce: [16]byte{}, SourceProtocol: "bedrock", SourceProtocolVersion: 776,
		PolicyRevision: 7,
	}
}

func validPayload() map[string]any {
	return map[string]any{
		"version": 2, "issuer": "minekube-connect-test",
		"trust_domain": "urn:minekube:connect:test:corpus-v2",
		"audience":     "urn:minekube:connect:test:bedrock-principal:v2",
		"subject_kind": "bedrock_xuid", "canonical_xuid": "1",
		"canonical_unlinked_uuid": "00000000-0000-0000-0000-000000000001",
		"bedrock_display_name":    "BedrockOne", "endpoint_id": "endpoint-test",
		"organization_id": "organization-test", "connect_session_id": "session-test",
		"connect_session_nonce": testNonce, "policy_revision": 7,
		"source_protocol": "bedrock", "source_protocol_version": 776,
		"iat": testNowUnix, "nbf": testNowUnix, "exp": testNowUnix + 30, "jti": testJTI,
		"verification_method": "minecraft_legacy_chain+client_jwt+ecdh_v1",
	}
}

func linkedPayload(payload map[string]any) map[string]any {
	payload["subject_kind"] = "bedrock_linked_java"
	payload["verification_method"] = "minecraft_full_jwks+client_jwt+ecdh_v1"
	payload["linked_java"] = map[string]any{
		"uuid": "123e4567-e89b-12d3-a456-426614174000",
		"name": "JavaOne",
		"provenance": map[string]any{
			"provider": "moxy_account_link_v1", "record_id": "record-test-1",
			"revision": int64(3), "verified_at": testNowUnix - 10,
		},
	}
	return payload
}

func signPayload(t *testing.T, payload map[string]any) string {
	t.Helper()
	header := []byte(`{"alg":"EdDSA","typ":"connect-bedrock-principal+jws;v=2","kid":"connect-v2-test"}`)
	payloadBytes, err := json.Marshal(payload)
	require.NoError(t, err)
	protected := base64.RawURLEncoding.EncodeToString(header)
	body := base64.RawURLEncoding.EncodeToString(payloadBytes)
	input := protected + "." + body
	signature := ed25519.Sign(testPrivateKey, []byte(input))
	return input + "." + base64.RawURLEncoding.EncodeToString(signature)
}

func mustEnvelope(t *testing.T, compact string) SignedPrincipalEnvelope {
	t.Helper()
	envelope, err := NewSignedPrincipalEnvelope([]byte(compact))
	require.NoError(t, err)
	return envelope
}

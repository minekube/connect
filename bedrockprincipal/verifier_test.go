package bedrockprincipal

import (
	"context"
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
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

func requirePrincipalError(t *testing.T, err error, want PrincipalError) {
	t.Helper()
	require.True(t, errors.Is(err, want), "got %v, want %s", err, want)
}

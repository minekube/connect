package bedrockprincipal

import (
	"bytes"
	"context"
	"crypto/ed25519"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestVerificationEventSchemaAdmitsOnlyNonSecretFields(t *testing.T) {
	typeInfo := reflect.TypeOf(VerificationEvent{})
	want := []string{"Error", "Correlation", "KID", "Ready", "Linked", "ReplaySize", "ActiveCount", "ProfileApplications"}
	var got []string
	for i := 0; i < typeInfo.NumField(); i++ {
		got = append(got, typeInfo.Field(i).Name)
	}
	require.Equal(t, want, got)
}

func TestVerifierCapturesNeverExposeInjectedPrincipalMaterial(t *testing.T) {
	payload, expected, sentinels := sentinelLinkedPayload()
	compact := signPayload(t, payload)
	const rawEnvelopeSentinel = "envelope-capture-sentinel"
	sentinels = append(sentinels, compact, rawEnvelopeSentinel)

	t.Run("success serialization and debug log", func(t *testing.T) {
		capture, err := captureVerification(t, verifierAt(testNowUnix), mustEnvelope(t, compact), expected)
		require.NoError(t, err)
		assertNoSentinels(t, capture, sentinels)
	})

	t.Run("rejection serialization and debug log", func(t *testing.T) {
		capture, err := captureVerification(t, verifierAt(testNowUnix), mustEnvelope(t, tamperSignature(t, compact)), expected)
		require.Equal(t, Signature, err)
		assertNoSentinels(t, capture, sentinels)
	})

	t.Run("malformed envelope debug log", func(t *testing.T) {
		capture, err := captureVerification(t, verifierAt(testNowUnix), mustEnvelope(t, rawEnvelopeSentinel), expected)
		require.Equal(t, Malformed, err)
		assertNoSentinels(t, capture, sentinels)
	})
}

func TestVerifierStripsSensitiveWrappedDependencyErrors(t *testing.T) {
	payload, expected, sentinels := sentinelLinkedPayload()
	compact := signPayload(t, payload)
	sentinels = append(sentinels, compact, "wrapped-provider-sentinel", "wrapped-replay-sentinel")

	t.Run("key provider", func(t *testing.T) {
		verifier := NewVerifier(VerifierConfiguration{
			KeyProvider: keyProviderFunc(func(context.Context, string, string) (ed25519.PublicKey, error) {
				return nil, fmt.Errorf("wrapped-provider-sentinel: %w", Trust)
			}),
			Now: func() time.Time { return time.Unix(testNowUnix, 0) },
		})
		capture, err := captureVerification(t, verifier, mustEnvelope(t, compact), expected)
		require.Equal(t, Trust, err)
		assertNoSentinels(t, capture, sentinels)
	})

	t.Run("replay consumer", func(t *testing.T) {
		verifier := NewVerifier(VerifierConfiguration{
			Keys: map[string]ed25519.PublicKey{"connect-v2-test": testPublicKey},
			Replay: replayConsumerFunc(func(context.Context, ReplayEntry) error {
				return fmt.Errorf("wrapped-replay-sentinel: %w", Replay)
			}),
			Now: func() time.Time { return time.Unix(testNowUnix, 0) },
		})
		capture, err := captureVerification(t, verifier, mustEnvelope(t, compact), expected)
		require.Equal(t, Replay, err)
		assertNoSentinels(t, capture, sentinels)
	})

	t.Run("unknown dependency error", func(t *testing.T) {
		verifier := NewVerifier(VerifierConfiguration{
			KeyProvider: keyProviderFunc(func(context.Context, string, string) (ed25519.PublicKey, error) {
				return nil, errors.New("wrapped-provider-sentinel")
			}),
			Now: func() time.Time { return time.Unix(testNowUnix, 0) },
		})
		capture, err := captureVerification(t, verifier, mustEnvelope(t, compact), expected)
		require.Equal(t, Internal, err)
		assertNoSentinels(t, capture, sentinels)
	})
}

type keyProviderFunc func(context.Context, string, string) (ed25519.PublicKey, error)

func (f keyProviderFunc) Eligible(ctx context.Context, trustDomain, kid string) (ed25519.PublicKey, error) {
	return f(ctx, trustDomain, kid)
}

type replayConsumerFunc func(context.Context, ReplayEntry) error

func (f replayConsumerFunc) Consume(ctx context.Context, entry ReplayEntry) error {
	return f(ctx, entry)
}

func sentinelLinkedPayload() (map[string]any, TrustedProposalContext, []string) {
	const (
		xuid        = "9223372036854775706"
		displayName = "BedrockCaptureSentinel"
		javaUUID    = "deadbeef-dead-4bad-8bad-feedfacecafe"
		javaName    = "JavaSentinel"
		recordID    = "record-capture-sentinel"
	)
	nonce := base64.RawURLEncoding.EncodeToString([]byte("nonce-sentinel01"))
	jti := base64.RawURLEncoding.EncodeToString([]byte("jti-sentinel-001"))
	payload := linkedPayload(validPayload())
	payload["canonical_xuid"] = xuid
	payload["canonical_unlinked_uuid"] = uuidFromXUID(9223372036854775706).String()
	payload["bedrock_display_name"] = displayName
	payload["connect_session_nonce"] = nonce
	payload["jti"] = jti
	linked := payload["linked_java"].(map[string]any)
	linked["uuid"] = javaUUID
	linked["name"] = javaName
	linked["provenance"].(map[string]any)["record_id"] = recordID
	linked["provenance"].(map[string]any)["verified_at"] = testNowUnix
	expected := trustedContext()
	copy(expected.ConnectSessionNonce[:], []byte("nonce-sentinel01"))
	return payload, expected, []string{xuid, displayName, javaUUID, javaName, recordID, jti, nonce}
}

func captureVerification(t *testing.T, verifier Verifier, envelope SignedPrincipalEnvelope, expected TrustedProposalContext) (string, error) {
	t.Helper()
	principal, err := verifier.VerifyAndConsume(context.Background(), envelope, expected)
	category := PrincipalError("")
	if err != nil {
		var principalError PrincipalError
		if errors.As(err, &principalError) {
			category = principalError
		}
	}
	linked := false
	if principal != nil {
		_, linked = principal.LinkedJava()
	}
	event := VerificationEvent{Error: category, Correlation: [16]byte{1}, KID: "connect-v2-test", Linked: linked}
	serializedEvent, marshalErr := json.Marshal(event)
	require.NoError(t, marshalErr)
	serializedPrincipal, marshalErr := json.Marshal(principal)
	require.NoError(t, marshalErr)

	var debug bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&debug, &slog.HandlerOptions{Level: slog.LevelDebug}))
	logger.Debug("bedrock principal verification", "event", event, "error", err)

	returnedError := ""
	if err != nil {
		returnedError = fmt.Errorf("verification failed: %w", err).Error()
	}
	return strings.Join([]string{string(serializedEvent), string(serializedPrincipal), debug.String(), returnedError}, "\n"), err
}

func tamperSignature(t *testing.T, compact string) string {
	t.Helper()
	parts := strings.Split(compact, ".")
	require.Len(t, parts, 3)
	signature, err := base64.RawURLEncoding.DecodeString(parts[2])
	require.NoError(t, err)
	signature[0] ^= 1
	parts[2] = base64.RawURLEncoding.EncodeToString(signature)
	return strings.Join(parts, ".")
}

func assertNoSentinels(t *testing.T, capture string, sentinels []string) {
	t.Helper()
	for _, sentinel := range sentinels {
		require.NotContains(t, capture, sentinel)
	}
}

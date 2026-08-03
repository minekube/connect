package bedrockprincipal

import (
	"context"
	"crypto/ed25519"
	"encoding/base64"
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

type coreVector struct {
	Name                 string                 `json:"name"`
	CompactJWS           string                 `json:"compact_jws"`
	TrustedContext       coreTrustedContext     `json:"trusted_context"`
	VerificationTimeUnix int64                  `json:"verification_time_unix"`
	ExpectedError        string                 `json:"expected_error"`
	ExpectedPrincipal    *coreExpectedPrincipal `json:"expected_principal"`
}

type coreTrustedContext struct {
	Issuer, TrustDomain, Audience                string
	EndpointID, OrganizationID, ConnectSessionID string
	ConnectSessionNonce, SourceProtocol          string
	SourceProtocolVersion                        int32
	PolicyRevision                               int64
}

func (c *coreTrustedContext) UnmarshalJSON(data []byte) error {
	type wire struct {
		Issuer                string `json:"issuer"`
		TrustDomain           string `json:"trust_domain"`
		Audience              string `json:"audience"`
		EndpointID            string `json:"endpoint_id"`
		OrganizationID        string `json:"organization_id"`
		ConnectSessionID      string `json:"connect_session_id"`
		ConnectSessionNonce   string `json:"connect_session_nonce"`
		SourceProtocol        string `json:"source_protocol"`
		SourceProtocolVersion int32  `json:"source_protocol_version"`
		PolicyRevision        int64  `json:"policy_revision"`
	}
	var value wire
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	*c = coreTrustedContext(value)
	return nil
}

type coreExpectedPrincipal struct {
	SubjectKind           SubjectKind             `json:"subject_kind"`
	CanonicalXUID         CanonicalXUID           `json:"canonical_xuid"`
	CanonicalUnlinkedUUID uuid.UUID               `json:"canonical_unlinked_uuid"`
	BedrockDisplayName    string                  `json:"bedrock_display_name"`
	EffectiveUUID         uuid.UUID               `json:"effective_uuid"`
	EffectiveName         string                  `json:"effective_name"`
	VerificationMethod    string                  `json:"verification_method"`
	KID                   string                  `json:"kid"`
	PolicyRevision        int64                   `json:"policy_revision"`
	LinkedJava            *coreExpectedLinkedJava `json:"linked_java,omitempty"`
}

type coreExpectedLinkedJava struct {
	UUID           uuid.UUID `json:"uuid"`
	Name           string    `json:"name"`
	Provider       string    `json:"provider"`
	RecordID       string    `json:"record_id"`
	Revision       int64     `json:"revision"`
	VerifiedAtUnix int64     `json:"verified_at_unix"`
}

func TestCoreVectorsVerifyAgainstLiteralOutcomes(t *testing.T) {
	raw, err := os.ReadFile("testdata/v2/core-vectors.json")
	require.NoError(t, err)
	var vectors []coreVector
	require.NoError(t, json.Unmarshal(raw, &vectors))
	for _, vector := range vectors {
		t.Run(vector.Name, func(t *testing.T) {
			nonce, err := base64.RawURLEncoding.Strict().DecodeString(vector.TrustedContext.ConnectSessionNonce)
			require.NoError(t, err)
			require.Len(t, nonce, 16)
			require.Equal(t, vector.TrustedContext.ConnectSessionNonce, base64.RawURLEncoding.EncodeToString(nonce))
			var nonceArray [16]byte
			copy(nonceArray[:], nonce)
			contextValue := TrustedProposalContext{
				Issuer: vector.TrustedContext.Issuer, TrustDomain: vector.TrustedContext.TrustDomain,
				Audience: vector.TrustedContext.Audience, EndpointID: vector.TrustedContext.EndpointID,
				OrganizationID: vector.TrustedContext.OrganizationID, ConnectSessionID: vector.TrustedContext.ConnectSessionID,
				ConnectSessionNonce: nonceArray, SourceProtocol: vector.TrustedContext.SourceProtocol,
				SourceProtocolVersion: vector.TrustedContext.SourceProtocolVersion, PolicyRevision: vector.TrustedContext.PolicyRevision,
			}
			verifier := NewVerifier(VerifierConfiguration{Keys: map[string]ed25519.PublicKey{"connect-v2-test": testPublicKey}, Replay: NewMemoryReplayConsumer(MaxReplayEntries, func() time.Time { return time.Unix(vector.VerificationTimeUnix, 0) }), Now: func() time.Time { return time.Unix(vector.VerificationTimeUnix, 0) }})
			principal, err := verifier.VerifyAndConsume(context.Background(), mustEnvelope(t, vector.CompactJWS), contextValue)
			if vector.ExpectedError != "OK" {
				require.ErrorIs(t, err, PrincipalError(vector.ExpectedError))
				require.Nil(t, principal)
				return
			}
			require.NoError(t, err)
			require.Equal(t, vector.ExpectedPrincipal.SubjectKind, principal.SubjectKind())
			require.Equal(t, vector.ExpectedPrincipal.CanonicalXUID, principal.XUID())
			require.Equal(t, vector.ExpectedPrincipal.CanonicalUnlinkedUUID, principal.CanonicalUnlinkedUUID())
			require.Equal(t, GameProfile{UUID: vector.ExpectedPrincipal.EffectiveUUID, Name: vector.ExpectedPrincipal.EffectiveName}, principal.EffectiveGameProfile())
			require.Equal(t, vector.ExpectedPrincipal.VerificationMethod, principal.Verification().VerificationMethod)
			require.Equal(t, vector.ExpectedPrincipal.KID, principal.Verification().KID)
			require.Equal(t, vector.ExpectedPrincipal.PolicyRevision, principal.Bindings().PolicyRevision)
			linked, hasLinked := principal.LinkedJava()
			if vector.ExpectedPrincipal.LinkedJava == nil {
				require.False(t, hasLinked)
				return
			}
			require.True(t, hasLinked)
			require.Equal(t, vector.ExpectedPrincipal.LinkedJava.UUID, linked.UUID)
			require.Equal(t, vector.ExpectedPrincipal.LinkedJava.Name, linked.Name)
			require.Equal(t, vector.ExpectedPrincipal.LinkedJava.Provider, linked.Provenance.Provider)
			require.Equal(t, vector.ExpectedPrincipal.LinkedJava.RecordID, linked.Provenance.RecordID)
			require.Equal(t, vector.ExpectedPrincipal.LinkedJava.Revision, linked.Provenance.Revision)
			require.Equal(t, time.Unix(vector.ExpectedPrincipal.LinkedJava.VerifiedAtUnix, 0), linked.Provenance.VerifiedAt)
		})
	}
}

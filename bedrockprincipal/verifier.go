package bedrockprincipal

import (
	"context"
	"crypto/ed25519"
	"encoding/base64"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Verifier interface {
	VerifyAndConsume(context.Context, SignedPrincipalEnvelope, TrustedProposalContext) (VerifiedBedrockPrincipal, error)
}

type verifier struct {
	keys   KeyProvider
	replay ReplayConsumer
	now    func() time.Time
}

type staticKeys map[string]ed25519.PublicKey

func (s staticKeys) Eligible(_ context.Context, _ string, kid string) (ed25519.PublicKey, error) {
	key, ok := s[kid]
	if !ok {
		return nil, Trust
	}
	return append(ed25519.PublicKey(nil), key...), nil
}

func NewVerifier(configuration VerifierConfiguration) Verifier {
	provider := configuration.KeyProvider
	if provider == nil {
		provider = staticKeys(configuration.Keys)
	}
	now := configuration.Now
	if now == nil {
		now = time.Now
	}
	replay := configuration.Replay
	if replay == nil {
		replay = NewMemoryReplayConsumer(MaxReplayEntries, now)
	}
	return &verifier{keys: provider, replay: replay, now: now}
}

type protectedHeader struct {
	Alg string `json:"alg"`
	Typ string `json:"typ"`
	KID string `json:"kid"`
}

type payloadClaims struct {
	Version               int               `json:"version"`
	Issuer                string            `json:"issuer"`
	TrustDomain           string            `json:"trust_domain"`
	Audience              string            `json:"audience"`
	SubjectKind           SubjectKind       `json:"subject_kind"`
	CanonicalXUID         string            `json:"canonical_xuid"`
	CanonicalUnlinkedUUID string            `json:"canonical_unlinked_uuid"`
	LinkedJava            *linkedJavaClaims `json:"linked_java,omitempty"`
	BedrockDisplayName    string            `json:"bedrock_display_name"`
	EndpointID            string            `json:"endpoint_id"`
	OrganizationID        string            `json:"organization_id"`
	ConnectSessionID      string            `json:"connect_session_id"`
	ConnectSessionNonce   string            `json:"connect_session_nonce"`
	PolicyRevision        int64             `json:"policy_revision"`
	SourceProtocol        string            `json:"source_protocol"`
	SourceProtocolVersion int32             `json:"source_protocol_version"`
	IAT                   int64             `json:"iat"`
	NBF                   int64             `json:"nbf"`
	EXP                   int64             `json:"exp"`
	JTI                   string            `json:"jti"`
	VerificationMethod    string            `json:"verification_method"`
}

type linkedJavaClaims struct {
	UUID       string               `json:"uuid"`
	Name       string               `json:"name"`
	Provenance linkProvenanceClaims `json:"provenance"`
}

type linkProvenanceClaims struct {
	Provider   string `json:"provider"`
	RecordID   string `json:"record_id"`
	Revision   int64  `json:"revision"`
	VerifiedAt int64  `json:"verified_at"`
}

func (v *verifier) VerifyAndConsume(ctx context.Context, env SignedPrincipalEnvelope, expected TrustedProposalContext) (VerifiedBedrockPrincipal, error) {
	header, claims, signingInput, signature, err := parseEnvelope(env)
	if err != nil {
		return nil, Malformed
	}
	if claims.Issuer != expected.Issuer || claims.TrustDomain != expected.TrustDomain || claims.Audience != expected.Audience {
		return nil, Trust
	}
	key, err := v.keys.Eligible(ctx, expected.TrustDomain, header.KID)
	if err != nil {
		return nil, err
	}
	if len(key) != ed25519.PublicKeySize || !ed25519.Verify(key, signingInput, signature) {
		return nil, Signature
	}
	nonce, err := decodeCanonical16(claims.ConnectSessionNonce)
	if err != nil {
		return nil, Malformed
	}
	if _, err := decodeCanonical16(claims.JTI); err != nil {
		return nil, Malformed
	}
	if nonce != expected.ConnectSessionNonce || claims.EndpointID != expected.EndpointID ||
		claims.OrganizationID != expected.OrganizationID || claims.ConnectSessionID != expected.ConnectSessionID ||
		claims.SourceProtocol != expected.SourceProtocol || claims.SourceProtocolVersion != expected.SourceProtocolVersion ||
		claims.PolicyRevision != expected.PolicyRevision || expected.PolicyRevision <= 0 {
		return nil, BindingMismatch
	}
	if err := validateTime(claims, v.now().Unix()); err != nil {
		return nil, err
	}
	principal, err := principalFromClaims(claims, header.KID, expected)
	if err != nil {
		return nil, err
	}
	if err := v.replay.Consume(ctx, ReplayEntry{TrustDomain: claims.TrustDomain, Issuer: claims.Issuer, JTI: claims.JTI, KID: header.KID, EndpointID: claims.EndpointID, SessionID: claims.ConnectSessionID, SessionNonce: nonce, ExpiresAt: time.Unix(claims.EXP+5, 0)}); err != nil {
		return nil, err
	}
	return principal, nil
}

func parseEnvelope(env SignedPrincipalEnvelope) (protectedHeader, payloadClaims, []byte, []byte, error) {
	parts := strings.Split(env.compact, ".")
	if len(parts) != 3 || parts[0] == "" || parts[1] == "" || parts[2] == "" {
		return protectedHeader{}, payloadClaims{}, nil, nil, Malformed
	}
	headerBytes, err := decodeCanonical(parts[0], MaxHeaderBytes)
	if err != nil {
		return protectedHeader{}, payloadClaims{}, nil, nil, err
	}
	payloadBytes, err := decodeCanonical(parts[1], MaxPayloadBytes)
	if err != nil {
		return protectedHeader{}, payloadClaims{}, nil, nil, err
	}
	signature, err := decodeCanonical(parts[2], ed25519.SignatureSize)
	if err != nil || len(signature) != ed25519.SignatureSize {
		return protectedHeader{}, payloadClaims{}, nil, nil, Malformed
	}
	var header protectedHeader
	if err := strictObject(headerBytes, &header); err != nil || header.Alg != "EdDSA" || header.Typ != WireType || header.KID == "" || len(header.KID) > 128 {
		return protectedHeader{}, payloadClaims{}, nil, nil, Malformed
	}
	if err := requireObjectMembers(headerBytes, []string{"alg", "typ", "kid"}); err != nil {
		return protectedHeader{}, payloadClaims{}, nil, nil, Malformed
	}
	var claims payloadClaims
	if err := strictObject(payloadBytes, &claims); err != nil {
		return protectedHeader{}, payloadClaims{}, nil, nil, Malformed
	}
	if err := requireObjectMembers(payloadBytes, []string{
		"version", "issuer", "trust_domain", "audience", "subject_kind", "canonical_xuid",
		"canonical_unlinked_uuid", "bedrock_display_name", "endpoint_id", "organization_id",
		"connect_session_id", "connect_session_nonce", "policy_revision", "source_protocol",
		"source_protocol_version", "iat", "nbf", "exp", "jti", "verification_method",
	}, "linked_java"); err != nil {
		return protectedHeader{}, payloadClaims{}, nil, nil, Malformed
	}
	return header, claims, []byte(parts[0] + "." + parts[1]), signature, nil
}

func decodeCanonical(value string, maxDecoded int) ([]byte, error) {
	if strings.Contains(value, "=") {
		return nil, Malformed
	}
	decoded, err := base64.RawURLEncoding.Strict().DecodeString(value)
	if err != nil || len(decoded) > maxDecoded || base64.RawURLEncoding.EncodeToString(decoded) != value {
		return nil, Malformed
	}
	return decoded, nil
}

func decodeCanonical16(value string) ([16]byte, error) {
	var result [16]byte
	if len(value) != 22 {
		return result, Malformed
	}
	decoded, err := decodeCanonical(value, 16)
	if err != nil || len(decoded) != 16 {
		return result, Malformed
	}
	copy(result[:], decoded)
	return result, nil
}

func validateTime(claims payloadClaims, now int64) error {
	for _, value := range []int64{claims.IAT, claims.NBF, claims.EXP} {
		if value < 0 || value > 253402300799 {
			return TimeInvalid
		}
	}
	if claims.NBF > claims.IAT || claims.IAT > claims.EXP || claims.EXP-claims.IAT > 30 ||
		claims.NBF > now+5 || claims.IAT > now+5 || claims.EXP < now-5 {
		return TimeInvalid
	}
	return nil
}

func principalFromClaims(claims payloadClaims, kid string, bindings PrincipalBindings) (VerifiedBedrockPrincipal, error) {
	if claims.Version != 2 || claims.CanonicalXUID == "" || len(claims.CanonicalXUID) > 19 || claims.CanonicalXUID[0] == '0' {
		return nil, IdentityInvalid
	}
	xuid, err := strconv.ParseUint(claims.CanonicalXUID, 10, 63)
	if err != nil || xuid == 0 || xuid > math.MaxInt64 {
		return nil, IdentityInvalid
	}
	unlinked, err := uuid.Parse(claims.CanonicalUnlinkedUUID)
	if err != nil || unlinked != uuidFromXUID(xuid) {
		return nil, IdentityInvalid
	}
	if claims.BedrockDisplayName == "" || len(claims.BedrockDisplayName) > 64 {
		return nil, IdentityInvalid
	}
	if claims.VerificationMethod != "minecraft_legacy_chain+client_jwt+ecdh_v1" && claims.VerificationMethod != "minecraft_full_jwks+client_jwt+ecdh_v1" {
		return nil, IdentityInvalid
	}
	p := &verifiedBedrockPrincipal{kind: claims.SubjectKind, xuid: CanonicalXUID(claims.CanonicalXUID), unlinkedUUID: unlinked, displayName: claims.BedrockDisplayName, bindings: bindings,
		verification: VerificationEvidence{KID: kid, VerificationMethod: claims.VerificationMethod, IssuedAt: time.Unix(claims.IAT, 0), NotBefore: time.Unix(claims.NBF, 0), ExpiresAt: time.Unix(claims.EXP, 0)}}
	switch claims.SubjectKind {
	case BedrockXUID:
		if claims.LinkedJava != nil {
			return nil, LinkInvalid
		}
	case BedrockLinkedJava:
		if claims.LinkedJava == nil {
			return nil, LinkInvalid
		}
		linkedUUID, err := uuid.Parse(claims.LinkedJava.UUID)
		if err != nil || claims.LinkedJava.Name == "" || len(claims.LinkedJava.Name) > 16 || claims.LinkedJava.Provenance.Provider != "moxy_account_link_v1" || claims.LinkedJava.Provenance.RecordID == "" || claims.LinkedJava.Provenance.Revision <= 0 {
			return nil, LinkInvalid
		}
		p.linked = VerifiedLinkedJavaIdentity{UUID: linkedUUID, Name: claims.LinkedJava.Name, Provenance: LinkProvenance{Provider: claims.LinkedJava.Provenance.Provider, RecordID: claims.LinkedJava.Provenance.RecordID, Revision: claims.LinkedJava.Provenance.Revision, VerifiedAt: time.Unix(claims.LinkedJava.Provenance.VerifiedAt, 0)}}
		p.hasLinked = true
	default:
		return nil, IdentityInvalid
	}
	return p, nil
}

func uuidFromXUID(xuid uint64) uuid.UUID {
	var id uuid.UUID
	for i := 15; i >= 8; i-- {
		id[i] = byte(xuid)
		xuid >>= 8
	}
	return id
}

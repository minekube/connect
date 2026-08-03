package bedrockprincipal

import (
	"context"
	"crypto/ed25519"
	"time"

	"github.com/google/uuid"
)

const (
	WireType            = "connect-bedrock-principal+jws;v=2"
	Capability          = "bedrock-verified-principal-v2"
	MaxEnvelopeBytes    = 16 << 10
	MaxHeaderBytes      = 2 << 10
	MaxPayloadBytes     = 12 << 10
	MaxReplayEntries    = 65536
	MaxEnvelopeLifetime = 30 * time.Second
	ClockSkew           = 5 * time.Second
)

type PrincipalError string

func (e PrincipalError) Error() string { return string(e) }

const (
	Malformed           PrincipalError = "MALFORMED"
	Trust               PrincipalError = "TRUST"
	Signature           PrincipalError = "SIGNATURE"
	BindingMismatch     PrincipalError = "BINDING_MISMATCH"
	TimeInvalid         PrincipalError = "TIME"
	IdentityInvalid     PrincipalError = "IDENTITY"
	LinkInvalid         PrincipalError = "LINK"
	Replay              PrincipalError = "REPLAY"
	Capacity            PrincipalError = "CAPACITY"
	MetadataUnavailable PrincipalError = "METADATA_UNAVAILABLE"
	KeyRevoked          PrincipalError = "KEY_REVOKED"
	Readiness           PrincipalError = "READINESS"
	Internal            PrincipalError = "INTERNAL"
)

type SubjectKind string

const (
	BedrockXUID       SubjectKind = "bedrock_xuid"
	BedrockLinkedJava SubjectKind = "bedrock_linked_java"
)

type CanonicalXUID string

type GameProfile struct {
	UUID uuid.UUID
	Name string
}

type LinkProvenance struct {
	Provider   string
	RecordID   string
	Revision   int64
	VerifiedAt time.Time
}

type VerifiedLinkedJavaIdentity struct {
	UUID       uuid.UUID
	Name       string
	Provenance LinkProvenance
}

type VerificationEvidence struct {
	KID                string
	VerificationMethod string
	IssuedAt           time.Time
	NotBefore          time.Time
	ExpiresAt          time.Time
}

type PrincipalBindings struct {
	Issuer                string
	TrustDomain           string
	Audience              string
	EndpointID            string
	OrganizationID        string
	ConnectSessionID      string
	ConnectSessionNonce   [16]byte
	SourceProtocol        string
	SourceProtocolVersion int32
	PolicyRevision        int64
}

type TrustedProposalContext = PrincipalBindings

type SignedPrincipalEnvelope struct{ compact string }

func NewSignedPrincipalEnvelope(compact []byte) (SignedPrincipalEnvelope, error) {
	if len(compact) == 0 || len(compact) > MaxEnvelopeBytes {
		return SignedPrincipalEnvelope{}, Malformed
	}
	return SignedPrincipalEnvelope{compact: string(append([]byte(nil), compact...))}, nil
}

type KeyProvider interface {
	Eligible(context.Context, string, string) (ed25519.PublicKey, error)
}

type ReplayConsumer interface {
	Consume(context.Context, ReplayEntry) error
}

type VerifierConfiguration struct {
	Keys        map[string]ed25519.PublicKey
	KeyProvider KeyProvider
	Replay      ReplayConsumer
	Now         func() time.Time
}

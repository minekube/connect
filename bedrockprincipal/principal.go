package bedrockprincipal

import "github.com/google/uuid"

type VerifiedPrincipal interface {
	SubjectKind() SubjectKind
	EffectiveGameProfile() GameProfile
	verifiedPrincipal()
}

type VerifiedBedrockPrincipal interface {
	VerifiedPrincipal
	XUID() CanonicalXUID
	CanonicalUnlinkedUUID() uuid.UUID
	LinkedJava() (VerifiedLinkedJavaIdentity, bool)
	BedrockDisplayName() string
	Verification() VerificationEvidence
	Bindings() PrincipalBindings
}

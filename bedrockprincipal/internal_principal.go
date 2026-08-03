package bedrockprincipal

import "github.com/google/uuid"

type verifiedBedrockPrincipal struct {
	kind         SubjectKind
	xuid         CanonicalXUID
	unlinkedUUID uuid.UUID
	linked       VerifiedLinkedJavaIdentity
	hasLinked    bool
	displayName  string
	verification VerificationEvidence
	bindings     PrincipalBindings
}

func (*verifiedBedrockPrincipal) verifiedPrincipal()                 {}
func (p *verifiedBedrockPrincipal) SubjectKind() SubjectKind         { return p.kind }
func (p *verifiedBedrockPrincipal) XUID() CanonicalXUID              { return p.xuid }
func (p *verifiedBedrockPrincipal) CanonicalUnlinkedUUID() uuid.UUID { return p.unlinkedUUID }
func (p *verifiedBedrockPrincipal) LinkedJava() (VerifiedLinkedJavaIdentity, bool) {
	return p.linked, p.hasLinked
}
func (p *verifiedBedrockPrincipal) BedrockDisplayName() string         { return p.displayName }
func (p *verifiedBedrockPrincipal) Verification() VerificationEvidence { return p.verification }
func (p *verifiedBedrockPrincipal) Bindings() PrincipalBindings        { return p.bindings }
func (p *verifiedBedrockPrincipal) EffectiveGameProfile() GameProfile {
	if p.hasLinked {
		return GameProfile{UUID: p.linked.UUID, Name: p.linked.Name}
	}
	return GameProfile{UUID: p.unlinkedUUID, Name: p.displayName}
}

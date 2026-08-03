package compilefail

import (
	"github.com/google/uuid"
	bp "go.minekube.com/connect/bedrockprincipal"
)

type forged struct{}

func (forged) SubjectKind() bp.SubjectKind          { return bp.BedrockXUID }
func (forged) EffectiveGameProfile() bp.GameProfile { return bp.GameProfile{UUID: uuid.Nil} }

var _ bp.VerifiedPrincipal = forged{}

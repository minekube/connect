package bedrockprincipal

// VerificationEvent is the complete typed observability surface for principal
// verification. It deliberately cannot carry an envelope or identity value.
type VerificationEvent struct {
	Error               PrincipalError `json:"error"`
	Correlation         [16]byte       `json:"correlation"`
	KID                 string         `json:"kid"`
	Ready               bool           `json:"ready"`
	Linked              bool           `json:"linked"`
	ReplaySize          uint32         `json:"replay_size"`
	ActiveCount         uint32         `json:"active_count"`
	ProfileApplications uint32         `json:"profile_applications"`
}

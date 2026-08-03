package bedrockprincipal

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadinessFailsClosedUnlessRequireAndEveryLocalCheckPasses(t *testing.T) {
	ready := ReadinessState{Mode: "require", MetadataFresh: true, ReplayAvailable: true, ReplayCapacityAvailable: true, SelfCheckPassed: true, EligibleKeyCount: 1}
	require.True(t, ready.Ready())
	cases := map[string]func(*ReadinessState){
		"warn":               func(s *ReadinessState) { s.Mode = "warn" },
		"stale-metadata":     func(s *ReadinessState) { s.MetadataFresh = false },
		"replay-unavailable": func(s *ReadinessState) { s.ReplayAvailable = false },
		"replay-capacity":    func(s *ReadinessState) { s.ReplayCapacityAvailable = false },
		"self-check":         func(s *ReadinessState) { s.SelfCheckPassed = false },
		"no-eligible-key":    func(s *ReadinessState) { s.EligibleKeyCount = 0 },
	}
	for name, mutate := range cases {
		t.Run(name, func(t *testing.T) {
			state := ready
			mutate(&state)
			require.False(t, state.Ready())
			require.ErrorIs(t, state.Err(), Readiness)
		})
	}
}

func TestReadinessRevisionIsStableAndInputBound(t *testing.T) {
	input := ReadinessRevisionInput{Issuer: "issuer", TrustDomain: "domain", Audience: "audience", SDKBuildIdentity: "sdk", Capability: Capability, Mode: "require", EligibleKeySetDigest: [32]byte{1}, ReplayConfigurationDigest: [32]byte{2}, ProfileApplierIdentity: "host", SelfCheckCorpusDigest: [32]byte{3}}
	first := input.Revision()
	require.Equal(t, first, input.Revision())
	input.Mode = "warn"
	require.NotEqual(t, first, input.Revision())
}

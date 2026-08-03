package bedrockprincipal

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"

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

func TestVerificationEventNeverSerializesPrincipalMaterial(t *testing.T) {
	event := VerificationEvent{Error: LinkInvalid, Correlation: [16]byte{1}, KID: "public-test-kid", Ready: false, Linked: true, ReplaySize: 2, ActiveCount: 3, ProfileApplications: 0}
	raw, err := json.Marshal(event)
	require.NoError(t, err)
	serialized := string(raw)
	for _, forbidden := range []string{"xuid-sentinel", "display-name-sentinel", "java-uuid-sentinel", "record-sentinel", "jti-sentinel", "nonce-sentinel", "envelope-sentinel"} {
		require.False(t, strings.Contains(serialized, forbidden))
	}
}

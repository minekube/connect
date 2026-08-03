package geyserliteabi

import (
	"os"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestFrozenIngressBoundsAndReturnCodes(t *testing.T) {
	require.Equal(t, uint32(1), VerifiedIngressVersion)
	require.Equal(t, 16, CorrelationBytes)
	require.Equal(t, 1, MinIngressFrameBytes)
	require.Equal(t, 4096, MaxIngressFrameBytes)
	require.Equal(t, 5*time.Second, MaxIngressLifetime)
	require.Equal(t, int32(0), CallbackRegistrationOK)
	require.Equal(t, int32(0), AssignmentOK)
	require.Equal(t, int32(-1), AssignmentUnknownOrClosedHandle)
	require.Equal(t, int32(-2), AssignmentDuplicateHandleOrCorrelation)
	require.Equal(t, int32(-3), AssignmentInvalidOrExpiredTime)
	require.Equal(t, int32(-4), AssignmentWrongConnectionState)
}

func TestVerifiedIngressV1SchemaIsClosedAndOrdered(t *testing.T) {
	raw, err := os.ReadFile("verified_ingress_v1.proto")
	require.NoError(t, err)

	message := regexp.MustCompile(`(?s)message VerifiedIngressV1\s*\{(.*?)\}`).FindSubmatch(raw)
	require.Len(t, message, 2)
	fieldPattern := regexp.MustCompile(`(?m)^\s*(?:repeated\s+)?([A-Za-z_][A-Za-z0-9_.]*)\s+([a-z][a-z0-9_]*)\s*=\s*([0-9]+)(?:\s*\[[^]]+\])?;\s*$`)
	fields := fieldPattern.FindAllStringSubmatch(string(message[1]), -1)
	require.Equal(t, [][]string{
		{"uint32", "version", "1"},
		{"bytes", "correlation_id", "2"},
		{"string", "canonical_xuid", "3"},
		{"string", "display_name", "4"},
		{"string", "source_protocol", "5"},
		{"int32", "source_protocol_version", "6"},
		{"string", "verification_method", "7"},
		{"uint64", "verified_at_unix_ms", "8"},
	}, stripFullMatches(fields))
	require.NotContains(t, string(message[1]), "repeated ")
	require.NotRegexp(t, `(?m)\b(link|token|jwt|key|ip|skin|verified)\b\s*=`, string(message[1]))
}

func TestNativeHeaderFreezesCallbackABI(t *testing.T) {
	raw, err := os.ReadFile("geyserlite_verified_ingress_v1.h")
	require.NoError(t, err)
	header := strings.Join(strings.Fields(string(raw)), " ")

	require.Contains(t, header, "typedef void (*gate_verified_ingress_v1_cb)( const uint8_t correlation[16], const uint8_t *frame, uint32_t frame_len, uint64_t expires_unix_ms);")
	require.Contains(t, header, "typedef void (*gate_ingress_connection_open_v1_cb)(uint64_t connection_handle);")
	require.Contains(t, header, "int32_t geyserlite_set_ingress_callbacks_v1( graal_isolatethread_t *thread, gate_ingress_connection_open_v1_cb open, gate_verified_ingress_v1_cb verified);")
	require.Contains(t, header, "int32_t geyserlite_assign_verified_ingress_v1( graal_isolatethread_t *thread, uint64_t connection_handle, const uint8_t correlation[16], uint64_t expires_unix_ms);")
	require.Contains(t, header, "Both callback arguments null is the only unregister operation")
	require.Contains(t, header, "registration and unregistration are linearization barriers")
	require.Contains(t, header, "validate pointers and frame_len before reading native memory")
	require.Contains(t, header, "Java/native owns callback memory after return")
	for _, definition := range []string{
		"GEYSERLITE_INGRESS_CORRELATION_BYTES 16",
		"GEYSERLITE_INGRESS_FRAME_MAX_BYTES 4096",
		"GEYSERLITE_INGRESS_LIFETIME_MAX_MS 5000",
		"GEYSERLITE_ASSIGN_UNKNOWN_OR_CLOSED_HANDLE -1",
		"GEYSERLITE_ASSIGN_DUPLICATE_HANDLE_OR_CORRELATION -2",
		"GEYSERLITE_ASSIGN_INVALID_OR_EXPIRED_TIME -3",
		"GEYSERLITE_ASSIGN_WRONG_CONNECTION_STATE -4",
	} {
		require.Contains(t, header, definition)
	}
}

func TestSubprocessFramingConstants(t *testing.T) {
	require.Equal(t, uint8(1), SubprocessFrameVersion)
	require.Equal(t, uint8(1), SubprocessAssignment)
	require.Equal(t, uint8(2), SubprocessAssignmentACK)
	require.Equal(t, uint8(3), SubprocessVerifiedIngress)
	require.Equal(t, 32, SubprocessIPCKeyBytes)
	require.Equal(t, 32, SubprocessMACBytes)
	require.Equal(t, 8192, MaxAuthenticatedPacketBytes)

	raw, err := os.ReadFile("SUBPROCESS_FRAMING.md")
	require.NoError(t, err)
	contract := string(raw)
	require.Contains(t, contract, "version_u8=1 || generation_u64_be || ipc_key_32")
	require.Contains(t, contract, "version=1 || type=1 || generation_u64_be || sequence_u64_be || connection_handle_u64_be || correlation_16 || expires_unix_ms_u64_be || HMAC-SHA256")
	require.Contains(t, contract, "type=2")
	require.Contains(t, contract, "type=3")
}

func stripFullMatches(matches [][]string) [][]string {
	result := make([][]string, len(matches))
	for i := range matches {
		result[i] = matches[i][1:]
	}
	return result
}

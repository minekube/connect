package geyserliteabi

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"errors"
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
		"GEYSERLITE_SUBPROCESS_CONNECTION_OPEN 4",
		"GEYSERLITE_SUBPROCESS_ACK_POSITIVE 0",
		"GEYSERLITE_SUBPROCESS_ACK_NEGATIVE 1",
		"GEYSERLITE_SUBPROCESS_ACK_STATUS_OFFSET 42",
		"GEYSERLITE_SUBPROCESS_ACK_PACKET_BYTES 75",
		"GEYSERLITE_SUBPROCESS_CONNECTION_OPEN_PACKET_BYTES 58",
	} {
		require.Contains(t, header, definition)
	}
}

func TestSubprocessFramingConstants(t *testing.T) {
	require.Equal(t, uint8(1), SubprocessFrameVersion)
	require.Equal(t, uint8(1), SubprocessAssignment)
	require.Equal(t, uint8(2), SubprocessAssignmentACK)
	require.Equal(t, uint8(3), SubprocessVerifiedIngress)
	require.Equal(t, uint8(4), SubprocessConnectionOpen)
	require.Equal(t, uint8(0), SubprocessACKPositive)
	require.Equal(t, uint8(1), SubprocessACKNegative)
	require.Equal(t, 32, SubprocessIPCKeyBytes)
	require.Equal(t, 32, SubprocessMACBytes)
	require.Equal(t, 8192, MaxAuthenticatedPacketBytes)
	require.Equal(t, 58, SubprocessConnectionOpenPacketBytes)
	require.Equal(t, 75, SubprocessAssignmentACKPacketBytes)

	raw, err := os.ReadFile("SUBPROCESS_FRAMING.md")
	require.NoError(t, err)
	contract := string(raw)
	require.Contains(t, contract, "version_u8=1 || generation_u64_be || ipc_key_32")
	require.Contains(t, contract, "version=1 || type=1 || generation_u64_be || sequence_u64_be || connection_handle_u64_be || correlation_16 || expires_unix_ms_u64_be || HMAC-SHA256")
	require.Contains(t, contract, "type=2")
	require.Contains(t, contract, "type=3")
	require.Contains(t, contract, "type=4")
}

func TestSubprocessConnectionOpenRoundTrip(t *testing.T) {
	key := subprocessTestKey()
	want := subprocessTestPacket{
		packetType:       SubprocessConnectionOpen,
		generation:       0x0102030405060708,
		sequence:         0x1112131415161718,
		connectionHandle: 0x2122232425262728,
	}

	wire := encodeSubprocessTestPacket(t, key, want)
	require.Len(t, wire, SubprocessConnectionOpenPacketBytes)
	require.Equal(t,
		"01040102030405060708111213141516171821222324252627283365040b57686da6d582b14d5cae21a4d52d337afcbddafafa09afb18f1389df",
		hex.EncodeToString(wire),
	)

	got, err := decodeSubprocessConnectionOpenTestPacket(key, wire)
	require.NoError(t, err)
	require.Equal(t, want, got)
}

func TestSubprocessAssignmentACKStatusRoundTrip(t *testing.T) {
	key := subprocessTestKey()
	tests := []struct {
		name   string
		status uint8
		hex    string
	}{
		{
			name:   "positive",
			status: SubprocessACKPositive,
			hex:    "0102010203040506070811121314151617182122232425262728303132333435363738393a3b3c3d3e3f00e34e1ea4f2729bae7994c1077eacba40f62f25fdf85a6bce159dbc999dc01a0a",
		},
		{
			name:   "negative",
			status: SubprocessACKNegative,
			hex:    "0102010203040506070811121314151617182122232425262728303132333435363738393a3b3c3d3e3f01225c9e9ec7386c89065bd036f9a427b28df11a36c72f5737c8444ab2eb428d8a",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			want := subprocessTestPacket{
				packetType:       SubprocessAssignmentACK,
				generation:       0x0102030405060708,
				sequence:         0x1112131415161718,
				connectionHandle: 0x2122232425262728,
				correlation:      subprocessTestCorrelation(),
				ackStatus:        tt.status,
			}

			wire := encodeSubprocessTestPacket(t, key, want)
			require.Len(t, wire, SubprocessAssignmentACKPacketBytes)
			require.Equal(t, tt.hex, hex.EncodeToString(wire))
			got, err := decodeSubprocessAssignmentACKTestPacket(key, wire)
			require.NoError(t, err)
			require.Equal(t, want, got)
		})
	}
}

func TestSubprocessAssignmentACKRejectsMissingOrMalformedStatus(t *testing.T) {
	key := subprocessTestKey()
	legacyACKWithoutStatus := mustDecodeHex(t,
		"0102010203040506070811121314151617182122232425262728303132333435363738393a3b3c3d3e3fc6f55900ddae71ead5dad40dd1d77ca6086a28c74346a790fbc5afc54360b0a2",
	)
	_, err := decodeSubprocessAssignmentACKTestPacket(key, legacyACKWithoutStatus)
	require.ErrorIs(t, err, errSubprocessTestPacketLength)

	malformed := encodeSubprocessTestPacket(t, key, subprocessTestPacket{
		packetType:       SubprocessAssignmentACK,
		generation:       1,
		sequence:         1,
		connectionHandle: 1,
		correlation:      subprocessTestCorrelation(),
		ackStatus:        2,
	})
	_, err = decodeSubprocessAssignmentACKTestPacket(key, malformed)
	require.ErrorIs(t, err, errSubprocessTestACKStatus)
}

func TestExistingSubprocessPacketWireEncodingRemainsFrozen(t *testing.T) {
	key := subprocessTestKey()
	tests := []struct {
		name   string
		packet subprocessTestPacket
		hex    string
	}{
		{
			name: "assignment type 1",
			packet: subprocessTestPacket{
				packetType:       SubprocessAssignment,
				generation:       0x0102030405060708,
				sequence:         0x1112131415161718,
				connectionHandle: 0x2122232425262728,
				correlation:      subprocessTestCorrelation(),
				expiresUnixMS:    0x4142434445464748,
			},
			hex: "0101010203040506070811121314151617182122232425262728303132333435363738393a3b3c3d3e3f4142434445464748af5225f34100d86a7cff80d08d11b2cc06ae9c22c8be08e259a56ee61f1bc525",
		},
		{
			name: "verified ingress type 3",
			packet: subprocessTestPacket{
				packetType:       SubprocessVerifiedIngress,
				generation:       0x0102030405060708,
				sequence:         0x1112131415161718,
				connectionHandle: 0x2122232425262728,
				payload:          mustDecodeHex(t, "08011210303132333435363738393a3b3c3d3e3f"),
			},
			hex: "010301020304050607081112131415161718212223242526272808011210303132333435363738393a3b3c3d3e3fe82e9c5cc7512273daf1a6b54537d49a3f59d0cdeef965f3215c7862007232e1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.hex, hex.EncodeToString(encodeSubprocessTestPacket(t, key, tt.packet)))
		})
	}

	// Type 2 retains every pre-existing field at its original byte offset. The
	// new authenticated status byte is appended immediately before the MAC.
	positiveACK := encodeSubprocessTestPacket(t, key, subprocessTestPacket{
		packetType:       SubprocessAssignmentACK,
		generation:       0x0102030405060708,
		sequence:         0x1112131415161718,
		connectionHandle: 0x2122232425262728,
		correlation:      subprocessTestCorrelation(),
		ackStatus:        SubprocessACKPositive,
	})
	require.Equal(t,
		"0102010203040506070811121314151617182122232425262728303132333435363738393a3b3c3d3e3f",
		hex.EncodeToString(positiveACK[:SubprocessAssignmentACKStatusOffset]),
	)
}

var (
	errSubprocessTestPacketLength = errors.New("invalid subprocess packet length")
	errSubprocessTestACKStatus    = errors.New("invalid subprocess ACK status")
)

type subprocessTestPacket struct {
	packetType       uint8
	generation       uint64
	sequence         uint64
	connectionHandle uint64
	correlation      [CorrelationBytes]byte
	expiresUnixMS    uint64
	ackStatus        uint8
	payload          []byte
}

func encodeSubprocessTestPacket(t *testing.T, key []byte, packet subprocessTestPacket) []byte {
	t.Helper()
	var authenticated []byte
	switch packet.packetType {
	case SubprocessAssignment:
		authenticated = make([]byte, SubprocessAssignmentMACOffset)
		copy(authenticated[SubprocessCorrelationOffset:], packet.correlation[:])
		binary.BigEndian.PutUint64(authenticated[SubprocessAssignmentExpiresOffset:], packet.expiresUnixMS)
	case SubprocessAssignmentACK:
		authenticated = make([]byte, SubprocessAssignmentACKMACOffset)
		copy(authenticated[SubprocessCorrelationOffset:], packet.correlation[:])
		authenticated[SubprocessAssignmentACKStatusOffset] = packet.ackStatus
	case SubprocessVerifiedIngress:
		authenticated = make([]byte, SubprocessVerifiedIngressPayloadOffset+len(packet.payload))
		copy(authenticated[SubprocessVerifiedIngressPayloadOffset:], packet.payload)
	case SubprocessConnectionOpen:
		authenticated = make([]byte, SubprocessConnectionOpenMACOffset)
	default:
		t.Fatalf("unsupported test packet type %d", packet.packetType)
	}
	authenticated[SubprocessVersionOffset] = SubprocessFrameVersion
	authenticated[SubprocessTypeOffset] = packet.packetType
	binary.BigEndian.PutUint64(authenticated[SubprocessGenerationOffset:], packet.generation)
	binary.BigEndian.PutUint64(authenticated[SubprocessSequenceOffset:], packet.sequence)
	binary.BigEndian.PutUint64(authenticated[SubprocessConnectionHandleOffset:], packet.connectionHandle)
	mac := hmac.New(sha256.New, key)
	_, err := mac.Write(authenticated)
	require.NoError(t, err)
	return append(authenticated, mac.Sum(nil)...)
}

func decodeSubprocessConnectionOpenTestPacket(key, wire []byte) (subprocessTestPacket, error) {
	if len(wire) != SubprocessConnectionOpenPacketBytes {
		return subprocessTestPacket{}, errSubprocessTestPacketLength
	}
	if !validSubprocessTestMAC(key, wire, SubprocessConnectionOpenMACOffset) {
		return subprocessTestPacket{}, errors.New("invalid subprocess packet MAC")
	}
	return decodeSubprocessTestPrefix(wire, SubprocessConnectionOpen), nil
}

func decodeSubprocessAssignmentACKTestPacket(key, wire []byte) (subprocessTestPacket, error) {
	if len(wire) != SubprocessAssignmentACKPacketBytes {
		return subprocessTestPacket{}, errSubprocessTestPacketLength
	}
	if !validSubprocessTestMAC(key, wire, SubprocessAssignmentACKMACOffset) {
		return subprocessTestPacket{}, errors.New("invalid subprocess packet MAC")
	}
	status := wire[SubprocessAssignmentACKStatusOffset]
	if status != SubprocessACKPositive && status != SubprocessACKNegative {
		return subprocessTestPacket{}, errSubprocessTestACKStatus
	}
	packet := decodeSubprocessTestPrefix(wire, SubprocessAssignmentACK)
	copy(packet.correlation[:], wire[SubprocessCorrelationOffset:SubprocessAssignmentACKStatusOffset])
	packet.ackStatus = status
	return packet, nil
}

func decodeSubprocessTestPrefix(wire []byte, packetType uint8) subprocessTestPacket {
	if wire[SubprocessVersionOffset] != SubprocessFrameVersion || wire[SubprocessTypeOffset] != packetType {
		return subprocessTestPacket{}
	}
	return subprocessTestPacket{
		packetType:       packetType,
		generation:       binary.BigEndian.Uint64(wire[SubprocessGenerationOffset:]),
		sequence:         binary.BigEndian.Uint64(wire[SubprocessSequenceOffset:]),
		connectionHandle: binary.BigEndian.Uint64(wire[SubprocessConnectionHandleOffset:]),
	}
}

func validSubprocessTestMAC(key, wire []byte, macOffset int) bool {
	mac := hmac.New(sha256.New, key)
	_, _ = mac.Write(wire[:macOffset])
	return hmac.Equal(wire[macOffset:], mac.Sum(nil))
}

func subprocessTestKey() []byte {
	key := make([]byte, SubprocessIPCKeyBytes)
	for i := range key {
		key[i] = byte(i)
	}
	return key
}

func subprocessTestCorrelation() [CorrelationBytes]byte {
	var correlation [CorrelationBytes]byte
	copy(correlation[:], "0123456789:;<=>?")
	return correlation
}

func mustDecodeHex(t *testing.T, value string) []byte {
	t.Helper()
	decoded, err := hex.DecodeString(value)
	require.NoError(t, err)
	return decoded
}

func stripFullMatches(matches [][]string) [][]string {
	result := make([][]string, len(matches))
	for i := range matches {
		result[i] = matches[i][1:]
	}
	return result
}

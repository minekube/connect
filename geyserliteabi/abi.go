// Package geyserliteabi freezes the source-level GeyserLite VerifiedIngressV1
// handoff. It defines no runtime implementation.
package geyserliteabi

import "time"

const (
	CallbackABIVersion     uint32 = 1
	VerifiedIngressVersion uint32 = 1

	CorrelationBytes     = 16
	MinIngressFrameBytes = 1
	MaxIngressFrameBytes = 4096
	MaxIngressLifetime   = 5 * time.Second
)

const (
	CallbackRegistrationOK int32 = 0

	AssignmentOK                           int32 = 0
	AssignmentUnknownOrClosedHandle        int32 = -1
	AssignmentDuplicateHandleOrCorrelation int32 = -2
	AssignmentInvalidOrExpiredTime         int32 = -3
	AssignmentWrongConnectionState         int32 = -4
)

const (
	SubprocessFrameVersion    uint8 = 1
	SubprocessAssignment      uint8 = 1
	SubprocessAssignmentACK   uint8 = 2
	SubprocessVerifiedIngress uint8 = 3
	SubprocessConnectionOpen  uint8 = 4

	SubprocessACKPositive uint8 = 0
	SubprocessACKNegative uint8 = 1

	SubprocessBootstrapPacketBytes = 41
	SubprocessIPCKeyBytes          = 32
	SubprocessMACBytes             = 32
	MaxAuthenticatedPacketBytes    = 8192
)

// Authenticated subprocess packet byte offsets and exact or bounded sizes.
// The ACK status is appended after the existing correlation field so every
// pre-existing type-2 field retains its frozen offset.
const (
	SubprocessVersionOffset          = 0
	SubprocessTypeOffset             = 1
	SubprocessGenerationOffset       = 2
	SubprocessSequenceOffset         = 10
	SubprocessConnectionHandleOffset = 18
	SubprocessCorrelationOffset      = 26

	SubprocessAssignmentExpiresOffset = 42
	SubprocessAssignmentMACOffset     = 50
	SubprocessAssignmentPacketBytes   = 82

	SubprocessAssignmentACKStatusOffset = 42
	SubprocessAssignmentACKMACOffset    = 43
	SubprocessAssignmentACKPacketBytes  = 75

	SubprocessVerifiedIngressPayloadOffset  = 26
	MinSubprocessVerifiedIngressPacketBytes = 59
	MaxSubprocessVerifiedIngressPacketBytes = 4154

	SubprocessConnectionOpenMACOffset   = 26
	SubprocessConnectionOpenPacketBytes = 58
)

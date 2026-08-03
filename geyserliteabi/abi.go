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

	SubprocessIPCKeyBytes       = 32
	SubprocessMACBytes          = 32
	MaxAuthenticatedPacketBytes = 8192
)

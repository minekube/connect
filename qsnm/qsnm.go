package qsnm

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"slices"
	"strings"
	"time"
)

const (
	EventConnectPublicSessionSuccess = "connect.public_session_success.v1"

	PropSource              = "source"
	PropOccurredAt          = "occurred_at"
	PropNetworkKey          = "network_key"
	PropEvidenceID          = "evidence_id"
	PropPlanState           = "plan_state"
	PropIntegrationFamily   = "integration_family"
	PropIntegrationVersion  = "integration_version"
	PropConnectTransport    = "connect_transport"
	PropPlayerCountBucket   = "player_count_bucket"
	PropSessionCountBucket  = "session_count_bucket"
	PropIdentityConfidence  = "identity_confidence"
	PropSchemaVersion       = "schema_version"
	PropPrivacyPosture      = "privacy_posture"
	PropQSNMQualifyingEvent = "qsnm_qualifying_event"
)

var ConnectPublicSessionSuccessAllowlist = []string{
	PropSource,
	PropOccurredAt,
	PropNetworkKey,
	PropEvidenceID,
	PropPlanState,
	PropIntegrationFamily,
	PropIntegrationVersion,
	PropConnectTransport,
	PropPlayerCountBucket,
	PropSessionCountBucket,
	PropIdentityConfidence,
	PropSchemaVersion,
	PropPrivacyPosture,
	PropQSNMQualifyingEvent,
}

type ConnectPublicSessionSuccess struct {
	Source             string
	OccurredAt         time.Time
	NetworkKey         string
	EvidenceID         string
	PlanState          string
	IntegrationFamily  string
	IntegrationVersion string
	ConnectTransport   string
	PlayerCountBucket  string
	SessionCountBucket string
	IdentityConfidence string
}

func MonthlyNetworkKey(secret []byte, canonicalNetworkFingerprint string, occurredAt time.Time) (string, error) {
	if len(secret) == 0 {
		return "", errors.New("qsnm secret must not be empty")
	}
	fingerprint := strings.TrimSpace(canonicalNetworkFingerprint)
	if fingerprint == "" {
		return "", errors.New("canonical network fingerprint must not be empty")
	}
	month := occurredAt.UTC().Format("2006-01")
	mac := hmac.New(sha256.New, secret)
	_, _ = mac.Write([]byte(fingerprint))
	_, _ = mac.Write([]byte(month))
	return "qsnm_" + month + "_" + hex.EncodeToString(mac.Sum(nil)), nil
}

func (e ConnectPublicSessionSuccess) Properties() (map[string]any, error) {
	if e.OccurredAt.IsZero() {
		return nil, errors.New("occurred_at must be set")
	}
	if strings.TrimSpace(e.NetworkKey) == "" {
		return nil, errors.New("network_key must be set")
	}
	if strings.TrimSpace(e.EvidenceID) == "" {
		return nil, errors.New("evidence_id must be set")
	}
	props := map[string]any{
		PropSource:              defaultString(e.Source, "connect"),
		PropOccurredAt:          e.OccurredAt.UTC().Format(time.RFC3339),
		PropNetworkKey:          e.NetworkKey,
		PropEvidenceID:          e.EvidenceID,
		PropPlanState:           e.PlanState,
		PropIntegrationFamily:   e.IntegrationFamily,
		PropIntegrationVersion:  e.IntegrationVersion,
		PropConnectTransport:    e.ConnectTransport,
		PropPlayerCountBucket:   e.PlayerCountBucket,
		PropSessionCountBucket:  e.SessionCountBucket,
		PropIdentityConfidence:  defaultString(e.IdentityConfidence, "strong"),
		PropSchemaVersion:       1,
		PropPrivacyPosture:      "hmac_monthly_network_key_no_raw_host_ip_domain_email_config",
		PropQSNMQualifyingEvent: true,
	}
	if err := validateAllowlist(props, ConnectPublicSessionSuccessAllowlist); err != nil {
		return nil, err
	}
	return props, nil
}

func validateAllowlist(props map[string]any, allowlist []string) error {
	for key := range props {
		if !slices.Contains(allowlist, key) {
			return fmt.Errorf("property %q is not in QSNM allowlist", key)
		}
	}
	return nil
}

func defaultString(value, fallback string) string {
	if strings.TrimSpace(value) == "" {
		return fallback
	}
	return value
}

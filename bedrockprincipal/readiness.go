package bedrockprincipal

import (
	"crypto/sha256"
	"encoding/binary"
	"hash"
)

type ReadinessState struct {
	Mode                    string
	MetadataFresh           bool
	ReplayAvailable         bool
	ReplayCapacityAvailable bool
	SelfCheckPassed         bool
	EligibleKeyCount        int
}

func (s ReadinessState) Ready() bool {
	return s.Mode == "require" && s.MetadataFresh && s.ReplayAvailable &&
		s.ReplayCapacityAvailable && s.SelfCheckPassed && s.EligibleKeyCount > 0
}

func (s ReadinessState) Err() error {
	if s.Ready() {
		return nil
	}
	return Readiness
}

type ReadinessRevisionInput struct {
	Issuer                    string
	TrustDomain               string
	Audience                  string
	SDKBuildIdentity          string
	Capability                string
	Mode                      string
	EligibleKeySetDigest      [32]byte
	ReplayConfigurationDigest [32]byte
	ProfileApplierIdentity    string
	SelfCheckCorpusDigest     [32]byte
}

func (i ReadinessRevisionInput) Revision() [32]byte {
	digest := sha256.New()
	for _, value := range []string{i.Issuer, i.TrustDomain, i.Audience, i.SDKBuildIdentity, i.Capability, i.Mode} {
		writeLengthPrefixed(digest, []byte(value))
	}
	writeLengthPrefixed(digest, i.EligibleKeySetDigest[:])
	writeLengthPrefixed(digest, i.ReplayConfigurationDigest[:])
	writeLengthPrefixed(digest, []byte(i.ProfileApplierIdentity))
	writeLengthPrefixed(digest, i.SelfCheckCorpusDigest[:])
	var result [32]byte
	copy(result[:], digest.Sum(nil))
	return result
}

func writeLengthPrefixed(digest hash.Hash, value []byte) {
	var length [4]byte
	binary.BigEndian.PutUint32(length[:], uint32(len(value)))
	_, _ = digest.Write(length[:])
	_, _ = digest.Write(value)
}

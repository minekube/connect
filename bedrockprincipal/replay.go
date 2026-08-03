package bedrockprincipal

import (
	"context"
	"sync"
	"time"
)

type ReplayEntry struct {
	TrustDomain  string
	Issuer       string
	JTI          string
	KID          string
	EndpointID   string
	SessionID    string
	SessionNonce [16]byte
	ExpiresAt    time.Time
}

type MemoryReplayConsumer struct {
	maxEntries int
	now        func() time.Time
	mu         sync.Mutex
	entries    map[string]ReplayEntry
}

func NewMemoryReplayConsumer(maxEntries int, now func() time.Time) *MemoryReplayConsumer {
	if maxEntries <= 0 || maxEntries > MaxReplayEntries {
		maxEntries = MaxReplayEntries
	}
	if now == nil {
		now = time.Now
	}
	return &MemoryReplayConsumer{maxEntries: maxEntries, now: now, entries: make(map[string]ReplayEntry, maxEntries)}
}

func (c *MemoryReplayConsumer) Consume(_ context.Context, entry ReplayEntry) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := c.now()
	key := entry.TrustDomain + "\x00" + entry.Issuer + "\x00" + entry.JTI
	if existing, exists := c.entries[key]; exists {
		if !now.After(existing.ExpiresAt) {
			return Replay
		}
		delete(c.entries, key)
	}
	if len(c.entries) >= c.maxEntries {
		removed := 0
		for existingKey, existing := range c.entries {
			if now.After(existing.ExpiresAt) {
				delete(c.entries, existingKey)
				removed++
				if removed == 64 {
					break
				}
			}
		}
	}
	if len(c.entries) >= c.maxEntries {
		return Capacity
	}
	c.entries[key] = entry
	return nil
}

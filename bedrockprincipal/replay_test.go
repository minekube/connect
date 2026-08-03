package bedrockprincipal

import (
	"context"
	"encoding/base64"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestReplaySameEnvelopeIsOKThenReplay(t *testing.T) {
	cache := replayCache()
	entry := replayEntry(1)
	require.NoError(t, cache.Consume(context.Background(), entry))
	require.ErrorIs(t, cache.Consume(context.Background(), entry), Replay)
}

func TestVerifierConsumesSameEnvelopeAsOKThenReplay(t *testing.T) {
	verifier := verifierAt(testNowUnix)
	envelope := mustEnvelope(t, signPayload(t, validPayload()))
	_, err := verifier.VerifyAndConsume(context.Background(), envelope, trustedContext())
	require.NoError(t, err)
	_, err = verifier.VerifyAndConsume(context.Background(), envelope, trustedContext())
	require.ErrorIs(t, err, Replay)
}

func TestReplayConcurrentConsumersAllowExactlyOne(t *testing.T) {
	cache := replayCache()
	entry := replayEntry(2)
	var successes atomic.Int32
	errorsByCaller := make(chan error, 64)
	var wg sync.WaitGroup
	for i := 0; i < 64; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := cache.Consume(context.Background(), entry)
			if err == nil {
				successes.Add(1)
				return
			}
			errorsByCaller <- err
		}()
	}
	wg.Wait()
	close(errorsByCaller)
	require.Equal(t, int32(1), successes.Load())
	for err := range errorsByCaller {
		require.ErrorIs(t, err, Replay)
	}
}

func TestReplayCapacityNeverEvictsUnexpiredEntries(t *testing.T) {
	cache := replayCache()
	for i := 0; i < MaxReplayEntries; i++ {
		require.NoError(t, cache.Consume(context.Background(), replayEntry(i)), i)
	}
	require.ErrorIs(t, cache.Consume(context.Background(), replayEntry(MaxReplayEntries)), Capacity)
	require.ErrorIs(t, cache.Consume(context.Background(), replayEntry(0)), Replay)
}

func TestReplayIdentityExcludesKID(t *testing.T) {
	cache := replayCache()
	entry := replayEntry(7)
	require.NoError(t, cache.Consume(context.Background(), entry))
	entry.KID = "rotated-kid"
	require.ErrorIs(t, cache.Consume(context.Background(), entry), Replay)
}

func TestReplayRetainsThroughExpiryBoundaryThenAllowsReuse(t *testing.T) {
	now := time.Unix(testNowUnix, 0)
	cache := NewMemoryReplayConsumer(MaxReplayEntries, func() time.Time { return now })
	entry := replayEntry(9)
	entry.ExpiresAt = now
	require.NoError(t, cache.Consume(context.Background(), entry))
	require.ErrorIs(t, cache.Consume(context.Background(), entry), Replay)
	now = now.Add(time.Second)
	require.NoError(t, cache.Consume(context.Background(), entry))
}

func replayCache() *MemoryReplayConsumer {
	return NewMemoryReplayConsumer(MaxReplayEntries, func() time.Time { return time.Unix(testNowUnix, 0) })
}

func replayEntry(i int) ReplayEntry {
	jti := base64.RawURLEncoding.EncodeToString([]byte(fmt.Sprintf("%016d", i)))
	return ReplayEntry{TrustDomain: "domain", Issuer: "issuer", JTI: jti, KID: "kid", ExpiresAt: time.Unix(testNowUnix+35, 0)}
}

package testutil

import (
	"fmt"
	"sync"
	"testing"
)

// TestSequenceKeepsConcurrentObservations covers the recorder used by the
// suite: TestWatchReject observes one exchange from the server and the client
// side at the same time, so observations are appended concurrently. A recorder
// that drops one of them makes the suite compare a truncated sequence and fail
// for a reason that has nothing to do with the exchange under test.
func TestSequenceKeepsConcurrentObservations(t *testing.T) {
	const observations = 200

	var seq sequence
	var wg sync.WaitGroup
	wg.Add(observations)
	for i := 0; i < observations; i++ {
		go func(i int) {
			defer wg.Done()
			seq.Add(fmt.Sprintf("observation %d", i))
		}(i)
	}
	wg.Wait()

	if got := len(seq.Get()); got != observations {
		t.Fatalf("recorded %d observations, want %d", got, observations)
	}
}

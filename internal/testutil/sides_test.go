package testutil

import (
	"errors"
	"strings"
	"testing"
	"time"
)

// TestSideCollectsItsChecksInsteadOfReportingThem covers the sink used by the
// concurrently running sides of a suite test method (the client side, the
// server side accepting the endpoint or tunnel, the call starting the test
// server): a side may outlive the method that started it, and testify's
// suite.Run installs the *testing.T and the assertions of the next method as
// soon as the current one returned (suite.Suite.SetT). A side that reports its
// checks through the suite therefore reads what that write changes - a data
// race - and a failure it reports is attributed to the next test, or to no test
// at all. The side collects them here instead, and the test method asserts them
// in its own goroutine.
func TestSideCollectsItsChecksInsteadOfReportingThem(t *testing.T) {
	var sides sides
	side := sides.add("test side")

	done := make(chan struct{})
	go func() {
		defer close(done)
		defer side.finish()

		check := side.check()
		check.NoError(nil)
		check.Equal(2, 1+1)
		check.NoError(errors.New("boom"))
	}()
	<-done

	failures := sides.wait(time.Second)
	if len(failures) != 1 {
		t.Fatalf("collected %d failures %q, want only the failing check", len(failures), failures)
	}
	if !strings.HasPrefix(failures[0], "test side: ") {
		t.Errorf("failure %q does not name the side that collected it", failures[0])
	}
	if !strings.Contains(failures[0], "boom") {
		t.Errorf("failure %q does not carry the check's message", failures[0])
	}
}

// TestSidesWaitDoesNotReturnBeforeTheSidesFinished covers the join a suite test
// method performs before it asserts anything: while a side is still running its
// collected checks are incomplete, and a side that still runs when the method
// returned is exactly what races with the next method. A side that never
// finishes must fail the test instead of hanging it.
func TestSidesWaitDoesNotReturnBeforeTheSidesFinished(t *testing.T) {
	var sides sides
	side := sides.add("blocked side")

	release := make(chan struct{})
	go func() {
		defer side.finish()
		<-release
		side.check().Equal(1, 2, "collected after the first wait")
	}()

	failures := sides.wait(time.Millisecond * 50)
	if len(failures) != 1 || !strings.Contains(failures[0], "did not finish") {
		t.Fatalf("wait collected %q, want a single 'did not finish' failure", failures)
	}

	close(release)
	failures = sides.wait(time.Second)
	if len(failures) != 2 || !strings.Contains(failures[1], "collected after the first wait") {
		t.Fatalf("wait collected %q, want the checks the side collected after the first wait", failures)
	}
}

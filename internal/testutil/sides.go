package testutil

import (
	"fmt"
	"slices"
	"sync"
	"time"

	"github.com/stretchr/testify/assert"
)

// sideWaitTimeout bounds how long a suite test method waits for its sides to
// finish. Every side is bound to the test's context and finishes right after
// the exchange under test is over; the bound only turns a side that never
// finishes into a failure of the test that waited for it instead of a hung
// test.
const sideWaitTimeout = time.Second * 5

// sides collects the concurrently running sides of one suite test method.
type sides struct {
	mu  sync.Mutex
	all []*side
}

// add registers a new side of the test method.
func (ss *sides) add(name string) *side {
	s := &side{name: name, done: make(chan struct{})}
	ss.mu.Lock()
	defer ss.mu.Unlock()
	ss.all = append(ss.all, s)
	return s
}

// wait blocks until every side finished, but no longer than d per side, and
// returns the checks they collected in the order they were collected. The
// method asserts them in its own goroutine; an empty result means every side
// made it through its checks.
func (ss *sides) wait(d time.Duration) []string {
	ss.mu.Lock()
	all := slices.Clone(ss.all)
	ss.mu.Unlock()

	var collected []string
	for _, s := range all {
		collected = append(collected, s.wait(d)...)
	}
	return collected
}

// side is one concurrently running side of a suite test method: the client
// side, the server side accepting the endpoint or tunnel, or the call that
// starts the test server.
//
// A side outlives the method that started it (the method returns as soon as the
// exchange under test is over, which is why the side exists), so it may NOT use
// the suite's assertions: testify's suite.Run installs the *testing.T and the
// assertions of the next method as soon as the current one returned
// (suite.Suite.SetT), leaving the side reading what that write changes - a data
// race - and attributing any failure it reports to the next test, or to no test
// at all (minekube/connect#172). A side collects its checks here instead, and
// the test method asserts them after waiting for every side (sides.wait).
type side struct {
	name string

	done chan struct{}
	once sync.Once

	mu       sync.Mutex
	failures []string
}

// finish marks the side as finished. It is called by the goroutine that runs
// the side, usually as a deferred call.
func (s *side) finish() { s.once.Do(func() { close(s.done) }) }

// check returns the assertions of the side. They may only be called from the
// goroutine that runs the side.
func (s *side) check() *assert.Assertions { return assert.New(s) }

// Errorf collects a failed check of the side instead of reporting it to the
// suite. It implements assert.TestingT.
func (s *side) Errorf(format string, args ...any) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.failures = append(s.failures, s.name+": "+fmt.Sprintf(format, args...))
}

// wait blocks until the side finished, but no longer than d, and returns the
// checks the side collected. A side that did not finish in time is collected as
// the first of them, so it fails the test that waited for it.
func (s *side) wait(d time.Duration) []string {
	select {
	case <-s.done:
	case <-time.After(d):
		s.Errorf("did not finish within %s", d)
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	return slices.Clone(s.failures)
}

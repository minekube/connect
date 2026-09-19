# Project agent memory

This file is the project's committed home for project-intrinsic agent knowledge: build, test, release, architecture, and sharp-edge notes that should travel with the code.

- Add durable project-specific notes here as they are discovered through real work.

## Maintaining this file

Keep this file for knowledge useful to almost every future agent session in this project.
Do not repeat what the codebase already shows; point to the authoritative file or command instead.
Prefer rewriting or pruning existing entries over appending new ones.
When updating this file, preserve this bar for all agents and keep entries concise.

## Suite test sides must not outlive their test method

The `testutil.Suite` test methods (`TestWatchReject`, `TestTunnel`) return as soon as the
exchange they test is over, so the sides that run it - the client, the server handler
accepting the endpoint/tunnel, the `Start*Server` call - keep running in goroutines that
outlive the method. A side must not use the suite's assertions: testify's `suite.Run`
installs the next method's `*testing.T` and assertions as soon as the current method
returned (`suite.Suite.SetT`), so a side that reads them races with that write and reports
its failures against the next test. Sides collect their checks with `side.check()`, and the
method joins every side with `sides.wait(sideWaitTimeout)` before it asserts anything.
See `internal/testutil/sides.go`.

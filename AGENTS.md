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

## Watch payloads must not be built with struct literals

The generated proto revision this module compiles against is decided by the
CONSUMER's build list (MVS), never by `go.mod`'s minimum, and revisions disagree
on how `WatchRequest`/`WatchResponse` expose their single payload: an older
generation declares plain optional fields (`session_rejection`, `session`), the
readiness-attestation generation declares them as members of a `payload` oneof.
Both encode the same wire field (`1`), so build payloads with
`connect.NewWatchRequestWithRejection` / `connect.NewWatchResponseWithSession`
(they set the field through the message descriptor and fail closed if it moves)
instead of a literal such as `&connect.WatchRequest{SessionRejection: ...}` - a
literal pins one Go shape and breaks every consumer that resolves the other
(that is how moxy's Renovate digest PRs and connect#42 went red in 2026-10).
Read side is `GetSession()`/`GetSessionRejection()`, which both shapes provide.
The `proto-compat` CI job runs the suite against the newest published
generation; do not raise `go.mod`'s generated-proto minimum to "fix" a shape
mismatch - consumers like Gate deliberately pin an older generation.

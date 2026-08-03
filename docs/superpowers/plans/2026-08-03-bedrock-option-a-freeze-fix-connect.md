# Bedrock Option A Freeze Fix (Connect) Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Complete Connect's captain-approved checkpoint-1 protocol freeze by landing the GeyserLite ingress ABI, enforcing omitted verifier bounds, proving privacy through real verifier paths, and recording the completed Moxy descriptor binding.

**Architecture:** Add a source-only `geyserliteabi` sibling package containing the canonical closed protobuf, native header, subprocess framing contract, and Go constants that downstream implementations can compile or validate against. Keep claim validation inside `bedrockprincipal` and collapse dependency errors to the existing closed `PrincipalError` enum. Update the wire follow-up only after the companion Moxy branch contains the real generated-descriptor proof.

**Tech Stack:** Go 1.21, protobuf source, C ABI header, standard-library testing, Testify.

## Global Constraints

- Preserve checkpoint correlation `3b2db3b29d50b1ac` and spec commit `78754829f610f54a8bbe3d043e48bd60296a1e9e`.
- This is source/interface-only: no release, deployment, credential, OAuth, production, customer-endpoint, or Moxy mutation.
- Correlation is exactly 16 bytes, ingress frames are 1–4096 bytes, and ingress lifetime is at most five seconds.
- The verifier exposes only the existing thirteen `PrincipalError` categories as errors.

---

### Task 1: Frozen GeyserLite ABI package

**Files:**
- Create: `geyserliteabi/abi_test.go`
- Create: `geyserliteabi/abi.go`
- Create: `geyserliteabi/verified_ingress_v1.proto`
- Create: `geyserliteabi/geyserlite_verified_ingress_v1.h`
- Create: `geyserliteabi/SUBPROCESS_FRAMING.md`

**Interfaces:**
- Produces: `VerifiedIngressV1` protobuf fields 1–8, callback ABI v1 symbols, assignment return codes 0/-1/-2/-3/-4, and frozen ingress/subprocess bounds.

- [ ] Write structural tests that require the exact eight-field protobuf, both native callback signatures, return codes, 16-byte correlation, 4096-byte frame maximum, five-second lifetime, and subprocess packet constants.
- [ ] Run `go test ./geyserliteabi -count=1` and confirm it fails because the package artifacts are absent.
- [ ] Add the minimal canonical source artifacts and constants.
- [ ] Run `go test ./geyserliteabi -count=1` and confirm it passes.

### Task 2: Frozen verifier claim limits

**Files:**
- Modify: `bedrockprincipal/verifier_test.go`
- Modify: `bedrockprincipal/verifier.go`
- Modify: `bedrockprincipal/internal/vectorgen/generator.go`
- Modify: `bedrockprincipal/testdata/v2/core-vectors.json`
- Modify: `bedrockprincipal/testdata/v2/connect-sdk-candidate.json`

**Interfaces:**
- Consumes: schema maxima in `bedrockprincipal/schema/v2.schema.json` and `ClockSkew`.
- Produces: rejection of issuer >128, trust domain >256, audience >256, and linked `verified_at` more than five seconds from `iat`.

- [ ] Add permanent negative tests for every omitted maximum and both one-second-outside linked-time directions, retaining equality acceptance.
- [ ] Run the focused tests and confirm the demonstrated invalid inputs are accepted before the fix.
- [ ] Add bounded trust-claim validation before key lookup and linked-time comparison during principal construction.
- [ ] Correct the valid-linked literal vector to the accepted five-second equality boundary, regenerate it deterministically, and update the pinned corpus digest.
- [ ] Run focused and full `bedrockprincipal` tests and confirm they pass.

### Task 3: End-to-end observability privacy proof

**Files:**
- Modify: `bedrockprincipal/observability_test.go`
- Modify: `bedrockprincipal/verifier.go`

**Interfaces:**
- Consumes: real signed linked envelopes, malformed envelopes, debug serialization/logging, and wrapped `KeyProvider`/`ReplayConsumer` errors.
- Produces: category-only returned errors and captures containing none of the injected XUID, display name, UUID, record ID, JTI, nonce, or envelope sentinels.

- [ ] Replace the tautological serialization test with success, rejection, debug-log, and wrapped-error path tests using unique sensitive sentinels.
- [ ] Run focused observability tests and confirm wrapped dependency errors leak sentinel text before the fix.
- [ ] Normalize all dependency failures to an unwrapped member of the closed `PrincipalError` set, using `INTERNAL` for unknown errors.
- [ ] Run focused and full tests and confirm every capture is sentinel-free.

### Task 4: Completed descriptor binding statement

**Files:**
- Modify: `bedrockprincipal/WIRE_FOLLOWUP.md`

**Interfaces:**
- Consumes: committed companion Moxy branch proof against actual generated Watch and libp2p descriptors.
- Produces: accurate source-only status with no still-deferred descriptor claim.

- [ ] Inspect the companion task status and branch commit for the real generated-descriptor proof.
- [ ] Once present, rewrite `WIRE_FOLLOWUP.md` to name the completed Moxy-owned binding and leave only genuinely downstream release/consumer work.

### Task 5: Verification and commit

**Files:**
- Verify all changed files and repository state.

- [ ] Run `gofmt` on changed Go files.
- [ ] Run `go test ./geyserliteabi ./bedrockprincipal/... -count=1`.
- [ ] Run `go test -race ./bedrockprincipal -count=1`.
- [ ] Run `go test ./... -count=1`.
- [ ] Run scoped forbidden-concept and release/deployment-path scans.
- [ ] Review `git diff --check`, `git diff`, and `git status --short`.
- [ ] Commit the complete fix on `fm/bedrock-option-a-freeze-fix-connect` and append the required terminal status.

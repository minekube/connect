# Bedrock principal v2 wire status

The repository-local protocol, verifier, replay, metadata, readiness, vector,
and structural privacy work is intentionally source-only and unreleased.

The real generated-descriptor binding is complete in Moxy PR
[`minekube/moxy#514`](https://github.com/minekube/moxy/pull/514), merged as
`821608c60e3372da13aa79856353937aad92d179`. Moxy's compatibility suite now
loads the actual generated Watch and libp2p descriptors and proves:

1. the exact closed field, type, label, and oneof sets for generated `Session`,
   `SessionOffer`, Watch request/response, readiness, and negotiation messages;
2. Watch request and response payload exclusivity; and
3. absence of raw XUID, Bedrock profile, or link-provenance fields throughout
   the actual generated Watch and libp2p descriptor trees.

This document therefore no longer defers the structural privacy proof to a
future generated-module release. Updating public Connect's generated module
dependency, exposing proposal accessors/shared codecs, and implementing
consumer transport behavior remain separate downstream work. No generated
dependency version, release, endpoint, credential, or production behavior is
invented by this branch.

## Endpoint-scoped organization binding

`organization_id` remains a required signed field, but its value may be empty
for an endpoint that has no dashboard organization. Empty does not mean
unbound: verifiers require an exact match with the trusted proposal context,
and continue to require non-empty endpoint and session bindings, a valid
signature, and replay protection. Consumers must therefore pass the actual
endpoint-scoped empty value rather than substituting an organization ID.

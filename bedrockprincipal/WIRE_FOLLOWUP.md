# Bedrock principal v2 wire follow-on

The repository-local protocol, verifier, replay, metadata, readiness, vector,
and structural privacy work is intentionally source-only and unreleased.

Consuming proposal messages is a separate follow-on gated on
`bedrock-option-a-moxy-wire-source` landing and producing its reviewed immutable
Buf release manifest. Only after that gate may Connect:

1. update the generated `buf.build/minekube/connect` dependency from the
   released manifest;
2. expose `SessionProposal` bindings for fields 7–12 and the Watch/libp2p
   readiness messages;
3. run `descriptorprivacy.ValidateV2ProposalAdditions` against the actual
   generated `Session` descriptor; and
4. add transport round-trip, legacy compatibility, cross-lease, reconnect,
   malformed-frame, and readiness-consumption tests.

No generated wire type, dependency version, release, endpoint, credential, or
production behavior is invented by this branch.

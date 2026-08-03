# GeyserLite verified-ingress subprocess framing v1

This is the immutable, source-only subprocess ABI paired with
`VerifiedIngressV1`. It defines an interface; this repository does not
implement or release GeyserLite.

Gate creates a Unix `SOCK_SEQPACKET` socketpair and a 32-byte random IPC key for
each nonzero, monotonically increasing process generation. The authenticated
bootstrap is:

```text
version_u8=1 || generation_u64_be || ipc_key_32
```

The bootstrap is exactly 41 bytes. Every authenticated packet is at most 8192
bytes, carries the generation, uses a sequence beginning at one and increasing
exactly, and ends in a 32-byte HMAC-SHA256 over every preceding packet byte.
Integers are unsigned big-endian. The common authenticated prefix is exactly 26
bytes:

| Field | Offset | Size |
| --- | ---: | ---: |
| version (`1`) | 0 | 1 |
| type | 1 | 1 |
| generation | 2 | 8 |
| sequence | 10 | 8 |
| connection handle | 18 | 8 |

Assignment is type 1 and remains exactly 82 bytes:

```text
version=1 || type=1 || generation_u64_be || sequence_u64_be || connection_handle_u64_be || correlation_16 || expires_unix_ms_u64_be || HMAC-SHA256
```

Its correlation begins at offset 26, expiry at offset 42, and MAC at offset 50.

The child acknowledges with type 2 and the same generation, handle, and
correlation. ACK is exactly 75 bytes:

```text
version=1 || type=2 || generation_u64_be || sequence_u64_be || connection_handle_u64_be || correlation_16 || status_u8 || HMAC-SHA256
```

The correlation remains at its original offset 26. Status is appended at offset
42: `0` is a positive ACK and `1` is a negative ACK. The MAC begins at offset
43 and authenticates the status. A missing status, any other status value, or
any packet length other than 75 is invalid and fail-closed. A negative ACK
closes the transport and clears the pending assignment.

Verified ingress uses type 3, repeats the generation and handle, and carries the
existing bounded `VerifiedIngressV1` payload:

```text
version=1 || type=3 || generation_u64_be || sequence_u64_be || connection_handle_u64_be || verified_ingress_v1_proto_1_to_4096 || HMAC-SHA256
```

The payload begins at offset 26; because `SOCK_SEQPACKET` preserves the packet
boundary, the MAC is the final 32 bytes. The packet is 59–4154 bytes. The
correlation inside `VerifiedIngressV1` is exactly 16 bytes, the ingress frame is
1–4096 bytes, and its expiry is at most five seconds after callback.

The child publishes a newly minted, nonzero connection handle with type 4
before translation or authentication. Connection-open is exactly 58 bytes:

```text
version=1 || type=4 || generation_u64_be || sequence_u64_be || connection_handle_u64_be || HMAC-SHA256
```

Its MAC begins at offset 26. A connection-open packet carries no correlation:
after accepting the authenticated event, Gate creates the 16-byte correlation
and binds it to the exact `(generation, connection_handle)` in the subsequent
type-1 assignment. GeyserLite must not release the connection until it has sent
a positive type-2 ACK for that same generation, handle, and correlation.

MAC or peer-owner failure, sequence gap/reuse, stale generation, unknown or
duplicate/reused handle or correlation, expiry, EOF, restart, or backlog
overflow is fail-closed and clears pending state. The IPC key is never placed
in arguments, environment, disk, logs, or artifacts.

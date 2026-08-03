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

Every authenticated packet is at most 8192 bytes, carries the generation, uses
a sequence beginning at one and increasing exactly, and ends in a 32-byte
HMAC-SHA256. Integers are unsigned big-endian. Assignment is:

```text
version=1 || type=1 || generation_u64_be || sequence_u64_be || connection_handle_u64_be || correlation_16 || expires_unix_ms_u64_be || HMAC-SHA256
```

The child acknowledges with `type=2` and the same generation, handle, and
correlation. Verified ingress uses `type=3`, repeats the generation and handle,
and carries the existing bounded `VerifiedIngressV1` payload. The correlation
is exactly 16 bytes, the ingress frame is 1–4096 bytes, and its expiry is at
most five seconds after callback.

MAC or peer-owner failure, sequence gap/reuse, stale generation, unknown or
duplicate handle/correlation, expiry, EOF, restart, or backlog overflow is
fail-closed and clears pending state. The IPC key is never placed in arguments,
environment, disk, logs, or artifacts.

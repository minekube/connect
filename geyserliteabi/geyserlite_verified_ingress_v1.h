#ifndef MINEKUBE_GEYSERLITE_VERIFIED_INGRESS_V1_H
#define MINEKUBE_GEYSERLITE_VERIFIED_INGRESS_V1_H

#include <stdint.h>

#ifdef __cplusplus
extern "C" {
#endif

typedef struct __graal_isolatethread_t graal_isolatethread_t;

#define GEYSERLITE_CALLBACK_ABI_VERSION 1
#define GEYSERLITE_INGRESS_CORRELATION_BYTES 16
#define GEYSERLITE_INGRESS_FRAME_MIN_BYTES 1
#define GEYSERLITE_INGRESS_FRAME_MAX_BYTES 4096
#define GEYSERLITE_INGRESS_LIFETIME_MAX_MS 5000

#define GEYSERLITE_CALLBACK_REGISTRATION_OK 0
#define GEYSERLITE_ASSIGN_OK 0
#define GEYSERLITE_ASSIGN_UNKNOWN_OR_CLOSED_HANDLE -1
#define GEYSERLITE_ASSIGN_DUPLICATE_HANDLE_OR_CORRELATION -2
#define GEYSERLITE_ASSIGN_INVALID_OR_EXPIRED_TIME -3
#define GEYSERLITE_ASSIGN_WRONG_CONNECTION_STATE -4

/*
 * Authenticated subprocess framing v1. The bootstrap is exactly 41 bytes with
 * a 32-byte key; authenticated packets are at most 8192 bytes. All integers
 * are unsigned big-endian, and every packet ends in a 32-byte HMAC-SHA256 over
 * all preceding bytes. Connection-open is type 4 and exactly 58 bytes.
 * Assignment ACK is type 2 and exactly 75 bytes; its authenticated status byte
 * is 0 (positive) or 1 (negative). Any missing or other status is fail-closed.
 */
#define GEYSERLITE_SUBPROCESS_FRAME_VERSION 1
#define GEYSERLITE_SUBPROCESS_ASSIGNMENT 1
#define GEYSERLITE_SUBPROCESS_ASSIGNMENT_ACK 2
#define GEYSERLITE_SUBPROCESS_VERIFIED_INGRESS 3
#define GEYSERLITE_SUBPROCESS_CONNECTION_OPEN 4

#define GEYSERLITE_SUBPROCESS_ACK_POSITIVE 0
#define GEYSERLITE_SUBPROCESS_ACK_NEGATIVE 1

#define GEYSERLITE_SUBPROCESS_BOOTSTRAP_PACKET_BYTES 41
#define GEYSERLITE_SUBPROCESS_IPC_KEY_BYTES 32
#define GEYSERLITE_SUBPROCESS_MAC_BYTES 32
#define GEYSERLITE_SUBPROCESS_MAX_PACKET_BYTES 8192

#define GEYSERLITE_SUBPROCESS_VERSION_OFFSET 0
#define GEYSERLITE_SUBPROCESS_TYPE_OFFSET 1
#define GEYSERLITE_SUBPROCESS_GENERATION_OFFSET 2
#define GEYSERLITE_SUBPROCESS_SEQUENCE_OFFSET 10
#define GEYSERLITE_SUBPROCESS_CONNECTION_HANDLE_OFFSET 18
#define GEYSERLITE_SUBPROCESS_CORRELATION_OFFSET 26

#define GEYSERLITE_SUBPROCESS_ASSIGNMENT_EXPIRES_OFFSET 42
#define GEYSERLITE_SUBPROCESS_ASSIGNMENT_MAC_OFFSET 50
#define GEYSERLITE_SUBPROCESS_ASSIGNMENT_PACKET_BYTES 82

#define GEYSERLITE_SUBPROCESS_ACK_STATUS_OFFSET 42
#define GEYSERLITE_SUBPROCESS_ACK_MAC_OFFSET 43
#define GEYSERLITE_SUBPROCESS_ACK_PACKET_BYTES 75

#define GEYSERLITE_SUBPROCESS_VERIFIED_INGRESS_PAYLOAD_OFFSET 26
#define GEYSERLITE_SUBPROCESS_VERIFIED_INGRESS_MIN_PACKET_BYTES 59
#define GEYSERLITE_SUBPROCESS_VERIFIED_INGRESS_MAX_PACKET_BYTES 4154

#define GEYSERLITE_SUBPROCESS_CONNECTION_OPEN_MAC_OFFSET 26
#define GEYSERLITE_SUBPROCESS_CONNECTION_OPEN_PACKET_BYTES 58

typedef void (*gate_verified_ingress_v1_cb)(
    const uint8_t correlation[16], const uint8_t *frame,
    uint32_t frame_len, uint64_t expires_unix_ms);
typedef void (*gate_ingress_connection_open_v1_cb)(uint64_t connection_handle);

/*
 * Both callback arguments null is the only unregister operation; a partial-null
 * pair is invalid. Callback registration and unregistration are linearization barriers.
 * A verified callback may enter from a foreign native thread. Consumers must
 * validate pointers and frame_len before reading native memory, synchronously
 * copy the 16-byte correlation and 1..4096-byte frame, and retain no pointer.
 * Java/native owns callback memory after return. Any nonzero result is fail-closed.
 */
int32_t geyserlite_set_ingress_callbacks_v1(
    graal_isolatethread_t *thread, gate_ingress_connection_open_v1_cb open,
    gate_verified_ingress_v1_cb verified);
int32_t geyserlite_assign_verified_ingress_v1(
    graal_isolatethread_t *thread,
    uint64_t connection_handle, const uint8_t correlation[16],
    uint64_t expires_unix_ms);

#ifdef __cplusplus
}
#endif

#endif

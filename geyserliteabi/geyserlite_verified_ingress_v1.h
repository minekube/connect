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

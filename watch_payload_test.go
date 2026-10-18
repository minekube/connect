package connect

import (
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protowire"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// Watch payloads must be encoded on their documented wire field numbers, no
// matter which generated proto revision this module builds against: older
// revisions expose the payload as a plain optional field, the
// readiness-attestation revision as a member of a `payload` oneof. Both encode
// the same tag, so the reflected constructor and the wire contract are what
// keep the library compatible with either shape.
func TestNewWatchRequestWithRejectionWireContract(t *testing.T) {
	const id = "session-1"
	reason := status.New(codes.Aborted, "not now").Proto()

	req, err := NewWatchRequestWithRejection(&SessionRejection{Id: id, Reason: reason})
	require.NoError(t, err)

	// The accessors read the payload in every generated revision.
	require.NotNil(t, req.GetSessionRejection())
	require.Equal(t, id, req.GetSessionRejection().GetId())
	require.Equal(t, reason.GetMessage(), req.GetSessionRejection().GetReason().GetMessage())

	// The payload stays on field 1 of WatchRequest.
	fd := req.ProtoReflect().Descriptor().Fields().ByName("session_rejection")
	require.NotNil(t, fd, "WatchRequest lost its session_rejection field")
	require.Equal(t, protoreflect.FieldNumber(1), fd.Number())

	// ... and the encoded frame is exactly that field.
	rejection, err := proto.Marshal(&SessionRejection{Id: id, Reason: reason})
	require.NoError(t, err)
	want := protowire.AppendTag(nil, 1, protowire.BytesType)
	want = protowire.AppendBytes(want, rejection)

	got, err := proto.Marshal(req)
	require.NoError(t, err)
	require.Equal(t, want, got)

	// A frame survives the round trip the ws transport performs.
	decoded := new(WatchRequest)
	require.NoError(t, proto.Unmarshal(got, decoded))
	require.Equal(t, id, decoded.GetSessionRejection().GetId())
}

func TestNewWatchResponseWithSessionWireContract(t *testing.T) {
	session := &Session{Id: "session-2", TunnelServiceAddr: "127.0.0.1:25565"}

	res, err := NewWatchResponseWithSession(session)
	require.NoError(t, err)
	require.Equal(t, session.GetId(), res.GetSession().GetId())

	fd := res.ProtoReflect().Descriptor().Fields().ByName("session")
	require.NotNil(t, fd, "WatchResponse lost its session field")
	require.Equal(t, protoreflect.FieldNumber(1), fd.Number())

	encodedSession, err := proto.Marshal(session)
	require.NoError(t, err)
	want := protowire.AppendTag(nil, 1, protowire.BytesType)
	want = protowire.AppendBytes(want, encodedSession)

	got, err := proto.Marshal(res)
	require.NoError(t, err)
	require.Equal(t, want, got)

	decoded := new(WatchResponse)
	require.NoError(t, proto.Unmarshal(got, decoded))
	require.Equal(t, session.GetId(), decoded.GetSession().GetId())
}

// A payload that cannot be carried must fail loudly: silently sending an empty
// frame would look like a sessionless watch update to the peer.
func TestWatchPayloadConstructorsFailClosed(t *testing.T) {
	var nilSession *Session
	var nilRejection *SessionRejection

	_, err := NewWatchResponseWithSession(nil)
	require.Error(t, err, "nil session")
	_, err = NewWatchResponseWithSession(nilSession)
	require.Error(t, err, "typed nil session")
	_, err = NewWatchRequestWithRejection(nil)
	require.Error(t, err, "nil rejection")
	_, err = NewWatchRequestWithRejection(nilRejection)
	require.Error(t, err, "typed nil rejection")

	// A field this revision does not have, and a field of the wrong kind, are
	// reported instead of panicking.
	require.Error(t, setMessageField(new(WatchRequest), "does_not_exist", &SessionRejection{}))
	require.Error(t, setMessageField(new(WatchResponse), "session", &SessionRejection{}),
		"a field of another message type must not be set")

	// An untouched message still encodes to an empty frame: the constructors
	// are the only supported way to build a payload.
	empty, err := proto.Marshal(new(WatchRequest))
	require.NoError(t, err)
	require.Empty(t, empty)
	require.Nil(t, new(WatchRequest).GetSessionRejection())
}

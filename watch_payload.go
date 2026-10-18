package connect

import (
	"fmt"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// The watch protocol carries exactly one payload per frame: a session
// proposal on a WatchResponse and a rejection on a WatchRequest. Generated
// proto revisions disagree on how that payload is exposed to Go code:
//
//   - older revisions declare it as a plain optional message field
//     (`WatchRequest.session_rejection`, `WatchResponse.session`);
//   - the readiness-attestation revision declares it as a member of a
//     `payload` oneof, alongside new sibling members of the same oneof.
//
// Both revisions use the SAME wire field numbers, so this package must not
// depend on either Go shape: a struct literal written for one of them stops
// compiling for every consumer that pins a different generated revision (a
// library consumer's build list decides the revision via MVS, not this
// module). The constructors below set the payload through the message
// descriptor, which is stable across both shapes, and fail closed if a future
// revision renames or drops the field instead of silently sending an empty
// frame.

// NewWatchResponseWithSession returns a WatchResponse proposing session.
func NewWatchResponseWithSession(session *Session) (*WatchResponse, error) {
	res := new(WatchResponse)
	if err := setMessageField(res, "session", session); err != nil {
		return nil, err
	}
	return res, nil
}

// NewWatchRequestWithRejection returns a WatchRequest rejecting a session.
func NewWatchRequestWithRejection(rejection *SessionRejection) (*WatchRequest, error) {
	req := new(WatchRequest)
	if err := setMessageField(req, "session_rejection", rejection); err != nil {
		return nil, err
	}
	return req, nil
}

// setMessageField sets the message field name on msg. The field may be a plain
// optional field or a member of a oneof; both encode identically.
func setMessageField(msg proto.Message, name protoreflect.Name, value proto.Message) error {
	m := msg.ProtoReflect()
	fd := m.Descriptor().Fields().ByName(name)
	if fd == nil {
		return fmt.Errorf("%s has no %q field", m.Descriptor().FullName(), name)
	}
	if fd.Message() == nil {
		return fmt.Errorf("%s.%s is not a message field", m.Descriptor().FullName(), name)
	}
	if value == nil {
		return fmt.Errorf("%s.%s: payload must not be nil", m.Descriptor().FullName(), name)
	}
	vm := value.ProtoReflect()
	if !vm.IsValid() {
		return fmt.Errorf("%s.%s: payload must not be nil", m.Descriptor().FullName(), name)
	}
	if got, want := vm.Descriptor().FullName(), fd.Message().FullName(); got != want {
		return fmt.Errorf("%s.%s expects a %s payload, got %s",
			m.Descriptor().FullName(), name, want, got)
	}
	m.Set(fd, protoreflect.ValueOfMessage(vm))
	return nil
}

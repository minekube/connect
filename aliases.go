package connect

import (
	"google.golang.org/genproto/googleapis/rpc/status"

	pb "go.minekube.com/connect/internal/api/minekube/connect/v1alpha1"
)

// Type alias to better support updating versions.
// See the referenced type for documentation.
//
// Other go files should only ever use the provided
// alias type and never import a specific proto version.
type (
	Session        = pb.Session
	Authentication = pb.Authentication

	SessionRejection = pb.SessionRejection
	RejectionReason  = status.Status // The reason why a session proposal is rejected.

	Player              = pb.Player
	GameProfile         = pb.GameProfile
	GameProfileProperty = pb.GameProfileProperty

	WatchRequest  = pb.WatchRequest
	WatchResponse = pb.WatchResponse
)

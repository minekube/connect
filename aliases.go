package connect

import (
	pb "go.minekube.com/connect/internal/api/minekube/connect/v1alpha1"
	"google.golang.org/genproto/googleapis/rpc/status"
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

	WatchServiceClient       = pb.WatchServiceClient
	WatchServiceServer       = pb.WatchServiceServer
	WatchService_WatchClient = pb.WatchService_WatchClient
	WatchService_WatchServer = pb.WatchService_WatchServer
	WatchRequest             = pb.WatchRequest
	WatchResponse            = pb.WatchResponse

	TunnelServiceClient        = pb.TunnelServiceClient
	TunnelServiceServer        = pb.TunnelServiceServer
	TunnelService_TunnelClient = pb.TunnelService_TunnelClient
	TunnelService_TunnelServer = pb.TunnelService_TunnelServer
	TunnelRequest              = pb.TunnelRequest
	TunnelResponse             = pb.TunnelResponse

	UnimplementedWatchServiceServer struct {
		pb.UnimplementedWatchServiceServer
	}
	UnimplementedTunnelServiceServer struct {
		pb.UnimplementedTunnelServiceServer
	}
)

// Alias to better support updating versions.
// See the referenced type for documentation.
//
// Other go files should only ever use the provided
// alias type and never import a specific proto version.
var (
	NewWatchServiceClient  = pb.NewWatchServiceClient
	NewTunnelServiceClient = pb.NewTunnelServiceClient

	RegisterWatchServiceServer  = pb.RegisterWatchServiceServer
	RegisterTunnelServiceServer = pb.RegisterTunnelServiceServer
)

// Well-known metadata keys
const (
	MDPrefix   = "connect-"            // The prefix of metadata keys.
	MDSession  = MDPrefix + "session"  // Metadata key specifying the session id when calling the TunnelService.
	MDEndpoint = MDPrefix + "endpoint" // Metadata key specifying the endpoint when calling the WatchService.
)

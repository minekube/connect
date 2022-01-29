package connect

import (
	pb "go.minekube.com/connect/api/minekube/connect/v1alpha1"
	"google.golang.org/genproto/googleapis/rpc/status"
)

// Type alias to better support updating versions.
// See the referenced type for documentation.
type (
	Session  = pb.Session
	Endpoint = pb.Endpoint

	DenySession = pb.DenySession
	// DenyReason is the reason why a session is denied to connect.
	DenyReason = status.Status

	WatchServiceClient       = pb.WatchServiceClient
	WatchService_WatchClient = pb.WatchService_WatchClient
	WatchService_WatchServer = pb.WatchService_WatchServer
	WatchRequest             = pb.WatchRequest
	WatchResponse            = pb.WatchResponse

	TunnelService_TunnelServer = pb.TunnelService_TunnelServer
	TunnelService_TunnelClient = pb.TunnelService_TunnelClient
	TunnelServiceClient        = pb.TunnelServiceClient
	TunnelRequest              = pb.TunnelRequest
	TunnelResponse             = pb.TunnelResponse
)

// Alias to better support updating versions.
// See the referenced type for documentation.
var (
	NewWatchServiceClient  = pb.NewWatchServiceClient
	NewTunnelServiceClient = pb.NewTunnelServiceClient
)

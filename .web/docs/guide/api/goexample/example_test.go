package goexample

import (
	"context"
	"net/http"

	minekube "buf.build/gen/go/minekube/connect/bufbuild/connect-go/minekube/connect/v1alpha1/connectv1alpha1connect"
	connectpb "buf.build/gen/go/minekube/connect/protocolbuffers/go/minekube/connect/v1alpha1"

	// You can read more about Buf's Connect for Go here https://connect.build/docs/go
	"github.com/bufbuild/connect-go"
)

const (
	// This is the official Connect API endpoint.
	baseURL = "https://connect-api.minekube.com"

	// These are the headers you need to set to authenticate with the Connect API.
	endpointHeader = "Connect-Endpoint"
	tokenHeader    = "Authorization"
)

// ExampleClient_ListEndpoints shows how to list endpoints you have access to.
// It uses the default http.Client and sets the endpoint and token headers
// manually.
func ExampleClient_ListEndpoints() {
	// Set up the client.
	client := minekube.NewConnectServiceClient(http.DefaultClient, baseURL)

	// Set up a request to list endpoints you have access to.
	ctx := context.TODO()
	req := connect.NewRequest(&connectpb.ListEndpointsRequest{})
	req.Header().Set(endpointHeader, "my-endpoint")
	req.Header().Set(tokenHeader, "Bearer"+"my-token")

	// Fetch all endpoints until the server returns an empty page.
	for {
		// Send the request.
		res, err := client.ListEndpoints(ctx, req)
		if err != nil {
			panic(err)
		}

		// Print the endpoints.
		for _, endpoint := range res.Msg.GetEndpoints() {
			// Do something with the endpoint.
			println(endpoint)
		}

		// Prepare the next request.
		req.Msg.PageToken = res.Msg.GetNextPageToken()
		if req.Msg.PageToken == "" {
			// No more pages.
			break
		}
	}
}

// ExampleClient_ListEndpoints_WithHeadersTransport shows how to connect players
// to an endpoint. It uses a custom http.Client that adds the endpoint and token
// headers to every request automatically.
func ExampleClient_ConnectEndpoint_WithHeadersTransport() {
	// Set up the client.
	httpClient := &http.Client{Transport: &headersTransport{
		headers: map[string]string{
			endpointHeader: "my-endpoint",
			tokenHeader:    "my-token",
		},
	}}
	client := minekube.NewConnectServiceClient(httpClient, baseURL)

	// Set up a request to connect a players to an endpoint you have access to.
	ctx := context.TODO()
	req := connect.NewRequest(&connectpb.ConnectEndpointRequest{
		Endpoint: "my-endpoint",
		Players: []string{
			// example player uuids,
			// the players must be online and on another endpoint you have access to.
			"11111111-1111-1111-1111-111111111111",
			"22222222-2222-2222-2222-222222222222",
			"33333333-3333-3333-3333-333333333333",
		},
	})

	// Send the request.
	_, err := client.ConnectEndpoint(ctx, req)
	if err != nil {
		panic(err)
	}
}

// headersTransport is a http.RoundTripper that adds headers to requests
// before sending them so that we don't have to add them to every request
// manually.
type headersTransport struct {
	headers map[string]string
	base    http.RoundTripper
}

// RoundTrip implements http.RoundTripper. It adds the headers to the request.
func (h *headersTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	for k, v := range h.headers {
		req.Header.Add(k, v)
	}
	base := h.base
	if base == nil {
		base = http.DefaultTransport
	}
	return base.RoundTrip(req)
}

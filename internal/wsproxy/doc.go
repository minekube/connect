// Package wsproxy implements a websocket proxy for bidirectional streaming grpc-gateway backed services.
// Inspired by https://github.com/tmc/grpc-websocket-proxy but with fixed issues and performance optimizations
// by using a more efficient websocket library and using protobuf binary messages instead of JSON like done
// in https://github.com/stackrox/go-grpc-http1.
package wsproxy

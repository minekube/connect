# grpc

This package was already tested to work, but we moved over to maintaining WebSocket
firstly because of wider support by web proxies.

These are the steps to introduce the grpc package back:
- don't move this `grpc` package to root
- instead, create a new `grpc` package and delete this one
- the new `grpc` package should take package `ws` as best practice reference

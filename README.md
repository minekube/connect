# Minekube Connect

[![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/minekube/connect?sort=semver)](https://github.com/minekube/connect/releases)
[![Doc](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go)](https://pkg.go.dev/go.minekube.com/connect)
[![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/minekube/connect?logo=go)](https://golang.org/doc/devel/release.html)
[![Go Report Card](https://goreportcard.com/badge/go.minekube.com/connect)](https://goreportcard.com/report/go.minekube.com/connect)
[![test](https://github.com/minekube/connect/workflows/test/badge.svg)](https://github.com/minekube/connect/actions?query=workflow%3Atest)
[![Discord](https://img.shields.io/discord/633708750032863232?logo=discord)](https://discord.gg/6vMDqWE)

Minekube Connect allows you to connect any Minecraft server,
whether online mode, public, behind your protected home network or anywhere else in the world,
with our highly available, performant and low latency edge proxies network nearest to you.

> WIP: This project is in active development and subject to major changes!

## Features

- [x] ProtoBuf typed
- [x] Streaming transport protocols
  - [x] WebSocket support
    - equally or more efficient than gRPC
    - better web proxy support: cloudflared, nginx, ...
  - [ ] gRPC support (improved developer experience)
    - No immediate support planned, [see](internal/grpc)
- [x] Minekube Connect plugin support for:
  - [x] [Gate](https://github.com/minekube/gate)
  - [x] [Spigot/PaperMC](https://github.com/minekube/connect-java)
  - [x] [Velocity](https://github.com/minekube/connect-java)
  - [x] [BungeeCord](https://github.com/minekube/connect-java)
  - [ ] Sponge
  - [ ] Minestom
- [x] Client side service tooling in Go
- [x] Server side service tooling in Go
- [x] Client- & service-side tests implementation in Go
- [ ] Awesome documentation website

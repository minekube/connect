---
title: Changelog
description: A dated ledger of user-visible changes to Minekube Connect, Gate, the Connect plugin, GeyserLite, and Craftless. One line per change, each linked to the release or pull request that shipped it.
---

# Changelog

User-visible changes across the Minekube platform, newest first. One line per
change, each linked to the release or pull request that shipped it.

Only changes that are actually released appear here. Dependency bumps are
omitted unless they carry a security fix. If a change requires you to do
something — upgrade, reconfigure, migrate — it is also posted to Discord
[`#announcements`](https://minekube.com/discord).

## 2026-07-26

- **Gate** — Raised the packet frame cap to vanilla's 2^21-1, so oversized packets no longer disconnect players. [gate#931](https://github.com/minekube/gate/pull/931) · [v0.68.30](https://github.com/minekube/gate/releases/tag/v0.68.30)
- **Gate** — Modern Forge token is appended to the backend handshake host, fixing modern Forge backends behind Gate. [gate#934](https://github.com/minekube/gate/pull/934) · [v0.68.30](https://github.com/minekube/gate/releases/tag/v0.68.30)
- **Gate** — Startup now warns about config settings that Lite mode ignores. [gate#933](https://github.com/minekube/gate/pull/933) · [v0.68.30](https://github.com/minekube/gate/releases/tag/v0.68.30)
- **Connect plugin** — Spigot connector handles ViaVersion 5.x legacy initializer unwrapping. [connect-java#67](https://github.com/minekube/connect-java/pull/67) · [v0.12.4](https://github.com/minekube/connect-java/releases/tag/0.12.4)

## 2026-07-25

- **Gate** — gRPC updated to v1.82.1 for a security fix. Upgrade Gate to pick it up. [gate#924](https://github.com/minekube/gate/pull/924) · [v0.68.29](https://github.com/minekube/gate/releases/tag/v0.68.29)

## 2026-07-24

- **Gate** — Fixed chat acknowledgement desync across commands on Java clients. [gate#922](https://github.com/minekube/gate/pull/922) · [v0.68.28](https://github.com/minekube/gate/releases/tag/v0.68.28)
- **Connect plugin** — Velocity 4 Guice injection is supported. [connect-java#64](https://github.com/minekube/connect-java/pull/64) · [v0.12.3](https://github.com/minekube/connect-java/releases/tag/0.12.3)

## 2026-07-22

- **Connect plugin** — Fixed a Paper startup config binding collision. [connect-java#59](https://github.com/minekube/connect-java/pull/59) · [v0.12.1](https://github.com/minekube/connect-java/releases/tag/0.12.1)

## 2026-07-19

- **Gate** — Lite status cache expiry is preserved, so cached status responses expire when they should. [gate#916](https://github.com/minekube/gate/pull/916) · [v0.68.27](https://github.com/minekube/gate/releases/tag/v0.68.27)
- **Gate** — Stale Via bridge requests are skipped on retry. [gate#911](https://github.com/minekube/gate/pull/911) · [v0.68.27](https://github.com/minekube/gate/releases/tag/v0.68.27)
- **Gate** — Reduced status response log churn. [gate#910](https://github.com/minekube/gate/pull/910) · [v0.68.27](https://github.com/minekube/gate/releases/tag/v0.68.27)
- **Gate** — Status response write failures are now logged. [gate#886](https://github.com/minekube/gate/pull/886) · [v0.68.27](https://github.com/minekube/gate/releases/tag/v0.68.27)

## 2026-07-15

- **Connect plugin** — Trusted Bedrock identity sessions are supported. [connect-java#57](https://github.com/minekube/connect-java/pull/57) · [v0.12.0](https://github.com/minekube/connect-java/releases/tag/0.12.0)

## 2026-07-13

- **GeyserLite** — Added opt-in Bedrock packet tracing. [geyserlite#76](https://github.com/minekube/geyserlite/pull/76) · [v0.4.0](https://github.com/minekube/geyserlite/releases/tag/v0.4.0)

## 2026-07-12

- **Connect plugin** — Exposed a watcher health endpoint. [connect-java#29](https://github.com/minekube/connect-java/pull/29) · [v0.11.0](https://github.com/minekube/connect-java/releases/tag/0.11.0)

## 2026-07-06

- **Connect plugin** — Bedrock identity verifier metadata is loaded from a configurable `metadata-url`. [v0.10.0](https://github.com/minekube/connect-java/releases/tag/0.10.0)
- **Craftless** — Official Fabric clients now stop cleanly. [craftless#13](https://github.com/minekube/craftless/pull/13) · [v0.3.5](https://github.com/minekube/craftless/releases/tag/v0.3.5)
- **Craftless** — CLI install includes the latest official lane. [v0.3.5](https://github.com/minekube/craftless/releases/tag/v0.3.5)

## 2026-07-05

- **Connect plugin** — Bedrock identity enforcement is configurable, and invalid Bedrock identity sessions are rejected. [v0.9.0](https://github.com/minekube/connect-java/releases/tag/0.9.0)
- **Connect plugin** — Added the Bedrock identity verifier API. [v0.8.0](https://github.com/minekube/connect-java/releases/tag/0.8.0)

## 2026-07-02

- **Gate** — Config plugin messages are queued during backend login instead of being dropped. [v0.68.22](https://github.com/minekube/gate/releases/tag/v0.68.22)

## 2026-07-01

- **Craftless** — Fabric runtime support targets are reported and exposed in the loader matrix. [v0.3.0](https://github.com/minekube/craftless/releases/tag/v0.3.0)

---

This ledger starts at 2026-07-01. For anything older, see the release notes on
[GitHub](https://github.com/minekube).

---
layout: Post
title: 'Managed Bedrock Support for Minekube Connect and Gate'
category: Engineering
date: 2026-06-08
imageUrl: '/blog/managed-bedrock-support/preview.svg'
author:
  name: Minekube Team
  role: Engineering
  href: 'https://github.com/minekube'
  imageUrl: 'https://github.com/minekube.png'
---

Minecraft has two huge player worlds: Java Edition on PC, and Bedrock Edition on mobile, console, and Windows. Server owners have usually had to pick one, or run extra infrastructure to bridge them.

Today, Minekube Connect starts making that split feel a lot smaller.

If your server is connected through the [Connect Plugin](/guide/connectors/plugin) or through standard [Gate](/guide/connectors/gate), Bedrock players can join through the same Connect network entry points your Java players already use. Your endpoint name still works. Your `play.minekube.net` subdomain still works. Your custom domains still work. The tunnel still works. The difference is that the edge can now speak Bedrock too.

No Geyser plugin to install. No extra proxy to run. No separate Bedrock hosting box to keep updated.

## What changed?

Connect has always been a global Minecraft edge network: players connect to a nearby edge, Connect proposes a session to your Connector, and your Connector opens an outbound tunnel back into the network. That is why servers can be reachable without port forwarding, without exposing their real IP, and without losing the domains and endpoint names that make them discoverable.

Now that edge is protocol-aware for both Java and Bedrock.

For server owners using the Connect Plugin, this means the same plugin already available for **PaperMC/Spigot**, **Velocity**, and **BungeeCord** can expose the server to Bedrock players through Connect. For Gate users, standard Gate can now enable Bedrock support directly in its own configuration.

The practical result is simple:

- Java players continue joining as before.
- Bedrock players can join the same endpoint.
- Connect domains and custom domains route the same way.
- Updates to the Bedrock translation layer are managed for you.

## Why this is exciting

The Minecraft community has used [GeyserMC](https://github.com/GeyserMC/Geyser) for years to translate between Bedrock and Java protocols. It is an impressive project, but the operational shape has usually been very Java-shaped too: run another JVM process, keep it updated, wire it to your proxy, configure Floodgate, and make sure Bedrock protocol updates land before players start seeing "outdated proxy" messages.

That is a lot to ask from someone who just wants mobile players to join their server.

So we changed the shape of the problem.

Instead of asking every server owner to operate their own Bedrock translation stack, Connect can run it as part of the network. If you are using Connect, the translation layer becomes another managed capability of the edge, like routing, tunneling, ping handling, endpoint domains, and protection.

## The technical bit: GeyserMC, but native

The engine behind this is [geyserlite](https://github.com/minekube/geyserlite), our open source project that takes GeyserMC's Bedrock-Java translation and compiles it into a small native artifact.

That sentence hides a lot of fun machinery.

Geyserlite does not maintain a long-lived fork of GeyserMC. Instead, its build pipeline checks out the official upstream [GeyserMC/Geyser](https://github.com/GeyserMC/Geyser) repository at a pinned commit, applies a small build overlay, and compiles Geyser Standalone with GraalVM `native-image`. The result is a native binary and shared library that can be embedded into Go and Rust programs.

In other words: Geyser's protocol work stays upstream, while the Minekube side focuses on packaging it into a form that fits modern proxy infrastructure.

The resource profile is the part we are especially proud of:

| Runtime shape | Typical memory profile |
| --- | ---: |
| Geyser Standalone on the JVM | about 440 MB RSS while idle |
| Geyserlite native runtime | about 110 MB idle / 175 MB peak |

That is the difference between "run a separate service for Bedrock" and "embed Bedrock support beside the proxy".

Geyserlite ships as:

- a native CLI for self-hosting,
- a container image,
- a Go module at [`go.minekube.com/geyserlite`](https://pkg.go.dev/go.minekube.com/geyserlite),
- a Rust crate at [`geyserlite`](https://crates.io/crates/geyserlite),
- and the integration layer used by Gate.

If you are building Minecraft infrastructure in Go or Rust, you can use the same engine directly.

## A release loop that follows upstream

Bedrock updates fast. Any bridge between Bedrock and Java is only useful if it stays current.

So this is automated from the top of the stack down:

```text
GeyserMC/Geyser updates
  -> Renovate opens a geyserlite update
  -> CI rebuilds the native binary and libraries
  -> release-please publishes a geyserlite release
  -> Gate sees the new Go module version
  -> Gate releases with the updated embedded Bedrock runtime
  -> Connect proxies roll forward globally
```

On the happy path, that loop does not need a human to remember that Bedrock changed again. CI builds, smoke tests, release automation, dependency updates, and deployment automation move the update through the chain.

That matters because a managed feature is only actually managed if it also stays managed after the next Minecraft protocol bump.

## Using it through Connect

If you already use Minekube Connect through the Java plugin, there is no new Bedrock-specific plugin to install.

The Connect Plugin supports:

- PaperMC / Spigot
- Velocity
- BungeeCord

Once your server or proxy is linked to the Connect Network, Bedrock players can use the same public address style Java players use:

```text
<endpoint>.play.minekube.net
```

Custom domains configured for your endpoint continue to route to the same endpoint. Endpoint names continue to behave the same way. Bedrock support is added at the network layer rather than becoming a second setup path that server owners need to keep in their head.

Start here if you are new to Connect:

- [Quick Start](/guide/quick-start)
- [Connect Plugin guide](/guide/connectors/plugin)
- [Custom domains](/guide/domains)

## Using it in self-hosted Gate

If you self-host standard Gate, Bedrock support is open source and built in there too.

For the simplest managed setup, enable Bedrock with one line:

```yaml
config:
  bind: 0.0.0.0:25565
  onlineMode: true
  servers:
    lobby: localhost:25566
  try:
    - lobby

  bedrock: true

connect:
  enabled: true
  name: your-endpoint-name
```

With `bedrock: true`, Gate handles the Bedrock translation runtime for you. Advanced users can still override Geyser configuration, run their own external Geyser instance, or tune Floodgate paths when they need to. The simple path is intentionally small.

One note: this is for standard Gate. Gate Lite is a thinner protocol proxy designed for lightweight Java forwarding, so Bedrock translation belongs in standard Gate or in the managed Connect edge path.

Read more:

- [Gate Bedrock guide](https://gate.minekube.com/guide/bedrock/)
- [Gate Connector guide](/guide/connectors/gate)
- [Gate on GitHub](https://github.com/minekube/gate)

## Scaling the edge for Bedrock

Adding Bedrock support means the Connect edge now does more protocol work. We also used this rollout to give the Connect proxy fleet more room for ping stability and login bursts.

Across our active edge regions in Ashburn, Dallas, Los Angeles, London, Amsterdam, Frankfurt, Stockholm, Singapore, and Sao Paulo, we doubled the baseline allocation per Connect proxy:

| Before | After |
| --- | --- |
| 1 vCPU, 256 MB RAM | 2 vCPU, 512 MB RAM |

That sounds small because it is small. And that is the point.

Gate's Go runtime and geyserlite's native build make the Bedrock translation path efficient enough that adding cross-play does not require a giant sidecar service. The Connect network can become smarter without becoming bloated.

## A smarter network, not another chore

The best infrastructure features are the ones that make the network better without making every server owner become an infrastructure operator.

This is the shape we want for Connect:

- install one Connector,
- keep using your endpoint name and domains,
- let Java and Bedrock players find you,
- and let the network handle the parts that change underneath.

Bedrock support is a big step in that direction. It makes every connected server easier to reach, gives self-hosted Gate users the same open source building block, and shows what is possible when Minecraft infrastructure is treated like a real edge platform.

Try it with [Minekube Connect](/guide/quick-start), or dig into [geyserlite](https://github.com/minekube/geyserlite) if you want to see the native Bedrock engine underneath.

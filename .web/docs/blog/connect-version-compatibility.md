---
layout: Post
title: 'Java Version Compatibility Is Now Managed in Minekube Connect'
category: Engineering
date: 2026-06-23
imageUrl: '/blog/connect-version-compatibility/preview.png'
author:
  name: Robin Brämer
  role: Founder
  href: 'https://github.com/robinbraemer'
  imageUrl: 'https://github.com/robinbraemer.png'
---

Minecraft Java server networks rarely move in lockstep.

Some players update immediately. Some stay on older clients because their mods, resource packs, or favorite server still need time. Server owners also move at different speeds: one backend might already be on the newest Paper release while another still runs the version its gameplay depends on.

That gap used to be annoying for Minekube Connect users. If your backend version and player version did not line up, you either had to update everything together or run extra version-translation infrastructure yourself.

That is changing today.

Minekube Connect now handles Java version compatibility at the edge through Gate-managed Via. Players can join Connect-backed servers across supported Minecraft Java versions without each server owner setting up a separate Via proxy path.

## What changed?

Connect edge proxies now run [Gate v0.68.0](https://github.com/minekube/gate/releases/tag/v0.68.0) with [vialite v0.2.8](https://github.com/minekube/vialite/releases/tag/v0.2.8).

Gate manages Via for the Java proxy path, and Connect enables that path by default for edge traffic. When a player joins a Connect endpoint, Gate registers the temporary session backend with Via, routes the player through the translation path, and removes the backend when the session ends.

The important part is the shape: Connect still owns endpoint auth, tunnel selection, ping behavior, and encrypted backend streams. Via is now part of the managed Java edge path instead of another process that every Connect user has to operate.

## What server owners get

For Connect users, this is intentionally boring:

- older and newer Java clients can join backend servers through the Connect edge,
- no manual ViaVersion sidecar is required for the Connect proxy path,
- endpoint names and custom domains keep working the same way,
- Connect still handles endpoint auth and encrypted tunnels,
- and the Via runtime can update automatically unless a version is pinned by operators.

If you already use the [Connect Plugin](/guide/connectors/plugin) or the [Gate Connector](/guide/connectors/gate), there is no new plugin to install for this edge feature. The public Connect network has already been updated.

## Bedrock still works

This does not replace the [managed Bedrock support](/blog/managed-bedrock-support) Connect already added.

The two layers sit in different places. Bedrock translation happens first at the Connect edge, before the connection reaches the Java proxy path. By the time Via is involved, the connection is already speaking the Java protocol, so Via can focus on Java version compatibility between the player session and the backend server.

In plain terms:

```text
Bedrock player -> Connect Bedrock edge -> Gate -> managed Via -> Connect tunnel -> backend server
Java player    -> Connect Java edge    -> Gate -> managed Via -> Connect tunnel -> backend server
```

That means Bedrock players keep working through the same managed Connect path, and Java players now get managed version compatibility on the Gate/Via side of that path.

## Why this belongs at the edge

Version compatibility is one of those Minecraft infrastructure chores that looks small until every network has to repeat it.

Running Via yourself is still a great option when you want full local control, especially on a self-hosted proxy. But Connect's job is to make the common path simpler: expose your server, keep your address stable, route players through a nearby edge, and avoid asking every server owner to become an infrastructure operator.

That is why managed version compatibility belongs beside routing, tunneling, protection, ping handling, Bedrock support, and endpoint domains.

## The technical bit: Gate-managed Via

The core feature landed in [Gate v0.68.0](https://github.com/minekube/gate/releases/tag/v0.68.0). Gate can now run managed Via through [vialite](https://github.com/minekube/vialite), the native and embeddable Via runtime wrapper used by Minekube.

For static Gate configs, managed Via is straightforward: Gate starts Via and points it at the configured backend servers. Connect needed one more piece because Connect backends are dynamic. A session backend only exists after an endpoint accepts a join.

Gate now supports that runtime shape. Connect can register a temporary session backend with Gate; Gate adds it to managed Via; Via translates packets; and Gate bridges Via's backend connection back into the original Connect tunnel dialer.

In plain terms:

```text
Java player -> Connect edge / Gate -> managed Via -> Connect tunnel -> backend server
```

That lets Via see and translate the Minecraft protocol while Connect keeps controlling the tunnel and backend stream.

## Self-hosted Gate config

If you run Gate directly, you can enable the same managed Via support in your Gate config:

```yaml
config:
  # Java protocol/version compatibility powered by Via
  via:
    enabled: true
    mode: subprocess
```

Read the full Gate compatibility guide here:

- [Gate Java Version Compatibility](https://gate.minekube.com/guide/compatibility)

## Release links

- [Gate v0.68.0 release](https://github.com/minekube/gate/releases/tag/v0.68.0)
- [vialite v0.2.8 release](https://github.com/minekube/vialite/releases/tag/v0.2.8)
- [vialite on GitHub](https://github.com/minekube/vialite)

## A smarter Connect network

Connect is becoming more than a tunnel.

It already gives servers stable public addresses, custom domains, Anycast edge routing, protection, and managed Bedrock support. Java version compatibility is another step in the same direction: make the network absorb the boring protocol work so server owners can spend less time wiring infrastructure and more time running communities.

Try it with [Minekube Connect](/guide/quick-start), or dig into [Gate](https://github.com/minekube/gate) and [vialite](https://github.com/minekube/vialite) if you want to see the open source pieces underneath.

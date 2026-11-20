---
title: "Gate connector"
description: "Set up the Gate proxy as a Minekube Connect Connector: the best supported option with the most routing features, including Gate Lite mode."
---

This page shows how to use Gate as the Connector for your Minecraft server or proxy.
Gate is the best supported, most frequently updated Connector with the most capabilities.

-> New to Gate? Start with [Gate's quick start](https://gate.minekube.com/guide/quick-start): Gate is a single small binary.

## Enable Connect in Gate configuration

In the Gate configuration:

1. Enable Connect mode.
2. Choose your endpoint name.

-> See [Enabling Connect](https://gate.minekube.com/guide/connect) in the Gate docs for how to do it.

## Gate (standard)

For an introduction to Gate, see [Gate Docs](https://gate.minekube.com/guide).

As one of the most capable open source Minecraft proxies in the world, Gate has first class support for Connect with
online-mode.

Connect-routed Bedrock players are translated by the Connect edge before they reach your Gate connector. You do not
need `bedrock: true` in your Gate config for that path. Enable `bedrock: true` only when Bedrock players should connect
directly to your self-hosted Gate instance instead of through a Connect address.

-> To customize Gate, start from the [starter plugin template](https://github.com/minekube/gate-plugin-template)

## Gate Lite mode

For an introduction to Gate Lite mode, see [Gate Lite Docs](https://gate.minekube.com/guide/lite).

Gate Lite remains a thin Java protocol reverse proxy. It can sit behind Connect as a connector for Java traffic, while
Connect-routed Bedrock translation happens at the Connect edge. Use standard Gate, not Gate Lite, if you want your own
direct Bedrock listener.

### Current behavior

In Lite mode, Gate does not terminate a player login as a full Minecraft proxy. It reads enough of the Java protocol to
select a backend route, then forwards the connection, including ping and player authentication, to that backend.

When Gate Lite is used behind Connect, the Connect edge is still the public ingress point. Gate Lite only sees the
connection that Connect forwards to your connector and then pipes it to the selected backend route.

This works well when the target behind Lite does not need to run its own Mojang online-mode login for
Connect-routed players. For example, Lite can route to an offline-mode backend, or to a proxy/server setup that is
already designed to trust the connector path you configured.

### Not supported today

`Connect -> Gate Lite -> Online Mode Backend`

That topology is not supported today by changing Gate Lite configuration alone.

Gate Lite forwards the backend's authentication flow; it does not translate Connect's edge-authenticated player session
into a backend online-mode session. Connect's [AuthSession API](../auth-api) and passthrough support are still
unavailable, so an online-mode Velocity, Paper, Waterfall, or similar backend behind Gate Lite should not be expected to
accept Connect-routed players through AuthSession yet.

If you need online-mode players through Connect today, use standard Gate with Connect enabled or the Connect Java Plugin
on a supported server/proxy. Use Gate Lite behind Connect only when the backend auth model is compatible with a thin
reverse-proxy hop.

## Example setups

Two common ways to run Gate as the Connector.

<span id="_1-example-online-mode-velocity"></span>

### Online-mode Velocity behind Gate Lite

`Connect -> [ Gate Lite -> [ Velocity -> Papers ] ]`

:::danger[Not supported today]
This example depends on Connect passthrough/AuthSession support for Lite backend routes. That support is not available
yet. Use standard Gate with Connect enabled, or install the Connect Java Plugin on Velocity, for online-mode Velocity
behind Connect.
:::

> Gate -> Paper without Velocity would also work, but this topology would be useful for existing Velocity networks
> once Connect passthrough/AuthSession support exists.

- We have Velocity in online mode running on `localhost:25577` and want to use Connect through Gate.
- Because Velocity is a proxy itself, we need to use Gate Lite mode for simple reverse proxying.
- We enable `lite` in Gate config and add a wildcard route `*` to forward all traffic to Velocity at `localhost:25577`.
- We enable `connect` mode in Gate config and choose a name for our endpoint.
- This will require Connect passthrough/AuthSession support before online-mode players can join reliably through Lite.

> Velocity here could have been replaced by another proxy like BungeeCord/Waterfall or even Gate standard
> (to avoid stacking proxy layers, use [Paper behind Gate](#paper-behind-gate) instead).

<span id="_2-example-paper-online-or-offline"></span>

### Paper behind Gate

`Connect -> [ Gate -> Papers ]`

- We have Paper in online or offline mode running on `localhost:25565` and want to use Connect through Gate.
- This time we chose to not enable `lite` mode in Gate config, because we want to use Gate's full proxying capabilities
  like switching servers with `/server`.
- We enable `connect` mode in Gate config and choose a name for our endpoint.
- Done! We can now join our Paper server at `<endpoint>.play.minekube.net`, and online-mode players from the Connect
  Network can join as well thanks to Gate's first-class Connect online-mode support. This needs no port forwarding,
  no [Java plugin](plugin.md) and no [AuthSession API](../auth-api).

## Getting support

If you have questions or need help, post a support request in the **#support** forum of the
[Minekube Community Discord](https://minekube.com/discord).

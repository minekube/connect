---
description: Allow offline or cracked Java players through Connect without weakening authenticated Bedrock identity.
---

# Offline Mode

_This page explains how to allow players to join your server in Offline Mode._

## What is Offline Mode?

Offline mode is a feature of Minecraft Vanilla that allows players to join a server without needing an Internet connection
to authenticate their Minecraft account. This is useful for LAN parties or for players who don't own a Minecraft account.
There are also many public offline-mode servers that allow players to join without a paid Minecraft account.

Offline mode servers and unauthenticated players are often referred to as _cracked_ servers and players.

::: warning Java only
Connect offline mode applies to Java Edition. It does not create an unauthenticated Bedrock path. Connect-managed
Bedrock always requires a valid Microsoft/Xbox session.
:::

This is different from Connect-managed Bedrock identity. Bedrock players who join through Connect can be authenticated
with official Microsoft/Xbox Bedrock auth and still travel through a Java connector. That trusted Bedrock/XUID path does
not require allowing arbitrary offline-mode Java players.

## Joining the [Connect Browser Hub](/guide/advertising#browser-hub)

Offline Java players use either canonical cracked entry point:

- `cracked.minekube.net`
- `cracked.minekube.com`

Use `cracked.*`, not `offline.*`. From the Browser Hub, select an endpoint or run
`/browser join <your-server-name>`. The target endpoint must explicitly allow offline Java players.

## Allowing Offline Java on an Endpoint

It is possible to join the [Connect Network](/guide/#the-connect-network) without
a valid Minecraft account. By default, Connect ensures that only Mojang authenticated players can join your online mode server.

To allow unauthenticated Java players, enable the option on the connector that owns the endpoint:

::: code-group
```yaml [plugins/connect/config.yml]
allow-offline-mode-players: true
```

```yaml [Gate config.yml]
connect:
  enabled: true
  allowOfflineModePlayers: true
```
:::

The default is `false`. You do not need to make the backend publicly reachable or create a second direct connection
path. Keep the rest of the topology's normal forwarding and online-mode settings unless you intentionally operate a
fully offline network.

::: danger Offline identities are not verified
An offline Java name is self-asserted and can be impersonated. Do not grant administrative permissions by name alone.
Use an authentication plugin if your offline community needs persistent accounts.
:::

Offline-mode player connections are not encrypted between the player and the [Connect Network](/guide/#the-connect-network) edge.
Player connections are always encrypted between the Connect edge and [Connect Endpoints](/guide/#connect-endpoints) - thanks to [Connect Tunnels](/guide/tunnels).

Connect also marks its own authenticated player connections offline mode at pre-login, because those players were
already authenticated at the Connect edge. That is an internal detail of the login flow and is unrelated to this
option. It does matter for login plugins - see [Login and Auth Plugins](/guide/login-plugins).

Do not enable either offline-mode option just because a Bedrock player is joining through Connect. Use the
[Bedrock support guide](/guide/bedrock) to decide whether the player is on the Connect-managed Bedrock path, a direct
self-hosted Gate Bedrock path, or a real offline-mode Java path.

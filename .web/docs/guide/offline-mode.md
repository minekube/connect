# Offline Mode

_This page explains how to allow players to join your server in Offline Mode._

## What is Offline Mode?

Offline mode is a feature of Minecraft Vanilla that allows players to join a server without needing an Internet connection
to authenticate their Minecraft account. This is useful for LAN parties or for players who don't own a Minecraft account.
There are also many public offline-mode servers that allow players to join without a paid Minecraft account.

Offline mode servers and unauthenticated players are often referred to as _cracked_ servers and players.

## Joining the [Connect Browser Hub](/guide/advertising#browser-hub)

To join the Browser Hub as a cracked player, you can use `cracked.minekube.net` to join the [Connect Network](/guide/#the-connect-network).

## Enabling Offline Mode

It is possible to join the [Connect Network](/guide/#the-connect-network) without
a valid Minecraft account. By default, Connect ensures that only Mojang authenticated players can join your online mode server.

To allow unauthenticated players to join your offline-mode server, you have to set the
`allow-offline-mode-players` option in the [Connect Plugin](/guide/connectors/plugin) configuration
to `true`.

::: code-group
```yaml [plugins/connect/config.yml]
allow-offline-mode-players: true
```
:::

Offline-mode player connections are not encrypted between the player and the [Connect Network](/guide/#the-connect-network) edge.
Player connections are always encrypted between the Connect edge and [Connect Endpoints](/guide/#connect-endpoints) - thanks to [Connect Tunnels](/guide/tunnels).

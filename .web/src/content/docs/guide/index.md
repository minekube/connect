---
title: "What is Minekube Connect?"
description: "What Minekube Connect is: a DDoS-protected edge network for Minecraft with free public domains, Connectors, the Connect Dashboard and the Browser Hub."
---

Minekube Connect gives any Minecraft server or proxy a public address and DDoS protection, without port forwarding.
You run a small **Connector** next to your server. Players join through the Connect edge network, which carries them to you.

It works the same for a server on your own PC, a VPS, or a whole network behind Gate, Velocity or BungeeCord.
The core features are free and will stay free.

**Start here:** [Quick start](/guide/quick-start) · [Open the dashboard](https://app.minekube.com/)

## How it works

1. **Install a Connector.** The Connect plugin for Paper, Spigot, Velocity or BungeeCord, or Gate with Connect switched on.
2. **Your server links to the edge.** The Connector opens an outbound tunnel to the nearest Connect edge and gets a name, for example `live-beru`.
3. **Players join by name.** They type `live-beru.play.minekube.net` or your own domain. Connect filters attacks and routes them through the tunnel to your server.

The tunnel is outbound, so no ports need to be opened and your server's IP address is never published.

## The Connect Platform

| Part | What it does |
| --- | --- |
| [Connect Network](#the-connect-network) | Global edge proxies that accept players and run the tunnels |
| [Connectors](/guide/connectors/) | The plugin or Gate integration that links your server to the network |
| [Connect Dashboard](#the-connect-dashboard) | The web app for your endpoints, connectors, domains and statistics |
| [Connect Browser](#the-connect-browser) | The in-game server list where players discover servers |

## The Connect Network

The Connect Network is a set of highly available edge proxies in [locations around the world](/guide/locations).
Every running Connector keeps a tunnel open to it, so players can reach your server from anywhere with little overhead.

The network is responsible for:

- Running the [tunnels](/guide/tunnels) between players and your servers
- Giving every server a free public [domain](/guide/domains), including servers on [localhost](/guide/localhost)
- Filtering [attacks and bots](/guide/protections) before they reach you
- [Advertising](/guide/advertising) your server to players, for example in the [Browser Hub](/guide/advertising#browser-hub)

## The Connect Dashboard

The [Connect Dashboard](https://app.minekube.com/) is where you manage your endpoints, connectors and custom domains,
and see statistics about your network.

## The Connect Browser

The Connect Browser is how players discover servers on the network.
Most notably it is the in-game server browser: players open it with the global `/browser` command or by joining `minekube.net`.
See the [advertising guide](/guide/advertising) to get listed.

## Connect endpoints

An endpoint is a Minecraft server or proxy linked to the Connect Network through a [Connector](/guide/connectors/).
Each endpoint has a globally unique, human-readable name, and the docs often just call it a _server_.

If you do not set a name in the Connector's configuration, it picks one from the [Random Name Service](https://randomname.minekube.net/).

## Connect tunnels

As soon as your server starts, its Connector links with the Connect Network and players can join,
even when the server runs [on your own PC](/guide/localhost).
[How tunnels work](/guide/tunnels) explains the connection in detail.

## Connect sessions

Connect uses two kinds of sessions:

- **Endpoint sessions** run from your Connector to the Connect edge and wait for incoming players.
- **Player sessions** run from the Connect edge to an endpoint with an active Connector, one per joining player.

## Next steps

- [Quick start](/guide/quick-start): link your first server in a few minutes
- [Connectors](/guide/connectors/): choose between the Java plugin and Gate
- [Host from your own PC](/guide/localhost): let friends join a server running at home
- [Use a custom domain](/guide/domains): put `play.yourserver.com` in front of your endpoint
- [Bedrock players](/guide/bedrock): let Bedrock players join on the same address

# Introduction

_Connect is a platform for developers, that allows you to
connect Minecraft (Java Edition) servers and proxies very easily, and
build awesome Minecraft networks._

---

## The Connect platform

The Connect platform consists of the following components:

- [The Connect Plugin](#the-connect-plugin) - A plugin that links your servers with the Connect Network
- [The Connect Network](#the-connect-network) - A global network of Minecraft endpoints
- [The Connect API](#the-connect-api) - A public API for developers
- [The Connect Browser](#the-connect-browser) - Server discovery for players

## The Connect Plugin

The Connect Plugin is a powerful multi-platform Minecraft plugin that links
your servers with the global [Connect Network](#the-connect-network).

It supports PaperMC, Gate, BungeeCord and Velocity platforms.

Your servers are automatically advertised to players on the open network.
Players can discover your server from the in-game server browser while you
stay in control of your Minecraft servers/proxies and hosting location.

Connect also makes localhost servers publicly joinable and provides a
free domain for every server. _Feels like magic!_

You may already have questions - don't worry. We will cover every little detail in the rest of the documentation.
For now, please read along, so you can have a high-level understanding of what the Connect platform offers.

## The Connect Network

The Connect Network consists of highly available, scalable and self-healing edge proxies
built on top of best-of-breed cloud native technologies and securely interconnects all active
Minecraft server and proxy endpoints through the [Connect Plugin](#the-connect-plugin)
without noticeable overhead.

The **Connect Network** is vital to the **Connect platform** and responsible for:

- Establishing [Tunnels](/guide/tunnels) between players and your servers
- Providing free public [domains](/guide/domains) to your public or [localhost](/guide/localhost) servers
- [Advertising](/guide/advertising) your servers to players through multiple channels like the [Browser Hub](/guide/advertising#browser-hub)
- Providing features to the public [Connect API](/guide/api/) for developers

Connect can replace traditional proxies like BungeeCord or Velocity and
become the largest shared Minecraft network in the world.

_Core features of Connect are free and will always be free for everyone._

## The Connect API

The public Connect API allows developers to build awesome
Minecraft networks and applications on top of the Connect platform.

It powers the [Connect Browser](#the-connect-browser) and is used
to move players between endpoints and retrieve server information.

Check out the [Developers API guide](/guide/api/) to learn more!

## The Connect Browser

The Connect Browser encompasses all the features that allow players to discover
and join endpoints on the Connect Network. **Most notably, it's the in-game server
browser that players can access with the `/browser` global command or by joining `minekube.net`.**

Check out the [Advertising guide](/guide/advertising) to learn more!

## Connect Endpoints

A Connect endpoint is a Minecraft server or proxy that is linked with the [Connect Network](#the-connect-network)
through the [Connect Plugin](#the-connect-plugin) and is identified by a unique name like.

Endpoints are also referred to as _servers_ for simplicity.

Endpoints are advertised to players and can be joined
by using the [Connect Browser](#the-connect-browser) or moving players with the [Connect API](#the-connect-api).
Any server that can serve Minecraft clients and is linked with the Connect Network can be an Endpoint.

If you do not specify the Endpoint name in your [Connect Plugin](#the-connect-plugin) configuration,
the plugin will ask the [Random Name Service](https://randomname.minekube.net/).

**Related Guides:**
- [Joining guide](/guide/joining) - Learn how to join endpoints.
- [Advertising guide](/guide/advertising) - Learn about advertising your endpoints.
- [Endpoint Domains guide](/guide/domains) - Learn how to use free domains for your endpoints.
- [Offline Mode guide](/guide/offline-mode) - Learn how to allow offline mode players on your server.


## Connect Tunnels

As soon as your server is started the [Connect Plugin](/guide/#the-connect-plugin) links with
Connect networking services and players can start joining your Minecraft server even if it's running
[locally on your PC](/guide/localhost).

Check out [How Tunnel connections work!](/guide/tunnels) for a technical explanation.

## Connect Sessions

_This is a placeholder for the Connect Sessions. Please check back later._

[//]: # (- TODO document player sessions & endpoint sessions)

## Let's Speedrun the Quick Start!

You have a Minecraft server locally or somewhere else?
Let's get started and link it with the **Connect Network** for the first time!
Just click the `Next page` button or click [Quick Start](quick-start)!

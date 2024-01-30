# Minekube Connect Introduction

-> The Ingress Tunnel for Minecraft Servers

Connect is a global gateway to a DDoS-protected Minecraft native Superproxy.
A developer platform with public Anycast domains for localhost allowing admins
to secure and scale servers/proxies to build durable Minecraft networks in
any environment easily.

_Sounds like techno marketing to you? Let's break the promises down, shall we?_

-> No time to read! Get me started! [Quick Start](quick-start) / [Dashboard](https://app.minekube.com/)

## The Connect Platform

The Connect Platform consists of the following components:

- [The Connect Network](#the-connect-network) - A global edge network of router proxies
- [The Connectors](/guide/connectors/) - Tunnels players <-> to servers through Connect Network
- [The Connect Dashboard](#the-connect-dashboard) - A web interface for organizing your Minecraft network
- [The Connect API](#the-connect-api) - A public API for developers

[//]: # (- [The Connect Browser]&#40;#the-connect-browser&#41; - Server discovery for players)

## The Connect Network

The Connect Network consists of highly available, scalable and self-healing edge proxies
built on top of best-of-breed cloud native technologies and securely interconnects all active
Minecraft server and proxy endpoints through [Connectors](/guide/connectors/)
without noticeable overhead.

The **Connect Network** is vital to the **Connect platform** and responsible for:

- Establishing [Tunnels](/guide/tunnels) between players and your servers
- Providing free public [domains](/guide/domains) to your public or [localhost](/guide/localhost) servers
- [Advertising](/guide/advertising) your servers to players through multiple channels like
  the [Browser Hub](/guide/advertising#browser-hub)
- Providing features to the public [Connect API](/guide/api/) for developers

Connect can replace traditional proxies like BungeeCord or Velocity and
become the largest shared Minecraft network in the world.

_Core features of Connect are free and will always be free for everyone._

## The Connect Dashboard

The Connect Dashboard is a web interface for managing all
resources of your Minecraft network in one place.

-> https://app.minekube.com/

It allows you to manage your Endpoints, Connectors, Custom Domains,
analyze statistics about your network, and much more.

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
through a [Connector](/guide/connectors/) and is identified by a globally unique human-readable name.

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

As soon as your server is started [Connectors](/guide/connectors/) link with
Connect networking services and players can start joining your Minecraft server even if it's running
[locally on your PC](/guide/localhost).

Check out [How Tunnel connections work!](/guide/tunnels) for a technical explanation.

## Connect Sessions

Connect defines two specific types of sessions:

- **Endpoint Sessions**: These sessions are established by your [Connector](/guide/connectors/)'s Minecraft server/proxy
  to the Connect edge proxy to watch for incoming player sessions.
- **Player Sessions**: These sessions are established by Connect edge proxies between players and available Endpoints
  that have an active Connector.

## Let's Speedrun the Quick Start!

You have a Minecraft server locally or somewhere else?
Let's get started and link it with the **Connect Network** for the first time!
Just click the `Next page` button or click [Quick Start](quick-start)!

## The Connect Plugin

Moved to [Connectors / Plugin](connectors/plugin).
# About Connect Tunnels
<VPBadge>Networking knowledge recommended</VPBadge>

_Don't believe in magic? This page explains how [Connect Tunnels](/guide/#connect-tunnels) enable player connections beyond private networks._

[[toc]]

## Introduction

Publicly routable IP addresses are required to allow inbound connections from the Internet (outside your home network).
Your ISP (Internet Service Provider) assigns you a public IP address, but sometimes it's dynamic and changes over time.

Further, your router has a NAT (Network Address Translation) that allows multiple devices to share the same
public IP address in your home network. This is why you can have multiple devices connected to the Internet
at the same time, their private IP addresses are translated to your public IP address when communicating with the Internet.

**However, your home firewall blocks inbound connections by default, including Minecraft connections on port `25565`.**

## Connect to the Rescue

Connect bypasses this whole problem by establishing outbound only tunnel connections from your server
to publicly routable edge proxies in the [Connect Network](/guide/#the-connect-network).
This way, you don't have to configure port forwarding on your router.

Connect edge proxies accept inbound connections from players and routes them through encrypted tunnels to your server.
The [Connectors](/guide/connectors/) establish a tunnel connection to the edge proxy for every player that joins your
server.

Attentive readers are now wondering how the Connector knows when a player wants to join and from which edge proxy.
Connect has two transport paths for that:

- **libp2p transport:** modern Connectors register a peer identity with the Connect Edge Network. The edge opens a
  short-lived libp2p session stream to the Connector for the player join or status request.
- **WatchService fallback:** older Connectors continue to watch for player session proposals and then open the existing
  websocket tunnel to the correct edge proxy.

The fallback remains part of Connect so existing servers keep working. The libp2p path is the newer default for updated
Connectors because the edge can initiate the session stream directly and avoid the extra proposal-watch round trip.

## The Join Flow

This is how the join flow looks for every player:

1. A player connects to a nearby Connect edge proxy.
2. The edge resolves the target endpoint and looks for a live libp2p Connector peer.
3. If a libp2p peer is available, the edge opens a session stream to that Connector and the Connector attaches the player
   session to your server.
4. If libp2p is unavailable or the Connector is older, Connect falls back to the WatchService proposal flow.

This flow is quite fast and happens within a few milliseconds.
The player doesn't notice any difference and can play on the performant tunnel to your server as usual.

For updated Connectors, server-list pings can also use the libp2p status projection instead of opening a full player-like
tunnel just to ask the backend server for its MOTD. This reduces noisy reconnect-style messages in Minecraft server
consoles and gives Connect a cleaner path for future latency-aware routing.

## Connectors

Check out the [Connectors Overview](/guide/connectors/) to learn more about Connectors.

---

**If you find these concepts intimidating, don't worry! The [Quick Start](/guide/quick-start) only requires basic
knowledge, and you should be able to follow along without being an expert in any of these.**

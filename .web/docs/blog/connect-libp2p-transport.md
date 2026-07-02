---
layout: Post
title: 'A Faster Tunnel Foundation for Minekube Connect'
category: Engineering
date: 2026-06-19
imageUrl: '/blog/connect-libp2p-transport/preview.jpeg'
author:
  name: Robin Brämer
  role: Founder
  href: 'https://github.com/robinbraemer'
  imageUrl: 'https://github.com/robinbraemer.png'
---

Minekube Connect has always had a simple promise: install a Connector, keep your server behind your own firewall, and let
players join through the global Connect edge.

Today we are laying down a new transport foundation underneath that promise.

Updated Connectors can now register as libp2p peers, which lets the Connect edge open session streams directly to the
Connector when a player joins or when the network needs fresh status. The old WatchService path stays available for
compatibility, but it no longer has to be the hot path for modern Connectors.

That sounds like a small implementation detail. It is not. It changes how quickly Connect can attach sessions, how much
console noise server owners see, and what kind of routing intelligence we can build next.

## The old shape

The original Connect tunnel architecture was intentionally conservative.

1. A Connector kept a WatchService connection open.
2. A player joined through a Connect edge proxy.
3. The edge proposed a session through the WatchService.
4. The Connector noticed the proposal and opened a websocket tunnel back to the edge.
5. The edge attached the player to that tunnel.

This model made Connect useful on day one. It worked from home networks, private machines, rented servers, school
networks, and places where inbound Minecraft ports were not an option.

It also had a cost: every join and many status checks needed a proposal-watch step before the real tunnel could exist.
For some servers that showed up as reconnect-looking messages in Minecraft logs, especially when pings needed to ask the
backend for a fresh MOTD. Nothing was necessarily broken, but it was noisy and harder to reason about than it should be.

## The new shape

With the new libp2p transport, an updated Connector has its own peer identity and registers live reachability with the
Connect Edge Network.

When a player joins, the edge can do this instead:

1. Resolve the endpoint.
2. Find a live Connector peer.
3. Open a libp2p session stream to that peer.
4. Carry the Minecraft session over that stream.

The edge starts the session. The Connector answers it. The compatibility WatchService path remains as a fallback for old
Connectors or temporary transport failures.

| Area | WatchService path | libp2p path |
| --- | --- | --- |
| Session setup | Proposal first, tunnel second | Edge opens the session stream directly |
| Status checks | Often looked like short tunnel activity | Can use live status projection |
| Connector identity | Endpoint token plus active watch session | Endpoint token plus per-Connector peer identity |
| Compatibility | Works for every existing Connector | Preferred for updated Connectors, fallback remains |
| Future routing | Limited by proposal flow | Ready for direct, relay, and hole-punched paths |

## Why libp2p?

Minecraft traffic is latency-sensitive, bursty, and very intolerant of mysterious pauses. A good Connect transport needs
to be fast on happy paths and still forgiving on weird networks.

libp2p gives us a toolbox built for that shape:

- peer identities instead of anonymous tunnel sockets,
- multiplexed streams for short control work and longer sessions,
- circuit relay support for networks that cannot be reached directly,
- DCUtR and hole punching as a path toward more direct connections,
- transport observability that tells us whether a session used direct or relayed connectivity.

The important part is not that Connect uses a fashionable networking library. The important part is that Connect can now
model every Connector as a real network peer and let the edge choose a better path.

## Less noise, cleaner pings

One of the most practical changes is status handling.

Server-list pings are constant. Every Minecraft server owner has seen them. In the older architecture, fresh ping data
could require a short tunnel-like interaction with the backend. On some platforms that made the console look busier than
the server actually was.

With the libp2p transport, updated Connectors can publish short-lived status information to the Connect edge. The edge can
answer many pings from that live projection, and only fall back to the heavier path when it needs to.

This is the kind of change that should feel boring when it works: fewer strange reconnect-looking messages, faster status
answers, and less work for the server just because someone refreshed a multiplayer list.

## A new identity file, not a new setup burden

The Connect Plugin still uses `plugins/connect/token.json` to prove ownership of your endpoint name.

The libp2p transport adds one more local file:

```text
plugins/connect/libp2p-identity.key
```

It is generated automatically. It is not a second token to paste into a dashboard. It is the private peer identity for
that one running Connector instance.

That distinction matters for load balancing. You can run multiple Connectors for the same endpoint and let Connect
distribute traffic between them, but each running Connector should keep its own libp2p identity. Sharing the endpoint
token is normal. Sharing the peer identity key between concurrently running Connectors is not.

## Compatibility is still part of the design

Minecraft server owners do not all update plugins at the same speed, and they should not need to.

The WatchService tunnel architecture remains in Connect. Existing Connectors keep working. Newer Connectors get the
libp2p transport first, and the edge can fall back when needed.

That lets us move the network forward without breaking servers that are simply running quietly and doing their job.

## What this enables next

The immediate win is a cleaner tunnel path. The more interesting win is what this lets us build afterward.

Once Connectors are live peers, the edge can make better choices:

- prefer already-warm paths,
- understand direct versus relayed sessions,
- prepare streams before the player is fully attached,
- collect better latency and failure signals,
- route around unhealthy paths faster,
- and eventually use hole punching where the network allows it.

This is the same direction Connect has been moving with managed Bedrock support: the edge network becomes smarter, while
server owners keep the setup small.

Install a Connector. Keep your endpoint name. Keep your custom domains. Let the network keep getting better underneath.

Read more:

- [Connectors overview](/guide/connectors/)
- [About Connect tunnels](/guide/tunnels)
- [Authentication and endpoint identity](/guide/api/authentication)

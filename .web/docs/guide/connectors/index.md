# Connectors Overview

Learn about Connectors, their function, and available options here.

## What is a Connector?

A Connector at Minekube Connect is a software component that facilitates the secure
communication between your Minecraft server/proxy and the Connect Edge Network by
creating secure outbound tunnels for receiving player connections.

## Available Connectors

**[Gate Connector](gate.md):** <VPBadge type="info">Recommended</VPBadge>

- Enable `connect` mode in the configuration to use it as a Connector for your
  server. Gate is updated most frequently and has the most capabilities.
- If you have an existing Java proxy, switch to [Gate Lite mode](gate.md#gate-lite-mode) to use it as a Connector
  without installing the Connect Java Plugin.

-> Continue with [Gate Connector Setup Guide](gate.md)

**[Java Plugin](plugin.md):** <VPBadge type="info">Fast set-up</VPBadge>

- The Spigot/Velocity/Bungee Connect Plugin can be installed on your Minecraft server/proxy to
  use it as a Connector for your endpoints.

-> Continue wit [Java Plugin Setup Guide](plugin.md)

## Advantages using a Connector

- You want to protect and hide your public IP address from attackers? Use a Connector.
- You don't have a public IP address, like running your server on a private network or at home? Use a Connector.

### Anycast Public IP

In any case, Connect gives your endpoints a shared public Anycast IP address reachable from
anywhere to route your player traffic through the nearest Connect Edge region to your Minecraft server.

## How does a Connector work?

The operation of a Connector can be simplified into three steps:

1. **Receiving the Player Session Proposal**: The Connector starts its operation when it receives a player session
   proposal from the Edge.

2. **Establishing the Outbound Tunnel Connection**: Upon receiving the session proposal, the Connector establishes an
   outbound tunnel connection back to the Edge.

3. **Forwarding the Player's Connection**: This outbound tunnel connection is then used by the Connector to forward the
   player's connection to the Minecraft server.

![https://mermaid.live/edit#pako:eNp1kk1uAyEMha-CvE4vwCKb_uwqVcqWjQtWijqDqYFWUZS71xMyaStSVpb53tMz5gieA4GFQh-NkqeHiHvB2SWjJ6PU6GPGVM3LhAeSsf8Y9jR27zkl8pVvCHYkn6tRN73bbhcXu6pM5V-2S6XE1dGqQwomn6WmUCmRk8nCmQtOXXOFf6yFsJLhVl-5qbw2JSbjO6gOo7AHteaJ5QslDGy__xttZYVK5lTovzi3ucuo_VVGCDYwk8wYgy7suEgc1DeayYHVMqC8O3DppBy2yrtD8mCrNNpAy0HHvyx3bVKImum5f4DzPzh9A86ity4](diagram.svg)

-> Why all that? Checkout [Connect Tunnels](/guide/tunnels) explained!

## Load Balancing multiple Connectors

You can run multiple Connectors for the same endpoint to load balance the player traffic between them.
The Connect Edge will automatically distribute the player traffic between the available endpoint's Connectors randomly.

::: tip Regional Load Balancing

Note that the Connect Edge will **currently** not take into account the accurate load or regional latency of the
Connectors when
distributing the player connections. If you want this behavior, let us know in
our [Discord](https://minekube.com/discord)!

::: 
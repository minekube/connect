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

To better understand how a Connector works, let's consider an example. Imagine a player named Alex from the US wants to
join a Minecraft server running in Singapore.

1. **Initial Connection to Nearest Edge**: Alex's request to join the server first reaches the closest Connect Edge,
   located in Ashburn. This step ensures Alex's connection benefits from an efficient data center network
   infrastructure, offering faster routing compared to a potentially slower ISP route. This optimized network path
   reduces latency from the outset.

2. **Session Proposal to Connector**: The Edge in Ashburn then sends a session proposal to the Connector that's
   associated with the target server endpoint in Singapore.

3. **Connector Connects to Local Edge**: Upon receiving the proposal, the Connector in Singapore connects to its own
   nearest Connect Edge in Singapore. This local Edge connection optimizes latency by ensuring server outbound
   connections use the most efficient network paths.

4. **Secure Tunnel Established Through Both Edges**: A secure tunnel is created from the Connector in Singapore to Alex
   through both the Singapore and Ashburn Edges. This tunnel ensures Alex's connection to the Minecraft server is secure
   and efficient, minimizing latency and enhancing the gaming experience.

This process, utilizing both local and player-nearby Edge servers, ensures a high-quality connection for players
globally, effectively reducing latency and securing the data path.

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
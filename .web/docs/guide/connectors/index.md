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
   located in Ashburn, for example. This step ensures Alex's connection benefits from an efficient data center network
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

![https://mermaid.live/edit#pako:eNptkrtuwzAMRX9F0BwvHjV46WMrWiAdvbASmwq1KYeSWgRB_r2ymcgpEg9-yIfUuTaP2gaH2uiI-4xk8dHDjmHsSZVjAk7e-gkoqbcBDsi3609uhwqiXD24-0RbkejpFnkIRGhTuNN_i_yDvJSTm4KnJIz4NF03dzWXDioFRQi87CWgnOfnM9ya0pScihijD6QmDlOIMKxgW8iqZNRz4F_gewVSUtGm6US3-hRcmNVkpavOa04fIRellMu7QRJgTSBKzTnpuzD2aoOriE0nH-a_gd7oEXkE78qvPs54r9MXjthrU24d8HevezoVDnIK2wNZbRJn3Og8OUiXsdDmE4ZYVtH5kuBFZmcZodMfAHfFrw](diagram.svg)

::: info When the Player and Connector are in the same region
Note that if both the Player and Connector are in the same region, they will likely be routed to the same Edge, thus
the Connector will create the Tunnel directly to the Edge the Player is connected. Only one Edge would be involved in the diagram.
::: 

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
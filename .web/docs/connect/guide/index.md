# Introduction

_Connect is an ecosystem that makes it very easy to
build awesome and future-proof Minecraft server networks._

## What is Connect?

Connect is a powerful Minecraft plugin that links servers with the global Connect Network.
It allows server owners to advertise their server through multiple channels
and get more players that are already playing on the shared open network.
Players can discover your server from the in-game server list while you are the complete
owner of your Minecraft server and decide where to host it.

Connect also makes localhost servers running on your PC publicly joinable and provides a
free domain for every server. _Feels like magic, you should try the quick start!_

You may already have questions - don't worry. We will cover every little detail in the rest of the documentation.
For now, please read along, so you can have a high-level understanding of what Connect offers.

## The Connect Network

Connect is also a highly available, scalable and self-healing network built on top of
best-of-breed cloud native technologies and securely interconnects all running Minecraft servers and proxies
that are using the **Connect plugin** without noticeable overhead.

The **Connect Network** is vital for the **Connect plugin** and is responsible for the following:

- Establishing a tunnel connection between players and servers
- Providing a free public domain for servers ([localhost servers too!](#public-localhost-servers))
- Advertising servers to players through multiple channels like the [Browser Hub](advertising-your-server#browser-hub)
- Providing public APIs for developers to integrate Connect into their own projects

Connect has the potential to obsolete traditional proxies like BungeeCord or Velocity and
become the largest shared Minecraft server network in the world.

_Core features of Connect are free and will always be free for everyone.
However, in the future we also offer premium features like [multiple custom domains](custom-domains),
[server advertising boosts](advertising-your-server) and [DDoS & Bot protection](protections)._

## Public Localhost Servers

_Host Minecraft servers anywhere!_

Running a localhost Minecraft server is awesome because it's free and quickly set up.
The problem is only players in your local network can join. This is where Connect comes in
and makes your localhost servers publicly joinable through the Internet.

Yes! As soon as your server is started the Connect Plugin links with the the Connect Network
and players can start joining your Minecraft server even if it's running locally
on your PC without configuring any port forwarding.


::: details Technical Details

Public localhost severs are possible because players connect to a public edge proxy in the
Connect Network that then establishes an output only tunnel session to your listing Minecraft server.

In fact, the server join flow looks the same for every server and works as follows:
1. Connect plugin watches for player session proposals from the Connect Network
2. Player pings and joins Connect Network at e.g. `your-server-name.play.minekube.net`
3. The edge proxy holding the player connection proposes a session to your server
4. Your server accepts session and establishes a tunnel to the edge proxy to take the player connection

Outbound tunnel connections are always established from the server to an edge proxy,
though players are allowed to join to the target server directly or through the Connect Network.

**If you find these concepts intimidating, don't worry! The tutorial and guide only require basic
knowledge, and you should be able to follow along without being an expert in any of these.**

:::

## Let's Speedrun the Quick Start!

You have a Minecraft server locally or somewhere else?
Let's get started and link it with the **Connect Network** for the first time!
Just click the `Next page` button or click [Quick Start](quick-start)!
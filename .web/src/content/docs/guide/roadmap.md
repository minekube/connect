---
title: "Roadmap"
description: "The Minekube Connect development roadmap: what is being built next for server owners, developers and players."
---

This page shows where Minekube is heading and which Connect features are planned next. It is a living document and
will be updated as Connect progresses.

## Where Minekube is heading

Connect runs our open source Gate proxy as an edge network for Minecraft. It is the first product; these are next.

| Product | Status | What it does |
| --- | --- | --- |
| Minekube Connect | Live now | You run the Connector, we run the edge network behind it. |
| Minekube Games | Coming soon | Bring the game, we bring the players, matchmaking and capacity. |
| Minekube Analytics | Coming soon | Learn who your players are and what they want to play next. |
| Minekube Developer Platform | Coming soon | Server search, reviews and votes, a marketplace and ads. |

General-purpose CDNs don't carry Minecraft connections for free. Connect does, with Gate deployed at every edge:

- **Edge tunnels:** protected entry points close to your players, wherever they are.
- **Anycast addresses:** one public IP that routes every player to the nearest edge.
- **Quick to set up:** one binary, container or plugin and a few lines of config. Any network layout works.

## Project roadmap

We are just getting started with the open Connect platform, and we are excited to see what the future holds.
We are working on many ideas to extend Connect and make it even more useful for creators and players.

Our major goal is to keep Connect a stable and reliable service that is
always on, always developing and the core features free to use for everyone.

## Additional components for the Connect Platform
<Badge type="caution">These are Drafts and not yet available</Badge>

Beyond what is already available (see the [introduction](/guide/)), these are ideas for extending Connect.

### Search

- Find the top-rated game servers and quickly join them
- See the history of servers you were connected on to join back later
- Search for any game server registered with Connect Network and filter by rating, category, tags, and more

### Reviews

Players can rate and review game servers they played on for a certain amount of time, and get rewarded for it.
This helps other players to find the best game servers when browsing server categories.

### Maps

**As a player:**
You walk through the world and explore previews of every game server.
You like what you see? Jump into its portal, and you get connected to the server in a blink of an eye!

**As a server creator:**
You design the gate of your game server to make it appealing to players walking by.
You get an area, limited only by the Minecraft blocks available and your imagination,
to capture the spirit of your games and advertise your server to new players.

### Games

Developers can create their own games and integrate them with APIs in the Connect Network.
Players are able to browse, play and rate these game servers. Developers host their games
but will be able to let Connect do the deployment and scaling for them to focus on game development only.

### Cloud

The Minekube Cloud will be a hosting platform that allows players and creators to create and manage
servers on-demand and pay only for the time they use it. Players are able to control their servers
and resources in-game from the Connect Hub and Dashboard.

## Other ideas

- We would like to support [Skript](https://forums.skunity.com/resources/skript.323/) with an addon
  to easily integrate Connect into your Skript scripts and manage your servers and players without needing to
  maintain a separate proxy.

---

**There are plenty more ideas to come. Stay tuned, and join the [Discord](https://minekube.com/discord) to talk about yours!**

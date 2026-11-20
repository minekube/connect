---
title: "Host from your own PC"
description: "Host a Minecraft server on your own PC and let friends join it publicly through Minekube Connect, without port forwarding or exposing your IP."
---

Host a Minecraft server on your own PC and let friends anywhere join it, without port forwarding.

## Use cases

Running a localhost Minecraft server is awesome because it's free and quickly set up.

- **Local Development** - Run your Minecraft server locally on your PC and test with friends.
- **Local Multiplayer** - Play with friends beyond your local network without port forwarding.
- **Local Hosting** - Host a Minecraft server on your local PC and make it publicly joinable.
- **Local Testing** - Test your Minecraft server with players from all over the world.

## Why Connect for a localhost server?

- You want to host a free Minecraft server without paying a provider.
- You don't want to mess with your router settings and configure port forwarding.
- Or, might not even have the admin password to your router.
- You want to let players join from anywhere through the Internet.

## How it works

Without help, only players on your local network can join a server running on your PC.
Connect makes it publicly joinable: as soon as the server starts, the Connect plugin opens an outbound tunnel to the
Connect Network, and players anywhere join through it. Your router stays as it is.

<div class="mk-arch-wrap not-content">
<div class="mk-arch" role="img" aria-label="Friends anywhere join the Connect edge, which tunnels them to the Minecraft server on your PC; your router needs no port forwarding.">
<div class="a-node"><svg viewBox="0 0 24 24" aria-hidden="true" fill="none" stroke="currentColor" stroke-width="1.75" stroke-linecap="round" stroke-linejoin="round"><path d="M16 21v-2a4 4 0 0 0-4-4H6a4 4 0 0 0-4 4v2"/><circle cx="9" cy="7" r="4"/><path d="M22 21v-2a4 4 0 0 0-3-3.87"/><path d="M16 3.13a4 4 0 0 1 0 7.75"/></svg><strong>Friends</strong><small>anywhere online</small></div>
<div class="a-link"><span>join</span></div>
<div class="a-node a-edge"><img src="/minekube-logo.png" alt="" width="28" height="28"><strong>Connect edge</strong><small>name.play.minekube.net</small></div>
<div class="a-link a-tunnel"><span>tunnel</span></div>
<div class="a-zone"><p class="a-zone-label">Your PC · behind your router, no port forwarding</p><div class="a-stack"><div class="a-node"><svg viewBox="0 0 24 24" aria-hidden="true" fill="none" stroke="currentColor" stroke-width="1.75" stroke-linecap="round" stroke-linejoin="round"><rect width="20" height="14" x="2" y="3" rx="2"/><path d="M8 21h8"/><path d="M12 17v4"/></svg><strong>Your server</strong><small>localhost:25565</small><em class="a-chip">Connect plugin</em></div></div></div>
</div>
</div>

Check out [About Tunnels](/guide/tunnels) for a technical explanation.

## Get started

Check out the [Quick Start guide](/guide/quick-start) to learn how to get started with Connect.

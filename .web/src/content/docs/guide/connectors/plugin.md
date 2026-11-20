---
title: "Java plugin connector"
description: "Install the Connect Java plugin on Paper, Spigot, Velocity or BungeeCord to put a Minecraft server on Minekube Connect with a public address."
---

Using the Connect Java Plugin as a Connector for your Minecraft server or proxy.
If you have a Minecraft Java server or proxy, this the most convenient Connector for you,
but it is not as capable as the [Gate Proxy Connector](gate.md) in terms of routing features and performance.

The Connect Plugin is a powerful multi-platform Minecraft plugin that tunnels
your players through the global [Connect Network](/guide/#the-connect-network) to your Minecraft server/proxy.

-> It supports PaperMC, BungeeCord and Velocity platforms.

<!--@include: ../includes/downloads.md-->

<!--@include: ../includes/joining.md-->

## Example Setups

Minekube Connect advances the way players connect and developers architect secure Minecraft servers and networks.
Let's take a look at some common example setups.

### #1 Example: Velocity

<div class="mk-arch-wrap not-content">
<div class="mk-arch" role="img" aria-label="Players join the Connect edge, which tunnels them to Velocity with the Connect plugin; Velocity forwards them to two Paper servers.">
<div class="a-node"><svg viewBox="0 0 24 24" aria-hidden="true" fill="none" stroke="currentColor" stroke-width="1.75" stroke-linecap="round" stroke-linejoin="round"><path d="M16 21v-2a4 4 0 0 0-4-4H6a4 4 0 0 0-4 4v2"/><circle cx="9" cy="7" r="4"/><path d="M22 21v-2a4 4 0 0 0-3-3.87"/><path d="M16 3.13a4 4 0 0 1 0 7.75"/></svg><strong>Players</strong><small>Java and Bedrock</small></div>
<div class="a-link"><span>join</span></div>
<div class="a-node a-edge"><img src="/minekube-logo.png" alt="" width="28" height="28"><strong>Connect edge</strong><small>name.play.minekube.net</small></div>
<div class="a-link a-tunnel"><span>tunnel</span></div>
<div class="a-zone"><p class="a-zone-label">Your machine · no open ports</p><div class="a-stack"><div class="a-node"><img src="/images/logos/velocity.png" alt="" width="28" height="28"><strong>Velocity</strong><small>localhost:25577</small><em class="a-chip">Connect plugin</em></div><div class="a-fan" aria-hidden="true"></div><div class="a-col"><div class="a-node sm"><img src="/images/logos/paper.png" alt="" width="28" height="28"><strong>Paper</strong></div><div class="a-node sm"><img src="/images/logos/paper.png" alt="" width="28" height="28"><strong>Paper</strong></div></div></div></div>
</div>
</div>

- We have Velocity in online mode running on `localhost:25577` and want to use Connect.
- We install the Connect plugin to Velocity's `plugins` folder.
- We choose a name for our endpoint in the Connect plugin config `plugins/connect/config.yml`.
- We start Velocity and Connect Plugin will automatically tunnel players from the Connect Network.
- Done! We can now join our Velocity server at `<endpoint>.play.minekube.net`. Online mode players from Connect Network
  can
  join thanks to Plugins's auth session injection mechanisms.

### #2 Example: Paper

<div class="mk-arch-wrap not-content">
<div class="mk-arch" role="img" aria-label="Players join the Connect edge, which tunnels them to a Paper server with the Connect plugin.">
<div class="a-node"><svg viewBox="0 0 24 24" aria-hidden="true" fill="none" stroke="currentColor" stroke-width="1.75" stroke-linecap="round" stroke-linejoin="round"><path d="M16 21v-2a4 4 0 0 0-4-4H6a4 4 0 0 0-4 4v2"/><circle cx="9" cy="7" r="4"/><path d="M22 21v-2a4 4 0 0 0-3-3.87"/><path d="M16 3.13a4 4 0 0 1 0 7.75"/></svg><strong>Players</strong><small>Java and Bedrock</small></div>
<div class="a-link"><span>join</span></div>
<div class="a-node a-edge"><img src="/minekube-logo.png" alt="" width="28" height="28"><strong>Connect edge</strong><small>name.play.minekube.net</small></div>
<div class="a-link a-tunnel"><span>tunnel</span></div>
<div class="a-zone"><p class="a-zone-label">Your machine · no open ports</p><div class="a-stack"><div class="a-node"><img src="/images/logos/paper.png" alt="" width="28" height="28"><strong>Paper</strong><small>localhost:25565</small><em class="a-chip">Connect plugin</em></div></div></div>
</div>
</div>

- We have Paper running on `localhost:25565` and want to use Connect.
- We install the Connect plugin to Paper's `plugins` folder.
- We choose a name for our endpoint in the Connect plugin config `plugins/connect/config.yml`.
- We start Paper and Connect Plugin will automatically tunnel players from the Connect Network.
- Done! We can now join our Paper server at `<endpoint>.play.minekube.net`. Online mode players from Connect Network can
  join thanks to Plugin's auth session injection mechanisms.

### #3 Example: Connect `passthrough` + AuthSession API

Passthrough mode and AuthSession API are <Badge type="caution">Coming soon</Badge>.

<div class="mk-arch-wrap not-content soon">
<div class="mk-arch" role="img" aria-label="Planned: players join the Connect edge, which passes them through Gate Lite to an online-mode backend that authenticates them with the AuthSession API.">
<div class="a-node"><svg viewBox="0 0 24 24" aria-hidden="true" fill="none" stroke="currentColor" stroke-width="1.75" stroke-linecap="round" stroke-linejoin="round"><path d="M16 21v-2a4 4 0 0 0-4-4H6a4 4 0 0 0-4 4v2"/><circle cx="9" cy="7" r="4"/><path d="M22 21v-2a4 4 0 0 0-3-3.87"/><path d="M16 3.13a4 4 0 0 1 0 7.75"/></svg><strong>Players</strong><small>Java and Bedrock</small></div>
<div class="a-link"><span>join</span></div>
<div class="a-node a-edge"><img src="/minekube-logo.png" alt="" width="28" height="28"><strong>Connect edge</strong><small>name.play.minekube.net</small></div>
<div class="a-link a-tunnel"><span>tunnel</span></div>
<div class="a-zone"><p class="a-zone-label">Coming soon</p><div class="a-stack"><div class="a-node"><img src="/minekube-logo.png" alt="" width="28" height="28"><strong>Gate Lite</strong><small>passthrough</small></div><div class="a-link a-inner"></div><div class="a-node"><img src="/images/logos/paper.png" alt="" width="28" height="28"><strong>Backend</strong><small>online mode</small><em class="a-chip">AuthSession API</em></div></div></div>
</div>
</div>

This is not required for the normal Connect Java Plugin setup shown above. The plugin integrates with Connect's current
auth/session handling on supported platforms.

Passthrough/AuthSession is planned for setups where the backend or proxy needs to run its own online-mode authentication
for Connect-routed players. That includes topologies such as `Connect -> Gate Lite -> online-mode backend`, which are
not supported today through Gate Lite configuration alone.

## Endpoint Token

The endpoint **name** is what you set in `plugins/connect/config.yml`. The endpoint **token** does
not belong there: the plugin reads it from `token.json` in the plugin's data directory (for example
`plugins/connect/token.json`) or from the `CONNECT_TOKEN` environment variable, which takes
precedence over the file. On first launch the plugin generates a token and saves it to `token.json`.

- Resetting the token in the Dashboard invalidates the previous token. Copy the new token
  byte-for-byte into `token.json` (a trailing space or newline is a different token) or set
  `CONNECT_TOKEN`, then restart the server or proxy.
- A Watch service authentication error (`the endpoint token ... does not match the token currently
  stored for this endpoint`) always reports this token mismatch, not a dashboard problem. If the
  endpoint name belongs to an organization, the token must have been created for this endpoint name
  inside that organization; a token from another organization, or for a different endpoint name,
  cannot take the name over. If the name is not yours, choose a different endpoint name.

## Bedrock Identity

Connect-managed Bedrock is handled at the Connect edge. The plugin can verify the Bedrock identity that Connect already
checked before forwarding the player to Paper, Velocity, or BungeeCord.

Current plugin releases supply compatible identity defaults automatically, including for older configuration files
that omit the legacy identity section. For the managed Connect service, install or update the plugin and keep its
defaults. No identity URLs, public keys, or capability settings are required. Explicit operator overrides are preserved.

The endpoint policy accepts Microsoft/Xbox-authenticated Bedrock players without a linked Java account by default.
Connector identity settings verify what the edge forwarded; they cannot change an edge rejection such as
`policy_linked_java_only`.

See the [Bedrock support guide](/guide/bedrock#bedrock-identity-enforcement) for the current defaults, rejection reasons,
and support checklist.

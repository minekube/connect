# Java Plugin - Connector Guide

Using the Connect Java Plugin as a Connector for your Minecraft server or proxy.
If you have a Minecraft Java server or proxy, this the most convenient Connector for you,
but it is not as capable as the [Gate Proxy Connector](gate.md) in terms of routing features and performance.

The Connect Plugin is a powerful multi-platform Minecraft plugin that tunnels
your players through the global [Connect Network](#the-connect-network) to your Minecraft server/proxy.

-> It supports PaperMC, BungeeCord and Velocity platforms.

<!--@include: ../includes/downloads.md-->

<!--@include: ../includes/joining.md-->

## Example Setups

Minekube Connect advances the way players connect and developers architect secure Minecraft servers and networks.
Let's take a look at some common example setups.

### #1 Example: Velocity

`Connect -> [ Velocity -> Papers ]`

- We have Velocity in online mode running on `localhost:25577` and want to use Connect.
- We install the Connect plugin to Velocity's `plugins` folder.
- We choose a name for our endpoint in the Connect plugin config `plugins/connect/config.yml`.
- We start Velocity and Connect Plugin will automatically tunnel players from the Connect Network.
- Done! We can now join our Velocity server at `<endpoint>.play.minekube.net`. Online mode players from Connect Network
  can
  join thanks to Plugins's auth session injection mechanisms.

### #2 Example: Paper

`Connect -> Paper`

- We have Paper running on `localhost:25565` and want to use Connect.
- We install the Connect plugin to Paper's `plugins` folder.
- We choose a name for our endpoint in the Connect plugin config `plugins/connect/config.yml`.
- We start Paper and Connect Plugin will automatically tunnel players from the Connect Network.
- Done! We can now join our Paper server at `<endpoint>.play.minekube.net`. Online mode players from Connect Network can
  join thanks to Plugin's auth session injection mechanisms.

### #3 Example: Connect `passthrough` + AuthSession API

Passthrough mode and AuthSession API are <VPBadge type="warning">Coming soon</VPBadge>.

This is not required for the normal Connect Java Plugin setup shown above. The plugin integrates with Connect's current
auth/session handling on supported platforms.

Passthrough/AuthSession is planned for setups where the backend or proxy needs to run its own online-mode authentication
for Connect-routed players. That includes topologies such as `Connect -> Gate Lite -> online-mode backend`, which are
not supported today through Gate Lite configuration alone.

## Bedrock Identity

Connect-managed Bedrock is handled at the Connect edge. The plugin can verify the Bedrock identity that Connect already
checked before forwarding the player to Paper, Velocity, or BungeeCord.

::: code-group
```yaml [plugins/connect/config.yml]
bedrock-identity:
  enforcement: warn
  metadata-url: "https://connect.minekube.com/.well-known/minekube-connect/bedrock-identity-keys.json"
  metadata-cache-seconds: 300
  public-key: ""
  public-keys: []
  expected-policy: "bedrock-xuid"
```
:::

Start with `enforcement: warn`, review connector logs for Bedrock joins, then switch to `require` when the expected joins
carry valid identities. See the [Bedrock support guide](/guide/bedrock#bedrock-identity-enforcement) for the rollout and
support checklist.

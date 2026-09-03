---
description: Supported Connect platforms and known proxy, login-plugin, modded-server, and Bedrock compatibility constraints.
---

# Compatibility Matrix

Connect works best when the Minecraft ingress path is simple and each layer has one responsibility. This matrix captures
the combinations that most often need extra care.

## Requirements

The Connect plugin requires **Minecraft 1.13 or newer** and **Java 17 or newer**:

- **Minecraft 1.13+** — the declared `api-version` floor in `plugin.yml`, and the version from which servers ship the
  Netty 4.1+ the plugin's bundled Netty (4.2.x, unrelocated) requires. Servers 1.8–1.11 ship Netty 4.0.x, where the
  packet listener injection fails with `AbstractMethodError` (missing `newChild` implementation) and the plugin
  disables itself.
- **Java 17+** — the plugin is compiled for Java 17; older server JVMs (Java 8/11, typical for Minecraft ≤1.12
  servers) cannot load it.

On Minecraft 1.12 and older, use [Gate standalone](/guide/connectors/gate) as the Connector instead — it includes
Connect tunneling and supports 1.8.8 backends.

## Platforms

| Platform | Connect support | Notes |
| --- | --- | --- |
| Paper and Spigot | Supported with the Connect plugin | Recommended for single-server setups. Keep direct backend access closed unless intentionally public. |
| Velocity | Supported with the Connect plugin | Use stable Velocity builds when possible. Snapshot builds may change internals that proxy plugins rely on. |
| BungeeCord | Supported with the Connect plugin | Make sure backend forwarding and online-mode behavior match your BungeeCord setup. |
| Standard Gate | Supported as a Connect connector | Connect-routed Bedrock is handled at the Connect edge. Direct Gate Bedrock needs standard Gate with `bedrock: true`. |
| Gate Lite | Can be used as a Java reverse-proxy connector behind Connect | Lite forwards ping and player authentication to the selected backend route. `Connect -> Gate Lite -> online-mode backend` is not supported today through AuthSession/passthrough. Use standard Gate or the Connect Java Plugin for online-mode Connect ingress. Gate Lite is not the direct Bedrock listener. |
| Sponge or Minestom | Not a primary Connect plugin target | Put a supported proxy or server in front, or test carefully before using in production. |

## Proxy and Login Plugins

| Component | Risk | Recommended support response |
| --- | --- | --- |
| Velocity snapshots | Medium | Ask for the exact Velocity build and Connect plugin version. Reproduce on a stable Velocity build if packet or login behavior changed. |
| FastLogin, AuthMe, NLogin, and similar login plugins | Medium to high | Ask whether the server is online-mode, offline-mode, or mixed. Confirm whether the player joined through Connect, TCPShield, direct proxy, or direct backend. Do not treat a Connect-managed Bedrock report as a reason to enable direct Gate Bedrock. See [Login and Auth Plugins](/guide/login-plugins) for the rule that decides whether a login plugin conflicts with Connect. |
| LibreLogin with premium autologin enabled | Version-dependent | Use Connect v0.13.1 or newer and keep premium autologin enabled. See [Login and Auth Plugins](/guide/login-plugins#librelogin-with-premium-autologin) for the default skin handling and optional profile settings. |
| PacketEvents-based Velocity plugins (including Sonar and some nLogin builds) | Product incompatibility | A tunneled player can reach post-login and then disconnect because PacketEvents expects Velocity's ordinary socket-channel state while Connect supplies a local tunnel channel. Removing only the auth plugin is not a durable fix when it depends on PacketEvents. Compare a direct Velocity join, collect the PacketEvents/plugin versions and exception, and follow [connect-java#141](https://github.com/minekube/connect-java/issues/141). |
| Gate Lite with an online-mode backend | High | AuthSession/passthrough for Lite backend routes is not available today. Do not suggest `-Dmojang.sessionserver=` or a Gate Lite config change as a working fix. Recommend standard Gate with Connect enabled or the Connect Java Plugin on the online-mode proxy/server. |
| MultiProxySync and profile/skin sync plugins | Medium | Check whether UUID/profile data is expected from the proxy, backend, or plugin. Compare behavior between the direct proxy path and the Connect path. |
| TCPShield in front of Java while Connect handles Bedrock | Medium | Treat as two ingress paths. Ask which hostname the player used and whether forwarding is configured consistently on both paths. |
| Backend direct public access | High | Players can bypass forwarding and authentication assumptions. Recommend closing direct backend access or documenting it as a separate path. |

## Modded and Hybrid Servers

| Server type | Risk | Notes |
| --- | --- | --- |
| Vanilla-compatible Paper or Spigot | Low | Best supported path. |
| Forge or NeoForge behind a supported proxy | Medium | Test login and plugin-message behavior. Some modded handshakes assume a direct client-to-server path. |
| Arclight/Ketting/Forge hybrids, Mohist, Magma, and similar servers | High | Hybrid server internals vary. Ask for exact server type/version and logs before recommending a Connect or proxy setting. |
| Fabric servers with FabricProxy-Lite | Medium | Verify the proxy forwarding mod expected by the backend is installed, enabled, and configured with the same forwarding secret as the proxy. |
| Fabric servers with CrossStitch | Medium | CrossStitch is required when the Fabric backend needs Velocity modern forwarding support that FabricProxy-Lite alone does not provide. It is recommended when users report login/profile mismatch with Velocity-style forwarding. Verify by checking the backend log for the forwarding mod loading and by confirming direct Velocity-to-backend login works before adding Connect. |
| Polymer and server-side mod stacks | Medium | Usually compatible when the server remains vanilla-protocol compatible. Ask whether the issue also occurs without Connect and collect the mod list when resource-pack, profile, or login behavior differs. |
| NeoForge 1.21.x / Proxy-Compatible-Forge through Connect | Product investigation | Gate v0.69.1 fixed the open-bundle transition defect that was proven, and the Connect edge consumes that fix with regression coverage. A complete end-to-end NeoForge path has not yet been proven, and later failures may have a different cause. Preserve the join address, Connect plugin version, proxy type/version, NeoForge version, Proxy-Compatible-Forge version, exact kick text, and status/login/configuration-phase logs, then follow [connect#111](https://github.com/minekube/connect/issues/111). |
| Modpacks requiring custom client handshakes | Medium to high | Confirm whether the proxy and backend both support the required handshake. |

Use a product bug when a NeoForge or Forge-compatible path works through the same proxy without Connect but still fails
only through the current Connect edge. The report should focus on whether Connect preserves the status, login,
plugin-message, and configuration-phase metadata the backend forwarding mod expects.

## Voice Chat

Connect does not currently tunnel the separate UDP transport used by proximity-voice plugins such as Simple Voice
Chat. Minecraft can still join through the Connect hostname, but voice packets cannot use that hostname because the
Connect edge has no UDP voice listener or route to the backend yet. This limitation is tracked in
[connect#159](https://github.com/minekube/connect/issues/159).

Server owners can use one of these paths today:

- expose the voice plugin's UDP port (Simple Voice Chat defaults to `24454`) on the backend and configure its
  `voice_host` to that public address; Minecraft traffic can continue through Connect
- use a voice system that does not require a direct game UDP path, such as a browser-based voice plugin

Installing Simple Voice Chat beside the Connect plugin does not create a tunnel automatically. The SVC client mod is
still required, and opening local UDP `19132` does not help: that port is for Bedrock Minecraft traffic, not voice.

## Bedrock and Account Linking

Connect-routed Bedrock players use the same Connect endpoint names and custom domains as Java players. The Connect edge
handles Bedrock translation before traffic reaches your connector.

If a Bedrock player can ping but cannot join:

- confirm the player joined a Connect address, not a direct Gate address
- read the exact rejection reason; `policy_linked_java_only` is an explicit Connect edge override, while the default
  endpoint policy accepts Microsoft/Xbox-authenticated Bedrock identities without a linked Java account
- check whether the kick is from the Connect edge, the connector, the proxy, or the backend
- avoid recommending backend Geyser, Floodgate, UDP port opening, or `bedrock: true` unless the user wants direct
  Bedrock clients to hit their own standard Gate instance

For backend Floodgate API issues, use the [Bedrock support matrix](/guide/bedrock#support-behavior-matrix) and collect
plugin versions before recommending topology changes.

## What to Ask For

For compatibility reports, ask for the smallest useful set of facts:

- player address used to join
- Java or Bedrock client
- connector type and version
- backend/proxy type and version
- whether there is TCPShield, Velocity, BungeeCord, Gate, or direct backend access in the path
- exact kick text and logs around the timestamp
- whether the issue reproduces through another address

That information usually tells support whether the issue belongs to Connect ingress, proxy forwarding, backend auth,
modded handshake compatibility, or local server configuration.

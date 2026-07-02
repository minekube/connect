# Compatibility Matrix

Connect works best when the Minecraft ingress path is simple and each layer has one responsibility. This matrix captures
the combinations that most often need extra care.

## Platforms

| Platform | Connect support | Notes |
| --- | --- | --- |
| Paper and Spigot | Supported with the Connect plugin | Recommended for single-server setups. Keep direct backend access closed unless intentionally public. |
| Velocity | Supported with the Connect plugin | Use stable Velocity builds when possible. Snapshot builds may change internals that proxy plugins rely on. |
| BungeeCord | Supported with the Connect plugin | Make sure backend forwarding and online-mode behavior match your BungeeCord setup. |
| Standard Gate | Supported as a Connect connector | Connect-routed Bedrock is handled at the Connect edge. Direct Gate Bedrock needs standard Gate with `bedrock: true`. |
| Gate Lite | Can be used as a Java connector behind Connect | Gate Lite is not the direct Bedrock listener. Use standard Gate for direct self-hosted Bedrock. |
| Sponge or Minestom | Not a primary Connect plugin target | Put a supported proxy or server in front, or test carefully before using in production. |

## Proxy and Login Plugins

| Component | Risk | Recommended support response |
| --- | --- | --- |
| Velocity snapshots | Medium | Ask for the exact Velocity build and Connect plugin version. Reproduce on a stable Velocity build if packet or login behavior changed. |
| FastLogin, AuthMe, NLogin, and similar login plugins | Medium to high | Ask whether the server is online-mode, offline-mode, or mixed. Confirm whether the player joined through Connect, TCPShield, direct proxy, or direct backend. Do not treat a Connect-managed Bedrock report as a reason to enable direct Gate Bedrock. |
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
| NeoForge 1.21.x / Proxy-Compatible-Forge through Connect | Product investigation | Classify as a Connect compatibility investigation, not generic local Gate misconfiguration. Preserve the join address, Connect plugin version, proxy type/version, NeoForge version, Proxy-Compatible-Forge version, exact kick text, and logs from the status/login/configuration phases. Track product regression coverage in [connect#111](https://github.com/minekube/connect/issues/111). |
| Modpacks requiring custom client handshakes | Medium to high | Confirm whether the proxy and backend both support the required handshake. |

Use a product bug when a NeoForge or Forge-compatible path works through the same proxy without Connect but fails only
through Connect. The report should focus on whether Connect preserves the status, login, plugin-message, and
configuration-phase metadata the backend forwarding mod expects.

## Bedrock and Account Linking

Connect-routed Bedrock players use the same Connect endpoint names and custom domains as Java players. The Connect edge
handles Bedrock translation before traffic reaches your connector.

If a Bedrock player can ping but cannot join:

- confirm the player joined a Connect address, not a direct Gate address
- check whether the target server requires a linked Java account
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

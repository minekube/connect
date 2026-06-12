# Forwarding and Topologies

Connect can sit in front of many Minecraft server layouts. The important first step is to identify which address the
player uses and which component is responsible for authentication, forwarding, and Bedrock translation.

## Choose the Right Path

| Player address | Connector | Backend sees | Bedrock support | Use this when |
| --- | --- | --- | --- | --- |
| `<endpoint>.play.minekube.net`, `minekube.net`, or a Connect custom domain | Connect plugin on Paper, Spigot, Velocity, or BungeeCord | The server or proxy running the plugin | Managed by the Connect edge | You want the simplest managed Connect setup |
| `<endpoint>.play.minekube.net`, `minekube.net`, or a Connect custom domain | Standard Gate with `connect` enabled | Your Gate instance, then its configured Java backends | Managed by the Connect edge for Connect-routed players | You already operate Gate and want Connect ingress |
| A hostname or IP that points directly at your Gate instance | Standard Gate | Your Gate instance, then its configured Java backends | Gate handles it only when `bedrock: true` is enabled | You self-host Gate and want direct Bedrock clients outside Connect |
| A TCPShield Java address plus a Connect Bedrock address | TCPShield for Java, Connect for Bedrock | Your backend may see two separate ingress paths | Connect handles Bedrock only on the Connect address | You intentionally keep Java and Bedrock ingress separate |
| A public backend IP or port | None | Backend directly | Not managed by Connect | Development only, or a deliberate direct public setup |

If the player joined through a Connect address, do not diagnose it as direct Gate Bedrock or direct backend traffic.

## Paper or Spigot With the Connect Plugin

This is the smallest setup. Install the Connect plugin and let players join the endpoint name, `play.minekube.net`
subdomain, or custom domain attached in the dashboard.

Use this topology when:

- the backend is a single Paper or Spigot server
- you want Connect to manage public ingress
- Bedrock players should join through the same Connect address as Java players

Keep Paper `online-mode=true` unless you intentionally run an offline-mode server. If you allow offline-mode players,
use the Connect plugin's `allow-offline-mode-players` setting instead of opening a second direct backend path.

## Velocity or BungeeCord With the Connect Plugin

Install the Connect plugin on the proxy when the proxy is the public entry point for the network. The proxy then keeps
its normal forwarding relationship with backend servers.

Check these pieces first when players join but have the wrong UUID, skin, permissions, or login state:

- the Connect plugin is installed on the proxy, not only on one backend server
- backend servers are configured for the proxy forwarding mode you already use
- Velocity, BungeeCord, Paper, and login plugins agree on online-mode behavior
- players are not bypassing the proxy through a direct backend address

Connect does not replace Velocity or BungeeCord forwarding between your proxy and backend servers. It gives the proxy a
managed public ingress path.

## Standard Gate With Connect Enabled

Use this topology when Gate is your Java proxy and you want Connect to expose it through the Connect network. Connect
players reach the Connect edge first, then your Gate connector, then Gate routes to its configured backend servers.

For Connect-routed Bedrock players, translation is already handled before traffic reaches your Gate connector. You do
not need `bedrock: true` just because players use a Connect address.

Enable `bedrock: true` only when Bedrock clients also connect directly to the Gate address you host yourself.

## TCPShield and Connect Together

Some networks keep TCPShield for Java traffic while adding Connect for managed Bedrock or additional endpoint routing.
That can work, but treat it as two ingress paths:

- Java players using the TCPShield address follow the TCPShield and proxy forwarding path.
- Java or Bedrock players using a Connect address follow the Connect path.
- Account, UUID, skin, and IP-forwarding behavior can differ between those paths if the backend/proxy is not configured
  consistently.

When debugging, always ask which address the player used. A kick that happens through `minekube.net` may not reproduce
through a TCPShield hostname, and the reverse is also true.

## Direct Backend Connections

Direct backend connections are useful for local debugging but should not be mixed casually with Connect or proxy
ingress. If players can bypass the proxy and join a backend directly, the backend may see different online-mode,
forwarding, UUID, skin, permission, or IP information.

For production, prefer one public path per player group and make the backend reject unintended direct traffic.

## Quick Debug Checklist

1. Ask for the exact address the player joined.
2. Identify the connector: Paper/Spigot plugin, Velocity plugin, BungeeCord plugin, or Gate connector.
3. Identify whether the player is Java or Bedrock.
4. Identify whether Bedrock is Connect-routed or direct Gate Bedrock.
5. Check whether the backend is reachable directly and whether players might bypass the intended proxy.
6. For UUID, skin, permission, or login issues, check forwarding and login-plugin configuration before changing Connect.


# Bedrock Support

Connect lets Bedrock players join connected endpoints through the same public addresses Java players already use:

- `<endpoint>.play.minekube.net`
- custom domains attached to the endpoint
- endpoint routing through the Connect network

For Connect-routed players, Bedrock translation is handled by the Connect edge before traffic reaches your connector.
That means the usual Connect setup stays the same for Paper/Spigot, Velocity, BungeeCord, and Gate connectors.

## Support Behavior Matrix

| Setup | What to configure | Bedrock translation | Backend Floodgate API | Online-mode and account behavior | Custom domain coexistence |
| --- | --- | --- | --- | --- | --- |
| Connect-managed Bedrock with Paper/Velocity/Bungee connector | Install the Connect plugin normally. Do not enable backend Gate `bedrock: true` for Connect-managed Bedrock. | Connect edge | Supported only for the compatibility data Connect forwards. If a plugin expects a local Floodgate/Geyser runtime, collect logs and plugin names. | Java online-mode stays on the Java path. Bedrock players may be rejected by servers that require a linked Java account. | Supported. Bedrock and Java can use the same endpoint subdomain or custom domain. |
| Connect-managed Bedrock with standard Gate `connect` enabled | Enable the Gate Connect connector normally. Do not add local Bedrock settings only for Connect players. | Connect edge | Same Connect-managed compatibility behavior as plugin connectors. | Treat kicks as Connect-managed first, then Gate/backend auth. Do not switch the backend to offline mode unless the whole topology requires it. | Supported. Connect addresses and custom domains stay shared. |
| self-hosted Gate direct Bedrock | Add `bedrock: true` to the standard Gate instance that Bedrock clients join directly. | Your Gate instance | Local Gate/Floodgate behavior belongs to that Gate setup, not Connect-managed Bedrock. | Follow Gate's direct Bedrock and account-linking requirements. | Use a separate direct Bedrock hostname or port unless you intentionally route that hostname outside Connect. |
| Standard Gate with both Connect and direct Bedrock | Enable Connect for Connect addresses and `bedrock: true` only for the separate direct Bedrock listener. | Connect edge for Connect addresses, Gate for the direct Bedrock address | Diagnose by join address. Connect-managed joins use Connect compatibility data; direct joins use local Gate behavior. | Ask which address the player used before changing online-mode or linking settings. | Supported when DNS clearly separates Connect-routed names from direct Gate names. |
| Gate Lite behind Connect | Configure Gate Lite as a Java connector behind Connect. | Connect edge | Treat backend Floodgate API reports as Connect-managed compatibility reports. | Gate Lite is not the direct Bedrock authority; keep auth decisions in the backend/proxy path. | Supported for Connect-routed endpoint domains. |
| Gate Lite as a direct Bedrock listener | Use standard Gate instead. | Not supported in Gate Lite | Not supported. | Not supported. | Not supported. |

## Connect-Routed Bedrock

Use this path when players join through Connect addresses such as `<endpoint>.play.minekube.net`, `minekube.net`, or a
custom domain attached to your endpoint.

You do not need to:

- install Geyser on your backend
- install Floodgate on your backend
- open a local UDP Bedrock port
- set `bedrock: true` in your backend Gate just for Connect players

This applies to the Connect Java plugin and to Gate used as a Connect connector.

Do not enable backend Gate `bedrock: true` for Connect-managed Bedrock. That setting is for self-hosted Gate direct
Bedrock listeners and can create a second Bedrock path that support then has to diagnose separately.

## Direct Self-Hosted Gate Bedrock

Use this path when Bedrock players connect directly to a Gate address you operate yourself, outside the Connect network.
In that case, standard Gate can manage Bedrock locally:

```yaml
bedrock: true
```

With that setting, Gate starts and manages the Bedrock translation runtime for direct Bedrock clients. This local setting
does not change how Connect-routed Bedrock players are handled.

## Domains

Connect domains are shared by Java and Bedrock players. If a custom domain already routes Java players to your endpoint,
Bedrock players use the same domain after it is attached to the endpoint in the dashboard.

Use a separate direct Bedrock address only when you intentionally run local Bedrock support on your own standard Gate
instance.

## Account Linking and Online Mode

Connect-routed Bedrock players arrive through the managed Bedrock path and are represented through the compatibility
layer before they reach your connector. Servers that require a linked Java account may still reject Bedrock players who
have not linked their account. In that case, the player needs to link their Bedrock account with the Java account expected
by the target server.

If you are troubleshooting online-mode behavior, first identify whether the player joined through a Connect address or a
direct self-hosted Gate Bedrock address. The required configuration is different for those two paths.

Keep online-mode decisions tied to the whole Java forwarding topology. A Connect-managed Bedrock report is not, by
itself, a reason to disable online-mode or allow offline-mode players on the backend.

## Backend Floodgate API Compatibility

Some backend plugins query Floodgate-style player metadata to decide whether a player is from Bedrock, linked to Java,
or allowed past a login step. Connect-managed Bedrock forwards compatibility data for supported paths, but servers can
still fail when a plugin assumes it is talking to a local Floodgate/Geyser installation.

When Floodgate-dependent behavior fails, ask for:

- the exact join address the Bedrock player used
- connector type and version
- backend server type and version
- Floodgate, Geyser, auth, or login plugin names and versions
- the kick text and logs from Connect, the connector, proxy, and backend around the same timestamp
- whether the same player succeeds through a direct self-hosted Gate Bedrock address, if one exists

Do not recommend installing backend Geyser/Floodgate or enabling Gate `bedrock: true` unless the user wants to operate a
separate self-hosted Gate direct Bedrock path.

## Discord support response draft

> Connect-managed Bedrock is separate from self-hosted Gate Bedrock. If your player joins through
> `<endpoint>.play.minekube.net`, `minekube.net`, or a custom domain attached to the endpoint, do not enable backend Gate
> `bedrock: true` just for that player. Connect handles Bedrock translation at the edge. If a Floodgate-dependent plugin
> rejects the player, send the join hostname, connector type/version, backend type/version, Floodgate/Geyser/auth plugin
> versions, the exact kick text, and logs from the Connect connector/proxy/backend at the same timestamp. Only use
> `bedrock: true` when Bedrock clients connect directly to your own standard Gate instance outside Connect.

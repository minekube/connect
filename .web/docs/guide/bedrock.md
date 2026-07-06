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

For Bedrock clients, add the same Connect hostname as the server address and use the normal Bedrock port `19132`.
For example:

- **Server Address:** `<endpoint>.play.minekube.net`
- **Port:** `19132`

If you attached a custom domain to the endpoint, Bedrock players can use that custom domain with port `19132` too.
This is a player-facing Connect edge port. It does not mean you need to open UDP `19132` on your home network, VPS,
backend server, or connector host.

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
Java clients normally use the Java port `25565`; Bedrock clients normally use the Bedrock port `19132`.

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

Connect-managed Bedrock identity uses official Microsoft/Xbox Bedrock authentication at the Connect edge. It is not a
Minekube password system, it is not backend Geyser authentication, and it is not the same thing as allowing generic
offline-mode Java players. The connector can enforce the Bedrock identity that Connect already verified, while Java
online-mode players continue to use the normal Java session path.

## Bedrock Identity Enforcement

For Connect-routed Bedrock, the Connect edge verifies the player's Microsoft/Xbox Bedrock identity and signs a short-lived
identity envelope before forwarding the session to your connector. The Connect Java Plugin can then verify that envelope
before the player reaches the backend.

The signed identity contains the Bedrock XUID, username, policy, issue time, and expiry time. Newer connectors also bind
the identity to the endpoint and organization that the player joined. This prevents a signed identity captured for one
endpoint from being replayed against another endpoint in a different organization.

Configure the connector with `bedrock-identity` when you want the backend to reject unsigned or invalid Connect-managed
Bedrock sessions:

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

Use `warn` first. In warn mode, invalid or missing identities are logged but the player is not rejected. After logs show
that expected Bedrock joins carry valid identities, switch to `require`.

Use `metadata-url` for normal production rollouts. It lets the connector fetch the current Ed25519 verifier key and the
previous verifier keys that remain valid during key rotation. Static `public-key` and `public-keys` are useful for
self-hosted, staged, or emergency rollouts. The key metadata is cacheable; keep `metadata-cache-seconds` aligned with the
HTTP cache time served by the metadata endpoint.

The key metadata endpoint is public and contains verifier public keys only:

```text
/.well-known/minekube-connect/bedrock-identity-keys.json
```

It must never expose private keys, signed player identity envelopes, access tokens, or player-specific data.

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
- whether `bedrock-identity.enforcement` is `disabled`, `warn`, or `require`
- Bedrock identity warnings such as missing envelope, invalid signature, expired identity, replayed nonce, policy
  mismatch, endpoint mismatch, or organization mismatch

Do not recommend installing backend Geyser/Floodgate or enabling Gate `bedrock: true` unless the user wants to operate a
separate self-hosted Gate direct Bedrock path.

Do not ask users to paste signed Bedrock identity envelopes, private keys, account tokens, or full player profile dumps
into public support channels. The useful support data is the join address, connector version, endpoint ownership, relevant
timestamps, sanitized logs, and the exact rejection reason.

## Discord support response draft

> Connect-managed Bedrock is separate from self-hosted Gate Bedrock. If your player joins through
> `<endpoint>.play.minekube.net`, `minekube.net`, or a custom domain attached to the endpoint, Bedrock clients should use
> that hostname with port `19132`. Do not install backend Geyser/Floodgate, open local UDP `19132`, or enable backend Gate
> `bedrock: true` just for that player. Connect handles Bedrock translation at the edge. If a Floodgate-dependent plugin
> rejects the player, send the join hostname, connector type/version, backend type/version, Floodgate/Geyser/auth plugin
> versions, the exact kick text, and logs from the Connect connector/proxy/backend at the same timestamp. Only use
> `bedrock: true` when Bedrock clients connect directly to your own standard Gate instance outside Connect.
